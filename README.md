# HAPROXY-CONFIG

Daemon which creates a REST-ful API which can be used to dynamically control
HAproxy configuration. It uses an htpasswd-formatted file to control
authorization for the API.

## BUILDING

	go get github.com/abbot/go-http-auth
	go get github.com/bradfitz/gomemcache/memcache
	go get github.com/gorilla/mux
	go build

## BUILD STATUS

[![Status](https://secure.travis-ci.org/jbuchbinder/haproxy-config.png)](http://travis-ci.org/jbuchbinder/haproxy-config)

[![Gobuild Download](http://gobuild.io/badge/github.com/jbuchbinder/haproxy-config/downloads.svg)](http://gobuild.io/github.com/jbuchbinder/haproxy-config)

## TODO

* ~~Persist configuration to NoSQL database (Memcache, plugins).~~
* ~~Mutex locking of config object, since maps aren't thread safe.~~
* jQuery-based UI.
* ~~Authentication layer/security, maybe ACLs?~~
* Implement global configuration
* Implement frontend configuration
* Finish implementing configuration for backend servers

