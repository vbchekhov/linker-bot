package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var config = NewConfig()

type Config struct {
	Debug  bool   `yaml:"debug"`
	Token  string `yaml:"token"`
	Folder string `yaml:"folder"`
	Port   string `yaml:"port"`
}

func NewConfig() *Config {

	c := new(Config)

	bytes, err := ioutil.ReadFile("./app.yaml")
	if err != nil {
		log.Printf("Error load read config %v", err)
	}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Printf("Error unmarhsl %v", err)
	}

	return c
}
