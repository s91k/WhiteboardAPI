package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		File     string `yaml:"file" envconfig:"DB_FILE"`
		Username string `yaml:"sql-user" envconfig:"DB_USERNAME"`
		Password string `yaml:"sql-pass" envconfig:"DB_PASSWORD"`
		Database string `yaml:"sql-database" envconfig:"DB_DATABASE"`
		Server   string `yaml:"sql-server" envconfig:"DB_SERVER"`
		Port     int    `yaml:"sql-port" envconfig:"DB_PORT"`
	} `yaml:"database"`

	Whiteboard struct {
		Width  int `yaml:"width" envconfig:"WHITEBOARD_WIDTH"`
		Height int `yaml:"height" envconfig:"WHITEBOARD_HEIGHT"`
	}
}

func readConfig(cfg *Config) {
	readFile(cfg)
	readEnv(cfg)
	fmt.Printf("%+v", cfg)
}

func readFile(cfg *Config) {
	fileName := "config.yml"
	s := os.Getenv("RUNENVIRONMENT")
	if len(s) > 0 {
		fileName = "config" + s + ".yml"
	}

	f, _ := os.Open(fileName)
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.Decode(cfg)
}

func readEnv(cfg *Config) {
	envconfig.Process("", cfg)
}
