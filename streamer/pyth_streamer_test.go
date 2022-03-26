package streamer

import (
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

var testRPC = "https://api.devnet.solana.com"
var testWS = "wss://api.devnet.solana.com"

func TestStreamerCanInitialize(t *testing.T) {
	require := require.New(t)

	avaxKey := solana.MustPublicKeyFromBase58("DDdPuysfkxPq5Y1ZtTSk1H5n7iBKc9wtEKUwd1TNu3Gc")
	products := make(map[solana.PublicKey]PythProduct)
	products[avaxKey] = PythProduct{
		Symbol: "AVAX/USD",
		Key:    avaxKey,
	}

	myStreamer := NewPythStreamer(products, testRPC, testWS)

	go myStreamer.StreamProducts()

	select {
	case <-myStreamer.Initialized:
	case <-time.After(6 * time.Second):
	}

	prices := myStreamer.GetPrices()

	require.Len(prices, 1, "Expected to receive back one price!")

	require.Greater(prices[0].Price, int64(0))
}
