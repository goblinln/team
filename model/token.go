package model

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

const (
	// AutoLoginCookieKey key to store auto login data in cookie
	AutoLoginCookieKey = "login_token"
	// AutoLoginSecret to sign login token
	AutoLoginSecret = "@team.auto_login_secret_01"
)

// AutoLoginToken holds auto login data.
type AutoLoginToken struct {
	ID   int64
	IP   string
	Sign string
}

// TryAutoLogin returns logined user id by token. <0 means failed.
func TryAutoLogin(tokenStr, ip string) int64 {
	j, err := base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		return -1
	}

	token := &AutoLoginToken{}
	err = json.Unmarshal(j, token)
	if err != nil || ip != token.IP {
		return -1
	}

	code := fmt.Sprintf("%d|%s|%s", token.ID, token.IP, AutoLoginSecret)
	hash := md5.New()

	hash.Write([]byte(code))
	sign := fmt.Sprintf("%X", hash.Sum(nil))
	if sign != token.Sign {
		return -1
	}

	user := FindUser(token.ID)
	if user == nil {
		return -1
	}

	now := time.Now()
	if user.IsLocked || user.AutoLoginExpire <= now.Unix() {
		return -1
	}

	return user.ID
}
