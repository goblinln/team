package controller

import (
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

	group.PATCH("/:id/name", t.setName)
	group.PATCH("/:id/creator", t.setCreator)
	group.PATCH("/:id/developer", t.setDeveloper)
	group.PATCH("/:id/tester", t.setTester)
	group.PATCH("/:id/weight", t.setWeight)
	group.PATCH("/:id/time", t.setTime)
	group.PATCH("/:id/content", t.setContent)
	group.POST("/:id/comment", t.addComment)
}

func (t *Task) create(c *web.Context) {
	name := c.PostFormValue("name")
	proj := atoi(c.PostFormValue("proj"))
	branch := atoi(c.PostFormValue("branch"))
	weight := atoi(c.PostFormValue("weight"))
	creator := atoi(c.PostFormValue("creator"))
	developer := atoi(c.PostFormValue("developer"))
	tester := atoi(c.PostFormValue("tester"))
	startTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("startTime"))
	endTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("endTime"))
	tags := intArray(c.PostFormArray("tags[]"))
	content := c.PostFormValue("content")

	assert(len(name) > 0, "任务名不可为空")
	assert(len(content) > 0, "任务详情不可为空")
	assert(developer > 0, "开发人员未指定")
	assert(tester > 0, "测试/验收人员未指定")

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
		Tags:        tags,
		Content:     content,
	}

	rs, err := orm.Insert(task)
	assert(err == nil, "写入任务信息失败")

	tid, err := rs.LastInsertId()
	assert(err == nil, "读取新任务ID失败")

	task.ID = tid

	for i := 0; i < len(uploaded); i++ {
		uploaded[i].TID = tid
		orm.Insert(uploaded[i])
	}

	model.AfterTaskOperation(task, uid, model.TaskEventCreate, "")
	c.JSON(200, web.Map{})
}

func (t *Task) moveBack(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := c.Session.Get("uid").(int64)
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

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
		assert(false, "任务不可回退了")
	}

	if !validOperator {
		rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=? AND `isadmin`=1", task.PID, uid)
		assert(err == nil, "您不可回退当前状态")
		defer rows.Close()

		assert(rows.Next(), "你无权回退该任务")

		count := 0
		assert(rows.Scan(&count) == nil && count > 0, "你无权回退该任务")
	}

	task.State--
	err = orm.Update(task)
	assert(err == nil, "修改任务失败")

	model.AfterTaskOperation(task, uid, model.TaskEventMoveBack, "")
	c.JSON(200, web.Map{})
}

func (t *Task) moveNext(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := c.Session.Get("uid").(int64)
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

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
		assert(false, "任务无下一步流程")
	}

	if !validOperator {
		rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=? AND `isadmin`=1", task.PID, uid)
		assert(err == nil, "您不可修改该任务状态")
		defer rows.Close()

		assert(rows.Next(), "您无权修改该任务")

		count := 0
		assert(rows.Scan(&count) == nil && count > 0, "您无权修改该任务")
	}

	task.State++
	err = orm.Update(task)
	assert(err == nil, "修改任务失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), ev, "")
	c.JSON(200, web.Map{})
}

func (t *Task) info(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	c.JSON(200, model.MakeTaskDetail(tid))
}

func (t *Task) delete(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := c.Session.Get("uid").(int64)
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	if task.Creator != uid {
		rows, err := orm.Query("SELECT COUNT(*) FROM `projectmember` WHERE `pid`=? AND `uid`=? AND `isadmin`=1", task.PID, uid)
		assert(err == nil, "只有创建者或管理员可以删除任务")
		defer rows.Close()

		assert(rows.Next(), "您无权删除任务")

		count := 0
		assert(rows.Scan(&count) == nil && count > 0, "只有创建者或管理员可以删除任务")
	}

	orm.Delete("task", tid)
	c.JSON(200, web.Map{})
}

