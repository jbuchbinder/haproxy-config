package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// Main config object

func configHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(ConfigObj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "false")
		return
	}
	fmt.Fprint(w, string(b))
}

func configReloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var err error
	ConfigObj, err = PersistenceObj.GetConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "false")
		return
	}
	fmt.Fprint(w, "true")
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
