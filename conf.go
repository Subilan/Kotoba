package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	JwtPrivateKey string `yaml:"jwt_private_key"`
	JwtExpiration int `yaml:"jwt_expiration"`
}

func (c *config) get() *config {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		return nil
	}
	
	ymlErr := yaml.Unmarshal(file, c)

	if ymlErr != nil {
		return nil
	}

	return c
}
