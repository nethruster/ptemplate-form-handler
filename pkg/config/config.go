package config

// Package config is the package that manages the functions related to the config file of ptemplate-form-handler.

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/ptemplate-form-handler/pkg/sender"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"strconv"
)

// config represents the structure of the config file of ptemplate-form-handler.
type config struct {
	WebName string `toml:"web_name"`
	RecaptchaSecret string `toml:"recaptcha_secret"`
	Mail struct {
		Mailto string `toml:"mailto"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		SmtpServer string `toml:"smtp_server"`
		Port int `toml:"port"`
	} `toml:"mail"`
}

// LoadConfig will read the config from the path provided and return a sender.Mail object.
func LoadConfig(path string) (*sender.Mail, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file from path \"%s\": %w", path, err)
	}

	var c config
	if err = toml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("error parsing config file from path \"%s\": %w", path, err)
	}

	if err = checkValidInput(&c); err != nil {
		return nil, fmt.Errorf("invalid configuration file: %w", err)
	}

	return &sender.Mail{
		WebName:         c.WebName,
		RecaptchaSecret: c.RecaptchaSecret,
		Mailto:          c.Mail.Mailto,
		Username:        c.Mail.Username,
		Password:        c.Mail.Password,
		Hostname:        c.Mail.SmtpServer,
		Port:            strconv.Itoa(c.Mail.Port),
	}, nil
}

// checkValidInput checks if all the fields in the config provided are valid
func checkValidInput(c *config) error {
	if c.WebName == "" {
		return errors.New("empty web_name")
	}
	if c.RecaptchaSecret == "" {
		return errors.New("empty recaptcha_secret")
	}
	if c.Mail.Mailto == "" {
		return errors.New("empty mailto")
	}
	if c.Mail.Username == "" {
		return errors.New("empty username")
	}
	if c.Mail.Password == "" {
		return errors.New("empty password")
	}
	if c.Mail.SmtpServer == "" {
		return errors.New("empty smtp_server")
	}
	if c.Mail.Port < 1 || c.Mail.Port > 65535 {
		return errors.New("invalid port")
	}
	return nil
}
