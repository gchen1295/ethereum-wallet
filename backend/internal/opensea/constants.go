package piratesea

import (
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ORDER_BOOK_VERSION = "1"
	API_VERSION        = "1"

	ORDER_BOOK_PATH = fmt.Sprintf("wyvern/v%s", ORDER_BOOK_VERSION)
	API_PATH        = fmt.Sprintf("api/v%s", API_VERSION)

	API_BASE_URL_MAIN = &url.URL{
		Scheme: "https",
		Host:   "api.opensea.io",
	}
	API_BASE_URL_RINKEBY = &url.URL{
		Scheme: "https",
		Host:   "testnets-api.opensea.io",
	}

	OPENSEA_FEE_RECIPIENT                = common.HexToAddress("0x5b3256965e7c3cf26e11fcaf296dfc8807c01073")
	OPENSEA_ADDRESS                      = common.HexToAddress("0x7be8076f4ea4a4ad08075c2508e481d6c946d12b")
	RINKEBY_OPENSEA_CONTRACT_ADDRESS     = common.HexToAddress("0x5206e78b21Ce315ce284FB24cf05e0585A93B1d9")
	ENJIN_ADDRESS                        = common.HexToAddress("0xfaaFDc07907ff5120a76b34b731b278c38d6043C")
	CHEEZE_WIZARDS_GUILD_ADDRESS         = common.HexToAddress("0x0000000000000000000000000000000000000000")
	CHEEZE_WIZARDS_GUILD_RINKEBY_ADDRESS = common.HexToAddress("0x095731b672b76b00A0b5cb9D8258CD3F6E976cB2")
	DECENTRALAND_ESTATE_ADDRESS          = common.HexToAddress("0x095731b672b76b00A0b5cb9D8258CD3F6E976cB2")
	NULL_ADDRESS                         = common.HexToAddress("0x0000000000000000000000000000000000000000")

	OPENSEA_TIMESTAMP_FORMAT = "2006-01-02T15:04:05"

	INVERSE_BASIS_POINT                = big.NewInt(10000)
	DEFAULT_SELLER_FEE_BASIS_POINTS    = big.NewInt(250)
	OPENSEA_SELLER_BOUNTY_BASIS_POINTS = big.NewInt(250)
	DEFAULT_BUYER_FEE_BASIS_POINTS     = big.NewInt(0)
)
