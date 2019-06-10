package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// MySQL configuration.
type MySQL struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	MaxConns int    `json:"maxConns"`
}

// Env holds configuration for this server.
type Env struct {
	Installed bool   `json:"-"`
	AppName   string `json:"appName"`
	MySQL     MySQL  `json:"mysql"`
}

// Save configuration into ./team.json
func (e *Env) Save() error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("./team.json", data, 0777)
}

// Environment in runtime.
var Environment = &Env{
	Installed: false,
	AppName:   "协作系统",
	MySQL: MySQL{
		Host:     "127.0.0.1:3306",
		User:     "root",
		Password: "root",
		Database: "team",
		MaxConns: 64,
	},
}

func init() {
	raw, err := ioutil.ReadFile("./team.json")
	if err != nil {
		fmt.Println("Can NOT find team.json. Using default instead.")
		return
	}

	err = json.Unmarshal(raw, Environment)
	if err != nil {
		fmt.Println("Failed to load team.json. Using default instead. Error: " + err.Error())
	}

	Environment.Installed = true
}
