package goid

import (
	"crypto/rand"
	"math/big"
	"sync/atomic"
	"time"
)

func NewID2() *ID2 {
	return &ID2{
		delta: 1,
	}
}

var _id2 = NewID2()

func GetID2() *ID2 {
	return _id2
}

func GenID2() int64 {
	return _id2.Generate()
}

func ResolveID2(id int64, oid *ID2) (timestamp int64, counter uint32) {
	return id >> 20, uint32(id) & uint32((1<<(20-oid.nodeBits))-1)
}

type ID2 struct {
	id          int64
	delta       uint32
	randomDelta uint32
	node        uint32
	nodeBits    uint8
}

func (i *ID2) Generate() int64 {
	for {
		old := atomic.LoadInt64(&i.id)
		nt := time.Now().Unix()
		lt := (old >> 20) & ((1 << 33) - 1)
		cBits := 20 - i.nodeBits
		mask := uint32((1 << cBits) - 1)
		ct := uint32(old) & mask
		if nt < lt {
			time.Sleep(time.Millisecond)
			continue
		}
		if nt == lt {
			ct += i.getDelta()
			if ct > mask {
				time.Sleep(10 * time.Millisecond)
				continue
			}
		} else {
			ct = i.getDelta()
		}

		now := (nt << 20) | int64(ct)
		if i.nodeBits > 0 {
			now |= int64(i.node) << cBits
		}
		if atomic.CompareAndSwapInt64(&i.id, old, now) {
			return now
		}
	}
}

func (i *ID2) getDelta() uint32 {
	if i.randomDelta > 0 {
		if de, err := rand.Int(rand.Reader, big.NewInt(int64(i.randomDelta))); err == nil {
			return uint32(de.Int64())
		}
	}
	return i.delta
}

func (i *ID2) SetDelta(d uint32) {
	if d == 0 || d >= (1<<(20-i.nodeBits)-1) {
		panic("delta too large or invalid")
	}
	i.delta = d
}

func (i *ID2) GetDelta() uint32 {
	return i.delta
}

func (i *ID2) SetRandomDelta(r uint32) {
	if r == 0 || r >= (1<<(20-i.nodeBits)-1) {
		panic("random delta too large or invalid")
	}
	i.randomDelta = r
}

func (i *ID2) GetRandomDelta() uint32 {
	return i.randomDelta
}

func (i *ID2) SetNode(node uint32, nodeBits uint8) {
	if nodeBits < 2 || nodeBits > 18 ||
		node == 0 || node > (1<<nodeBits-1) ||
		i.delta >= (1<<(20-nodeBits)-1) ||
		i.randomDelta >= (1<<(20-nodeBits)-1) {
		panic("node or nodeBits is invalid")
	}
	i.node, i.nodeBits = node, nodeBits
}

func (i *ID2) GetNode() (node uint32, nodeBits uint8) {
	return i.node, i.nodeBits
}