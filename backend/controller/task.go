package controller

import (
	"time"

	"team/model"
	"team/web"
	"team/web/orm"
)

// Task controller
type Task int

// Register implements web.Controller interface.
func (t *Task) Register(group *web.Router) {
	group.POST("", t.create)
	group.POST("/:id/next", t.moveNext)
	group.GET("/:id", t.info)
	group.DELETE("/:id", t.delete)

	group.GET("/mine", t.mine)
	group.GET("/project/:id", t.project)

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

	if len(name) == 0 || len(content) == 0 {
		c.JSON(200, &web.JObject{"err": "任务名及详情不可为空"})
		return
	}

	uid := c.Session.Get("uid").(int64)
	me := model.FindUser(uid)

	if creator == 0 {
		creator = uid
	}

	if developer == 0 || tester == 0 {
		c.JSON(200, &web.JObject{"err": "人员配置出错"})
		return
	}

	uploaded := []*model.TaskAttachment{}
	fhs, ok := c.MultipartForm().File["files[]"]
	if ok {
		helper := new(File)
		for _, fh := range fhs {
			url, _, err := helper.save(fh, uid)
			if err != nil {
				c.JSON(200, &web.JObject{"err": err.Error()})
				return
			}

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
	if err != nil {
		c.JSON(200, &web.JObject{"err": "写入任务信息失败"})
		return
	}

	tid, err := rs.LastInsertId()
	if err != nil {
		c.JSON(200, &web.JObject{"err": "读取新任务ID失败"})
		return
	}

	for i := 0; i < len(uploaded); i++ {
		uploaded[i].TID = tid
		orm.Insert(uploaded[i])
	}

	orm.Insert(&model.TaskEvent{
		TID:   tid,
		UID:   uid,
		Event: model.TaskEventCreate,
		Time:  time.Now(),
		Extra: "",
	})

	c.JSON(200, &web.JObject{})
}

func (t *Task) moveNext(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "任务不存在或已被删除"})
		return
	}

	ev := model.TaskEventUnderway
	switch task.State {
	case 0:
		ev = model.TaskEventUnderway
	case 1:
		ev = model.TaskEventTesting
	case 2:
		ev = model.TaskEventFinished
	case 3:
		ev = model.TaskEventArchived
		task.ArchiveTime = time.Now()
	}

	task.State++
	err = orm.Update(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "修改任务失败"})
		return
	}

	orm.Insert(&model.TaskEvent{
		TID:   tid,
		UID:   c.Session.Get("uid").(int64),
		Event: ev,
		Time:  time.Now(),
		Extra: "",
	})

	c.JSON(200, &web.JObject{})
}

func (t *Task) info(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	c.JSON(200, model.MakeTaskDetail(tid))
}

func (t *Task) delete(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	orm.Delete("task", tid)
}

func (t *Task) mine(c *web.Context) {
	uid := c.Session.Get("uid").(int64)

	rows, err := orm.Query(
		"SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringTop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `state`<4 AND (`creator`=? OR `developer`=? OR `tester`=?)",
		uid, uid, uid)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "读取数据库错误"})
		return
	}

	defer rows.Close()

	ret := []interface{}{}
	for rows.Next() {
		task := &model.Task{}
		err = orm.Scan(rows, task)
		if err != nil {
			c.JSON(200, &web.JObject{"err": "读取数据库错误"})
			return
		}

		ret = append(ret, model.MakeTaskBrief(task))
	}

	c.JSON(200, &web.JObject{"data": ret})
}

func (t *Task) project(c *web.Context) {
	pid := atoi(c.RouteValue("id"))
	rows, err := orm.Query(
		"SELECT `id`,`pid`,`branch`,`creator`,`developer`,`tester`,`name`,`bringTop`,`weight`,`state`,`starttime`,`endtime` FROM `task` WHERE `state`<4 AND `pid`=?",
		pid)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "读取数据库错误"})
		return
	}

	defer rows.Close()

	ret := []interface{}{}
	for rows.Next() {
		task := &model.Task{}
		err = orm.Scan(rows, task)
		if err != nil {
			c.JSON(200, &web.JObject{"err": "读取数据库错误"})
			return
		}

		ret = append(ret, model.MakeTaskBrief(task))
	}

	c.JSON(200, &web.JObject{"data": ret})
}

