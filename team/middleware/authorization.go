package middleware

import (
	"net/http"

	"team/model"
	"team/web"
)

// AutoLogin tries login user using cookie token.
func AutoLogin(next web.Handler) web.Handler {
	return func(c *web.Context) {
		if !c.Session.Has("uid") {
			cookie, err := c.Cookie(model.AutoLoginCookieKey)
			if err == nil {
				uid := model.TryAutoLogin(cookie.Value, c.RemoteIP())
				if uid < 0 {
					c.SetCookie(&http.Cookie{
						Name:   model.AutoLoginCookieKey,
						Value:  "",
						MaxAge: -1,
					})
				} else {
					c.Session.Set("uid", uid)
				}
			}
		}

		next(c)
	}
}

// MustLogined makes sure the client has logined in this server.
func MustLogined(next web.Handler) web.Handler {
	return func(c *web.Context) {
		if c.Session.Has("uid") {
			next(c)
		} else {
			c.JSON(http.StatusUnauthorized, web.Map{"err": "请先登录后操作"})
		}
	}
}

// MustLoginedAsAdmin makes sure the client has logined as Administrator
func MustLoginedAsAdmin(next web.Handler) web.Handler {
	return func(c *web.Context) {
		if !c.Session.Has("uid") {
			c.JSON(http.StatusUnauthorized, web.Map{"err": "请先登录后操作"})
			return
		}

		me := model.FindUser(c.Session.Get("uid").(int64))
		if me == nil {
			c.JSON(http.StatusUnauthorized, web.Map{"err": "请先登录后操作"})
			return
		}

		if me.IsSu {
			next(c)
		} else {
			c.JSON(http.StatusUnauthorized, web.Map{"err": "权限不足"})
		}
	}
}
