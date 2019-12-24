package user

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"team/orm"
)

const (
	// AutoLoginCookieKey key to store auto login data in cookie
	AutoLoginCookieKey = "login_token"
	// AutoLoginSecret to sign login token
	AutoLoginSecret = "@team.auto_login_secret_01"
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

	// AutoLoginData holds auto login data.
	AutoLoginData struct {
		ID   int64
		IP   string
		Sign string
	}

	// AutoLoginCookie holds cookie data needs to send back to client
	AutoLoginCookie struct {
		Name    string
		Value   string
		Expires time.Time
	}
)

// Global user cache.
var userCache = sync.Map{}

// GetAll returns all users from db.
func GetAll() ([]*User, error) {
	users := []*User{}

	rows, err := orm.Query("SELETE * FROM `user`")
	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		one := &User{}
		if err = orm.Scan(rows, one); err != nil {
			return users, err
		}

		userCache.Store(one.ID, one)
		users = append(users, one)
	}

	return users, nil
}

// Find user by ID.
func Find(ID int64) *User {
	exists, ok := userCache.Load(ID)
	if ok {
		return exists.(*User)
	}

	user := &User{ID: ID}
	if err := orm.Read(user); err != nil {
		return nil
	}

	userCache.Store(ID, user)
	return user
}

// FindInfo returns name and avatar of user by ID.
func FindInfo(ID int64) (string, string) {
	exists := Find(ID)
	if exists == nil {
		return "Unknown", ""
	}

	return exists.Name, exists.Avatar
}

// Add a new user.
func Add(account, name, pswd string, isSu bool) error {
	rows, err := orm.Query("SELECT COUNT(*) FROM `user` WHERE `account`=? OR `name`=?", account, name)
	if err != nil {
		return err
	}

	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	if count != 0 {
		return errors.New("帐号或昵称已存在")
	}

	one := &User{
		Account:  account,
		Name:     name,
		Avatar:   "",
		Password: pswd,
		IsSu:     isSu,
		IsLocked: false,
	}

	result, err := orm.Insert(one)
	if err != nil {
		return err
	}

	one.ID, _ = result.LastInsertId()
	userCache.Store(one.ID, one)
	return nil
}

// Delete an existed user.
func Delete(uid int64) {
	orm.Delete("user", uid)
	userCache.Delete(uid)
}

// AutoLogin by cookie data. <0 means failed.
func AutoLogin(data, ip string) int64 {
	j, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return -1
	}

	param := &AutoLoginData{}
	err = json.Unmarshal(j, param)
	if err != nil || ip != param.IP {
		return -1
	}

	code := fmt.Sprintf("%d|%s|%s", param.ID, param.IP, AutoLoginSecret)
	hash := md5.New()

	hash.Write([]byte(code))
	sign := fmt.Sprintf("%X", hash.Sum(nil))
	if sign != param.Sign {
		return -1
	}

	user := Find(param.ID)
	if user == nil {
		return -1
	}

	now := time.Now()
	if user.IsLocked || user.AutoLoginExpire <= now.Unix() {
		return -1
	}

	return user.ID
}

// Login process.
func Login(account, pswd, ip string, remember bool) (*User, *AutoLoginCookie) {
	hash := md5.New()
	hash.Write([]byte(pswd))

	try := &User{
		Account:  account,
		Password: fmt.Sprintf("%X", hash.Sum(nil)),
	}

	if err := orm.Read(try, "account", "password"); err != nil {
		return nil, nil
	}

	if remember {
		expire := time.Now().Add(30 * 24 * time.Hour)
		code := fmt.Sprintf("%d|%s|%s", try.ID, ip, AutoLoginSecret)
		sign := md5.New()
		sign.Write([]byte(code))

		token := &AutoLoginData{
			ID:   try.ID,
			IP:   ip,
			Sign: fmt.Sprintf("%X", sign.Sum(nil)),
		}

		j, err := json.Marshal(token)
		if err == nil {
			try.AutoLoginExpire = expire.Unix()
			orm.Update(try)
			userCache.Store(try.ID, try)

			return try, &AutoLoginCookie{
				Name:    AutoLoginCookieKey,
				Value:   base64.StdEncoding.EncodeToString(j),
				Expires: expire,
			}
		}
	}

	userCache.Store(try.ID, try)
	return try, nil
}

// Save user data to database.
func (u *User) Save() error {
	return orm.Update(u)
}
