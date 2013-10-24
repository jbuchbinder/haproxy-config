package main

import (
	"flag"
	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"log/syslog"
	"net/http"
	"time"
)

var (
	bind                = flag.String("bind", ":8888", "Port/IP for binding interface")
	haproxyBinary       = flag.String("haproxyBinary", "/usr/sbin/haproxy", "Path to haproxy binary")
	haproxyConfigFile   = flag.String("haproxyConfig", "/etc/haproxy.cfg", "Configuration file for haproxy")
	haproxyPidFile      = flag.String("haproxyPidFile", "/var/run/haproxy.pid", "Location of haproxy PID file")
	haproxyTemplateFile = flag.String("template", "haproxy.cfg.template", "Template file to build haproxy config")
	htpasswd            = flag.String("htpasswd", "haproxy.htpasswd", "htpasswd-formatted authentication file")
	persistence         = flag.String("persist", "dummy", "Persistence plugin to use")
	persistenceOpts     = flag.String("persistOpt", "", "Options to pass to the active persistence plugin")
	ConfigObj           *Config
	PersistenceObj      PersistenceLayer
	log, _              = syslog.New(syslog.LOG_DEBUG, "haproxy-config")
)

func main() {
	flag.Parse()

	PersistenceObj = GetPersistenceLayer(*persistence)
	if PersistenceObj == nil {
		log.Err("Unable to load persistence plugin " + *persistence)
		panic("Dying")
	}
	var err error
	ConfigObj, err = PersistenceObj.GetConfig()
	if err != nil {
		log.Err("Unable to load config from persistence plugin " + *persistence)
		panic("Dying")
	}

	if ConfigObj.PidFile != *haproxyPidFile {
		ConfigObj.PidFile = *haproxyPidFile
	}

	r := mux.NewRouter()

	// Define paths
	sub := r.PathPrefix("/api").Subrouter()

	// Display handlers
	sub.HandleFunc("/config", configHandler).Methods("GET")
	sub.HandleFunc("/reload", configReloadHandler).Methods("GET")
	sub.HandleFunc("/backend/{backend}", backendHandler).Methods("GET")
	sub.HandleFunc("/backend/{backend}", backendAddHandler).Methods("POST")
	sub.HandleFunc("/backend/{backend}", backendDeleteHandler).Methods("DELETE")
	sub.HandleFunc("/backend/{backend}/server/{server}", backendServerHandler).Methods("GET")
	sub.HandleFunc("/backend/{backend}/server/{server}", backendServerAddHandler).Methods("POST")
	sub.HandleFunc("/backend/{backend}/server/{server}", backendServerDeleteHandler).Methods("DELETE")

	s := &http.Server{
		Addr:           *bind,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	h := auth.HtpasswdFileProvider(*htpasswd)
	a := auth.NewBasicAuthenticator("haproxy config", h)
	http.Handle("/", a.Wrap(func(w http.ResponseWriter, ar *auth.AuthenticatedRequest) {
		r.ServeHTTP(w, &ar.Request)
	}))
	log.Err(s.ListenAndServe().Error())
}

// On configuration change, call this to reload config
func configChangeHook() {
	// Attempt to serialize back to persistence layer
	_ = PersistenceObj.SetConfig(ConfigObj)

	err := RenderConfig(*haproxyConfigFile, *haproxyTemplateFile, ConfigObj)
	if err != nil {
		log.Err("Error rendering config file")
		return
	}

	err = HaproxyReload(*haproxyBinary, *haproxyConfigFile, *haproxyPidFile)
	if err != nil {
		log.Err("Error rendering config file")
		return
	}
}
