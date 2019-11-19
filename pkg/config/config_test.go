package config

import (
	"fmt"
	"strconv"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	checkValid("testdata/valid.toml", config{
		WebName:         "ptemplate.nethruster.com",
		RecaptchaSecret: "xkmBhVrYaB0NhtHpHgAWeTnLZpTSxCKs0gigByk5",
		Mail: struct {
			Mailto     string `toml:"mailto"`
			Username   string `toml:"username"`
			Password   string `toml:"password"`
			SmtpServer string `toml:"smtp_server"`
			Port       int    `toml:"port"`
		}{
			Mailto: "personal@gmail.com",
			Username: "no-reply@nethruster.com",
			Password: "bNRxxIPxX7kLrbN8WCG22VUmpBqVBGgLTnyLdjob",
			SmtpServer: "smtp.nethruster.com",
			Port: 587,
		},
	}, t)

	checkValid("testdata/extra-info.toml", config{
		WebName:         "ptemplate.nethruster.com",
		RecaptchaSecret: "xkmBhVrYaB0NhtHpHgAWeTnLZpTSxCKs0gigByk5",
		Mail: struct {
			Mailto     string `toml:"mailto"`
			Username   string `toml:"username"`
			Password   string `toml:"password"`
			SmtpServer string `toml:"smtp_server"`
			Port       int    `toml:"port"`
		}{
			Mailto: "personal@gmail.com",
			Username: "no-reply@nethruster.com",
			Password: "bNRxxIPxX7kLrbN8WCG22VUmpBqVBGgLTnyLdjob",
			SmtpServer: "smtp.nethruster.com",
			Port: 587,
		},
	}, t)

	checkInvalid("testdata/invalid.toml", config{
		WebName:         "ptemplate.nethruster.com",
		RecaptchaSecret: "xkmBhVrYaB0NhtHpHgAWeTnLZpTSxCKs0gigByk5",
		Mail: struct {
			Mailto     string `toml:"mailto"`
			Username   string `toml:"username"`
			Password   string `toml:"password"`
			SmtpServer string `toml:"smtp_server"`
			Port       int    `toml:"port"`
		}{
			Mailto: "personal@gmail.com",
			Username: "no-reply@nethruster.com",
			Password: "bNRxxIPxX7kLrbN8WCG22VUmpBqVBGgLTnyLdjob",
			SmtpServer: "smtp.nethruster.com",
			Port: 832429423,
		},
	}, t)

	checkInvalid("testdata/incomplete.toml", config{
		RecaptchaSecret: "xkmBhVrYaB0NhtHpHgAWeTnLZpTSxCKs0gigByk5",
		Mail: struct {
			Mailto     string `toml:"mailto"`
			Username   string `toml:"username"`
			Password   string `toml:"password"`
			SmtpServer string `toml:"smtp_server"`
			Port       int    `toml:"port"`
		}{
			Password: "bNRxxIPxX7kLrbN8WCG22VUmpBqVBGgLTnyLdjob",
			SmtpServer: "smtp.nethruster.com",
		},
	}, t)

	checkInvalid("testdata/empty.toml", config{}, t)
	checkInvalid("testdata/nonexistent.toml", config{}, t)
}

func checkValid(path string, expectedConfig config, t *testing.T) {
	if err := testConfig(path, expectedConfig); err != nil {
		t.Errorf("unexpected error in path %s: %s", path, err)
	}
}

func checkInvalid(path string, expectedConfig config, t *testing.T) {
	if err := testConfig(path, expectedConfig); err == nil {
		t.Errorf("not error found in path %s", path)
	}
}

func testConfig(path string, expectedConfig config) error {
	actualConfig, err := LoadConfig(path)
	if err != nil {
		return err
	}

	if actualConfig.WebName != expectedConfig.WebName {
		return fmt.Errorf("web_name dont match: expected (%s) - found (%s)", expectedConfig.WebName, actualConfig.WebName)
	}
	if actualConfig.RecaptchaSecret != expectedConfig.RecaptchaSecret {
		return fmt.Errorf("recaptcha secret dont match: expected (%s) - found (%s)", expectedConfig.RecaptchaSecret, actualConfig.RecaptchaSecret)
	}
	if actualConfig.Mailto != expectedConfig.Mail.Mailto {
		return fmt.Errorf("mailto dont match: expected (%s) - found (%s)", expectedConfig.Mail.Mailto, actualConfig.Mailto)
	}
	if actualConfig.Username != expectedConfig.Mail.Username {
		return fmt.Errorf("username dont match: expected (%s) - found (%s)", expectedConfig.Mail.Username, actualConfig.Username)
	}
	if actualConfig.Password != expectedConfig.Mail.Password {
		return fmt.Errorf("password dont match: expected (%s) - found (%s)", expectedConfig.Mail.Password, actualConfig.Password)
	}
	if actualConfig.Hostname != expectedConfig.Mail.SmtpServer {
		return fmt.Errorf("server dont match: expected (%s) - found (%s)", expectedConfig.Mail.SmtpServer, actualConfig.Hostname)
	}
	if actualConfig.Port != strconv.Itoa(expectedConfig.Mail.Port) {
		return fmt.Errorf("port dont match: expected (%d) - found (%s)", expectedConfig.Mail.Port, actualConfig.Port)
	}

	return nil
}
