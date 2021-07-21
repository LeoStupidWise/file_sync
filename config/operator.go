package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
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

// 存放每一个目标目录的详情
type TargetDir struct {
	// 根目录
	BaseDir string
	Dirs [] DirInfo
}

// 每一个目录的信息
type DirInfo struct {
	// 地址
	Path string
	// 是否是文件夹
	IsDir bool
	// 修改时间
	UpdatedAt time.Time
}
