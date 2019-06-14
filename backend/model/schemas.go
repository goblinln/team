package model

import (
	"team/orm"
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
		ID       int64     `json:"id" mysql:"BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Time     time.Time `json:"time" mysql:"TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
		UID      int64     `json:"uid" mysql:"BIGINT NOT NULL"`
		TID      int64     `json:"tid" mysql:"BIGINT NOT NULL"`
		TName    string    `json:"tname" mysql:"VARCHAR(128) NOT NULL"`
		Operator int64     `json:"operator" mysql:"BIGINT NOT NULL"`
		Event    int16     `json:"event" mysql:"INTEGER NOT NULL"`
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
					"event":    ev.Event,
					"extra":    ev.Extra,
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

// AfterTaskOperation generates event and notify
func AfterTaskOperation(task *Task, operator int64, ev int, extra string) {
	event := &TaskEvent{
		TID:   task.ID,
		UID:   operator,
		Event: ev,
		Time:  time.Now(),
		Extra: extra,
	}

	notice := &Notice{
		Time:     time.Now(),
		TID:      task.ID,
		TName:    task.Name,
		Operator: operator,
		Event:    int16(ev),
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
