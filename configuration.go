package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Database struct {
		User     string
		Password string
		Database string
	}
	Mailer struct {
		Server string
	}
	Crypto struct {
		Application string
		Session     string
	}
	Web struct {
		Root       string
		CookieName string `yaml:"cookie_name"`
		XsrfToken  string `yaml:"xsrf_token"`
	}
}

func LoadConfig(path string) (*Config, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
