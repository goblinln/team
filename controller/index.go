package controller

import (
	"team/model"
	"team/web"
)

// Index handler.
func Index(c *web.Context) {
	if !model.Environment.Installed {
		c.Redirect(302, "/install")
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

		c.Redirect(302, "/login")
		return
	}

	c.HTML(200, model.Index)
}
