package controller

import (
	"strconv"
	"time"

	"team/model"
	"team/orm"
	"team/web"
)

// Task controller
type Task int

// Register implements web.Controller interface.
func (t *Task) Register(group *web.Router) {
	group.POST("", t.create)
	group.POST("/:id/back", t.moveBack)
	group.POST("/:id/next", t.moveNext)
	group.GET("/:id", t.info)
	group.DELETE("/:id", t.delete)

	group.GET("/mine", t.mine)
	group.GET("/project/:id", t.project)

	group.PUT("/:id/name", t.setName)
	group.PUT("/:id/creator", t.setCreator)
	group.PUT("/:id/developer", t.setDeveloper)
	group.PUT("/:id/tester", t.setTester)
	group.PUT("/:id/weight", t.setWeight)
	group.PUT("/:id/time", t.setTime)
	group.PUT("/:id/content", t.setContent)
	group.PUT("/:id/status", t.setStatus)
	group.POST("/:id/comment", t.addComment)
}

func (t *Task) create(c *web.Context) {
	name := c.PostFormValue("name").MustString("任务名不可空")
	proj := c.PostFormValue("proj").MustInt("无效的项目ID")
	branch := c.PostFormValue("branch").MustInt("无效的分支")
	weight := c.PostFormValue("weight").MustInt("无效的优先级")
	creator, _ := c.PostFormValue("creator").Int()
	developer := c.PostFormValue("developer").MustInt("开发人员未指定")
	tester := c.PostFormValue("tester").MustInt("测试人员未指定")
	startTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("startTime").MustString("开始时间未指定"))
	endTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("endTime").MustString("结束时间未指定"))
	tags, _ := c.PostFormValue("tags[]").Ints()
	content := c.PostFormValue("content").MustString("任务内容不可空")

	uid := c.Session.Get("uid").(int64)
	me := model.FindUser(uid)

	if creator == 0 {
		creator = uid
	}

	uploaded := []*model.TaskAttachment{}
	fhs, ok := c.MultipartForm().File["files[]"]
	if ok {
		helper := new(File)
		for _, fh := range fhs {
			url, _ := helper.save(fh, uid)
			uploaded = append(uploaded, &model.TaskAttachment{
				Name: fh.Filename,
				Path: url,
			})
		}
	}

	task := &model.Task{
		PID:         proj,
		Branch:      int8(branch),
		Creator:     creator,
		Developer:   developer,
		Tester:      tester,
		Name:        name,
		BringTop:    me.IsSu,
		Weight:      int8(weight),
		State:       0,
		StartTime:   startTime,
		EndTime:     endTime,
		ArchiveTime: model.TaskTimeInfinite,
		Tags:        []int{},
		Content:     content,
	}
	for _, t := range tags {
		task.Tags = append(task.Tags, int(t))
	}

	rs, err := orm.Insert(task)
	web.Assert(err == nil, "写入任务信息失败")

	tid, err := rs.LastInsertId()
	web.Assert(err == nil, "读取新任务ID失败")

	task.ID = tid

	for i := 0; i < len(uploaded); i++ {
		uploaded[i].TID = tid
		orm.Insert(uploaded[i])
	}

	model.AfterTaskOperation(task, uid, model.TaskEventCreate, "")
	c.JSON(200, web.Map{})
}

func (t *Task) moveBack(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	tid := c.RouteValue("id").MustInt("")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	validOperator := false
	switch task.State {
	case 1:
		validOperator = uid == task.Developer
	case 2:
		validOperator = uid == task.Tester || uid == task.Developer
	case 3:
		validOperator = uid == task.Creator || uid == task.Tester
	case 4:
		validOperator = uid == task.Creator
		task.ArchiveTime = model.TaskTimeInfinite
	default:
		web.Assert(false, "任务不可回退了")
	}

	if !validOperator {
		rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=? AND `isadmin`=1", task.PID, uid)
		web.Assert(err == nil, "您不可回退当前状态")
		defer rows.Close()

		web.Assert(rows.Next(), "你无权回退该任务")

		count := 0
		web.Assert(rows.Scan(&count) == nil && count > 0, "你无权回退该任务")
	}

	task.State--
	err = orm.Update(task)
	web.Assert(err == nil, "修改任务失败")

	model.AfterTaskOperation(task, uid, model.TaskEventMoveBack, "")
	c.JSON(200, web.Map{})
}

