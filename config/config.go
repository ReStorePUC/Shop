package config

import (
	"github.com/restore/shop/repository"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const EmailHeader = "X-Consumer-Username"

type Configuration struct {
	Mysql repository.Config `yaml:"mysql"`
}

var config Configuration

func Init() {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
}

func NewDBConfig() *repository.Config {
	return &config.Mysql
}
