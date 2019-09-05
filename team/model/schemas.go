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
	TaskEventRename       = 14
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
		ID              int64  `json:"id"`
		Account         string `json:"account" orm:"type=VARCHAR(64),unique,notnull"`
		Name            string `json:"name" orm:"type=VARCHAR(32),unique,notnull"`
		Avatar          string `json:"avatar" orm:"type=VARCHAR(128)"`
		Password        string `json:"-" orm:"type=CHAR(32),notnull"`
		IsSu            bool   `json:"isSu"`
		IsLocked        bool   `json:"isLocked"`
		AutoLoginExpire int64  `json:"-"`
	}
	// Project schema
	Project struct {
		ID       int64    `json:"id"`
		Name     string   `json:"name" orm:"type=VARCHAR(64),unique,notnull"`
		Branches []string `json:"branches"`
	}
	// ProjectMember schema
	ProjectMember struct {
		ID      int64 `json:"id"`
		UID     int64 `json:"uid"`
		PID     int64 `json:"pid"`
		Role    int8  `json:"role"`
		IsAdmin bool  `json:"isAdmin"`
	}
	// Share schema
	Share struct {
		ID   int64     `json:"id"`
		Name string    `json:"name" orm:"type=VARCHAR(128),notnull"`
		Path string    `json:"url" orm:"type=VARCHAR(128),notnull"`
		UID  int64     `json:"uid"`
		Time time.Time `json:"time" orm:"default=CURRENT_TIMESTAMP"`
		Size int64     `json:"size"`
	}
	// Document schema
	Document struct {
		ID       int64     `json:"id"`
		Parent   int64     `json:"parent" orm:"default='-1'"`
		Title    string    `json:"title" orm:"type=VARCHAR(64),notnull"`
		Author   int64     `json:"author"`
		Modifier int64     `json:"modifier"`
		Time     time.Time `json:"time" orm:"default=CURRENT_TIMESTAMP"`
		Content  string    `json:"content"`
	}
	// Notice schema
	Notice struct {
		ID       int64     `json:"id"`
		Time     time.Time `json:"time" orm:"default=CURRENT_TIMESTAMP"`
		UID      int64     `json:"uid"`
		TID      int64     `json:"tid"`
		Operator int64     `json:"operator"`
		Event    int16     `json:"event"`
	}
	// Task schema
	Task struct {
		ID          int64     `json:"id"`
		PID         int64     `json:"pid"`
		Branch      int8      `json:"branch"`
		Creator     int64     `json:"creator"`
		Developer   int64     `json:"developer"`
		Tester      int64     `json:"tester"`
		Name        string    `json:"name" orm:"type=VARCHAR(128),notnull"`
		BringTop    bool      `json:"bringTop"`
		Weight      int8      `json:"weight"`
		State       int8      `json:"state"`
		StartTime   time.Time `json:"startTime" orm:"default=CURRENT_TIMESTAMP"`
		EndTime     time.Time `json:"endTime" orm:"notnull"`
		ArchiveTime time.Time `json:"archiveTime" orm:"notnull"`
		Tags        []int     `json:"tags"`
		Content     string    `json:"content"`
	}
	// TaskAttachment schema
	TaskAttachment struct {
		ID   int64  `json:"id"`
		TID  int64  `json:"tid"`
		Name string `json:"name" orm:"type=VARCHAR(128),notnull"`
		Path string `json:"url" orm:"type=VARCHAR(128),notnull"`
	}
	// TaskComment schema
	TaskComment struct {
		ID      int64     `json:"id"`
		TID     int64     `json:"tid"`
		UID     int64     `json:"uid"`
		Time    time.Time `json:"time" orm:"default=CURRENT_TIMESTAMP"`
		Comment string    `json:"comment"`
	}
	// TaskEvent schema
	TaskEvent struct {
		ID    int64     `json:"id"`
		TID   int64     `json:"tid"`
		UID   int64     `json:"uid"`
		Event int       `json:"ev"`
		Time  time.Time `json:"time" orm:"default=CURRENT_TIMESTAMP"`
		Extra string    `json:"extra"`
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
		Operator: operator,
		Event:    int16(ev),
	}

	orm.Insert(event)

	noticed := map[int64]bool{}
	noticed[task.Developer] = false
	noticed[task.Tester] = false

	if task.Creator != operator {
		notice.UID = task.Creator
		noticed[task.Creator] = true
		orm.Insert(notice)
	}

	if task.Developer != operator && !noticed[task.Developer] {
		notice.UID = task.Developer
		noticed[task.Developer] = true
		orm.Insert(notice)
	}

	if task.Tester != operator && !noticed[task.Tester] {
		notice.UID = task.Tester
		orm.Insert(notice)
	}
}
