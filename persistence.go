package main

import (
	"strings"
)

var (
	PersistenceLayerMap = map[string]func() PersistenceLayer{}
)

type PersistenceLayer interface {
	Configure(c string) error
	GetConfig() (*Config, error)
	SetConfig(c *Config) error
}

// Resolves PersistenceLayer objects based on their string names.
func GetPersistenceLayer(p string) PersistenceLayer {
	log.Info("Selecting persistence layer using: '" + p + "'")
	pName := strings.TrimSpace(p)
	if _, exists := PersistenceLayerMap[pName]; exists {
		return PersistenceLayerMap[pName]()
	} else {
		log.Err("Unable to resolve persistence layer " + pName)
	}
	return nil
}
