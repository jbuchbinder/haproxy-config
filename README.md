# HAPROXY-CONFIG

Daemon which creates a REST-ful API which can be used to dynamically control
HAproxy configuration.

## BUILDING

	go get github.com/gorilla/mux
	go get github.com/bradfitz/gomemcache/memcache
	go build

## BUILD STATUS

[![Status](https://secure.travis-ci.org/jbuchbinder/haproxy-config.png)](http://travis-ci.org/jbuchbinder/haproxy-config)

## TODO

* ~~Persist configuration to NoSQL database (Memcache, plugins).~~
* ~~Mutex locking of config object, since maps aren't thread safe.~~
* jQuery-based UI.
* Authentication layer/security, maybe ACLs?
* Implement global configuration
* Implement frontend configuration
* Finish implementing configuration for backend servers

