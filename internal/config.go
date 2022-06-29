package internal

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Creds struct {
	Username string
	Key      string
	Password string
}

type Config struct {
	Source    string
	Inventory string
	Creds     Creds
	Cooldown  int
}

func ParseConfig(fileName string) (*Config, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading config file %s: %s\n", fileName, err)
		return nil, err
	}

	config := &Config{}

	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}
