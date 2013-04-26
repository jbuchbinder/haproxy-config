package main

import (
	"testing"
)

func TestRenderConfig(t *testing.T) {
	outFile := "test.render"
	template := "haproxy.cfg.template"
	config := &Config{
		Backends: map[string]*Backend{
			"a": &Backend{
				Name: "a",
				BackendServers: map[string]*BackendServer{
					"x": &BackendServer{
						Name:          "serverX",
						Bind:          "10.0.1.11:8080",
						Weight:        1,
						MaxConn:       1000,
						Check:         true,
						CheckInterval: 1000,
					},
				},
			},
		},
	}
	err := RenderConfig(outFile, template, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Done")
}
