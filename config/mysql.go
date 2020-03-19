package config

import "fmt"

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
