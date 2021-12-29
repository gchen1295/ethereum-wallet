package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	minter "nft-engine/internal/deth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
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

func Log(msg string) {
	fmt.Print(colorGreen)
	log.Println(msg + colorReset)
}

func main() {
	// sk := flag.String("secret", "", "Very secret stuff")
	// flag.Parse()

	// contractAdd := common.HexToAddress("0x262f103bb7ab38c6e4dcd848586dfebe0bde2959")
	// bot, err := minter.NewBot(minter.BotOptions{
	// 	ClientOptions: minter.ClientOptions{
	// 		EtherscanToken: "WFVX7JFYSNW3D636591ZX83G1S55H8RT12",
	// 		RelayLink:      "wss://mainnet.infura.io/ws/v3/abbd89cc93b24947b0f96a096fbdffd9",
	// 	},
	// 	Contract: &contractAdd,
	// }, context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	eventChan := make(chan core.NewTxsEvent)

	txPool := core.NewTxPool(core.DefaultTxPoolConfig, core.DefaultGenesisBlock().Config, &core.BlockChain{})
	sub := txPool.SubscribeNewTxsEvent(eventChan)

	defer sub.Unsubscribe()
	for i := range eventChan {
		// log.Println(fmt.Sprintf("%+v\n", i))

		log.Println("Number Txns: ", len(i.Txs))

		for _, tx := range i.Txs {
			eventJson, err := tx.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(eventJson))
		}
	}

}

func ostest() {

	sk := flag.String("secret", "", "Very secret stuff")
	flag.Parse()

	contractAdd := common.HexToAddress("0x262f103bb7ab38c6e4dcd848586dfebe0bde2959")
	bot, err := minter.NewBot(minter.BotOptions{
		ClientOptions: minter.ClientOptions{
			EtherscanToken: "WFVX7JFYSNW3D636591ZX83G1S55H8RT12",
			RelayLink:      "wss://mainnet.infura.io/ws/v3/abbd89cc93b24947b0f96a096fbdffd9",
		},
		Contract: &contractAdd,
	}, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	Log("Loading wallet")
	Log(*sk)

}
