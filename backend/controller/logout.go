package controller

import (
	"team/web"
)

// Logout handler.
func Logout(c *web.Context) {
	c.EndSession()
	c.JSON(200, web.JObject{})
}
