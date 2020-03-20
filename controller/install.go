package controller

import (
	"crypto/md5"
	"fmt"

	"team/common/orm"
	"team/common/web"
	"team/config"
	"team/model/document"
	"team/model/notice"
	"team/model/project"
	"team/model/share"
	"team/model/task"
	"team/model/user"
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
	appName := c.FormValue("name").MustString("无效的站点名称")
	appPort := c.FormValue("port").MustInt("无效的监听端口")
	appLoginType := c.FormValue("loginType").MustInt("无效的认证方式")
	mysqlHost := c.FormValue("mysqlHost").MustString("无效MySQL地址")
	mysqlUser := c.FormValue("mysqlUser").MustString("无效MySQL用户")
	mysqlPswd := c.FormValue("mysqlPswd").String()
	mysqlDB := c.FormValue("mysqlDB").MustString("数据名不可为空")

	web.Assert(appPort > 0 && appPort < 65535, "无效的端口参数")

	i.Done = false
	i.IsError = false
	i.Status = []string{"连接数据库..."}

	config.App.Name = appName
	config.App.Port = int(appPort)
	config.App.Auth = config.AuthKind(appLoginType)
	config.MySQL = &config.MySQLInfo{
		Host:     mysqlHost,
		User:     mysqlUser,
		Password: mysqlPswd,
		Database: mysqlDB,
	}

	switch config.App.Auth {
	case config.AuthKindSMTP:
		smtpHost := c.FormValue("smtpLoginHost").MustString("无效的SMTP地址")
		smtpPort := int(c.FormValue("smtpLoginPort").MustInt("端口号不可为空"))
		smtpPlain := c.FormValue("smtpLoginKind").MustInt("未选择SMTP登录方式") == 0
		smtpTLS, _ := c.FormValue("smtpLoginTLS").Bool()
		smtpSkipVerfiy, _ := c.FormValue("smtpLoginSkipVerify").Bool()
		config.UseSMTPAuth(smtpHost, smtpPort, smtpPlain, smtpTLS, smtpSkipVerfiy)
	}

	go func() {
		err := orm.OpenDB("mysql", config.MySQL.URL())
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
			{Table: "project milestone", Schema: &project.Milestone{}},
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
	web.Assert(config.Save() == nil, "保存配置失败")

	config.Installed = true
	c.JSON(200, web.Map{})
}
