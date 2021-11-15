package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	UserFilePath    string `yaml:"user-file-path"`
	Email           string `yaml:"email"`
	EmailPassword   string `yaml:"email-password"`
	PendingFilePath string `yaml:"pending-file-path"`
}

func (c *Config) Load() error {
	fmt.Println("loading config")
	userFilePath := os.Getenv("user-file-path")
	pendingFilePath := os.Getenv("pending-file-path")
	email := os.Getenv("email")
	psw := os.Getenv("email-password")
	if pendingFilePath == "" || userFilePath == "" || email == "" || psw == "" {
		fmt.Println("opening config.yaml")
		dat, err := ioutil.ReadFile("config.yaml")
		err = yaml.Unmarshal([]byte(dat), &c)
		if err != nil {
			return err
		}
	} else {
		c.UserFilePath = userFilePath
		c.Email = email
		c.EmailPassword = psw
	}
	return nil
}
