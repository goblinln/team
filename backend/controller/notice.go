package controller

import (
	"team/model"
	"team/orm"
	"team/web"
	"time"
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

	rows, err := orm.Query("SELECT `notice`.`id` AS id, `notice`.`tid` AS tid, `task`.`name` AS tname, `notice`.`operator` AS operator, `notice`.`time` AS time, `notice`.`event` AS ev FROM `notice` LEFT JOIN `task` ON `notice`.`tid`=`task`.`id` WHERE `uid`=?", uid)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "拉取通知信息失败"})
		return
	}

	defer rows.Close()

	type msg struct {
		ID       int64
		TID      int64
		TName    string
		Operator int64
		Time     time.Time
		Ev       int16
	}

	notices := []map[string]interface{}{}
	for rows.Next() {
		notice := &msg{}
		err = orm.Scan(rows, notice)
		if err == nil {
			operator, _ := model.FindUserInfo(notice.Operator)
			notices = append(notices, map[string]interface{}{
				"id":       notice.ID,
				"tid":      notice.TID,
				"tname":    notice.TName,
				"operator": operator,
				"time":     notice.Time.Format("2006-01-02 15:04:05"),
				"ev":       notice.Ev,
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
