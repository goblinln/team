package project

import (
	"errors"
	"sync"
	"time"

	"team/model/user"
	"team/orm"
)

type (
	// Project schema
	Project struct {
		ID   int64  `json:"id"`
		Name string `json:"name" orm:"type=VARCHAR(64),unique,notnull"`

		// Runtime data.
		Milestones []*Milestone `json:"-" orm:"-"`
		Members    []*Member    `json:"-" orm:"-"`
	}

	// Milestone schema
	Milestone struct {
		ID        int64     `json:"id"`
		PID       int64     `json:"pid"`
		Name      string    `json:"name" orm:"type=VARCHAR(64),unique,notnull"`
		StartTime time.Time `json:"startTime" orm:"notnull,default='2000-01-01'"`
		EndTime   time.Time `json:"endTime" orm:"notnull,default='2000-01-01'"`
	}

	// Member schema
	Member struct {
		ID      int64 `json:"id"`
		PID     int64 `json:"pid"`
		UID     int64 `json:"uid"`
		Role    int8  `json:"role"`
		IsAdmin bool  `json:"isAdmin"`
	}
)

var (
	// Global project cache.
	projectCache = sync.Map{}
)

// GetAll returns all projects in db
func GetAll() ([]*Project, error) {
	rows, err := orm.Query("SELECT * FROM `project`")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projs := []*Project{}
	for rows.Next() {
		one := &Project{}
		if err = orm.Scan(rows, one); err != nil {
			return nil, err
		}

		projs = append(projs, one)
	}

	return projs, nil
}

// GetAllByUser returns all project that the given user has joined.
func GetAllByUser(uid int64) ([]*Project, error) {
	rows, err := orm.Query("SELECT * FROM `member` WHERE `uid`=?", uid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projs := []*Project{}
	for rows.Next() {
		one := &Member{}
		if err = orm.Scan(rows, one); err != nil {
			return nil, err
		}

		proj := Find(one.PID)
		if proj != nil {
			projs = append(projs, proj)
		}
	}

	return projs, nil
}

// Find project
func Find(ID int64) *Project {
	exists, ok := projectCache.Load(ID)
	if ok {
		return exists.(*Project)
	}

	proj := &Project{ID: ID}
	if err := orm.Read(proj); err != nil {
		return nil
	}

	proj.FetchMilestones()
	proj.FetchMembers()

	projectCache.Store(proj.ID, proj)
	return proj
}

// Add a new project.
func Add(name string, uid int64, role int8) error {
	rows, err := orm.Query("SELECT COUNT(*) FROM `project` WHERE `name`=?", name)
	if err != nil {
		return err
	}

	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	if count != 0 {
		return errors.New("同名项目已存在")
	}

	admin := user.Find(uid)
	if admin == nil {
		return errors.New("默认管理员不存在或已被删除")
	}

	if admin.IsLocked {
		return errors.New("默认管理员当前被禁止登录")
	}

	proj := &Project{Name: name}
	rs, err := orm.Insert(proj)
	if err != nil {
		return errors.New("写入新项目失败")
	}

	proj.ID, _ = rs.LastInsertId()
	orm.Insert(&Member{PID: proj.ID, UID: uid, Role: role, IsAdmin: true})
	projectCache.Store(proj.ID, proj)
	return nil
}

// Delete an existed project by ID
func Delete(ID int64) {
	orm.Delete("project", ID)

	orm.Exec("DELETE FROM `milestone` WHERE `pid`=?", ID)
	orm.Exec("DELETE FROM `member` WHERE `pid`=?", ID)
	orm.Exec("DELETE FROM `task` WHERE `pid`=?", ID)

	projectCache.Delete(ID)
}

// Info returns detail information
func (p *Project) Info() map[string]interface{} {
	if p.Members == nil {
		p.FetchMembers()
	}

	if p.Milestones == nil {
		p.FetchMilestones()
	}

	members := []map[string]interface{}{}
	for _, one := range p.Members {
		u := user.Find(one.UID)
		if u == nil || u.IsLocked {
			continue
		}

		members = append(members, map[string]interface{}{
			"user":    u,
			"role":    one.Role,
			"isAdmin": one.IsAdmin,
		})
	}

	milestones := []*Milestone{}
	now := time.Now()
	for _, one := range p.Milestones {
		if one.EndTime.After(now) {
			milestones = append(milestones, one)
		}
	}

	return map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"milestones": milestones,
		"members":    members,
	}
}

// FetchMilestones preloads all valid(NOT ended) milestones for this project.
func (p *Project) FetchMilestones() {
	rows, err := orm.Query("SELECT * FROM `milestone` WHERE `pid`=?", p.ID)
	if err != nil {
		return
	}

	defer rows.Close()

	p.Milestones = []*Milestone{}
	for rows.Next() {
		one := &Milestone{}
		if err = orm.Scan(rows, one); err != nil {
			continue
		}

		p.Milestones = append(p.Milestones, one)
	}
}

// AddMilestone adds new milestone to this project.
func (p *Project) AddMilestone(name string, startTime, endTime time.Time) error {
	add := &Milestone{
		PID:       p.ID,
		Name:      name,
		StartTime: startTime,
		EndTime:   endTime,
	}

	rs, err := orm.Insert(add)
	if err == nil {
		add.ID, _ = rs.LastInsertId()
		p.Milestones = append(p.Milestones, add)
	}

	return err
}

// FindMilestone returns milestone by ID
func (p *Project) FindMilestone(mid int64) *Milestone {
	for _, one := range p.Milestones {
		if one.ID == mid {
			return one
		}
	}

	return nil
}

// FetchMembers preloads all members for this project.
func (p *Project) FetchMembers() {
	rows, err := orm.Query("SELECT * FROM `member` WHERE `pid`=?", p.ID)
	if err != nil {
		return
	}

	defer rows.Close()

	p.Members = []*Member{}
	for rows.Next() {
		one := &Member{}
		if err = orm.Scan(rows, one); err != nil {
			continue
		}

		p.Members = append(p.Members, one)
	}
}

// AddMember adds member to this project.
func (p *Project) AddMember(uid int64, role int8, isAdmin bool) error {
	add := &Member{
		PID:     p.ID,
		UID:     uid,
		Role:    role,
		IsAdmin: isAdmin,
	}

	rs, err := orm.Insert(add)
	if err == nil {
		add.ID, _ = rs.LastInsertId()
		p.Members = append(p.Members, add)
	}

	return err
}

// IsAdmin returns true if given user is administrator of this project.
func (p *Project) IsAdmin(uid int64) bool {
	for _, one := range p.Members {
		if one.UID == uid {
			return one.IsAdmin
		}
	}

	return false
}

// DelMember remove given user from this project.
func (p *Project) DelMember(uid int64) {
	orm.Exec("DELETE FROM `member` WHERE `pid`=? AND `uid`=?", p.ID, uid)
}

// Save project.
func (p *Project) Save() error {
	return orm.Update(p)
}

// Save modified member.
func (m *Member) Save() error {
	return orm.Update(m)
}