func (t *Task) setCreator(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := atoi(c.PostFormValue("member"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "任务不存在或已被删除"})
		return
	}

	old := task.Creator
	task.Creator = uid
	err = orm.Update(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "更新数据库失败"})
		return
	}

	oldCreator, _ := model.FindUserInfo(old)
	orm.Insert(&model.TaskEvent{
		TID:   tid,
		UID:   c.Session.Get("uid").(int64),
		Event: model.TaskEventModCreator,
		Time:  time.Now(),
		Extra: oldCreator,
	})

	c.JSON(200, &web.JObject{})
}

func (t *Task) setDeveloper(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := atoi(c.PostFormValue("member"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "任务不存在或已被删除"})
		return
	}

	old := task.Developer
	task.Developer = uid
	err = orm.Update(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "更新数据库失败"})
		return
	}

	oldDeveloper, _ := model.FindUserInfo(old)
	orm.Insert(&model.TaskEvent{
		TID:   tid,
		UID:   c.Session.Get("uid").(int64),
		Event: model.TaskEventModDeveloper,
		Time:  time.Now(),
		Extra: oldDeveloper,
	})

	c.JSON(200, &web.JObject{})
}

func (t *Task) setTester(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := atoi(c.PostFormValue("member"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "任务不存在或已被删除"})
		return
	}

	old := task.Tester
	task.Tester = uid
	err = orm.Update(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "更新数据库失败"})
		return
	}

	oldTester, _ := model.FindUserInfo(old)
	orm.Insert(&model.TaskEvent{
		TID:   tid,
		UID:   c.Session.Get("uid").(int64),
		Event: model.TaskEventModTester,
		Time:  time.Now(),
		Extra: oldTester,
	})

	c.JSON(200, &web.JObject{})
}

func (t *Task) setWeight(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	weight := atoi(c.PostFormValue("weight"))
	old := c.PostFormValue("old")
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "任务不存在或已被删除"})
		return
	}

	task.Weight = int8(weight)
	err = orm.Update(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "更新数据库失败"})
		return
	}

	orm.Insert(&model.TaskEvent{
		TID:   tid,
		UID:   c.Session.Get("uid").(int64),
		Event: model.TaskEventModWeight,
		Time:  time.Now(),
		Extra: old,
	})

	c.JSON(200, &web.JObject{})
}

func (t *Task) setTime(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := c.Session.Get("uid").(int64)
	startTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("startTime"))
	endTime, _ := time.Parse(model.TaskTimeFormat, c.PostFormValue("endTime"))
	task := &model.Task{ID: tid}
	err := orm.Read(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "任务不存在或已被删除"})
		return
	}

	if !task.StartTime.Equal(startTime) {
		orm.Insert(&model.TaskEvent{
			TID:   tid,
			UID:   uid,
			Event: model.TaskEventModStartTime,
			Time:  time.Now(),
			Extra: task.StartTime.Format(model.TaskTimeFormat),
		})

		task.StartTime = startTime
	}

	if !task.EndTime.Equal(endTime) {
		orm.Insert(&model.TaskEvent{
			TID:   tid,
			UID:   uid,
			Event: model.TaskEventModEndTime,
			Time:  time.Now(),
			Extra: task.EndTime.Format(model.TaskTimeFormat),
		})

		task.EndTime = endTime
	}

	err = orm.Update(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "修改任务时间失败"})
		return
	}

	c.JSON(200, &web.JObject{})
}

func (t *Task) setContent(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	uid := c.Session.Get("uid").(int64)
	content := c.PostFormValue("content")

	task := &model.Task{ID: tid}
	err := orm.Read(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "任务不存在或已被删除"})
		return
	}

	task.Content = content
	err = orm.Update(task)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "修改任务时间失败"})
		return
	}

	orm.Insert(&model.TaskEvent{
		TID:   tid,
		UID:   uid,
		Event: model.TaskEventModContent,
		Time:  time.Now(),
		Extra: "",
	})

	c.JSON(200, &web.JObject{})
}

func (t *Task) addComment(c *web.Context) {
	tid := atoi(c.RouteValue("id"))
	_, err := orm.Insert(&model.TaskComment{
		TID:     tid,
		UID:     c.Session.Get("uid").(int64),
		Time:    time.Now(),
		Comment: c.PostFormValue("content"),
	})

	if err != nil {
		c.JSON(200, &web.JObject{"err": "发送评论失败"})
		return
	}

	c.JSON(200, &web.JObject{})
}
