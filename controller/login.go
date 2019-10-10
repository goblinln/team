package controller

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"team/model"
	"team/orm"
	"team/web"
)

// Login controller
type Login int

// Register implements web.Controller interface.
func (l *Login) Register(group *web.Router) {
	group.GET("", l.index)
	group.POST("", l.doLogin)
}

func (l *Login) index(c *web.Context) {
	if c.Session.Has("uid") {
		c.Redirect(302, "/")
		return
	}

	c.HTML(200, model.Index)
}

func (l *Login) doLogin(c *web.Context) {
	account := c.PostFormValue("account").MustString("登录帐号未填写")
	password := c.PostFormValue("password").MustString("登录密码未填写")
	remember, _ := c.PostFormValue("remember").Bool()

	hash := md5.New()
	hash.Write([]byte(password))

	user := &model.User{Account: account, Password: fmt.Sprintf("%X", hash.Sum(nil))}
	err := orm.Read(user, "account", "password")
	web.Assert(err == nil, "帐号或密码不正确")
	web.Assert(!user.IsLocked, "帐号已被禁止登录，请联系管理员解除锁定！")

	if remember {
		expire := time.Now().Add(30 * 24 * time.Hour)
		code := fmt.Sprintf("%d|%s|%s", user.ID, c.RemoteIP(), model.AutoLoginSecret)
		sign := md5.New()
		sign.Write([]byte(code))

		token := &model.AutoLoginToken{
			ID:   user.ID,
			IP:   c.RemoteIP(),
			Sign: fmt.Sprintf("%X", sign.Sum(nil)),
		}

		j, err := json.Marshal(token)
		if err == nil {
			c.SetCookie(&http.Cookie{
				Name:    model.AutoLoginCookieKey,
				Value:   base64.StdEncoding.EncodeToString(j),
				Expires: expire,
			})

			user.AutoLoginExpire = expire.Unix()
			orm.Update(user)
		}
	}

	model.Cache.SetUser(user)

	c.Session.Set("uid", user.ID)
	c.JSON(200, web.Map{})
}
