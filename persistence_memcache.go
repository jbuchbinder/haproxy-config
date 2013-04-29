package main

import (
	"github.com/bradfitz/gomemcache/memcache"
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
	self.Memcache = memcache.New(c)
	return nil
}
