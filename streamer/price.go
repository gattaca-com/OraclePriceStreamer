package streamer

import (
	"container/list"
	"fmt"
)

type Price struct {
	price int64
	slot uint64
	symbol string
	decimals uint
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

func (fifo *PriceBuffer) Append(elt Price) {

	defer fifo.lock.Release()
	fifo.lock.Acquire()

	if fifo.buffer.Len() == int(fifo.size) {
		last := fifo.buffer.Back()
		fifo.buffer.Remove(last)
	}
	fifo.buffer.PushFront(elt)

}

func (fifo *PriceBuffer) IsValidPrice(price Price) bool {

	defer fifo.lock.Release()
	fifo.lock.Acquire()


	for elt := fifo.buffer.Front(); elt != nil; elt = elt.Next() {
		cachedPrice := elt.Value.(Price)

		if cachedPrice.price == price.price && cachedPrice.slot == price.slot && price.symbol == price.symbol{
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