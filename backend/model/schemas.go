package model

import (
	"team/web/orm"
	"time"
)

// Task events.
const (
	TaskEventCreate       = 0
	TaskEventUnderway     = 1
	TaskEventTesting      = 2
	TaskEventFinished     = 3
	TaskEventArchived     = 4
	TaskEventModStartTime = 5
	TaskEventModEndTime   = 6
	TaskEventModCreator   = 7
	TaskEventModDeveloper = 8
	TaskEventModTester    = 9
	TaskEventModWeight    = 10
	TaskEventModContent   = 11
	TaskEventComment      = 12
	TaskEventMoveBack     = 13
)

var (
	// TaskTimeFormat --
	TaskTimeFormat = "2006-01-02"
	// TaskTimeInfinite is the time never reached.
	TaskTimeInfinite, _ = time.Parse(TaskTimeFormat, "2000-01-01")
)

type (
	// User schema.
	User struct {
		ID              int64  `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Account         string `json:"account" mysql:"VARCHAR(64) UNIQUE NOT NULL"`
		Name            string `json:"name" mysql:"VARCHAR(32) UNIQUE NOT NULL"`
		Avatar          string `json:"avatar" mysql:"VARCHAR(128)"`
		Password        string `json:"-" mysql:"CHAR(32) NOT NULL"`
		IsSu            bool   `json:"isSu" mysql:"TINYINT(1) DEFAULT '0'"`
		IsLocked        bool   `json:"isLocked" mysql:"TINYINT(1) DEFAULT '0'"`
		AutoLoginExpire int64  `json:"-" mysql:"BIGINT UNSIGNED DEFAULT '0'"`
	}
	// Project schema
	Project struct {
		ID       int64    `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Name     string   `json:"name" mysql:"VARCHAR(64) UNIQUE NOT NULL"`
		Branches []string `json:"branches" mysql:"TEXT"`
	}
	// ProjectMember schema
	ProjectMember struct {
		ID      int64 `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		UID     int64 `json:"uid" mysql:"BIGINT NOT NULL"`
		PID     int64 `json:"pid" mysql:"BIGINT NOT NULL"`
		Role    int8  `json:"role" mysql:"INTEGER NOT NULL"`
		IsAdmin bool  `json:"isAdmin" mysql:"TINYINT(1) DEFAULT '0'"`
	}
	// Share schema
	Share struct {
		ID   int64     `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Name string    `json:"name" mysql:"VARCHAR(128) NOT NULL"`
		Path string    `json:"url" mysql:"VARCHAR(128) NOT NULL"`
		UID  int64     `json:"uid" mysql:"BIGINT NOT NULL"`
		Time time.Time `json:"time" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		Size int64     `json:"size" mysql:"BIGINT DEFAULT '0'"`
	}
	// Document schema
	Document struct {
		ID       int64     `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Parent   int64     `json:"parent" mysql:"BIGINT DEFAULT '-1'"`
		Title    string    `json:"title" mysql:"VARCHAR(64) NOT NULL"`
		Author   int64     `json:"author" mysql:"BIGINT NOT NULL"`
		Modifier int64     `json:"modifier" mysql:"BIGINT NOT NULL"`
		Time     time.Time `json:"time" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		Content  string    `json:"content" mysql:"TEXT"`
	}
	// Notice schema
	Notice struct {
		ID      int64     `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Time    time.Time `json:"time" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		UID     int64     `json:"uid" mysql:"BIGINT NOT NULL"`
		Content string    `json:"content" mysql:"TEXT"`
	}
	// Task schema
	Task struct {
		ID          int64     `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		PID         int64     `json:"pid" mysql:"BIGINT NOT NULL"`
		Branch      int8      `json:"branch" mysql:"INTEGER DEFAULT '0'"`
		Creator     int64     `json:"creator" mysql:"BIGINT NOT NULL"`
		Developer   int64     `json:"developer" mysql:"BIGINT NOT NULL"`
		Tester      int64     `json:"tester" mysql:"BIGINT NOT NULL"`
		Name        string    `json:"name" mysql:"VARCHAR(128) NOT NULL"`
		BringTop    bool      `json:"bringTop" mysql:"TINYINT(1) DEFAULT '0'"`
		Weight      int8      `json:"weight" mysql:"INTEGER DEFAULT '0'"`
		State       int8      `json:"state" mysql:"INTEGER DEFAULT '0'"`
		StartTime   time.Time `json:"startTime" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		EndTime     time.Time `json:"endTime" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		ArchiveTime time.Time `json:"archiveTime" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		Tags        []int     `json:"tags" mysql:"TEXT"`
		Content     string    `json:"content" mysql:"TEXT"`
	}
	// TaskAttachment schema
	TaskAttachment struct {
		ID   int64  `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		TID  int64  `json:"tid" mysql:"BIGINT NOT NULL"`
		Name string `json:"name" mysql:"VARCHAR(128) NOT NULL"`
		Path string `json:"url" mysql:"VARCHAR(128) NOT NULL"`
	}
	// TaskComment schema
	TaskComment struct {
		ID      int64     `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		TID     int64     `json:"tid" mysql:"BIGINT NOT NULL"`
		UID     int64     `json:"uid" mysql:"BIGINT NOT NULL"`
		Time    time.Time `json:"time" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		Comment string    `json:"comment" mysql:"TEXT"`
	}
	// TaskEvent schema
	TaskEvent struct {
		ID    int64     `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		TID   int64     `json:"tid" mysql:"BIGINT NOT NULL"`
		UID   int64     `json:"uid" mysql:"BIGINT NOT NULL"`
		Event int       `json:"ev" mysql:"INTEGER DEFAULT '0'"`
		Time  time.Time `json:"time" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		Extra string    `json:"extra" mysql:"TEXT"`
	}
)

