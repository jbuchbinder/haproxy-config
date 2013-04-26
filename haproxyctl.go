package main

import (
	"bytes"
	"log"
	"net"
	"os/exec"
)

// Execute haproxy command over administrative/stats socket
func HaproxyCtl(socket string, command []byte) error {
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(command)
	if err != nil {
		return err
	}
	return nil
}

// Configuration reload
func HaproxyReload(binary, config, pidfile string) error {
	// Read pid
	pid := ""
	cmd := exec.Command(binary, "-f", config, "-p", pidfile, "-sf", pid)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print("HaproxyReload: " + out.String())
	return nil
}
