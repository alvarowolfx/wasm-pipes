package infra

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
)

func NewCache(uri string) (*cache.Cache, error) {
	if strings.HasPrefix(uri, "mem") {
		storeImpl, err := newInMemStore()
		if err != nil {
			return nil, fmt.Errorf("err creating cache: %v", err)
		}
		return cache.New(storeImpl), nil
	}
	return nil, fmt.Errorf("unknown cache type: %s", uri)
}

func newInMemStore() (store.StoreInterface, error) {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("err creating cache: %v", err)
	}
	return store.NewRistretto(ristrettoCache, nil), nil
}
