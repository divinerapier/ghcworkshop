package main

import (
	"sort"
	"sync"
)

type Beers struct {
	beers []Beer
	mu    sync.Mutex
}

// sort beers by id
func (beers *Beers) Len() int {
	beers.mu.Lock()
	defer beers.mu.Unlock()
	return len(beers.beers)
}
func (beer *Beers) Less(i, j int) bool {
	beer.mu.Lock()
	defer beer.mu.Unlock()
	return beer.beers[i].ID < beer.beers[j].ID
}
func (beer *Beers) Swap(i, j int) {
	beer.mu.Lock()
	defer beer.mu.Unlock()
	beer.beers[i], beer.beers[j] = beer.beers[j], beer.beers[i]
}

func NewBeers(beers []Beer) *Beers {
	return &Beers{
		beers: beers,
	}
}

func (bs *Beers) GetRemains() []Beer {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	sort.Sort(bs)
	return bs.beers
}

func (bs *Beers) GetByID(id int) (*Beer, bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	for _, beer := range bs.beers {
		if beer.ID == id {
			return &beer, true
		}
	}
	return nil, false
}
