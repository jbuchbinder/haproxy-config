package main

import (
	"bytes"
	"io/ioutil"
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
	pid, err := ioutil.ReadFile(pidfile)
	args := make([]string, 1)
	args = append(args, "-f")
	args = append(args, config)
	args = append(args, "-p")
	args = append(args, pidfile)
	if pid != nil {
		args = append(args, "-sf")
		args = append(args, string(pid))
	}
	cmd := exec.Command(binary, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Err(err.Error())
		return err
	}
	log.Info("HaproxyReload: " + out.String())
	return nil
}
