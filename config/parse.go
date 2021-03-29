package config

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

const defaultConfig = `# 项目名称
project_name: pastebin
# 端口
port: :8080
# mongodb配置
mongodb:
  uri: mongodb://127.0.0.1:27017
  dbname: test
# https配置
https:
  # 是否开启https
  enable: false
  # 证书路径
  crt_path:
  key_path:
`

var (
	readFileError  = errors.New("")
	unmarshalError = errors.New("")
)

// Environments 项目主要配置项[子项] 如果需要扩展 在这里添加结构来实现yaml的解析
type Environments struct {
	ProjectName string `yaml:"project_name"` //项目名称
	Port        string `yaml:"port"`         //服务运行的 :port
	Https       HTTPS  `yaml:"https"`        //https配置
	Mongo       Mongo  `yaml:"mongodb"`      //mongo配置
}

type HTTPS struct {
	Enable  bool   `yaml:"enable"`
	CrtFile string `yaml:"crt_path"`
	KeyFile string `yaml:"key_path"`
}

type Mongo struct {
	URI    string `yaml:"uri"`
	Dbname string `yaml:"dbname"`
}

var (
	std      *Environments
	loadOnce sync.Once
)

func GetConfig() *Environments {
	defer func() {
		if err := recover(); err != nil {
			if err == readFileError {
				if err = ioutil.WriteFile("./config.yaml", []byte(defaultConfig), 0666); err != nil {
					log.Errorf("create config err: %v", err)
				} else {
					log.Infof("please rewrite config.yaml and restart")
				}
				os.Exit(1)
			}

			if err == unmarshalError {
				log.Errorf("unmarshal config file error, please see .config.yaml template file")
			}

			os.Exit(1)
		}
	}()

	loadOnce.Do(func() {
		std = new(Environments)
		yamlFile, err := ioutil.ReadFile("./config.yaml")
		if err != nil {
			log.Errorf("read file error: %v", err)
			panic(readFileError)
		}
		err = yaml.Unmarshal(yamlFile, std)
		if err != nil {
			//读取配置文件失败,停止执行
			log.Errorf("unmarshal file err: %v", err)
			panic(unmarshalError)
		}
	})
	return std
}
