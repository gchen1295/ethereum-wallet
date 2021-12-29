package piratesea

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	minter "nft-engine/internal/deth"
	queries "nft-engine/internal/opensea/queries"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Order struct {
	ID big.Int
}

func getUintArgs(order *queries.Order, buyFeeParams *FeeParameters) ([9]*big.Int, [9]*big.Int, error) {
	basePrice, ok := big.NewInt(0).SetString(order.BasePrice, 10)
	if !ok {
		return [9]*big.Int{}, [9]*big.Int{}, errors.New("failed to parse base price")
	}

	salt := getSalt()
	if salt == nil {
		return [9]*big.Int{}, [9]*big.Int{}, errors.New("failed to generate a valid salt")
	}

	buyerUints := [9]*big.Int{
		buyFeeParams.TakerRelayFee,          // Maker Relay Fee (devSellerFee + OSFee) taken from seller
		buyFeeParams.MakerRelayFee,          // Taker Relay Fee
		buyFeeParams.TakerProtocolFee,       // Maker Protocol Fee
		buyFeeParams.MakerProtocolFee,       // Taker Protocol Fee
		basePrice,                           // Base Price  (wei)
		big.NewInt(0),                       // Extra
		big.NewInt(time.Now().Unix() - 120), // Buy list time
		big.NewInt(0),                       // Buy expiration time
		salt,                                //salt
	}

	takerRelayFee, ok := big.NewInt(0).SetString(order.TakerRelayerFee, 10)
	if !ok {
		return [9]*big.Int{}, [9]*big.Int{}, errors.New("failed to parse Taker Relay Fee")
	}

	makerRelayFee, ok := big.NewInt(0).SetString(order.MakerRelayerFee, 10)
	if !ok {
		return [9]*big.Int{}, [9]*big.Int{}, errors.New("failed to parse Maker Relay Fee")
	}

	makerProtocolFee, ok := big.NewInt(0).SetString(order.MakerProtocolFee, 10)
	if !ok {
		return [9]*big.Int{}, [9]*big.Int{}, errors.New("failed to parse Maker Protocol Fee")
	}

	takerProtocolFee, ok := big.NewInt(0).SetString(order.TakerProtocolFee, 10)
	if !ok {
		return [9]*big.Int{}, [9]*big.Int{}, nil
	}
	//sellerFee := big.NewInt(int64(order.Asset.AssetContract.SellerFeeBasisPoints))
	m, ok := big.NewInt(0).SetString(order.Salt, 10)
	if !ok {
		return [9]*big.Int{}, [9]*big.Int{}, errors.New("failed to parse order salt")
	}
	sellerUints := [9]*big.Int{
		makerRelayFee,                           // Maker Relay Fee (devSellerFee + OSFee) taken from seller
		takerRelayFee,                           // Taker Relay Fee
		makerProtocolFee,                        // Maker Protocol Fee
		takerProtocolFee,                        // Taker Protocol Fee
		basePrice,                               // Base Price  (wei)
		big.NewInt(0),                           // Extra
		big.NewInt(int64(order.ListingTime)),    // Sell List Time
		big.NewInt(int64(order.ExpirationTime)), // Sell End Time
		m,
	}

	return buyerUints, sellerUints, nil
}

func getAddressArgs(order *queries.Order, fromAddress common.Address) ([7]common.Address, [7]common.Address, error) {
	buyerAddr := [7]common.Address{
		// verify if buyer
		OPENSEA_ADDRESS, // Exchange
		fromAddress,     // Maker (Buyer)
		common.HexToAddress(order.Asset.Owner.Address),                    // Taker (Seller)
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // Recipient Fee
		common.HexToAddress(order.Metadata.Asset.Address),                 // Target (Asset Contract)
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // Static target
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // Payment token
	}
	sellerAddr := [7]common.Address{
		OPENSEA_ADDRESS, // Exchange
		// Verify if seller
		common.HexToAddress(order.Asset.Owner.Address),                    // Maker (Seller)
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // Taker (Auction Winner / Buyer)
		OPENSEA_FEE_RECIPIENT,                                             // Fee recipient
		common.HexToAddress(order.Metadata.Asset.Address),                 // Target (Asset Contract)
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // Static target
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // Payment token
	}

	return buyerAddr, sellerAddr, nil
}

