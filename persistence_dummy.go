package main

import (
	"sync"
)

func init() {
	PersistenceLayerMap["dummy"] = func() PersistenceLayer {
		return new(PersistenceLayerDummy)
	}
}

type PersistenceLayerDummy struct {
}

func (self *PersistenceLayerDummy) Configure(c string) error {
	return nil
}

func (self *PersistenceLayerDummy) GetConfig() (*Config, error) {
	c := &Config{
		Backends: map[string]*Backend{},
		Mutex:    new(sync.RWMutex),
	}
	return c, nil
}

func (self *PersistenceLayerDummy) SetConfig(c *Config) error {
	return nil
}
