package main

import (
	"sync"
)

// Main configuration object. This contains all variables and is passed to
// the templating engine.
type Config struct {
	Backends map[string]*Backend `json:"backends"`
	PidFile  string              `json:"pidfile"`
	Mutex    *sync.RWMutex       `json:"-"`
}

// Defines a single haproxy "backend".
type Backend struct {
	Name           string                    `json:"name"`
	BackendServers map[string]*BackendServer `json:"servers"`
	Options        ProxyOptions              `json:"options"`
}

// Options which are common between frontends, backends, etc
type ProxyOptions struct {
	AbortOnClose    bool `json:"abortOnClose"`
	AllBackups      bool `json:"allBackups"`
	CheckCache      bool `json:"checkCache"`
	ForwardFor      bool `json:"forwardFor"`
	HttpClose       bool `json:"httpClose"`
	HttpCheck       bool `json:"httpCheck"`
	LdapCheck       bool `json:"ldapCheck"`
	MysqlCheck      bool `json:"mysqlCheck"`
	PgsqlCheck      bool `json:"pgsqlCheck"`
	RedisCheck      bool `json:"redisCheck"`
	SmtpCheck       bool `json:"smtpCheck"`
	SslHelloCheck   bool `json:"sslHelloCheck"`
	TcpKeepAlive    bool `json:"tcpKeepAlive"`
	TcpLog          bool `json:"tcpLog"`
	TcpSmartAccept  bool `json:"tcpSmartAccept"`
	TcpSmartConnect bool `json:"tcpSmartConnect"`
	Transparent     bool `json:"transparent"`
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
