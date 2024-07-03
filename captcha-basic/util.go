package basic

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Info struct {
	SlugName string `yaml:"slug_name"`
	Type     string `yaml:"type"`
	Version  string `yaml:"version"`
	Author   string `yaml:"author"`
	Link     string `yaml:"link"`
}

func (c *Info) getInfo() *Info {
	_, filename, _, _ := runtime.Caller(0)
	wd := filepath.Dir(filename)

	yamlFilePath := filepath.Join(wd, "info.yaml")
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err)
	}
	return c
}
