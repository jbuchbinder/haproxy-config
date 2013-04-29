package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
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

	r := mux.NewRouter()

	// Define paths
	sub := r.PathPrefix("/").Subrouter()

	// Display handlers
	sub.HandleFunc("/", configHandler)
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
	http.Handle("/", r)
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

// Main config object

func configHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(ConfigObj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}

// Backend functions

func backendHandler(w http.ResponseWriter, r *http.Request) {
	ConfigObj.Mutex.RLock()
	defer ConfigObj.Mutex.RUnlock()
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	backend := vars["backend"]
	if _, found := ConfigObj.Backends[backend]; !found {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	b, err := json.Marshal(ConfigObj.Backends[backend])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}

func backendAddHandler(w http.ResponseWriter, r *http.Request) {
	ConfigObj.Mutex.Lock()
	defer ConfigObj.Mutex.Unlock()
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	backend := vars["backend"]
	if _, found := ConfigObj.Backends[backend]; !found {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	obj := &Backend{}
	err = json.Unmarshal(body, obj)
	if err != nil {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ConfigObj.Backends[backend] = obj
	fmt.Fprint(w, "true")
	configChangeHook()
}

func backendDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ConfigObj.Mutex.Lock()
	defer ConfigObj.Mutex.Unlock()
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	backend := vars["backend"]
	if _, found := ConfigObj.Backends[backend]; !found {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	delete(ConfigObj.Backends, backend)
	fmt.Fprint(w, "true")
	configChangeHook()
}

// Backend server functions

func backendServerHandler(w http.ResponseWriter, r *http.Request) {
	ConfigObj.Mutex.RLock()
	defer ConfigObj.Mutex.RUnlock()
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	backend := vars["backend"]
	if _, found := ConfigObj.Backends[backend]; !found {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	b, err := json.Marshal(ConfigObj.Backends[backend])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}

func backendServerAddHandler(w http.ResponseWriter, r *http.Request) {
	ConfigObj.Mutex.Lock()
	defer ConfigObj.Mutex.Unlock()
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	backend := vars["backend"]
	server := vars["server"]
	if _, found := ConfigObj.Backends[backend]; !found {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if _, found := ConfigObj.Backends[backend].BackendServers[server]; found {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	obj := &BackendServer{}
	err = json.Unmarshal(body, obj)
	if err != nil {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ConfigObj.Backends[backend].BackendServers[server] = obj
	fmt.Fprint(w, "true")
	configChangeHook()
}

func backendServerDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ConfigObj.Mutex.Lock()
	defer ConfigObj.Mutex.Unlock()
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	backend := vars["backend"]
	server := vars["server"]
	if _, found := ConfigObj.Backends[backend]; !found {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if _, found := ConfigObj.Backends[backend].BackendServers[server]; !found {
		fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	delete(ConfigObj.Backends[backend].BackendServers, server)
	fmt.Fprint(w, "true")
	configChangeHook()
}
