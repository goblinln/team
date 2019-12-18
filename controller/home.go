package controller

import (
	"team/model"
	"team/web"
)

// Home handler.
func Home(c *web.Context) {
	if !model.Environment.Installed {
		c.JSON(200, web.Map{"data": "install"})
		return
	}

	for !c.Session.Has("uid") {
		cookie, err := c.Cookie(model.AutoLoginCookieKey)
		if err == nil {
			uid := model.TryAutoLogin(cookie.Value, c.RemoteIP())
			if uid < 0 {
				c.EndSession()
			} else {
				c.Session.Set("uid", uid)
				break
			}
		}

		c.JSON(200, web.Map{"data": "login"})
		return
	}

	c.JSON(200, web.Map{"data": "home"})
}
