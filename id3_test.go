package goid

import (
	"sync"
	"testing"
	"time"
)

func TestID3_Generate_duplicate(t *testing.T) {
	ll := 20000000
	idArr := make(map[int64]struct{}, ll)
	idChan := make(chan int64, ll)
	i := 0
	id := NewID3()
	wg := sync.WaitGroup{}
	wg.Add(ll + 1)
	go func() {
		defer wg.Done()
		for {
			idV := <-idChan
			idArr[idV] = struct{}{}
			i++
			if i == ll {
				break
			}
		}
	}()
	tt := time.Now()
	for j := 0; j < ll; j++ {
		go func() {
			defer wg.Done()
			idChan <- id.Generate()
		}()
	}
	wg.Wait()
	t.Logf("generate %v id cost %d ms", ll, time.Now().Sub(tt).Milliseconds())
	if len(idArr) != ll {
		t.Errorf("Duplicate ID generated, want %d, got %d", ll, len(idArr))
	}
}

func TestID3_Generate_increment(t *testing.T) {
	ll := 20000000
	var latestID int64
	id := NewID3()
	for i := 0; i < ll; i++ {
		idV := id.Generate()
		if idV > latestID {
			latestID = idV
			continue
		}
		t.Errorf("id (%d) <= latestID (%d) ", idV, latestID)
	}

}

func TestID3_Generate_43_duplicate(t *testing.T) {
	ll := 20000000
	idArr := make(map[int64]struct{}, ll)
	idChan := make(chan int64, ll)
	i := 0
	id := NewID3()
	id.SetBits(43)
	wg := sync.WaitGroup{}
	wg.Add(ll + 1)
	go func() {
		defer wg.Done()
		for {
			idV := <-idChan
			idArr[idV] = struct{}{}
			i++
			if i == ll {
				break
			}
		}
	}()
	tt := time.Now()
	for j := 0; j < ll; j++ {
		go func() {
			defer wg.Done()
			idChan <- id.Generate()
		}()
	}
	wg.Wait()
	t.Logf("generate %v id cost %d ms", ll, time.Now().Sub(tt).Milliseconds())
	if len(idArr) != ll {
		t.Errorf("Duplicate ID generated, want %d, got %d", ll, len(idArr))
	}
}

func TestID3_Generate_43_increment(t *testing.T) {
	ll := 20000000
	var latestID int64
	id := NewID3()
	id.SetBits(43)
	for i := 0; i < ll; i++ {
		idV := id.Generate()
		if idV > latestID {
			latestID = idV
			continue
		}
		t.Errorf("id (%d) <= latestID (%d) ", idV, latestID)
	}

}

func TestID3_SetDelta(t *testing.T) {
	id := NewID3()
	delta := uint16(1 << 5)
	id.SetDelta(delta)
	var lt int64
	var lc uint16
	for i := 0; i < 100000; i++ {
		idV := id.Generate()
		idt, idc := ResolveID3(idV, id)

		if lt > idt {
			t.Errorf("idt (%d) < lt (%d)", idt, lt)
			break
		}
		if lt == idt && idc-lc != delta {
			t.Errorf("idc-lc (%d) != delta (%d)", idc-lc, delta)
			break
		}
		if lt < idt {
			if idc != delta {
				t.Errorf("idc (%d) != delta (%d)", idc, delta)
				break
			}
			lt = idt
		}
		lc = idc
	}
}

func TestID3_SetRandomDelta(t *testing.T) {
	id := NewID3()
	delta := uint16(1 << 5)
	id.SetRandomDelta(delta)
	var lt int64
	var lc uint16

	for i := 0; i < 100000; i++ {
		idV := id.Generate()
		idt, idc := ResolveID3(idV, id)
		if lt > idt {
			t.Errorf("idt (%d) < lt (%d)", idt, lt)
			break
		}

		if lt == idt && idc-lc > delta {
			t.Errorf("idc-lc (%d) > delta (%d)", idc-lc, delta)
			break
		}

		if lt < idt {
			if idc > delta {
				t.Errorf("idc (%d) > delta (%d)", idc, delta)
				break
			}
			lt = idt
		}
		lc = idc
	}
}
