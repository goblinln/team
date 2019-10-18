package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"team/orm"
)

type (
	// MySQL configuration.
	MySQL struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	}

	// Env holds configuration for this server.
	Env struct {
		Installed bool   `json:"-"`
		AppPort   string `json:"appPort"`
		MySQL     *MySQL `json:"mysql"`
	}
)

// Addr returns DSN address of MYSQL server
func (m *MySQL) Addr() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?multiStatements=true&charset=utf8&collation=utf8_general_ci",
		m.User, m.Password, m.Host, m.Database,
	)
}

// Environment in runtime.
var Environment = &Env{
	Installed: false,
	AppPort:   ":8080",
	MySQL: &MySQL{
		Host:     "127.0.0.1:3306",
		User:     "root",
		Password: "root",
		Database: "team",
	},
}

// Save configuration into ./team.json
func (e *Env) Save() error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("./team.json", data, 0777)
}

// Prepare environment for this app.
func (e *Env) Prepare() {
	raw, err := ioutil.ReadFile("./team.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, e)
	if err != nil {
		return
	}

	e.Installed = true

	if err := orm.OpenDB("mysql", e.MySQL.Addr()); err != nil {
		log.Fatalf("Failed to prepare database : %s", err.Error())
	}
}
