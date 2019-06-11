package controller

import (
	"crypto/md5"
	"fmt"

	"team/model"
	"team/web"
	"team/web/orm"
)

// Install prepares databases
type Install struct {
	Done    bool
	IsError bool
	Status  []string
}

// Register implements web.Controller interface.
func (i *Install) Register(group *web.Router) {
	group.POST("/configure", i.configure)
	group.GET("/status", i.status)
	group.POST("/admin", i.createAdmin)
}

func (i *Install) configure(c *web.Context) {
	appName := c.FormValue("appName")
	mysqlHost := c.FormValue("mysqlHost")
	mysqlUser := c.FormValue("mysqlUser")
	mysqlPswd := c.FormValue("mysqlPswd")
	mysqlDB := c.FormValue("mysqlDB")

	i.Done = false
	i.IsError = false
	i.Status = []string{"连接数据库..."}

	go i.setup(appName, mysqlHost, mysqlUser, mysqlPswd, mysqlDB)

	c.JSON(200, web.JObject{})
}

func (i *Install) status(c *web.Context) {
	c.JSON(200, web.JObject{
		"data": web.JObject{
			"done":    i.Done,
			"isError": i.IsError,
			"status":  i.Status,
		},
	})
}

func (i *Install) createAdmin(c *web.Context) {
	account := c.FormValue("account")
	name := c.FormValue("name")
	pswd := c.FormValue("pswd")

	hash := md5.New()
	hash.Write([]byte(pswd))

	user := &model.User{
		Account:  account,
		Name:     name,
		Password: fmt.Sprintf("%X", hash.Sum(nil)),
		IsSu:     true,
	}

	_, err := orm.Insert(user)
	if err != nil {
		web.Logger.Error("Create default admin failed: %v", err)
		c.JSON(200, web.JObject{"err": "创建默认管理员失败"})
	} else {
		model.Environment.Installed = true
		model.Environment.Save()
		c.JSON(200, web.JObject{})
	}
}

func (i *Install) setup(appName, host, user, pswd, db string) {
	err := orm.ConnectDB(host, user, pswd, db, 64)
	if err != nil {
		i.IsError = true
		i.Status = append(i.Status, "无法连接数据："+err.Error())
		return
	}

	type Job struct {
		Table  string
		Schema interface{}
	}

	jobs := []Job{
		{Table: "user", Schema: &model.User{}},
		{Table: "project", Schema: &model.Project{}},
		{Table: "project member", Schema: &model.ProjectMember{}},
		{Table: "task", Schema: &model.Task{}},
		{Table: "task attachment", Schema: &model.TaskAttachment{}},
		{Table: "task event", Schema: &model.TaskEvent{}},
		{Table: "task comment", Schema: &model.TaskComment{}},
		{Table: "document", Schema: &model.Document{}},
		{Table: "notice", Schema: &model.Notice{}},
		{Table: "share", Schema: &model.Share{}},
	}

	for _, job := range jobs {
		i.Status = append(i.Status, "创建数据表: "+job.Table)

		err = orm.CreateTable(job.Schema)
		if err != nil {
			i.IsError = true
			i.Status = append(i.Status, "出错了："+err.Error())
			return
		}
	}

	i.Done = true
	i.Status = append(i.Status, "应用配置完成!")

	model.Environment = &model.Env{
		Installed: false,
		AppName:   appName,
		AppPort:   ":8080",
		MySQL: model.MySQL{
			Host:     host,
			User:     user,
			Password: pswd,
			Database: db,
			MaxConns: 64,
		},
	}
}
