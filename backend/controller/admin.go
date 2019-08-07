package controller

import (
	"crypto/md5"
	"fmt"

	"team/model"
	"team/orm"
	"team/web"
)

// Admin controller
type Admin int

// Register implements web.Controller interface.
func (a *Admin) Register(group *web.Router) {
	group.POST("/user", a.addUser)
	group.PUT("/user/:id", a.editUser)
	group.PUT("/user/:id/lock", a.lockUser)
	group.DELETE("/user/:id", a.deleteUser)
	group.GET("/user/list", a.users)

	group.POST("/project", a.addProject)
	group.PUT("/project/:id", a.editProject)
	group.DELETE("/project/:id", a.deleteProject)
	group.GET("/project/list", a.projects)
}

func (a *Admin) addUser(c *web.Context) {
	account := c.PostFormValue("account").MustString("请填写新用户帐号")
	name := c.PostFormValue("name").MustString("新用户名称不可为空")
	pswd := c.PostFormValue("pswd").MustString("初始密码不可为空")
	cfmPswd := c.PostFormValue("cfmPswd").String()
	isSu, _ := c.PostFormValue("isSu").Bool()

	web.Assert(cfmPswd == pswd, "两次输入的新密码不一致")

	pswdEncoder := md5.New()
	pswdEncoder.Write([]byte(pswd))
	pswdEncoded := fmt.Sprintf("%X", pswdEncoder.Sum(nil))

	rows, err := orm.Query("SELECT COUNT(*) FROM `user` WHERE `account`=? OR `name`=?", account, name)
	web.Assert(err == nil, "创建帐号失败，代码#1")
	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	web.Assert(count == 0, "帐号或角色名已存在")

	user := &model.User{
		Account:  account,
		Name:     name,
		Avatar:   "",
		Password: pswdEncoded,
		IsSu:     isSu,
		IsLocked: false,
	}

	rs, err := orm.Insert(user)
	web.Assert(err == nil, "写入新帐号失败")

	user.ID, _ = rs.LastInsertId()
	model.Cache.SetUser(user)

	c.JSON(200, web.Map{})
}

func (a *Admin) editUser(c *web.Context) {
	uid := c.RouteValue("id").MustInt("")
	account := c.PostFormValue("account").MustString("用户帐号不可为空")
	name := c.PostFormValue("name").MustString("用户显示名称不可为空")
	isSu, _ := c.PostFormValue("isSu").Bool()

	rows, err := orm.Query("SELECT * FROM `user` WHERE `account`=? OR `name`=?", account, name)
	web.Assert(err == nil, "修改帐号失败，代码#1")
	defer rows.Close()

	user := &model.User{}
	if rows.Next() {
		orm.Scan(rows, user)
		web.Assert(user.ID == uid, "帐号或角色名已存在")
	} else {
		user.ID = uid
		err = orm.Read(user)
		web.Assert(err == nil, "帐号不存在或已被删除")
	}

	user.Account = account
	user.Name = name
	user.IsSu = isSu
	err = orm.Update(user)
	web.Assert(err == nil, "写入帐号失败")

	model.Cache.SetUser(user)
	c.JSON(200, web.Map{})
}

func (a *Admin) lockUser(c *web.Context) {
	uid := c.RouteValue("id").MustInt("")
	user := model.FindUser(uid)
	web.Assert(user != nil, "帐号不存在或已被删除")

	user.IsLocked = !user.IsLocked
	err := orm.Update(user)
	web.Assert(err == nil, "写入帐号失败")

	c.JSON(200, web.Map{})
}

func (a *Admin) deleteUser(c *web.Context) {
	uid := c.RouteValue("id").MustInt("")
	orm.Delete("user", uid)
	model.Cache.DeleteUser(uid)
	c.JSON(200, web.Map{})
}

func (a *Admin) users(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `user`")
	web.Assert(err == nil, "拉取用户列表失败")
	defer rows.Close()

	users := []*model.User{}
	for rows.Next() {
		user := &model.User{}
		err = orm.Scan(rows, user)
		if err == nil {
			model.Cache.SetUser(user)
			users = append(users, user)
		}
	}

	c.JSON(200, web.Map{"data": users})
}

func (a *Admin) addProject(c *web.Context) {
	name := c.PostFormValue("name").MustString("非法项目名称")
	admin := c.PostFormValue("admin").MustInt("无效管理员ID")
	role, _ := c.PostFormValue("role").Int()

	rows, err := orm.Query("SELECT COUNT(*) FROM `project` WHERE `name`=?", name)
	web.Assert(err == nil, "创建项目失败，代码#1")
	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	web.Assert(count == 0, "同名项目已存在")

	user := &model.User{ID: admin}
	err = orm.Read(user)
	web.Assert(err == nil, "默认管理员不存在或已被删除")
	web.Assert(!user.IsLocked, "默认管理员当前被禁止登录")

	proj := &model.Project{
		Name:     name,
		Branches: []string{"默认"},
	}
	rs, err := orm.Insert(proj)
	web.Assert(err == nil, "写入新项目失败")

	proj.ID, _ = rs.LastInsertId()
	model.Cache.SetProject(proj)

	orm.Insert(&model.ProjectMember{
		UID:     admin,
		PID:     proj.ID,
		Role:    int8(role),
		IsAdmin: true,
	})

	c.JSON(200, web.Map{})
}

func (a *Admin) editProject(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	name := c.PostFormValue("name").MustString("无效项目名")

	rows, err := orm.Query("SELECT * FROM `project` WHERE `name`=?", name)
	web.Assert(err == nil, "修改项目失败，代码#1")
	defer rows.Close()

	proj := &model.Project{}
	if rows.Next() {
		orm.Scan(rows, proj)
		web.Assert(proj.ID == pid, "同名项目已存在")
	} else {
		proj.ID = pid
		err = orm.Read(proj)
		web.Assert(err == nil, "项目不存在或已被删除")
	}

	proj.Name = name
	err = orm.Update(proj)
	web.Assert(err == nil, "写入新项目失败")

	model.Cache.SetProject(proj)
	c.JSON(200, web.Map{})
}

func (a *Admin) deleteProject(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	orm.Delete("project", pid)
	orm.Exec("DELETE FROM `projectmember` WHERE `pid`=?", pid)
	orm.Exec("DELETE FROM `task` WHERE `pid`=?", pid)
	model.Cache.DeleteProject(pid)
	c.JSON(200, web.Map{})
}

func (a *Admin) projects(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `project`")
	web.Assert(err == nil, "拉取项目列表失败")
	defer rows.Close()

	projs := []*model.Project{}
	for rows.Next() {
		proj := &model.Project{}
		err = orm.Scan(rows, proj)
		if err == nil {
			model.Cache.SetProject(proj)
			projs = append(projs, proj)
		}
	}

	c.JSON(200, web.Map{"data": projs})
}
