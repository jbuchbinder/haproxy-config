package main

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"strings"
	"sync"
)

const (
	MEMCACHE_PERSISTENCE_KEY = "haproxyConfigMemcachePersist"
)

func init() {
	PersistenceLayerMap["memcache"] = func() PersistenceLayer {
		return new(PersistenceLayerMemcache)
	}
}

type PersistenceLayerMemcache struct {
	Memcache *memcache.Client
}

func (self *PersistenceLayerMemcache) Configure(c string) error {
	s := strings.Split(c, ",")
	self.Memcache = memcache.New(s...)
	return nil
}

func (self *PersistenceLayerMemcache) GetConfig() (*Config, error) {
	var c *Config
	s, err := self.Memcache.Get(MEMCACHE_PERSISTENCE_KEY)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(s.Value, c)
	if err == nil {
		c.Mutex = new(sync.RWMutex)
	}
	return c, err
}

func (self *PersistenceLayerMemcache) SetConfig(c *Config) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = self.Memcache.Set(&memcache.Item{
		Key:   MEMCACHE_PERSISTENCE_KEY,
		Value: b,
	})
	if err != nil {
		return err
	}
	return nil
}
