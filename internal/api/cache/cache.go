package cache

import (
	"errors"
	"sync"
	"time"
	"wbtech/level0/internal/api/models"
	"wbtech/level0/internal/db/repositories"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Data       models.Data
	Created    time.Time
	Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *Cache) Set(key string, data models.Data, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Data:       data,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

func (c *Cache) Get(key string) (models.Data, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return models.Data{}, false
	}

	if item.Expiration > 0 {

		if time.Now().UnixNano() > item.Expiration {
			return models.Data{}, false
		}

	}

	return item.Data, true
}

func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}

func (c *Cache) Fill(repositories repositories.Repositories) error {
	data, err := repositories.OrderRepository.GetAll()

	if err != nil {
		return err
	}

	for _, v := range data {
		c.Set(v.OrderUID, v, 0)
	}

	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {

	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
