package goid

import (
	"crypto/rand"
	"github.com/beevik/ntp"
	"math/big"
	"sync/atomic"
	"time"
)

func NewID() *ID {
	return &ID{
		delta:            1,
		maxBacktrackWait: 3 * time.Second,
	}
}

var _id = NewID()

func GetID() *ID {
	return _id
}

func GenID() int64 {
	return _id.Generate()
}

func ResolveID(id int64, oid *ID) (timestamp int64, counter uint32) {
	return id >> 21, uint32(id) & uint32((1<<(21-oid.nodeBits))-1)
}

type ID struct {
	id               int64
	delta            uint32
	randomDelta      uint32
	node             uint32
	nodeBits         uint8
	maxBacktrackWait time.Duration
	ntpServer        string
}

func (i *ID) Generate() int64 {
	for {
		old := atomic.LoadInt64(&i.id)
		nt := uint32(time.Now().Unix())
		lt := uint32(old >> 21)
		cBits := 21 - i.nodeBits
		mask := uint32((1 << cBits) - 1)
		ct := uint32(old) & mask
		if nt < lt {
			if time.Duration(lt-nt)*time.Second <= i.maxBacktrackWait {
				time.Sleep(time.Millisecond)
				continue
			}
			if i.ntpServer == "" {
				panic("ntp server not set")
			}
			ntTime, err := ntp.Time(i.ntpServer)
			if err != nil {
				panic(err)
			}
			nt = uint32(ntTime.Unix())
			if nt < lt {
				panic("ntp time error")
			}
		}
		if nt == lt {
			ct += i.getDelta()
			if ct > mask {
				time.Sleep(time.Millisecond)
				continue
			}
		} else {
			ct = i.getDelta()
		}

		now := (int64(nt) << 21) | int64(ct)
		if i.nodeBits > 0 {
			now |= int64(i.node) << cBits
		}
		if atomic.CompareAndSwapInt64(&i.id, old, now) {
			return now
		}
	}
}

func (i *ID) getDelta() uint32 {
	if i.randomDelta > 0 {
		if de, err := rand.Int(rand.Reader, big.NewInt(int64(i.randomDelta))); err == nil {
			return uint32(de.Int64() + 1)
		}
	}
	return i.delta
}

func (i *ID) SetDelta(d uint32) {
	if d == 0 || d >= (1<<(21-i.nodeBits)-1) {
		panic("delta too large or invalid")
	}
	i.delta = d
}

func (i *ID) GetDelta() uint32 {
	return i.delta
}

func (i *ID) SetRandomDelta(r uint32) {
	if r == 0 || r >= (1<<(21-i.nodeBits)-1) {
		panic("random delta too large or invalid")
	}
	i.randomDelta = r
}

func (i *ID) GetRandomDelta() uint32 {
	return i.randomDelta
}

func (i *ID) SetNode(node uint32, nodeBits uint8) {
	if nodeBits < 2 || nodeBits > 19 ||
		node > (1<<nodeBits-1) ||
		i.delta >= (1<<(21-nodeBits)-1) ||
		i.randomDelta >= (1<<(21-nodeBits)-1) {
		panic("node or nodeBits is invalid")
	}
	i.node, i.nodeBits = node, nodeBits
}

func (i *ID) GetNode() (node uint32, nodeBits uint8) {
	return i.node, i.nodeBits
}

func (i *ID) SetMaxBacktrackWait(d time.Duration) {
	if d < 0 {
		panic("invalid maxBacktrackWait")
	}
	i.maxBacktrackWait = d
}

func (i *ID) GetMaxBacktrackWait() time.Duration {
	return i.maxBacktrackWait
}

func (i *ID) SetNTPServer(s string) {
	i.ntpServer = s
}

func (i *ID) GetNTPServer() string {
	return i.ntpServer
}
