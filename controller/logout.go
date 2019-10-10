package controller

import (
	"team/web"
)

// Logout handler.
func Logout(c *web.Context) {
	c.EndSession()
	c.Redirect(302, "/login")
}
