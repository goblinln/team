package controller

import (
	"team/model"
	"team/orm"
	"team/web"
)

// Project controller
type Project int

// Register implements web.Controller interface.
func (p *Project) Register(group *web.Router) {
	group.GET("/:id", p.info)
	group.GET("/mine", p.mine)
	group.POST("/:id/branch", p.addBranch)
	group.GET("/:id/invites", p.getInviteList)
	group.POST("/:id/member", p.addMember)
	group.PUT(`/:id/member/{uid:[\d]+}`, p.editMember)
	group.DELETE(`/:id/member/{uid:[\d]+}`, p.deleteMember)
	group.GET(`/:id/report/{from:[\d]+}`, p.getReports)
}

func (p *Project) info(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	c.JSON(200, web.Map{"data": model.MakeProjectInfo(pid)})
}

func (p *Project) mine(c *web.Context) {
	uid := c.Session.Get("uid").(int64)

	rows, err := orm.Query("SELECT `pid` FROM `projectmember` WHERE `uid`=?", uid)
	if err != nil {
		c.JSON(200, web.Map{"err": "获取项目列表失败"})
		return
	}

	defer rows.Close()

	projs := []map[string]interface{}{}
	pid := int64(0)
	for rows.Next() {
		err := rows.Scan(&pid)
		if err != nil {
			continue
		}

		projs = append(projs, model.MakeProjectInfo(pid))
	}

	c.JSON(200, web.Map{"data": projs})
}

func (p *Project) addBranch(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	branch := c.PostFormValue("branch")
	proj := model.FindProject(pid)
	if proj == nil {
		c.JSON(200, web.Map{"err": "项目不存在或已被删除"})
		return
	}

	for _, b := range proj.Branches {
		if b == branch {
			c.JSON(200, web.Map{"err": "同名分支已存在"})
			return
		}
	}

	proj.Branches = append(proj.Branches, branch)
	err := orm.Update(proj)
	if err != nil {
		c.JSON(200, web.Map{"err": "写入修改失败"})
		return
	}

	c.JSON(200, web.Map{})
}

func (p *Project) getInviteList(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	proj := model.FindProject(pid)
	if proj == nil {
		c.JSON(200, web.Map{"err": "项目不存在或已被删除"})
		return
	}

	rows, err := orm.Query("SELECT `uid` FROM `projectmember` WHERE `pid`=?", pid)
	if err != nil {
		c.JSON(200, web.Map{"err": "获取成员列表失败"})
		return
	}

	defer rows.Close()

	members := make(map[int64]bool)
	for rows.Next() {
		uid := int64(0)
		rows.Scan(&uid)
		members[uid] = true
	}

	userRows, err := orm.Query("SELECT * FROM `user`")
	if err != nil {
		c.JSON(200, web.Map{"err": "获取用户列表失败"})
		return
	}

	defer userRows.Close()

	valids := []*model.User{}
	for userRows.Next() {
		user := &model.User{}
		orm.Scan(userRows, user)

		if _, ok := members[user.ID]; !ok && !user.IsLocked {
			valids = append(valids, user)
		}
	}

	c.JSON(200, web.Map{"data": valids})
}

func (p *Project) addMember(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	uid := atoi(c.PostFormValue("uid"))
	isAdmin := atoi(c.PostFormValue("isAdmin")) == 1
	role := atoi(c.PostFormValue("role"))

	proj := model.FindProject(pid)
	if proj == nil {
		c.JSON(200, web.Map{"err": "项目不存在或已被删除"})
		return
	}

	rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=?", pid, uid)
	if err != nil {
		c.JSON(200, web.Map{"err": "获取成员列表失败"})
		return
	}

	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	if count > 0 {
		c.JSON(200, web.Map{})
		return
	}

	user := model.FindUser(uid)
	if user == nil || user.IsLocked {
		c.JSON(200, web.Map{"err": "无效的成员ID"})
		return
	}

	_, err = orm.Insert(&model.ProjectMember{
		PID:     pid,
		UID:     uid,
		Role:    int8(role),
		IsAdmin: isAdmin,
	})
	if err != nil {
		c.JSON(200, web.Map{"err": "写入数据库失败"})
		return
	}

	c.JSON(200, web.Map{})
}

func (p *Project) editMember(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	uid := atoi(c.RouteValue("uid"))
	role := atoi(c.PostFormValue("role"))
	isAdmin := atoi(c.PostFormValue("isAdmin")) == 1

	member := &model.ProjectMember{
		PID: pid,
		UID: uid,
	}
	err := orm.Read(member, "pid", "uid")
	if err != nil {
		if err == orm.ErrNotFound {
			c.JSON(200, web.Map{"err": "参数错误"})
			return
		}

		c.JSON(200, web.Map{"err": "写入数据库失败"})
		return
	}

	member.Role = int8(role)
	member.IsAdmin = isAdmin
	err = orm.Update(member)
	if err != nil {
		c.JSON(200, web.Map{"err": "写入数据库失败"})
		return
	}

	c.JSON(200, web.Map{})
}

func (p *Project) deleteMember(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	uid := atoi(c.RouteValue("uid"))

	orm.Exec("DELETE from `projectmember` WHERE `pid`=? AND `uid`=?", pid, uid)
	c.JSON(200, web.Map{})
}

func (p *Project) getReports(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	from := atoi(c.RouteValue("from"))
	end := from + 3600*24*7

	unarchived := []map[string]interface{}{}
	archived := []map[string]interface{}{}

	unarchivedRows, err := orm.Query("SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringtop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `pid`=? AND `state`<4 AND UNIX_TIMESTAMP(`endtime`)<=? AND (`archivetime`=? OR UNIX_TIMESTAMP(`archivetime`)>?)", pid, end, model.TaskTimeInfinite.Format(orm.TimeFormat), end)
	if err != nil {
		c.JSON(200, web.Map{"err": "拉取该周内未归档任务列表出错"})
		return
	}

	defer unarchivedRows.Close()

	for unarchivedRows.Next() {
		task := &model.Task{}
		err = orm.Scan(unarchivedRows, task)
		if err == nil {
			unarchived = append(unarchived, model.MakeTaskBrief(task))
		}
	}

	archivedRows, err := orm.Query("SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringTop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `pid`=? AND `state`=4 AND UNIX_TIMESTAMP(`archivetime`)>=? AND UNIX_TIMESTAMP(`archivetime`)<=?", pid, from, end)
	if err != nil {
		c.JSON(200, web.Map{"err": "拉取该周内归档任务列表出错"})
		return
	}

	defer archivedRows.Close()

	for archivedRows.Next() {
		task := &model.Task{}
		err = orm.Scan(archivedRows, task)
		if err == nil {
			archived = append(archived, model.MakeTaskBrief(task))
		}
	}

	c.JSON(200, web.Map{
		"data": map[string]interface{}{
			"archived":   archived,
			"unarchived": unarchived,
		},
	})
}