func (t *Task) moveNext(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	uid := c.Session.Get("uid").(int64)
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	ev := model.TaskEventUnderway
	validOperator := false
	switch task.State {
	case 0:
		validOperator = uid == task.Developer
		ev = model.TaskEventUnderway
	case 1:
		validOperator = uid == task.Developer
		ev = model.TaskEventTesting
	case 2:
		validOperator = uid == task.Tester
		ev = model.TaskEventFinished
	case 3:
		validOperator = uid == task.Creator
		ev = model.TaskEventArchived
		task.ArchiveTime = time.Now()
	default:
		web.Assert(false, "任务无下一步流程")
	}

	if !validOperator {
		rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=? AND `isadmin`=1", task.PID, uid)
		web.Assert(err == nil, "您不可修改该任务状态")
		defer rows.Close()

		web.Assert(rows.Next(), "您无权修改该任务")

		count := 0
		web.Assert(rows.Scan(&count) == nil && count > 0, "您无权修改该任务")
	}

	task.State++
	err = orm.Update(task)
	web.Assert(err == nil, "修改任务失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), ev, "")
	c.JSON(200, web.Map{})
}

func (t *Task) info(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	c.JSON(200, model.MakeTaskDetail(tid))
}

func (t *Task) delete(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	uid := c.Session.Get("uid").(int64)
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	if task.Creator != uid {
		rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=? AND `isadmin`=1", task.PID, uid)
		web.Assert(err == nil, "只有创建者或管理员可以删除任务")
		defer rows.Close()

		web.Assert(rows.Next(), "您无权删除任务")

		count := 0
		web.Assert(rows.Scan(&count) == nil && count > 0, "只有创建者或管理员可以删除任务")
	}

	orm.Delete("task", tid)
	c.JSON(200, web.Map{})
}

func (t *Task) mine(c *web.Context) {
	uid := c.Session.Get("uid").(int64)

	rows, err := orm.Query(
		"SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringtop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `state`<4 AND (`creator`=? OR `developer`=? OR `tester`=?)",
		uid, uid, uid)
	web.Assert(err == nil, "读取数据库错误")
	defer rows.Close()

	ret := []interface{}{}
	for rows.Next() {
		task := &model.Task{}
		err = orm.Scan(rows, task)
		web.Assert(err == nil, "读取数据库错误")
		ret = append(ret, model.MakeTaskBrief(task))
	}

	c.JSON(200, web.Map{"data": ret})
}

func (t *Task) project(c *web.Context) {
	pid := c.RouteValue("id").MustInt("")
	rows, err := orm.Query(
		"SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringtop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `state`<4 AND `pid`=?",
		pid)
	web.Assert(err == nil, "读取数据库错误")
	defer rows.Close()

	ret := []interface{}{}
	for rows.Next() {
		task := &model.Task{}
		err = orm.Scan(rows, task)
		web.Assert(err == nil, "读取数据库错误")
		ret = append(ret, model.MakeTaskBrief(task))
	}

	c.JSON(200, web.Map{"data": ret})
}

func (t *Task) setName(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	name := c.PostFormValue("name").MustString("任务名不可为空")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	old := task.Name
	task.Name = name
	err = orm.Update(task)
	web.Assert(err == nil, "更新数据库失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventRename, old)
	c.JSON(200, web.Map{})
}

func (t *Task) setCreator(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	uid := c.PostFormValue("member").MustInt("人员ID无效")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	old := task.Creator
	task.Creator = uid
	err = orm.Update(task)
	web.Assert(err == nil, "更新数据库失败")

	oldCreator, _ := model.FindUserInfo(old)
	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModCreator, oldCreator)
	c.JSON(200, web.Map{})
}

func (t *Task) setDeveloper(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	uid := c.PostFormValue("member").MustInt("开发人员ID无效")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	old := task.Developer
	task.Developer = uid
	err = orm.Update(task)
	web.Assert(err == nil, "更新数据库失败")

	oldDeveloper, _ := model.FindUserInfo(old)
	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModDeveloper, oldDeveloper)
	c.JSON(200, web.Map{})
}