func (t *Task) mine(c *web.Context) {
	uid := c.Session.Get("uid").(int64)

	rows, err := orm.Query(
		"SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringtop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `state`<4 AND (`creator`=? OR `developer`=? OR `tester`=?)",
		uid, uid, uid)
	assert(err == nil, "读取数据库错误")
	defer rows.Close()

	ret := []interface{}{}
	for rows.Next() {
		task := &model.Task{}
		err = orm.Scan(rows, task)
		assert(err == nil, "读取数据库错误")
		ret = append(ret, model.MakeTaskBrief(task))
	}

	c.JSON(200, web.Map{"data": ret})
}

func (t *Task) project(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	rows, err := orm.Query(
		"SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringtop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `state`<4 AND `pid`=?",
		pid)
	assert(err == nil, "读取数据库错误")
	defer rows.Close()

	ret := []interface{}{}
	for rows.Next() {
		task := &model.Task{}
		err = orm.Scan(rows, task)
		assert(err == nil, "读取数据库错误")
		ret = append(ret, model.MakeTaskBrief(task))
	}

	c.JSON(200, web.Map{"data": ret})
}

func (t *Task) setName(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	name := c.PostFormValue("name")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	old := task.Name
	task.Name = name
	err = orm.Update(task)
	assert(err == nil, "更新数据库失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventRename, old)
	c.JSON(200, web.Map{})
}

func (t *Task) setCreator(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := atoi(c.PostFormValue("member"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	old := task.Creator
	task.Creator = uid
	err = orm.Update(task)
	assert(err == nil, "更新数据库失败")

	oldCreator, _ := model.FindUserInfo(old)
	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModCreator, oldCreator)
	c.JSON(200, web.Map{})
}

func (t *Task) setDeveloper(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := atoi(c.PostFormValue("member"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	old := task.Developer
	task.Developer = uid
	err = orm.Update(task)
	assert(err == nil, "更新数据库失败")

	oldDeveloper, _ := model.FindUserInfo(old)
	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModDeveloper, oldDeveloper)
	c.JSON(200, web.Map{})
}

func (t *Task) setTester(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := atoi(c.PostFormValue("member"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	old := task.Tester
	task.Tester = uid
	err = orm.Update(task)
	assert(err == nil, "更新数据库失败")

	oldTester, _ := model.FindUserInfo(old)
	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModTester, oldTester)
	c.JSON(200, web.Map{})
}

func (t *Task) setWeight(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	weight := atoi(c.PostFormValue("weight"))
	old := c.PostFormValue("old")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	task.Weight = int8(weight)
	err = orm.Update(task)
	assert(err == nil, "更新数据库失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventModWeight, old)
	c.JSON(200, web.Map{})
}

func (t *Task) setTime(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := c.Session.Get("uid").(int64)
	startTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("startTime"))
	endTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("endTime"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	if !task.StartTime.Equal(startTime) {
		model.AfterTaskOperation(task, uid, model.TaskEventModStartTime, task.StartTime.Format(model.TaskTimeFormat))
		task.StartTime = startTime
	}

	if !task.EndTime.Equal(endTime) {
		model.AfterTaskOperation(task, uid, model.TaskEventModEndTime, task.EndTime.Format(model.TaskTimeFormat))
		task.EndTime = endTime
	}

	err = orm.Update(task)
	assert(err == nil, "更新数据库失败")

	c.JSON(200, web.Map{})
}

func (t *Task) setContent(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := c.Session.Get("uid").(int64)
	content := c.PostFormValue("content")

	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	task.Content = content
	err = orm.Update(task)
	assert(err == nil, "更新数据库失败")

	model.AfterTaskOperation(task, uid, model.TaskEventModContent, "")
	c.JSON(200, web.Map{})
}

func (t *Task) addComment(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	assert(err == nil, "任务不存在或已被删除")

	_, err = orm.Insert(&model.TaskComment{
		TID:     tid,
		UID:     c.Session.Get("uid").(int64),
		Time:    time.Now(),
		Comment: c.PostFormValue("content"),
	})
	assert(err == nil, "发送评论失败")

	model.AfterTaskOperation(task, c.Session.Get("uid").(int64), model.TaskEventComment, "")
	c.JSON(200, web.Map{})
}
