package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type PathConf struct {
	BaseModel string `yaml:"base_model"`
}

func (c *PathConf) GetPathConf() *PathConf {
	yamlFile, err := ioutil.ReadFile("config/path.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}
