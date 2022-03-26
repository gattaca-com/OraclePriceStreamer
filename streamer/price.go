package streamer

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type Price struct {
	Price int64
	Slot uint64
	Symbol string
	Decimals uint
}

func MarshallPrice(price *Price) ([]byte, error) {
	var toHash [32]byte

	binary.LittleEndian.PutUint64(toHash[:8], uint64(price.Price))
	binary.LittleEndian.PutUint64(toHash[8:8+8], uint64(price.Slot))
	binary.LittleEndian.PutUint16(toHash[8+8:8+8+2], uint16(price.Decimals))

	copy(toHash[8+8+2:], []byte(price.Symbol))

	return toHash[:], nil
}

func PriceToHash(price *Price) common.Hash {
	b, err := MarshallPrice(price)
	if err != nil {
		// gattaca TODO is this a good idea?
		return common.Hash{}
	}

	return common.BytesToHash(b)
}

func UnmarshallPrice(data []byte) (*Price, error) {

	priceVal := binary.LittleEndian.Uint64(data[:8])
	slotVal := binary.LittleEndian.Uint64(data[8 : 8+8])
	decimalVal := binary.LittleEndian.Uint16(data[8+8 : 8+8+2])
	symbol := string(data[8+8+2:])

	price := Price {
		Price:    int64(priceVal),
		Slot:     slotVal,
		Symbol:   symbol,
		Decimals: uint(decimalVal),
	}

	return &price, nil
}


func PricesToBytes(prices []*Price) ([] byte, error) {
	var byteArr []byte
	for _, price := range prices {

		priceBytes, err := MarshallPrice(price)
		if err != nil {
			byteArr = append(byteArr, priceBytes...)
		} else {
			return nil, err
		}
		
	}
	return byteArr, nil
}

func BytesToPrices(priceBytes []byte) ([]*Price, error) {

	byteStep := 32
	offset := 0

	var prices []*Price

	if len(priceBytes) % byteStep != 0 {
		return nil, fmt.Errorf("Malformed byte array")
	}

	for {

		if offset + byteStep > len(priceBytes) {
			break
		}

		price, err := UnmarshallPrice(priceBytes[offset:offset+byteStep])

		if  err != nil {
			return nil, err
		}
		prices = append(prices, price)

		offset =+ byteStep

	}

	return prices, nil

}


func (tx *Price) EncodeRLP(w io.Writer) error {

	// panic("GTC encoding price")

	return fmt.Errorf("GTC encoding price")
	// if tx.Type() == LegacyTxType {
	// 	return rlp.Encode(w, tx.inner)
	// }
	// // It's an EIP-2718 typed TX envelope.
	// buf := encodeBufferPool.Get().(*bytes.Buffer)
	// defer encodeBufferPool.Put(buf)
	// buf.Reset()
	// if err := tx.encodeTyped(buf); err != nil {
	// 	return err
	// }
	// return rlp.Encode(w, buf.Bytes())
}



// // MarshalBinary returns the canonical encoding of the transaction.
// // For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// // transactions, it returns the type and payload.
// func (tx *Transaction) MarshalBinary() ([]byte, error) {
// 	if tx.Type() == LegacyTxType {
// 		return rlp.EncodeToBytes(tx.inner)
// 	}
// 	var buf bytes.Buffer
// 	err := tx.encodeTyped(&buf)
// 	return buf.Bytes(), err
// }

// DecodeRLP implements rlp.Decoder
func (px *Price) DecodeRLP(s *rlp.Stream) error {
	// panic("GTC decoding price")

	return fmt.Errorf("GTC decoding price")
}


type PriceBuffer struct {
	buffer *list.List
	size   uint64
	lock Lock
}

func NewPriceBuffer(size uint64) *PriceBuffer {

	fifo := &PriceBuffer{}
	fifo.buffer = list.New()
	fifo.size = size
	fifo.lock = NewLock()

	return fifo
}

func (buffer *PriceBuffer) GetLatest() *Price {
	defer buffer.lock.Release()
	buffer.lock.Acquire()

	latest := buffer.buffer.Front().Value.(Price)
	return &latest
}

func (fifo *PriceBuffer) Append(elt Price) {

	defer fifo.lock.Release()
	fifo.lock.Acquire()

	if fifo.buffer.Len() == int(fifo.size) {
		last := fifo.buffer.Back()
		fifo.buffer.Remove(last)
	}
	fifo.buffer.PushFront(elt)

}

func (fifo *PriceBuffer) IsValidPrice(price *Price) bool {

	defer fifo.lock.Release()
	fifo.lock.Acquire()


	for elt := fifo.buffer.Front(); elt != nil; elt = elt.Next() {
		cachedPrice := elt.Value.(Price)

		if cachedPrice.Price == price.Price && cachedPrice.Slot == price.Slot && price.Symbol == price.Symbol{
			return true
		}
	}
	return false
}

func (fifo *PriceBuffer) Len() int {
	defer fifo.lock.Release()
	fifo.lock.Acquire()
	return fifo.buffer.Len()
}


type Lock struct {
	lock chan interface{}
}

func NewLock() Lock {
	return Lock{
		lock: make(chan interface{}, 1),
	}
}

func (l *Lock) Acquire() {
	l.lock <- 1
}

func (l *Lock) Release() {
	<-l.lock
}