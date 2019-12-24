package controller

import (
	"crypto/md5"
	"fmt"

	"team/config"
	"team/model/document"
	"team/model/notice"
	"team/model/project"
	"team/model/share"
	"team/model/task"
	"team/model/user"
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
	group.POST("/configure", i.configure)
	group.GET("/status", i.status)
	group.POST("/admin", i.createAdmin)
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

	config.Default.AppPort = fmt.Sprintf(":%d", appPort)
	config.Default.MySQL = &config.MySQL{
		Host:     mysqlHost,
		User:     mysqlUser,
		Password: mysqlPswd,
		Database: mysqlDB,
	}

	go func() {
		err := orm.OpenDB("mysql", config.Default.GetMySQLAddr())
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
			{Table: "user", Schema: &user.User{}},
			{Table: "notice", Schema: &notice.Notice{}},
			{Table: "project", Schema: &project.Project{}},
			{Table: "project member", Schema: &project.Member{}},
			{Table: "task", Schema: &task.Task{}},
			{Table: "task attachment", Schema: &task.Attachment{}},
			{Table: "task event", Schema: &task.Event{}},
			{Table: "task comment", Schema: &task.Comment{}},
			{Table: "document", Schema: &document.Document{}},
			{Table: "share", Schema: &share.Share{}},
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
	}()

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

	user := &user.User{
		Account:  account,
		Name:     name,
		Password: fmt.Sprintf("%X", hash.Sum(nil)),
		IsSu:     true,
	}

	_, err := orm.Insert(user)
	web.Assert(err == nil, "创建默认管理员失败")

	config.Default.Installed = true
	config.Default.Save()
	c.JSON(200, web.Map{})
}