func (t *Task) setTester(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	uid := c.PostFormValue("member").MustInt("新人员未指定")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	old := task.Tester
	task.Tester = uid
	err = orm.Update(task)
	web.Assert(err == nil, "更新数据库失败")

	oldTester, _ := model.FindUserInfo(old)
	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModTester, oldTester)
	c.JSON(200, web.Map{})
}

func (t *Task) setWeight(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	weight := c.PostFormValue("weight").MustInt("无效的优先级")
	old := c.PostFormValue("old").MustString("无效的旧优先级")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	task.Weight = int8(weight)
	err = orm.Update(task)
	web.Assert(err == nil, "更新数据库失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModWeight, old)
	c.JSON(200, web.Map{})
}

func (t *Task) setTime(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	uid := c.Session.Get("uid").(int64)
	startTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("startTime").MustString("未指定开始时间"))
	endTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("endTime").MustString("未指定结束时间"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	if !task.StartTime.Equal(startTime) {
		model.AfterTaskOperation(task, uid, model.TaskEventModStartTime, task.StartTime.Format(model.TaskTimeFormat))
		task.StartTime = startTime
	}

	if !task.EndTime.Equal(endTime) {
		model.AfterTaskOperation(task, uid, model.TaskEventModEndTime, task.EndTime.Format(model.TaskTimeFormat))
		task.EndTime = endTime
	}

	err = orm.Update(task)
	web.Assert(err == nil, "更新数据库失败")

	c.JSON(200, web.Map{})
}

func (t *Task) setContent(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	uid := c.Session.Get("uid").(int64)
	content := c.PostFormValue("content").MustString("任务内容不可为空")

	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	task.Content = content
	err = orm.Update(task)
	web.Assert(err == nil, "更新数据库失败")

	model.AfterTaskOperation(task, uid, model.TaskEventModContent, "")
	c.JSON(200, web.Map{})
}

func (t *Task) setStatus(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	status := int8(c.FormValue("moveTo").MustInt("非法的任务状态"))
	uid := c.Session.Get("uid").(int64)
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	if task.State == status {
		c.JSON(200, &web.Map{})
		return
	}

	ev := model.TaskEventSetStatus
	validOperator := false
	switch task.State {
	case 0:
		validOperator = (uid == task.Developer && status == 1)
	case 1:
		validOperator = (uid == task.Developer && status < 3)
	case 2:
		validOperator = (uid == task.Tester && status > 0 && status < 4) || uid == task.Creator
	case 3:
		validOperator = (uid == task.Tester && status > 0 && status < 3) || uid == task.Creator
	case 4:
		validOperator = uid == task.Creator
	}

	if !validOperator {
		rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=? AND `isadmin`=1", task.PID, uid)
		web.Assert(err == nil, "您不可修改该任务状态")
		defer rows.Close()

		web.Assert(rows.Next(), "您无权修改该任务")

		count := 0
		web.Assert(rows.Scan(&count) == nil && count > 0, "您无权修改该任务")
	}

	if task.State == 4 {
		task.ArchiveTime = model.TaskTimeInfinite
	}

	if status == 4 {
		task.ArchiveTime = time.Now()
	}

	task.State = status
	err = orm.Update(task)
	web.Assert(err == nil, "修改任务状态失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), ev, strconv.Itoa(int(status)))
	c.JSON(200, web.Map{})
}

func (t *Task) addComment(c *web.Context) {
	tid := c.RouteValue("id").MustInt("")
	content := c.PostFormValue("content").MustString("任务内容不可为空")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	web.Assert(err == nil, "任务不存在或已被删除")

	_, err = orm.Insert(&model.TaskComment{
		TID:     tid,
		UID:     c.Session.Get("uid").(int64),
		Time:    time.Now(),
		Comment: content,
	})
	web.Assert(err == nil, "发送评论失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventComment, "")
	c.JSON(200, web.Map{})
}
