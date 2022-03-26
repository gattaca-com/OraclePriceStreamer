
package streamer

import (
	"fmt"

	"go.blockdaemon.com/pyth"
	"github.com/gagliardetto/solana-go"
)


type PythProduct struct{
	Key solana.PublicKey
	Symbol string
	Decimals uint
}

type PythStreamer struct{
	priceCache map[solana.PublicKey]PriceBuffer
	products map[solana.PublicKey]PythProduct
	rpcURL string
	wsURL string
}

func NewPythStreamer(products map[solana.PublicKey]PythProduct, rpcURL string, wsURL string) *PythStreamer {

	priceStreamer := &PythStreamer{
		priceCache: make(map[solana.PublicKey]PriceBuffer),
		products: products,
		rpcURL: rpcURL,
		wsURL: wsURL,
	}

	for key, _ := range products {
        priceStreamer.priceCache[key] = *NewPriceBuffer(1000)
    }

	return priceStreamer
}

func (streamer *PythStreamer) IsValidPrice(price *Price) bool {

	priceBuffer, err := streamer.GetPriceBuffer(price)

	if err != nil {
		return false
	}

	return priceBuffer.IsValidPrice(price)

}

func (streamer *PythStreamer) StreamProducts() {

	client := pyth.NewClient(pyth.Devnet, streamer.rpcURL, streamer.wsURL)
	stream := client.StreamPriceAccounts()
	

	// Print updates.
	for update := range stream.Updates() {

		if streamer.shouldDump(update.Product) {

			buffer := streamer.priceCache[update.Product]
			price := Price{
				Price: update.Agg.Price,
				Slot: update.Agg.PubSlot, 
				Symbol: streamer.products[update.Product].Symbol,
				Decimals: uint(update.Exponent),
			}

			buffer.Append(price)

			output := fmt.Sprintf("Symbol: %s, Price: %d, Product: %s Slot: %d", streamer.products[update.Product].Symbol, update.Agg.Price, update.Product, update.Agg.PubSlot)
			fmt.Println(output)
		}
		
	}

}

func (streamer *PythStreamer) GetPrices() []*Price {

	var prices []*Price

	for key, _ := range streamer.products {
        buffer := streamer.priceCache[key] 
		prices = append(prices, buffer.GetLatest())
    }


	return prices

}

func (streamer *PythStreamer) GetPricesBytes()([] byte, error) {
	return PricesToBytes(streamer.GetPrices())
}

func (streamer *PythStreamer) GetPriceBuffer(price *Price) (*PriceBuffer, error) {
	for key, _ := range streamer.products {

		if streamer.products[key].Symbol == price.Symbol {
			buffer := streamer.priceCache[key] 
			return &buffer, nil

		}
        
    }
	return nil, fmt.Errorf("failed to find price buffer for %s", price.Symbol)
}


func (streamer *PythStreamer)shouldDump(productKey solana.PublicKey) bool {

	if _, ok := streamer.products[productKey]; ok {
		return true
	}

    return false
}
