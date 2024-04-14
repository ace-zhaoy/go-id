package goid

import (
	"sync"
	"testing"
)

func TestID2_Generate_duplicate(t *testing.T) {
	ll := 10000000
	idArr := make(map[int64]struct{}, ll)
	idChan := make(chan int64, ll)
	i := 0
	id := NewID2()
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
	for j := 0; j < ll; j++ {
		go func() {
			defer wg.Done()
			idChan <- id.Generate()
		}()
	}
	wg.Wait()
	if len(idArr) != ll {
		t.Errorf("Duplicate ID generated, want %d, got %d", ll, len(idArr))
	}
}

func TestID2_Generate_increment(t *testing.T) {
	ll := 10000000
	var latestID int64
	id := NewID2()
	for i := 0; i < ll; i++ {
		idV := id.Generate()
		if idV > latestID {
			latestID = idV
			continue
		}
		t.Errorf("id (%d) <= latestID (%d) ", idV, latestID)
	}
}

func TestID2_SetDelta(t *testing.T) {
	id := NewID2()
	delta := uint32(1 << 10)
	id.SetDelta(delta)
	var lt int64
	var lc uint32
	for i := 0; i < 100000; i++ {
		idV := id.Generate()
		idt, idc := ResolveID2(idV, id)
		if lt > idt {
			t.Errorf("idt (%d) < lt (%d)", idt, lt)
			break
		}
		if lt == idt && idc-lc != delta {
			t.Errorf("idc-lc (%d) != delta (%d)", idc-lc, delta)
			break
		}
		if lt < idt {
			if lt != 0 && idt != lt+1 {
				t.Errorf("idt (%d) != lt+1 (%d)", idt, lt+1)
				break
			}
			if idc != delta {
				t.Errorf("idc (%d) != delta (%d)", idc, delta)
				break
			}
			lt = idt
		}
		lc = idc
	}
}

func TestID2_SetRandomDelta(t *testing.T) {
	id := NewID2()
	delta := uint32(1 << 10)
	id.SetRandomDelta(delta)
	var lt int64
	var lc uint32

	for i := 0; i < 100000; i++ {
		idV := id.Generate()
		idt, idc := ResolveID2(idV, id)
		if lt > idt {
			t.Errorf("idt (%d) < lt (%d)", idt, lt)
			break
		}

		if lt == idt && idc-lc > delta {
			t.Errorf("idc-lc (%d) > delta (%d)", idc-lc, delta)
			break
		}

		if lt < idt {
			if lt != 0 && idt != lt+1 {
				t.Errorf("idt (%d) != lt+1 (%d)", idt, lt+1)
				break
			}
			if idc > delta {
				t.Errorf("idc (%d) > delta (%d)", idc, delta)
				break
			}
			lt = idt
		}
		lc = idc
	}
}
