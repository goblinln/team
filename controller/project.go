package controller

import (
	"time"

	"team/model/project"
	"team/model/task"
	"team/model/user"
	"team/web"
)

// Project controller
type Project int

// Register implements web.Controller interface.
func (p *Project) Register(group *web.Router) {
	group.GET("/:id", p.info)
	group.GET("/:id/summary", p.summary)
	group.GET("/mine", p.mine)

	group.PUT("/:id/desc", p.setDesc)

	group.GET("/:id/invites", p.getInviteList)
	group.POST("/:id/member", p.addMember)
	group.PUT(`/:id/member/{uid:[\d]+}`, p.editMember)
	group.DELETE(`/:id/member/{uid:[\d]+}`, p.deleteMember)

	group.GET("/:id/milestone/list", p.getMilestones)
	group.POST("/:id/milestone", p.addMilestone)
	group.PUT(`/:id/milestone/{mid:[\d]+}`, p.editMilestone)
	group.DELETE(`/:id/milestone/{mid:[\d]+}`, p.delMilestone)

	group.GET(`/:id/week/{start:[\d]+}`, p.getWeekReport)
}

func (*Project) info(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	proj := project.Find(pid)
	web.Assert(proj != nil, "指定项目不存在或已被删除")

	c.JSON(200, web.Map{"data": proj.Info()})
}

func (*Project) summary(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	proj := project.Find(pid)
	web.Assert(proj != nil, "指定项目不存在或已被删除")

	c.JSON(200, web.Map{"data": proj.Summary()})
}

func (*Project) mine(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	projs, err := project.GetAllByUser(uid)
	web.AssertError(err)
	c.JSON(200, web.Map{"data": projs})
}

func (*Project) setDesc(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	desc := c.PostFormValue("desc").String()
	proj := project.Find(pid)
	web.Assert(proj != nil, "指定项目不存在或已被删除")
	web.AssertError(proj.SetDesc(desc))
	c.JSON(200, web.Map{})
}

func (*Project) getInviteList(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")

	members := map[int64]bool{}
	for _, one := range proj.Members {
		members[one.UID] = true
	}

	users, err := user.GetAll()
	web.AssertError(err)

	invites := []*user.User{}
	for _, one := range users {
		if _, ok := members[one.ID]; !ok {
			invites = append(invites, one)
		}
	}

	c.JSON(200, web.Map{"data": invites})
}

func (*Project) addMember(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	uid := c.PostFormValue("uid").MustInt("无效的用户ID")
	isAdmin, _ := c.PostFormValue("isAdmin").Bool()
	role := c.PostFormValue("role").MustInt("无效的职能")

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")

	for _, one := range proj.Members {
		if one.UID == uid {
			c.JSON(200, web.Map{})
			return
		}
	}

	u := user.Find(uid)
	web.Assert(u != nil && !u.IsLocked, "无效的成员ID")
	web.AssertError(proj.AddMember(uid, int8(role), isAdmin))

	c.JSON(200, web.Map{})
}

func (*Project) editMember(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	uid := c.RouteValue("uid").MustInt("")
	role := c.PostFormValue("role").MustInt("无效的职能")
	isAdmin, _ := c.PostFormValue("isAdmin").Bool()

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")

	for _, one := range proj.Members {
		if uid == one.UID {
			one.Role = int8(role)
			one.IsAdmin = isAdmin
			web.AssertError(one.Save())
			c.JSON(200, web.Map{})
			return
		}
	}

	c.JSON(200, web.Map{"err": "成员不存在或已被删除"})
}

func (*Project) deleteMember(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	uid := c.RouteValue("uid").MustInt("")

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")

	proj.DelMember(uid)
	c.JSON(200, web.Map{})
}

func (*Project) getMilestones(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")

	c.JSON(200, web.Map{"data": proj.GetMilestones()})
}

func (*Project) addMilestone(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	name := c.PostFormValue("name").MustString("里程碑名不可为空")
	startTime, _ := time.Parse("2006-01-02", c.PostFormValue("startTime").MustString("开始时间不可为空"))
	endTime, _ := time.Parse("2006-01-02", c.PostFormValue("endTime").MustString("终止时间不可为空"))
	desc := c.PostFormValue("desc").String()

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")
	web.AssertError(proj.AddMilestone(name, desc, startTime, endTime))

	c.JSON(200, web.Map{})
}

func (*Project) editMilestone(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	mid := c.RouteValue("mid").MustInt("")

	name := c.PostFormValue("name").MustString("里程碑名不可为空")
	startTime, _ := time.Parse("2006-01-02", c.PostFormValue("startTime").MustString("开始时间不可为空"))
	endTime, _ := time.Parse("2006-01-02", c.PostFormValue("endTime").MustString("终止时间不可为空"))
	desc := c.PostFormValue("desc").String()

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")
	web.AssertError(proj.EditMilestone(mid, name, desc, startTime, endTime))

	c.JSON(200, web.Map{})
}

func (*Project) delMilestone(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	mid := c.RouteValue("mid").MustInt("")

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")

	proj.DelMilestone(mid)
	c.JSON(200, web.Map{})
}

func (*Project) getWeekReport(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	start := c.RouteValue("start").MustInt("")

	proj := project.Find(pid)
	web.Assert(proj != nil, "项目不存在或已被删除")

	c.JSON(200, web.Map{"data": task.GetWeekReport(pid, start)})
}
