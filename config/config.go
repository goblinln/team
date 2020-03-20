package config

import (
	"fmt"

	"team/common/auth"
	"team/common/ini"
)

// Installed flag.
var Installed = false

// AuthKind for this app.
type AuthKind int

const (
	// AuthKindOnlyBuildin only uses build account.
	AuthKindOnlyBuildin AuthKind = iota
	// AuthKindSMTP uses SMTP auth
	AuthKindSMTP
	// AuthKindLDAP uses LDAP auth
	AuthKindLDAP
)

// ExtraAuth for this app.
var ExtraAuth auth.LoginProcessor = nil

// AppInfo holds basic information for this server.
type AppInfo struct {
	Name string
	Port int
	Auth AuthKind
}

// App information of this system.
var App = &AppInfo{
	Name: "Team",
	Port: 8080,
	Auth: AuthKindOnlyBuildin,
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

// MySQL information.
var MySQL = &MySQLInfo{
	Host:     "127.0.0.1:3306",
	User:     "root",
	Password: "root",
	Database: "team",
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
	App.Auth = AuthKind(setting.GetInt("app", "auth"))

	MySQL.Host = setting.GetString("mysql", "host")
	MySQL.User = setting.GetString("mysql", "user")
	MySQL.Password = setting.GetString("mysql", "password")
	MySQL.Database = setting.GetString("mysql", "database")

	switch App.Auth {
	case AuthKindSMTP:
		UseSMTPAuth(
			setting.GetString("smtp_login", "host"),
			setting.GetInt("smtp_login", "port"),
			setting.GetBool("smtp_login", "plain"),
			setting.GetBool("smtp_login", "tls"),
			setting.GetBool("smtp_login", "skip_verify"),
		)
	case AuthKindLDAP:
	}

	Installed = true
}

// Save configuration to file
func Save() error {
	setting := ini.New()

	setting.SetString("app", "name", App.Name)
	setting.SetInt("app", "port", App.Port)
	setting.SetInt("app", "login_type", int(App.Auth))

	setting.SetString("mysql", "host", MySQL.Host)
	setting.SetString("mysql", "user", MySQL.User)
	setting.SetString("mysql", "password", MySQL.Password)
	setting.SetString("mysql", "database", MySQL.Database)

	switch App.Auth {
	case AuthKindSMTP:
		smtp := ExtraAuth.(*auth.SMTPProcessor)
		setting.SetString("smtp_login", "host", smtp.Host)
		setting.SetInt("smtp_login", "port", smtp.Port)
		setting.SetBool("smtp_login", "plain", smtp.Plain)
		setting.SetBool("smtp_login", "tls", smtp.TLS)
		setting.SetBool("smtp_login", "skip_verify", smtp.SkipVerfiy)
	case AuthKindLDAP:
	}

	return setting.Save("./team.ini")
}

// UseSMTPAuth uses SMTP as extra auth method.
func UseSMTPAuth(host string, port int, plain, tls, skipVerify bool) {
	ExtraAuth = &auth.SMTPProcessor{
		Host:       host,
		Port:       port,
		Plain:      plain,
		TLS:        tls,
		SkipVerfiy: skipVerify,
	}
}
