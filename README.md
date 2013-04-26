# HAPROXY-CONFIG

Daemon which creates a REST-ful API which can be used to dynamically control
HAproxy configuration.

## BUILDING

`
go get github.com/gorilla/mux
go build
`

## BUILD STATUS

[![Status](https://secure.travis-ci.org/jbuchbinder/haproxy-config.png)](http://travis-ci.org/jbuchbinder/haproxy-config)

## TODO

* Persist configuration to NoSQL database (Redis).
* Mutex locking of config object, since maps aren't thread safe.
* jQuery-based UI.
* Authentication layer/security, maybe ACLs?

