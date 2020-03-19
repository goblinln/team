package controller

import (
	"net/http"

	"team/common/web"
	"team/config"
	"team/model/user"
)

// Login handler
func Login(c *web.Context) {
	account := c.PostFormValue("account").MustString("登录帐号未填写")
	password := c.PostFormValue("password").MustString("登录密码未填写")
	remember, _ := c.PostFormValue("remember").Bool()

	logined, cookie := user.Login(account, password, c.RemoteIP(), remember)
	if logined == nil {
		web.Assert(config.ExtraLoginProcessor != nil, "帐号或密码不正确")

		err := config.ExtraLoginProcessor.Login(account, password)
		if err != nil {
			web.Assert(false, "帐号验证失败")
		}

		logined, cookie = user.LoginViaCustom(account, password, c.RemoteIP(), remember)
	}

	web.Assert(logined != nil, "帐号或密码不正确")
	web.Assert(!logined.IsLocked, "帐号已被禁止登录，请联系管理员解除锁定！")

	if cookie != nil {
		c.SetCookie(&http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Expires: cookie.Expires,
		})
	}

	c.Session.Set("uid", logined.ID)
	c.JSON(200, web.Map{})
}

// Logout handler.
func Logout(c *web.Context) {
	c.EndSession()
	c.Redirect(302, "/")
}