func (c *Bot) AtomicMatch(key *ecdsa.PrivateKey, order *queries.Order) error {
	pubKey := key.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	sellerSide := order.Side
	buyerSide := (sellerSide + 1) % 2

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	sellerAddress := common.HexToAddress(order.Asset.Owner.Address)

	computedFees, err := computeFees(order, OrderSideBuy)
	if err != nil {
		return err
	}

	buyFeeParams, err := getBuyFeeParameters(INVERSE_BASIS_POINT, computedFees.TotalBuyerFeeBasisPoints, computedFees.TotalSellerFeeBasisPoints, order)
	if err != nil {
		return err
	}

	var replacementPattern, callData []byte
	var transferFn *abi.Method
	if order.Metadata.Schema == "ERC721" {
		replacementPattern, err = hexutil.Decode("0x00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		if err != nil {
			return err
		}
		for name, method := range c.Erc721Abi.Methods {
			if name == "transferFrom" {
				transferFn = &method
				break
			}
		}

		if transferFn != nil {
			tokenId, ok := big.NewInt(0).SetString(order.Metadata.Asset.ID, 10)
			if !ok {
				return errors.New("failed to match parse token ID")
			}

			callData, err = transferFn.Inputs.Pack(NULL_ADDRESS, fromAddress, tokenId)
			if err != nil {
				return err
			}
			callData = append(transferFn.ID, callData...)
		}
	} else {
		replacementPattern, err = hexutil.Decode("0x00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		if err != nil {
			return err
		}

		for name, method := range c.Erc1155Abi.Methods {
			if name == "safeTransferFrom" {
				transferFn = &method
				break
			}
		}
		if transferFn != nil {
			tokenId, ok := big.NewInt(0).SetString(order.Metadata.Asset.ID, 10)
			if !ok {
				return errors.New("failed to match parse token ID")
			}

			// #	Name	Type	Data
			// 0	_from	address	0xB819bf611D4051959BfFfB75b8F81e19127C3660
			// 1	_to	address	0x9A0cfD0fE1344FBAA01Cd5f8f40E5995a8f3789f
			// 2	_id	uint256	10392
			// 3	_value	uint256	1
			// 4	_data	bytes
			// TODO Check for quantity
			callData, err = transferFn.Inputs.Pack(sellerAddress, fromAddress, tokenId, big.NewInt(1))
			if err != nil {
				return errors.New("failed to packed call data for erc 1155 token")
			}
		}
	}

	if callData == nil {
		return errors.New("failed to pack call data")
	}

	buyerAddr, sellerAddr, err := getAddressArgs(order, fromAddress)
	if err != nil {
		return err
	}
	buyerUints, sellerUints, err := getUintArgs(order, buyFeeParams)
	if err != nil {
		return err
	}

	feeMethodsSidesKindsHowToCalls := [8]uint8{
		uint8(buyFeeParams.FeeMethod),
		uint8(buyerSide),
		uint8(order.SaleKind),
		uint8(order.HowToCall),
		uint8(order.FeeMethod),
		uint8(sellerSide),
		uint8(order.SaleKind),
		uint8(order.HowToCall),
	}

	sellCallData, err := hexutil.Decode(order.Calldata)
	if err != nil {
		return err
	}

	buyCallData := callData
	if err != nil {
		return err
	}

	var r, s, metadata [32]byte

	rssByte, err := hexutil.Decode(order.R)
	if err != nil {
		return err
	}
	copy(r[:], common.LeftPadBytes(rssByte, 32))

	rssByte, err = hexutil.Decode(order.S)
	if err != nil {
		return err
	}
	copy(s[:], common.LeftPadBytes(rssByte, 32))

	// order.Metadata.ReferrerAddress = "0x5c5321ae45550685308a405827575e3d6b4a84aa000000000000000000000000"
	if order.Metadata.ReferrerAddress == "" {
		order.Metadata.ReferrerAddress = "0x5c5321ae45550685308a405827575e3d6b4a84aa"
	}
	rssByte, err = hexutil.Decode(order.Metadata.ReferrerAddress)
	if err != nil {
		return err
	}
	copy(metadata[:], common.RightPadBytes(rssByte, 32))

	uints := [18]*big.Int{}
	tmp := []*big.Int{}
	tmp = append(tmp, buyerUints[:]...)
	tmp = append(tmp, sellerUints[:]...)
	copied := copy(uints[:], tmp)
	if copied == 0 {
		return errors.New("failed to create uints")
	}

	addr := [14]common.Address{}
	addrTmp := []common.Address{}
	addrTmp = append(addrTmp, buyerAddr[:]...)
	addrTmp = append(addrTmp, sellerAddr[:]...)
	copied = copy(addr[:], addrTmp)
	if copied == 0 {
		return errors.New("failed to create addr")
	}

	sellerReplacementPattern, err := hexutil.Decode(order.ReplacementPattern)
	if err != nil {
		return nil
	}

	empty, err := hexutil.Decode("0x")
	if err != nil {
		return err
	}

	// Validate that the sell order is valid
	valid, err := c.OpenSea.ValidateOrder(&bind.CallOpts{}, sellerAddr, sellerUints, uint8(order.FeeMethod), uint8(sellerSide), uint8(order.SaleKind), uint8(order.HowToCall), sellCallData, sellerReplacementPattern, empty, uint8(order.V), r, s)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("sell order is not valid")
	}

	txnData, err := c.wyvernAbi.Pack(
		"atomicMatch_",
		addr,
		uints,
		feeMethodsSidesKindsHowToCalls,
		buyCallData,
		sellCallData,
		replacementPattern,
		sellerReplacementPattern,
		empty,
		empty,
		[2]uint8{uint8(order.V), uint8(order.V)},
		[5][32]byte{r, s, r, s, metadata},
	)

	if err != nil {
		return err
	}

	ok, err = c.OpenSea.ValidateOrderParameters(&bind.CallOpts{}, buyerAddr, buyerUints, uint8(order.FeeMethod), uint8(buyerSide), uint8(order.SaleKind), uint8(order.HowToCall), buyCallData, replacementPattern, empty)
	if err != nil {
		return err
	}

	log.Println("Buy Order Params: ", ok)

	estimate, err := c.OpenSea.CalculateCurrentPrice(&bind.CallOpts{}, buyerAddr, buyerUints, uint8(order.FeeMethod), uint8(buyerSide), uint8(order.SaleKind), uint8(order.HowToCall), buyCallData, replacementPattern, empty)
	if err != nil {
		return err
	}

	basePrice, ok := big.NewInt(0).SetString(order.BasePrice, 10)
	if !ok {
		return errors.New("failed to parse base price")
	}

	if estimate.Cmp(basePrice) == 1 {
		buyerUints[4] = estimate
		sellerUints[4] = estimate
	}

	log.Println("Price Estimate: ", estimate)

	// Check if the call data can match
	ok, err = c.OpenSea.OrderCalldataCanMatch(&bind.CallOpts{}, buyCallData, replacementPattern, sellCallData, sellerReplacementPattern)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("calldata cannot match")
	}

	// Grab misc info we need for signing and sending a txn
	chainID, err := c.ethClient.NetworkID(context.Background())
	if err != nil {
		return err
	}

	nonce, err := c.ethClient.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return errors.New("failed to get account nonce")
	}

	gas, err := c.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	gasTipCap, err := c.ethClient.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}

	// Get our oracle estimate
	// TODO give an option to manually set gas overrides
	oracleEstimate, err := c.ethClient.GetGasEstimate()
	if err == nil {
		gas = big.NewInt(oracleEstimate.Fast)
	}

	log.Println("Asset Name: " + order.Asset.Name)
	log.Println("Asset ID: " + order.Asset.TokenID)
	log.Printf("URL: https://opensea.io/assets/%s/%s\n", order.Asset.AssetContract.Address, order.Asset.TokenID)

	log.Println("Suggested Gas: ", gas)
	log.Println("Suggest Gas Tip: ", gasTipCap)

	// Create a new eip-1155 txn
	txn := types.NewTx(&types.DynamicFeeTx{
		ChainID: chainID,
		Nonce:   nonce,
		Data:    txnData,
		Value:   basePrice,
		Gas:     210000,
		// GasPrice:  big.NewInt(130515038612),
		GasFeeCap: gas,
		GasTipCap: gasTipCap,
		To:        &OPENSEA_ADDRESS,
	})

	// Sign our txn with the lastest chain signer
	signedTx, err := types.SignTx(txn, types.LatestSignerForChainID(chainID), key)
	if err != nil {
		return err
	}

	// Simulate bundle on last block
	result, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		Data:      signedTx.Data(),
		Value:     basePrice,
		Gas:       210000,
		GasFeeCap: gas,
		GasTipCap: gasTipCap,
		To:        &OPENSEA_ADDRESS,
		From:      fromAddress,
	}, nil)
	if err != nil {
		return err
	}

	// Check results and adjust as needed
	log.Println(string(result))

	bundle := []*types.Transaction{
		signedTx,
	}

	log.Println(bundle)

	// err = c.ethClient.Client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// result, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
	// 	Data:  hd,
	// 	Value: basePrice,
	// 	Gas:   210000,
	// 	// GasPrice:  big.NewInt(130515038612),
	// 	GasFeeCap: gas,
	// 	GasTipCap: gasTipCap,
	// 	To:        &OPENSEA_ADDRESS,
	// 	From:      fromAddress,
	// }, nil)

	// log.Println("RESULT: " + string(result))
	// if err != nil {
	// 	log.Println(err)
	// 	//	return err
	// }
	return nil
}

