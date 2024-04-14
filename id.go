package goid

import (
	"crypto/rand"
	"math/big"
	"sync/atomic"
	"time"
)

func NewID() *ID {
	return &ID{
		delta: 1,
	}
}

var _id = NewID()

func GetID() *ID {
	return _id
}

func GenID() int64 {
	return _id.Generate()
}

type ID struct {
	id          int64
	randomDelta uint32
	delta       uint32
}

func (i *ID) Generate() int64 {
	for {
		old := atomic.LoadInt64(&i.id)
		nt := uint32(time.Now().Unix())
		lt := uint32(old >> 21)
		ct := uint32(old) & 0x1FFFFF
		if nt < lt {
			time.Sleep(time.Millisecond)
			continue
		}
		if nt == lt {
			ct += i.getDelta()
			if ct >= 0x1FFFFF {
				time.Sleep(time.Millisecond)
				continue
			}
		} else {
			ct = 1
		}

		now := (int64(nt) << 21) | (int64(ct))
		if atomic.CompareAndSwapInt64(&i.id, old, now) {
			return now
		}
	}
}

func (i *ID) getDelta() uint32 {
	if i.randomDelta > 0 {
		if de, err := rand.Int(rand.Reader, big.NewInt(int64(i.randomDelta))); err == nil {
			return uint32(de.Int64())
		}
	}
	return i.delta
}

func (i *ID) SetDelta(d uint32) {
	if d == 0 || d >= 0x1FFFFF {
		panic("delta too large or invalid")
	}
	i.delta = d
}

func (i *ID) GetDelta() uint32 {
	return i.delta
}

func (i *ID) SetRandomDelta(r uint32) {
	if r == 0 || r >= 0x1FFFFF {
		panic("random delta too large or invalid")
	}
	i.randomDelta = r
}

func (i *ID) GetRandomDelta() uint32 {
	return i.randomDelta
}
