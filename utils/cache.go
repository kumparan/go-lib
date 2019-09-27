package utils

import (
	"time"

	"github.com/kumparan/cacher"
	"github.com/kumparan/tapao"
)

// SetNotFoundCache set nil for an object that not found
func SetNotFoundCache(keeper cacher.Keeper, cacheKey string, ttl time.Duration) error {
	value, err := tapao.Marshal(nil, tapao.With(tapao.JSON))
	if err != nil {
		return err
	}

	item := cacher.NewItemWithCustomTTL(cacheKey, value, ttl)
	err = keeper.StoreWithoutBlocking(item)
	if err != nil {
		return err
	}

	return nil
}
