package goid

import (
	"sync"
	"testing"
	"time"
)

func TestID_Generate_duplicate(t *testing.T) {
	ll := 10000000
	idArr := make(map[int64]struct{}, ll)
	idChan := make(chan int64, ll)
	i := 0
	id := NewID()
	wg := sync.WaitGroup{}
	wg.Add(ll + 1)
	go func() {
		defer wg.Done()
		for {
			id := <-idChan
			idArr[id] = struct{}{}
			i++
			if i == ll {
				break
			}
		}
	}()
	tt := time.Now()
	for i := 0; i < ll; i++ {
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

func TestID_Generate_increment(t *testing.T) {
	ll := 10000000
	var latestID int64
	id := NewID()
	for i := 0; i < ll; i++ {
		id1 := id.Generate()
		if id1 > latestID {
			latestID = id1
			continue
		}
		t.Errorf("id (%d) <= latestID (%d) ", id1, latestID)
	}

}
