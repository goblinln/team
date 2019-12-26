package controller

import (
	"team/config"
	"team/model/user"
	"team/web"
)

// Home handler.
func Home(c *web.Context) {
	if !config.Default.Installed {
		c.JSON(200, web.Map{"data": "install"})
		return
	}

	for !c.Session.Has("uid") {
		cookie, err := c.Cookie(user.AutoLoginCookieKey)
		if err == nil {
			uid := user.AutoLogin(cookie.Value, c.RemoteIP())
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
