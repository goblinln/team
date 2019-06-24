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
	appPort := c.FormValue("port").MustInt("无效的监听端口")
	mysqlHost := c.FormValue("mysqlHost").MustString("无效MySQL地址")
	mysqlUser := c.FormValue("mysqlUser").MustString("无效MySQL用户")
	mysqlPswd := c.FormValue("mysqlPswd").String()
	mysqlDB := c.FormValue("mysqlDB").MustString("数据名不可为空")

	web.Assert(appPort > 0 && appPort < 65535, "无效的端口参数")

	i.Done = false
	i.IsError = false
	i.Status = []string{"连接数据库..."}

	go i.setup(appPort, &model.MySQL{Host: mysqlHost, User: mysqlUser, Password: mysqlPswd, Database: mysqlDB})
	c.JSON(200, web.Map{})
}

func (i *Install) status(c *web.Context) {
	c.JSON(200, web.Map{
		"data": web.Map{
			"done":    i.Done,
			"isError": i.IsError,
			"status":  i.Status,
		},
	})
}

func (i *Install) createAdmin(c *web.Context) {
	account := c.FormValue("account").MustString("帐号不可为空")
	name := c.FormValue("name").MustString("显示名称不可为空")
	pswd := c.FormValue("pswd").MustString("超级管理员必须设置密码")

	hash := md5.New()
	hash.Write([]byte(pswd))

	user := &model.User{
		Account:  account,
		Name:     name,
		Password: fmt.Sprintf("%X", hash.Sum(nil)),
		IsSu:     true,
	}

	_, err := orm.Insert(user)
	web.Assert(err == nil, "创建默认管理员失败")

	model.Environment.Installed = true
	model.Environment.Save()
	c.JSON(200, web.Map{})
}

func (i *Install) setup(appPort int64, mysql *model.MySQL) {
	err := orm.OpenDB("mysql", mysql.Addr())
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
		MySQL:     mysql,
	}
}
