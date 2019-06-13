package controller

import (
	"team/model"
	"team/orm"
	"team/web"
)

// Notice controller
type Notice int

// Register implements web.Controller interface.
func (n *Notice) Register(group *web.Router) {
	group.GET("/list", n.mine)
	group.DELETE("/:id", n.deleteOne)
	group.DELETE("/all", n.deleteAll)
}

func (n *Notice) mine(c *web.Context) {
	uid := c.Session.Get("uid").(int64)

	rows, err := orm.Query("SELECT * FROM `notice` WHERE `uid`=?", uid)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "拉取通知信息失败"})
		return
	}

	defer rows.Close()

	notices := []map[string]interface{}{}
	for rows.Next() {
		notice := &model.Notice{}
		err = orm.Scan(rows, notice)
		if err == nil {
			notices = append(notices, map[string]interface{}{
				"id":      notice.ID,
				"uid":     uid,
				"time":    notice.Time.Format("2006-01-02 15:04:05"),
				"content": notice.Content,
			})
		}
	}

	c.JSON(200, &web.JObject{"data": notices})
}

func (n *Notice) deleteOne(c *web.Context) {
	orm.Delete("notice", atoi(c.RouteValue("id")))
	c.JSON(200, &web.JObject{})
}

func (n *Notice) deleteAll(c *web.Context) {
	orm.Exec("DELETE FROM `notice` WHERE `uid`=?", c.Session.Get("uid").(int64))
	c.JSON(200, &web.JObject{})
}