// MakeProjectInfo returns detail information for project
func MakeProjectInfo(id int64) map[string]interface{} {
	proj := FindProject(id)
	if proj == nil {
		return map[string]interface{}{}
	}

	type memberInfo struct {
		User    *User `json:"user"`
		Role    int8  `json:"role"`
		IsAdmin bool  `json:"isAdmin"`
	}

	members := []*memberInfo{}
	ret := map[string]interface{}{
		"id":       proj.ID,
		"name":     proj.Name,
		"branches": proj.Branches,
	}

	rows, err := orm.Query("SELECT * FROM `projectmember` WHERE `pid`=?", id)
	if err != nil {
		ret["members"] = members
		return ret
	}

	defer rows.Close()

	for rows.Next() {
		member := &ProjectMember{}
		err = orm.Scan(rows, member)
		if err != nil {
			continue
		}

		user := FindUser(member.UID)
		if user == nil || user.IsLocked {
			continue
		}

		members = append(members, &memberInfo{
			User:    user,
			Role:    member.Role,
			IsAdmin: member.IsAdmin,
		})
	}

	ret["members"] = members
	return ret
}

// MakeTaskBrief return task brief information.
func MakeTaskBrief(task *Task) map[string]interface{} {
	creator := FindUser(task.Creator)
	developer := FindUser(task.Developer)
	tester := FindUser(task.Tester)
	proj := FindProject(task.PID)

	return map[string]interface{}{
		"id":        task.ID,
		"name":      task.Name,
		"proj":      proj,
		"branch":    task.Branch,
		"bringTop":  task.BringTop,
		"weight":    task.Weight,
		"state":     task.State,
		"creator":   creator,
		"developer": developer,
		"tester":    tester,
		"startTime": task.StartTime.Format(TaskTimeFormat),
		"endTime":   task.EndTime.Format(TaskTimeFormat),
	}
}

// MakeTaskDetail returns detail information for task.
func MakeTaskDetail(id int64) map[string]interface{} {
	task := &Task{ID: id}
	err := orm.Read(task)
	if err != nil {
		return map[string]interface{}{"err": "任务不存在"}
	}

	proj := MakeProjectInfo(task.PID)
	if _, ok := proj["id"]; !ok {
		return map[string]interface{}{"err": "任务所属项目已被删除"}
	}

	creator := FindUser(task.Creator)
	developer := FindUser(task.Developer)
	tester := FindUser(task.Tester)

	detail := map[string]interface{}{
		"id":        task.ID,
		"name":      task.Name,
		"proj":      proj,
		"branch":    task.Branch,
		"bringTop":  task.BringTop,
		"weight":    task.Weight,
		"state":     task.State,
		"creator":   creator,
		"developer": developer,
		"tester":    tester,
		"startTime": task.StartTime.Format(TaskTimeFormat),
		"endTime":   task.EndTime.Format(TaskTimeFormat),
		"tags":      task.Tags,
		"content":   task.Content,
	}

	comments := []map[string]string{}
	commRows, err := orm.Query("SELECT * FROM `taskcomment` WHERE `tid`=?", task.ID)
	if err == nil {
		defer commRows.Close()

		for commRows.Next() {
			comment := &TaskComment{}
			err = orm.Scan(commRows, comment)
			if err == nil {
				name, avatar := FindUserInfo(comment.UID)
				comments = append(comments, map[string]string{
					"time":    comment.Time.Format(TaskTimeFormat),
					"user":    name,
					"avatar":  avatar,
					"content": comment.Comment,
				})
			}
		}
	}

	events := []map[string]interface{}{}
	evRows, err := orm.Query("SELECT * FROM `taskevent` WHERE `tid`=? ORDER BY `time` DESC", task.ID)
	if err == nil {
		defer evRows.Close()

		for evRows.Next() {
			ev := &TaskEvent{}
			err = orm.Scan(evRows, ev)
			if err == nil {
				name, _ := FindUserInfo(ev.UID)
				events = append(events, map[string]interface{}{
					"time":     ev.Time.Format(TaskTimeFormat),
					"operator": name,
					"desc":     MakeTaskEventDesc(ev),
				})
			}
		}
	}

	attachments := []*TaskAttachment{}
	attachmentRows, err := orm.Query("SELECT * FROM `taskattachment` WHERE `tid`=?", task.ID)
	if err == nil {
		defer attachmentRows.Close()

		for attachmentRows.Next() {
			one := &TaskAttachment{}
			err = orm.Scan(attachmentRows, one)
			if err == nil {
				attachments = append(attachments, one)
			}
		}
	}

	detail["comments"] = comments
	detail["events"] = events
	detail["attachments"] = attachments

	return map[string]interface{}{"data": detail}
}

