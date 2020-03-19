package config

import (
	"fmt"

	"team/common/ini"
)

var (
	// Installed check if already installed.
	Installed bool = false

	// App information of this system.
	App = &AppInfo{
		Name:      "Team",
		Port:      8080,
		LoginType: LoginTypeSimple,
	}

	// MySQL information.
	MySQL = &MySQLInfo{
		Host:     "127.0.0.1:3306",
		User:     "root",
		Password: "root",
		Database: "team",
	}

	// ExtraLoginProcessor for this app.
	ExtraLoginProcessor LoginProcessor = nil
)

// Load configuration from file.
func Load() {
	setting, err := ini.Load("./team.ini")
	if err != nil {
		fmt.Printf("Load configuration failed. %v. Using default configuration.\n", err)
		return
	}

	App.Name = setting.GetString("app", "name")
	App.Port = setting.GetInt("app", "port")
	App.LoginType = LoginType(setting.GetInt("app", "login_type"))

	MySQL.Host = setting.GetString("mysql", "host")
	MySQL.User = setting.GetString("mysql", "user")
	MySQL.Password = setting.GetString("mysql", "password")
	MySQL.Database = setting.GetString("mysql", "database")

	switch App.LoginType {
	case LoginTypeSMTP:
		ExtraLoginProcessor = &SMTPLoginProcessor{
			Host:       setting.GetString("smtp_login", "host"),
			Port:       setting.GetInt("smtp_login", "port"),
			Plain:      setting.GetBool("smtp_login", "plain"),
			TLS:        setting.GetBool("smtp_login", "tls"),
			SkipVerfiy: setting.GetBool("smtp_login", "skip_verify"),
		}
	}

	Installed = true
}

// Save configuration to file
func Save() error {
	setting := ini.New()

	setting.SetString("app", "name", App.Name)
	setting.SetInt("app", "port", App.Port)
	setting.SetInt("app", "login_type", int(App.LoginType))

	setting.SetString("mysql", "host", MySQL.Host)
	setting.SetString("mysql", "user", MySQL.User)
	setting.SetString("mysql", "password", MySQL.Password)
	setting.SetString("mysql", "database", MySQL.Database)

	switch App.LoginType {
	case LoginTypeSMTP:
		smtp := ExtraLoginProcessor.(*SMTPLoginProcessor)
		setting.SetString("smtp_login", "host", smtp.Host)
		setting.SetInt("smtp_login", "port", smtp.Port)
		setting.SetBool("smtp_login", "plain", smtp.Plain)
		setting.SetBool("smtp_login", "tls", smtp.TLS)
		setting.SetBool("smtp_login", "skip_verify", smtp.SkipVerfiy)
	}

	return setting.Save("./team.ini")
}
