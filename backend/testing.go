package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"

	"nft-engine/internal/eclient"
	"nft-engine/internal/ui"
	"nft-engine/pkg/utils"

	"github.com/ethereum/go-ethereum/common"
)

var (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func Log(msg ...interface{}) {
	fmt.Print(colorGreen)
	log.Println(fmt.Sprintln(msg...), colorReset)
}

var (
	etherscanToken = "WFVX7JFYSNW3D636591ZX83G1S55H8RT12"
	relayLink      = "wss://mainnet.infura.io/ws/v3/abbd89cc93b24947b0f96a096fbdffd9"
	testAddress    = common.HexToAddress("0x5a934DD1967C9BB3a6Ff5A329de373A9dB920607")
)

func main() {
	bot, err := eclient.NewClient(etherscanToken, relayLink)
	if err != nil {
		log.Fatal(err)
	}

	form := ui.TaskForm{}
	err = ui.ContractAddressPrompt(&form)
	if err != nil {
		log.Fatal(err)
	}

	contract, err := bot.PullContractABI(form.ContractAddress)
	if err != nil {
		log.Fatal(err)
	}

	err = ui.SelectMintFunctionPrompt(contract, &form)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Contract Address: ", form.ContractAddress.String())
	log.Println("Mint Function: ", form.MintFunction.String())

	// convert to gwei
	payableAmt := big.NewFloat(0)
	payableAmt, ok := payableAmt.SetString(form.PayableAmount)
	if !ok {
		log.Fatal("Failed to unmarshal payable amount")
	}

	payableGwei := utils.EthToGwei(payableAmt)

	log.Println("Payable: ", payableGwei.String())

	args := []interface{}{}
	for _, arg := range form.Args {
		log.Println(fmt.Sprintf("%s (%s): %s", arg.Name, arg.Type.String(), arg.Value))
		if arg.Argument.Type.String() == "uint256" {
			val, ok := big.NewInt(0).SetString(arg.Value, 10)
			if ok {
				args = append(args, val)
				continue
			}
		}
		args = append(args, arg.Value)
	}

	results, err := bot.QueryContract(contract, form.MintFunction, &testAddress, payableGwei, args...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(results...)
}

func testing() {
	//	sk := flag.String("secret", "", "Very secret stuff")
	flag.Parse()

	contractAdd := common.HexToAddress("0xa5c0bd78d1667c13bfb403e2a3336871396713c5")

	bot, err := eclient.NewClient(etherscanToken, relayLink)
	if err != nil {
		log.Fatal(err)
	}

	contract, err := bot.PullContractABI(&contractAdd)
	if err != nil {
		log.Fatal(err)
	}

	for _, method := range contract.ABI.Methods {
		Log(method.String())
		Log("--------------------------------------")
	}

	saleState, err := bot.Read(contract, "saleActive", &testAddress)
	if err != nil {
		log.Fatal(err)
	}
	Log(saleState...)
	// log.Println(est)
}
