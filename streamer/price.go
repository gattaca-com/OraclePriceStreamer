package streamer

import (
	"container/list"
	"fmt"
)

type Price struct {
	Price int64
	Slot uint64
	Symbol string
	Decimals uint
}

type AggPrice struct {
	price int64
	sourceIds string
	symbol string
	decimals uint
}

func (p *AggPrice) setSourceId( pythSlot int, chainlinkSlot int) {

	p.sourceIds = fmt.Sprintf("%d:%d", pythSlot, chainlinkSlot)
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