// MakeTaskEventDesc return description for task event.
func MakeTaskEventDesc(ev *TaskEvent) string {
	switch ev.Event {
	case TaskEventCreate:
		return "创建了任务"
	case TaskEventUnderway:
		return "接手了任务"
	case TaskEventTesting:
		return "开启了任务测试/验收流程"
	case TaskEventFinished:
		return "将任务设置为完成"
	case TaskEventArchived:
		return "将任务归档"
	case TaskEventModStartTime:
		return "修改了任务计划开始时间，原时间：" + ev.Extra
	case TaskEventModEndTime:
		return "修改了任务计划截止时间，原时间：" + ev.Extra
	case TaskEventModCreator:
		return "移交了任务，原负责人：" + ev.Extra
	case TaskEventModDeveloper:
		return "修改了任务开发者，原开发者：" + ev.Extra
	case TaskEventModTester:
		return "修改了任务验收/测试者，原测试人：" + ev.Extra
	case TaskEventModWeight:
		return "修改了任务优先级，原优先级：" + ev.Extra
	case TaskEventModContent:
		return "修改了任务的具体内容"
	case TaskEventComment:
		return "评论了任务"
	case TaskEventMoveBack:
		return "回退了任务"
	}

	return "未定义操作"
}

// MakeTaskNotice returns notice content
func MakeTaskNotice(task *Task, operator string, event *TaskEvent) string {
	link := "<a href='#'>" + task.Name + "</a>"

	switch event.Event {
	case TaskEventCreate:
		return operator + "创建了任务: " + link
	case TaskEventUnderway:
		return operator + "接手了任务: " + link
	case TaskEventTesting:
		return operator + "开启了任务：" + link + "的测试/验收流程"
	case TaskEventFinished:
		return operator + "将任务：" + link + "设置为完成"
	case TaskEventArchived:
		return operator + "归档了任务：" + link
	case TaskEventModStartTime:
		return operator + "修改了任务：" + link + "计划开始时间，原时间：" + event.Extra
	case TaskEventModEndTime:
		return operator + "修改了任务计划截止时间，原时间：" + event.Extra
	case TaskEventModCreator:
		return operator + "移交了任务：" + link + "，原负责人：" + event.Extra
	case TaskEventModDeveloper:
		return operator + "修改了任务：" + link + "的开发者，原开发者：" + event.Extra
	case TaskEventModTester:
		return operator + "修改了任务：" + link + "的验收/测试者，原测试人：" + event.Extra
	case TaskEventModWeight:
		return operator + "修改了任务：" + link + "的优先级，原优先级：" + event.Extra
	case TaskEventModContent:
		return operator + "修改了任务：" + link + "的具体内容"
	case TaskEventComment:
		return operator + "评论了任务：" + link
	case TaskEventMoveBack:
		return operator + "回退了任务：" + link
	}

	return operator + "修改了任务：" + link
}

// AfterTaskOperation generates event and notify
func AfterTaskOperation(task *Task, operator int64, ev int, extra string) {
	operatorName, _ := FindUserInfo(operator)

	event := &TaskEvent{
		TID:   task.ID,
		UID:   operator,
		Event: ev,
		Time:  time.Now(),
		Extra: extra,
	}

	notice := &Notice{
		Time:    time.Now(),
		Content: MakeTaskNotice(task, "<b>"+operatorName+"</b>", event),
	}

	orm.Insert(event)

	if task.Creator != operator {
		notice.UID = task.Creator
		orm.Insert(notice)
	}

	if task.Developer != operator {
		notice.UID = task.Developer
		orm.Insert(notice)
	}

	if task.Tester != operator {
		notice.UID = task.Tester
		orm.Insert(notice)
	}
}
