package controller

import (
	"crypto/md5"
	"fmt"

	"team/model"
	"team/orm"
	"team/web"
)

// Install prepares databases
type Install struct {
	Done    bool
	IsError bool
	Status  []string
}

// Register implements web.Controller interface.
func (i *Install) Register(group *web.Router) {
	group.GET("", i.index)
	group.POST("/configure", i.configure)
	group.GET("/status", i.status)
	group.POST("/admin", i.createAdmin)
}

func (i *Install) index(c *web.Context) {
	c.ResponseHeader().Set("Content-Type", "text/html")
	c.File(200, "./assets/app.html")
}

func (i *Install) configure(c *web.Context) {
	appPort := atoi(c.FormValue("port"))
	mysqlHost := c.FormValue("mysqlHost")
	mysqlUser := c.FormValue("mysqlUser")
	mysqlPswd := c.FormValue("mysqlPswd")
	mysqlDB := c.FormValue("mysqlDB")

	if appPort <= 0 || appPort > 65535 {
		c.JSON(200, &web.JObject{"err": "无效的端口参数"})
		return
	}

	i.Done = false
	i.IsError = false
	i.Status = []string{"连接数据库..."}

	go i.setup(appPort, mysqlHost, mysqlUser, mysqlPswd, mysqlDB)

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

func (i *Install) setup(appPort int64, host, user, pswd, db string) {
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
		AppPort:   fmt.Sprintf(":%d", appPort),
		MySQL: model.MySQL{
			Host:     host,
			User:     user,
			Password: pswd,
			Database: db,
			MaxConns: 64,
		},
	}
}
