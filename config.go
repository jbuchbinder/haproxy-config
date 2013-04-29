package main

import (
	"sync"
)

// Main configuration object. This contains all variables and is passed to
// the templating engine.
type Config struct {
	Backends map[string]*Backend `json:"backends"`
	Mutex    *sync.RWMutex       `json:"-"`
}

// Defines a single haproxy "backend".
type Backend struct {
	HttpClose      bool                      `json:"httpClose"`
	Name           string                    `json:"name"`
	BackendServers map[string]*BackendServer `json:"servers"`
}

// Defines a server which exists in a backend.
type BackendServer struct {
	Name          string `json:"name"`
	Bind          string `json:"bind"`
	Weight        int    `json:"weight"`
	MaxConn       int    `json:"maxconn"`
	Check         bool   `json:"check"`
	CheckInterval int    `json:"checkInterval"`
}
