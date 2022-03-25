package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gattaca-pyth/streamer"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"go.blockdaemon.com/pyth"

)


var testRPC = "https://api.devnet.solana.com"
var testWS = "wss://api.devnet.solana.com"


func main() {
	fmt.Println("Hello from Gattaca")

	streamPythPrices()

	
}

func streamPythPrices() {

	// {
	// 	"first_price": "FVb5h1VmHPfVb1RfqZckchq18GxRv4iKt8T4eVTQAqdz",
	// 	"attrs": {
	// 	  "asset_type": "Crypto",
	// 	  "base": "AVAX",
	// 	  "description": "AVAX/USD",
	// 	  "generic_symbol": "AVAXUSD",
	// 	  "quote_currency": "USD",
	// 	  "symbol": "Crypto.AVAX/USD"
	// 	},
	// 	"pubkey": "DDdPuysfkxPq5Y1ZtTSk1H5n7iBKc9wtEKUwd1TNu3Gc",
	// 	"slot": 123458988
	//   }

	avaxKey := solana.MustPublicKeyFromBase58("DDdPuysfkxPq5Y1ZtTSk1H5n7iBKc9wtEKUwd1TNu3Gc")
	products := make(map[solana.PublicKey]streamer.PythProduct)
	products[avaxKey] = streamer.PythProduct{
		Symbol: "AVAX/USD",
		Key: avaxKey,
	}
	
	

	myStreamer := streamer.NewPythStreamer(products, testRPC, testWS)

	myStreamer.StreamProducts()

}



func PrintAllPythProductAccounts() {
	
	client := pyth .NewClient(pyth.Devnet, testRPC, testWS)
	products, _ := client.GetAllProductAccounts(context.TODO(), rpc.CommitmentProcessed)
	
	for _, product := range products{
		jsonData, _ := json.MarshalIndent(product, "", "  ")
	fmt.Println(string(jsonData))
	}

}

