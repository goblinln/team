package model

import (
	"sync"
	"team/orm"
)

// CacheManager --
type CacheManager struct {
	sync.Mutex

	Users    map[int64]*User
	Projects map[int64]*Project
}

// Cache is runtime singleton instance of CacheManager
var Cache = &CacheManager{
	Mutex:    sync.Mutex{},
	Users:    make(map[int64]*User),
	Projects: make(map[int64]*Project),
}

// GetUser return user data only from cached list.
func (c *CacheManager) GetUser(ID int64) *User {
	c.Lock()
	defer c.Unlock()
	return c.Users[ID]
}

// SetUser pushes user data into cache list
func (c *CacheManager) SetUser(user *User) {
	c.Lock()
	defer c.Unlock()
	c.Users[user.ID] = user
}

// DeleteUser removes user from this cache
func (c *CacheManager) DeleteUser(ID int64) {
	c.Lock()
	defer c.Unlock()
	delete(c.Users, ID)
}

// GetProject return user data only from cached list.
func (c *CacheManager) GetProject(ID int64) *Project {
	c.Lock()
	defer c.Unlock()
	return c.Projects[ID]
}

// SetProject pushes user data into cache list
func (c *CacheManager) SetProject(proj *Project) {
	c.Lock()
	defer c.Unlock()
	c.Projects[proj.ID] = proj
}

// DeleteProject removes project from this cache
func (c *CacheManager) DeleteProject(ID int64) {
	c.Lock()
	defer c.Unlock()
	delete(c.Projects, ID)
}

// FindUser returns user by ID
func FindUser(ID int64) *User {
	user := Cache.GetUser(ID)
	if user == nil {
		user = &User{ID: ID}
		if err := orm.Read(user); err != nil {
			return nil
		}
		Cache.SetUser(user)
	}

	return user
}

// FindUserInfo returns user's name and avatar.
func FindUserInfo(ID int64) (string, string) {
	user := FindUser(ID)
	if user == nil {
		return "未知者", ""
	}

	return user.Name, user.Avatar
}

// FindProject returns project by ID
func FindProject(ID int64) *Project {
	proj := Cache.GetProject(ID)
	if proj == nil {
		proj = &Project{ID: ID}
		if err := orm.Read(proj); err != nil {
			return nil
		}
		Cache.SetProject(proj)
	}

	return proj
}
