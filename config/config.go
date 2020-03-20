package config

import (
	"fmt"

	"team/common/auth"
	"team/common/ini"
)

var (
	// Installed check if already installed.
	Installed bool = false

	// App information of this system.
	App = &AppInfo{
		Name:      "Team",
		Port:      8080,
		LoginType: auth.KindNone,
	}

	// MySQL information.
	MySQL = &MySQLInfo{
		Host:     "127.0.0.1:3306",
		User:     "root",
		Password: "root",
		Database: "team",
	}

	// ExtraLoginProcessor for this app.
	ExtraLoginProcessor auth.LoginProcessor = nil
)

// AppInfo holds basic information for this server.
type AppInfo struct {
	Name      string
	Port      int
	LoginType auth.Kind
}

// Addr returns address for this service.
func (a *AppInfo) Addr() string {
	return fmt.Sprintf(":%d", a.Port)
}

// MySQLInfo host information of mysql server to be used.
type MySQLInfo struct {
	Host     string
	User     string
	Password string
	Database string
}

// URL returns MySQL service address.
func (m *MySQLInfo) URL() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?multiStatements=true&charset=utf8&collation=utf8_general_ci",
		m.User, m.Password, m.Host, m.Database)
}

// Load configuration from file.
func Load() {
	setting, err := ini.Load("./team.ini")
	if err != nil {
		fmt.Printf("Load configuration failed. %v. Using default configuration.\n", err)
		return
	}

	App.Name = setting.GetString("app", "name")
	App.Port = setting.GetInt("app", "port")
	App.LoginType = auth.Kind(setting.GetInt("app", "login_type"))

	MySQL.Host = setting.GetString("mysql", "host")
	MySQL.User = setting.GetString("mysql", "user")
	MySQL.Password = setting.GetString("mysql", "password")
	MySQL.Database = setting.GetString("mysql", "database")

	switch App.LoginType {
	case auth.KindSMTP:
		ExtraLoginProcessor = &auth.SMTPLoginProcessor{
			Host:       setting.GetString("smtp_login", "host"),
			Port:       setting.GetInt("smtp_login", "port"),
			Plain:      setting.GetBool("smtp_login", "plain"),
			TLS:        setting.GetBool("smtp_login", "tls"),
			SkipVerfiy: setting.GetBool("smtp_login", "skip_verify"),
		}
	case auth.KindLDAP:
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
	case auth.KindSMTP:
		smtp := ExtraLoginProcessor.(*auth.SMTPLoginProcessor)
		setting.SetString("smtp_login", "host", smtp.Host)
		setting.SetInt("smtp_login", "port", smtp.Port)
		setting.SetBool("smtp_login", "plain", smtp.Plain)
		setting.SetBool("smtp_login", "tls", smtp.TLS)
		setting.SetBool("smtp_login", "skip_verify", smtp.SkipVerfiy)
	case auth.KindLDAP:
	}

	return setting.Save("./team.ini")
}
