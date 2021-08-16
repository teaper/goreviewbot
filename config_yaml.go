package main

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Rss struct {
	ClientURL string `yaml:"client_url"`
	Pubdate string `yaml:"pubdate"`
}

type  Channels struct {
	ChatID int64 `yaml:"chat_id"` //群 id
	ChatUserName string `yaml:"chat_user_name"` //群名称
	ChatType string `yaml:"chat_type"` //群类型
	Creator string `yaml:"creator"` //群创建者
}

type Bot struct {
	Token string `yaml:"token"`   //机器人 token
}

type Config struct {
	Enabled  bool   `yaml:"enabled"` //yaml：yaml格式 enabled：属性的为enabled
	Bot      Bot
	Channels Channels //群组或超级组
	Rss      Rss
}

//go:embed conf.yaml
var cyaml string

func (c *Config) GetConf() *Config {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("加载 config.yaml 文件出错：#%v ", err)
		//从 github 仓库下载模板
		log.Println("没事，我送你一个 conf.yaml 模板")
		f, err := os.Create("conf.yaml")
		if err != nil {
			log.Println(err)
		}
		l, err := f.WriteString(cyaml)
		if err != nil {
			log.Println(err)
			err := f.Close()
			if err != nil {
				return nil
			}
		}
		log.Println(l, "您的 conf.yaml 已送到！")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("解析 config.yaml 失败：%v", err)
	}
	return c
}