type OrderSide int

var (
	OrderSideBuy  = OrderSide(0)
	OrderSideSell = OrderSide(1)
)

func computeFees(order *queries.Order, side OrderSide) (*ComputedFees, error) {
	openseaBuyerFeeBasisPoints := DEFAULT_BUYER_FEE_BASIS_POINTS
	openseaSellerFeeBasisPoints := DEFAULT_SELLER_FEE_BASIS_POINTS
	devBuyerFeeBasisPoints := big.NewInt(0)
	devSellerFeeBasisPoints := big.NewInt(0)
	transferFee := big.NewInt(0)
	transferFeeTokenAddress := common.Address{}
	maxTotalBountyBPS := DEFAULT_SELLER_FEE_BASIS_POINTS

	if order != nil {
		toAdd, ok := big.NewInt(0).SetString(order.Asset.Collection.OpenseaBuyerFeeBasisPoints, 10)
		if !ok {
			return nil, errors.New("failed to set OpenseaBuyerFeeBasisPoints")
		}
		openseaBuyerFeeBasisPoints = openseaBuyerFeeBasisPoints.Add(openseaBuyerFeeBasisPoints, toAdd)

		toAdd, ok = toAdd.SetString(order.Asset.Collection.OpenseaSellerFeeBasisPoints, 10)
		if !ok {
			return nil, errors.New("failed to set OpenseaSellerFeeBasisPoints")
		}
		openseaSellerFeeBasisPoints = openseaSellerFeeBasisPoints.Add(openseaSellerFeeBasisPoints, toAdd)

		toAdd, ok = toAdd.SetString(order.Asset.Collection.DevBuyerFeeBasisPoints, 10)
		if !ok {
			return nil, errors.New("failed to set DevBuyerFeeBasisPoints")
		}
		devBuyerFeeBasisPoints = devBuyerFeeBasisPoints.Add(devBuyerFeeBasisPoints, toAdd)

		toAdd, ok = toAdd.SetString(order.Asset.Collection.DevSellerFeeBasisPoints, 10)
		if !ok {
			return nil, errors.New("failed to set DevSellerFeeBasisPoints")
		}
		devSellerFeeBasisPoints = devSellerFeeBasisPoints.Add(devSellerFeeBasisPoints, toAdd)

		maxTotalBountyBPS = openseaSellerFeeBasisPoints
	}

	if side == OrderSideSell && order != nil {
		// TODO update transfer fee settings for sell orders
	}

	// TODO set to extraBountyBasisPoints for sell orders
	sellerBountyBasisPoints := big.NewInt(0)

	isBountyTooLarge := sellerBountyBasisPoints.Add(sellerBountyBasisPoints, OPENSEA_SELLER_BOUNTY_BASIS_POINTS).Cmp(maxTotalBountyBPS) == 1
	if sellerBountyBasisPoints.Cmp(big.NewInt(0)) == 1 && isBountyTooLarge {
		errMsg := fmt.Sprintf("total bounty exceeds the maximum for this asset type (%d%%.", maxTotalBountyBPS.Div(maxTotalBountyBPS, big.NewInt(100)))

		if maxTotalBountyBPS.Cmp(OPENSEA_SELLER_BOUNTY_BASIS_POINTS) != -1 {
			errMsg += fmt.Sprintf(" remember that OpenSea will add %d%% for referrers with OpenSea accounts!", OPENSEA_SELLER_BOUNTY_BASIS_POINTS.Div(OPENSEA_SELLER_BOUNTY_BASIS_POINTS, big.NewInt(100)))
		}

		return nil, fmt.Errorf("%s", errMsg)
	}

	return &ComputedFees{
		TotalBuyerFeeBasisPoints:    openseaBuyerFeeBasisPoints.Add(openseaBuyerFeeBasisPoints, devBuyerFeeBasisPoints),
		TotalSellerFeeBasisPoints:   openseaSellerFeeBasisPoints.Add(openseaSellerFeeBasisPoints, devSellerFeeBasisPoints),
		OpenseaBuyerFeeBasisPoints:  openseaBuyerFeeBasisPoints,
		OpenseaSellerFeeBasisPoints: openseaSellerFeeBasisPoints,
		DevBuyerFeeBasisPoints:      devBuyerFeeBasisPoints,
		DevSellerFeeBasisPoints:     devSellerFeeBasisPoints,
		TransferFee:                 transferFee,
		TransferFeeTokenAddress:     &transferFeeTokenAddress,
	}, nil
}

