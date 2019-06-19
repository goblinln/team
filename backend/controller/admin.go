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
	group.PATCH("/user/:id/lock", a.lockUser)
	group.DELETE("/user/:id", a.deleteUser)
	group.GET("/user/list", a.users)

	group.POST("/project", a.addProject)
	group.PUT("/project/:id", a.editProject)
	group.DELETE("/project/:id", a.deleteProject)
	group.GET("/project/list", a.projects)
}

func (a *Admin) addUser(c *web.Context) {
	account := c.PostFormValue("account")
	name := c.PostFormValue("name")
	pswd := c.PostFormValue("pswd")
	cfmPswd := c.PostFormValue("cfmPswd")
	isSu := atoi(c.PostFormValue("isSu"))

	if cfmPswd != pswd {
		c.JSON(200, web.Map{"err": "两次输入的新密码不一致"})
		return
	}

	pswdEncoder := md5.New()
	pswdEncoder.Write([]byte(pswd))
	pswdEncoded := fmt.Sprintf("%X", pswdEncoder.Sum(nil))

	rows, err := orm.Query("SELECT COUNT(*) FROM `user` WHERE `account`=? OR `name`=?", account, name)
	if err != nil {
		c.JSON(200, web.Map{"err": "创建帐号失败，代码#1"})
		return
	}

	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	if count > 0 {
		c.JSON(200, web.Map{"err": "帐号或角色名已存在"})
		return
	}

	user := &model.User{
		Account:  account,
		Name:     name,
		Avatar:   "",
		Password: pswdEncoded,
		IsSu:     isSu == 1,
		IsLocked: false,
	}

	rs, err := orm.Insert(user)
	if err != nil {
		c.JSON(200, web.Map{"err": "写入新帐号失败"})
		return
	}

	user.ID, _ = rs.LastInsertId()
	model.Cache.SetUser(user)

	c.JSON(200, web.Map{})
}

func (a *Admin) editUser(c *web.Context) {
	uid := atoi(c.RouteValue("id"))
	account := c.PostFormValue("account")
	name := c.PostFormValue("name")
	isSu := atoi(c.PostFormValue("isSu"))

	rows, err := orm.Query("SELECT * FROM `user` WHERE `account`=? OR `name`=?", account, name)
	if err != nil {
		c.JSON(200, web.Map{"err": "修改帐号失败，代码#1"})
		return
	}

	defer rows.Close()

	user := &model.User{}
	if rows.Next() {
		orm.Scan(rows, user)
		if user.ID != uid {
			c.JSON(200, web.Map{"err": "帐号或角色名已存在"})
			return
		}
	} else {
		user.ID = uid
		err = orm.Read(user)
		if err != nil {
			c.JSON(200, web.Map{"err": "帐号不存在或已被删除"})
			return
		}
	}

	user.Account = account
	user.Name = name
	user.IsSu = isSu == 1
	err = orm.Update(user)
	if err != nil {
		c.JSON(200, web.Map{"err": "写入帐号失败"})
		return
	}

	model.Cache.SetUser(user)
	c.JSON(200, web.Map{})
}

func (a *Admin) lockUser(c *web.Context) {
	uid := atoi(c.RouteValue("id"))
	user := model.FindUser(uid)
	if user == nil {
		c.JSON(200, web.Map{"err": "帐号不存在或已被删除"})
		return
	}

	user.IsLocked = !user.IsLocked
	err := orm.Update(user)
	if err != nil {
		c.JSON(200, web.Map{"err": "写入帐号失败"})
		return
	}

	c.JSON(200, web.Map{})
}

func (a *Admin) deleteUser(c *web.Context) {
	uid := atoi(c.RouteValue("id"))
	orm.Delete("user", uid)
	model.Cache.DeleteUser(uid)
	c.JSON(200, web.Map{})
}

func (a *Admin) users(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `user`")
	if err != nil {
		c.JSON(200, web.Map{"err": "拉取用户列表失败"})
		return
	}

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
	name := c.PostFormValue("name")
	admin := atoi(c.PostFormValue("admin"))
	role := atoi(c.PostFormValue("role"))

	rows, err := orm.Query("SELECT COUNT(*) FROM `project` WHERE `name`=?", name)
	if err != nil {
		c.JSON(200, web.Map{"err": "创建项目失败，代码#1"})
		return
	}

	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	if count > 0 {
		c.JSON(200, web.Map{"err": "同名项目已存在"})
		return
	}

	user := &model.User{ID: admin}
	err = orm.Read(user)
	if err != nil {
		c.JSON(200, web.Map{"err": "默认管理员不存在或已被删除"})
		return
	}

	if user.IsLocked {
		c.JSON(200, web.Map{"err": "默认管理员当前被禁止登录"})
		return
	}

	proj := &model.Project{
		Name:     name,
		Branches: []string{"默认"},
	}
	rs, err := orm.Insert(proj)
	if err != nil {
		c.JSON(200, web.Map{"err": "写入新项目失败"})
		return
	}

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
	pid := atoi(c.RouteValue("id"))
	name := c.PostFormValue("name")

	rows, err := orm.Query("SELECT * FROM `project` WHERE `name`=?", name)
	if err != nil {
		c.JSON(200, web.Map{"err": "修改项目失败，代码#1"})
		return
	}

	defer rows.Close()

	proj := &model.Project{}
	if rows.Next() {
		orm.Scan(rows, proj)
		if proj.ID != pid {
			c.JSON(200, web.Map{"err": "同名项目已存在"})
			return
		}
	} else {
		proj.ID = pid
		err = orm.Read(proj)
		if err != nil {
			c.JSON(200, web.Map{"err": "项目不存在或已被删除"})
			return
		}
	}

	proj.Name = name
	err = orm.Update(proj)
	if err != nil {
		c.JSON(200, web.Map{"err": "写入新项目失败"})
		return
	}

	model.Cache.SetProject(proj)
	c.JSON(200, web.Map{})
}

func (a *Admin) deleteProject(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	orm.Delete("project", pid)
	orm.Exec("DELETE FROM `projectmember` WHERE `pid`=?", pid)
	orm.Exec("DELETE FROM `task` WHERE `pid`=?", pid)
	model.Cache.DeleteProject(pid)
	c.JSON(200, web.Map{})
}

func (a *Admin) projects(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `project`")
	if err != nil {
		c.JSON(200, web.Map{"err": "拉取项目列表失败"})
		return
	}

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
