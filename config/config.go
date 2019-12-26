package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type (
	// MySQL configuration.
	MySQL struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	}

	// Configure holds configuration for this server.
	Configure struct {
		Installed bool   `json:"-"`
		AppPort   string `json:"appPort"`
		MySQL     *MySQL `json:"mysql"`
	}
)

// Default setting for this app.
var Default = &Configure{
	Installed: false,
	AppPort:   ":8080",
	MySQL: &MySQL{
		Host:     "127.0.0.1:3306",
		User:     "root",
		Password: "root",
		Database: "team",
	},
}

// Load configure file for this app.
func (c *Configure) Load() {
	raw, err := ioutil.ReadFile("./team.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, c)
	if err != nil {
		fmt.Printf("Failed to parse configure file : ./team.json. Error: %v\n", err)
		os.Exit(-1)
	}

	c.Installed = true
}

// GetMySQLAddr returns MySQL service address
func (c *Configure) GetMySQLAddr() string {
	m := c.MySQL
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?multiStatements=true&charset=utf8&collation=utf8_general_ci",
		m.User, m.Password, m.Host, m.Database)
}

// Save configure data of this app.
func (c *Configure) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("./team.json", data, 0777)
}
