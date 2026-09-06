package xcache

import (
	"fmt"
	gcache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Cache struct {
	// 单例对象的其他属性
	*gcache.Cache
}

var Instance *Cache
var once sync.Once

const defaultExpiration = 5 * time.Minute

func init() {
	once.Do(func() {
		c := gcache.New(defaultExpiration, 10*time.Minute)
		Instance = &Cache{Cache: c}
	})
}

type ValueLoader func() (any, error)

type ValueLoaderE[T any] func() (T, error)

func (c *Cache) Load(key string, loader ValueLoader) (any, error) {
	return c.LoadDuration(key, loader, defaultExpiration)
}

func LoadE[T any](key string, loader ValueLoaderE[T]) (T, error) {
	return LoadDurationE(key, loader, defaultExpiration)
}
func LoadDurationE[T any](key string, loader ValueLoaderE[T], d time.Duration) (T, error) {
	o, found := Instance.Get(key)
	if found {
		if v, ok := o.(T); ok {
			return v, nil
		} else {
			return *new(T), fmt.Errorf("convert faild: interface is not of type %T", new(T))
		}
	}
	r, err := loader()
	if err != nil {
		return *new(T), err
	}
	Instance.Set(key, r, d)
	return r, nil
}
func (c *Cache) LoadDuration(key string, loader ValueLoader, d time.Duration) (any, error) {
	o, found := c.Get(key)
	if found {
		return o, nil
	}
	r, err := loader()
	if err != nil {
		return nil, err
	}
	c.Set(key, r, d)
	return r, nil
}
