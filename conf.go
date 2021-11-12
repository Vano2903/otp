package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Email         string `yaml:"email"`
	EmailPassword string `yaml:"email-password"`
}

func (c *Config) Load() error {
	email := os.Getenv("email")
	psw := os.Getenv("email-password")
	if email == "" || psw == "" {
		dat, err := ioutil.ReadFile("config.yaml")
		err = yaml.Unmarshal([]byte(dat), &c)
		if err != nil {
			return err
		}
	} else {
		c.Email = email
		c.EmailPassword = psw
	}
	return nil
}
