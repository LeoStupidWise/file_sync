package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type PathConf struct {
	// base_model 项目的目录
	BaseModel string `yaml:"base_model"`
	// 目标目录
	TargetFiles []string `yaml:"target_files"`
	// 定时规则
	Cron string `yaml:"cron"`
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

type TargetDir struct {
	// 根目录
	BaseDir string
	Dirs [] string
}