type ComputedFees struct {
	TotalBuyerFeeBasisPoints    *big.Int
	TotalSellerFeeBasisPoints   *big.Int
	OpenseaBuyerFeeBasisPoints  *big.Int
	OpenseaSellerFeeBasisPoints *big.Int
	DevBuyerFeeBasisPoints      *big.Int
	DevSellerFeeBasisPoints     *big.Int
	SellerBountyBasisPoints     *big.Int
	TransferFee                 *big.Int
	TransferFeeTokenAddress     *common.Address
}

func getSalt() *big.Int {
	max := new(big.Int)
	max.SetString("999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)
	min := big.NewInt(10)
	min.Exp(min, big.NewInt(77), nil)

	n, err := rand.Int(rand.Reader, max.Sub(max, min))
	if err != nil {
		return nil
	}
	n.Add(n, min)

	return n
}

func getTransferFeeSettings(client *minter.Client, assetAddress common.Address) (*big.Int, common.Address, error) {
	if assetAddress.Hash() == ENJIN_ADDRESS.Hash() {
		// transferSettings(uint256 tokenID, address from)

		// TODO get ENJIN asset transfer settings
	}

	return nil, common.Address{}, nil
}

func validateFees(inverseBasisPoints, totalBuyerFeeBasisPoints, totalSellerFeeBasisPoints *big.Int) error {
	maxFeePercent := inverseBasisPoints.Div(inverseBasisPoints, big.NewInt(100))

	if totalBuyerFeeBasisPoints.Cmp(maxFeePercent) == 1 || totalSellerFeeBasisPoints.Cmp(maxFeePercent) == 1 {
		return fmt.Errorf("invalid buyer/seller fees: must be less than %d%%", maxFeePercent)
	}

	if totalBuyerFeeBasisPoints.Cmp(big.NewInt(0)) == -1 || totalSellerFeeBasisPoints.Cmp(big.NewInt(0)) == -1 {
		return fmt.Errorf("invalid buyer/seller fees: must be at least 0%%")
	}

	return nil
}

func getBuyFeeParameters(inverseBasisPoints, totalBuyerFeeBasisPoints, totalSellerFeeBasisPoints *big.Int, order *queries.Order) (*FeeParameters, error) {
	err := validateFees(inverseBasisPoints, totalBuyerFeeBasisPoints, totalBuyerFeeBasisPoints)
	if err != nil {
		return nil, err
	}

	var makerRelayFee = big.NewInt(0)
	var takerRelayFee = big.NewInt(0)
	var ok bool
	if order != nil {
		if order.WaitingForBO != nil {
			if makerRelayFee, ok = makerRelayFee.SetString(order.MakerRelayerFee, 10); !ok {
				makerRelayFee = big.NewInt(0)
			}
			if takerRelayFee, ok = takerRelayFee.SetString(order.TakerRelayerFee, 10); !ok {
				takerRelayFee = big.NewInt(0)
			}
		} else {
			if makerRelayFee, ok = makerRelayFee.SetString(order.TakerRelayerFee, 10); !ok {
				makerRelayFee = big.NewInt(0)
			}
			if takerRelayFee, ok = takerRelayFee.SetString(order.MakerRelayerFee, 10); !ok {
				takerRelayFee = big.NewInt(0)
			}
		}
	} else {
		makerRelayFee = totalBuyerFeeBasisPoints
		takerRelayFee = totalSellerFeeBasisPoints
	}

	return &FeeParameters{
		MakerRelayFee:    makerRelayFee,
		TakerRelayFee:    takerRelayFee,
		MakerProtocolFee: big.NewInt(0),
		TakerProtocolFee: big.NewInt(0),
		MakerRefererFee:  big.NewInt(0),
		FeeRecipient:     OPENSEA_FEE_RECIPIENT,
		FeeMethod:        1,
	}, nil
}

type FeeParameters struct {
	MakerRelayFee    *big.Int
	TakerRelayFee    *big.Int
	MakerProtocolFee *big.Int
	TakerProtocolFee *big.Int
	MakerRefererFee  *big.Int
	FeeRecipient     common.Address
	FeeMethod        int
}

type PriceParameters struct {
	BasePrice    *big.Int
	Extra        *big.Int
	PaymentToken common.Address
	ReservePrice *big.Int
}

func getPriceParameters(side int, tokenAddress common.Address, expirationTime int, startAmount, endAmount *big.Int, waitingForBO bool, reservePrice *big.Int) (*PriceParameters, error) {
	priceDiff := big.NewInt(0)
	if endAmount != nil {
		priceDiff = startAmount.Sub(startAmount, endAmount)
	}

	isEther := tokenAddress.Hash().String() == common.HexToAddress("0x0000000000000000000000000000000000000000000000000000000000000000").Hash().String()

	if !isEther && waitingForBO {
		return nil, errors.New("english auctions must use wrapped ETH or an ERC-20 token")
	}

	// buy side
	if isEther && side == 0 {
		return nil, errors.New("offers must use wrapped ETH or an ERC-20 token")
	}

	if priceDiff.Cmp(big.NewInt(0)) == -1 {
		return nil, errors.New("end price must be less than or equal to the start price")
	}

	if priceDiff.Cmp(big.NewInt(0)) == 1 && expirationTime == 0 {
		return nil, errors.New("expiration time must be set if order will change in price")
	}

	if reservePrice != nil && !waitingForBO {
		return nil, errors.New("reserve prices may only be set on English auctions")
	}

	if reservePrice != nil {
		if reservePrice.Cmp(startAmount) == -1 {
			return nil, errors.New("reserve price must be greater than or equal to the start amount")
		}
	}

	return &PriceParameters{
		BasePrice:    startAmount,
		ReservePrice: reservePrice,
		PaymentToken: tokenAddress,
		Extra:        big.NewInt(0),
	}, nil

}

// 0x7Be8076f4EA4A4AD08075C2508e481d6C946D12b
// 0x37998495c09662E26b021Cb29c6B7859E97Cdc90
// 0x69f8D754C5f4F73aad00f3C22EaFB77Aa57Ff1BC
// 0x0000000000000000000000000000000000000000
// 0x338866F8ba75bb9D7a00502E11b099a2636C2C18
// 0x0000000000000000000000000000000000000000
// 0x0000000000000000000000000000000000000000

// 0x7Be8076f4EA4A4AD08075C2508e481d6C946D12b
// 0x69f8D754C5f4F73aad00f3C22EaFB77Aa57Ff1BC
// 0x0000000000000000000000000000000000000000
// 0x5b3256965e7C3cF26E11FCAf296DfC8807C01073
// 0x338866F8ba75bb9D7a00502E11b099a2636C2C18
// 0x0000000000000000000000000000000000000000
// 0x0000000000000000000000000000000000000000
