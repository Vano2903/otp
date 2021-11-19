package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Email           string `yaml:"email"`
	EmailPassword   string `yaml:"email-password"`
	UserFilePath    string `yaml:"user-file-path"`
	PendingFilePath string `yaml:"pending-file-path"`
	OtpFilePath     string `yaml:"otp-file-path"`
}

func (c *Config) Load() error {
	email := os.Getenv("email")
	psw := os.Getenv("email-password")
	userFilePath := os.Getenv("user-file-path")
	pendingFilePath := os.Getenv("pending-file-path")
	otpFilePath := os.Getenv("otp-file-path")

	if otpFilePath == "" || pendingFilePath == "" || userFilePath == "" || email == "" || psw == "" {
		fmt.Println("opening config.yaml")
		dat, err := ioutil.ReadFile("config.yaml")
		err = yaml.Unmarshal([]byte(dat), &c)
		return err
	}

	c.Email = email
	c.EmailPassword = psw
	c.UserFilePath = userFilePath
	c.PendingFilePath = pendingFilePath
	c.OtpFilePath = otpFilePath
	return nil
}
