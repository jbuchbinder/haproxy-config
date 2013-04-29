package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

func init() {
	PersistenceLayerMap["file"] = func() PersistenceLayer {
		return new(PersistenceLayerFile)
	}
}

type PersistenceLayerFile struct {
	LocalFilename string
}

func (self *PersistenceLayerFile) Configure(c string) error {
	self.LocalFilename = c
	return nil
}

func (self *PersistenceLayerFile) GetConfig() (*Config, error) {
	var c *Config
	s, err := ioutil.ReadFile(self.LocalFilename)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(s, c)
	if err == nil {
		c.Mutex = new(sync.RWMutex)
	}
	return c, err
}

func (self *PersistenceLayerFile) SetConfig(c *Config) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(self.LocalFilename, b, 0666)
	if err != nil {
		return err
	}
	return nil
}
