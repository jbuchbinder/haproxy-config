package main

import (
	"io/ioutil"
	"os"
	"text/template"
)

func RenderConfig(outFile string, templateFile string, config *Config) error {
	//configMap := map[string]interface{}{
	//	"config": config,
	//}

	f,err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	t := template.Must(template.New("haproxy.cfg.template").Parse(string(f)))
	err = t.Execute(os.Stdout, &config)
	if err != nil {
		return err
	}

	return nil
}
