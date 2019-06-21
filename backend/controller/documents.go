package controller

import (
	"time"

	"team/model"
	"team/orm"
	"team/web"
)

// Document controller
type Document int

// Register implements web.Controller interface.
func (d *Document) Register(group *web.Router) {
	group.POST("", d.create)
	group.GET("/list", d.getAll)
	group.GET("/:id", d.detail)
	group.PATCH("/:id/title", d.rename)
	group.PATCH("/:id/content", d.edit)
	group.DELETE("/:id", d.delete)
}

func (d *Document) create(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	title := c.PostFormValue("title").MustString("无效文档名")
	parent := c.PostFormValue("parent").MustInt("无效父节点")

	rows, err := orm.Query("SELECT COUNT(*) FROM `document` WHERE `parent`=? AND `title`=?", parent, title)
	web.Assert(err == nil, "数据库连接失败")
	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	web.Assert(count == 0, "同级目录下已存在同名文档")

	_, err = orm.Insert(&model.Document{
		Parent:   parent,
		Title:    title,
		Author:   uid,
		Modifier: uid,
		Time:     time.Now(),
		Content:  "",
	})
	web.Assert(err == nil, "写入数据失败")
	c.JSON(200, web.Map{})
}

func (d *Document) getAll(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `document`")
	web.Assert(err == nil, "数据库连接失败")
	defer rows.Close()

	list := []map[string]interface{}{}
	for rows.Next() {
		one := &model.Document{}
		err = orm.Scan(rows, one)
		if err == nil {
			creator, _ := model.FindUserInfo(one.Author)
			modifier, _ := model.FindUserInfo(one.Modifier)

			list = append(list, map[string]interface{}{
				"id":       one.ID,
				"parent":   one.Parent,
				"title":    one.Title,
				"creator":  creator,
				"modifier": modifier,
				"tile":     one.Time.Format(model.TaskTimeFormat),
			})
		}
	}

	c.JSON(200, web.Map{"data": list})
}

func (d *Document) detail(c *web.Context) {
	id := c.RouteValue("id").MustInt("")
	doc := &model.Document{ID: id}
	err := orm.Read(doc)
	web.Assert(err == nil, "读取文档信息失败")

	creator, _ := model.FindUserInfo(doc.Author)
	modifier, _ := model.FindUserInfo(doc.Modifier)

	c.JSON(200, web.Map{
		"data": map[string]interface{}{
			"id":       doc.ID,
			"parent":   doc.Parent,
			"title":    doc.Title,
			"creator":  creator,
			"modifier": modifier,
			"tile":     doc.Time.Format(model.TaskTimeFormat),
			"content":  doc.Content,
		},
	})
}

func (d *Document) rename(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	did := c.RouteValue("id").MustInt("")
	title := c.PostFormValue("title").MustString("无效文档名")

	doc := &model.Document{ID: did}
	err := orm.Read(doc)
	web.Assert(err == nil, "读取文档信息失败")

	doc.Title = title
	doc.Modifier = uid
	doc.Time = time.Now()
	err = orm.Update(doc)
	web.Assert(err == nil, "写入数据失败")

	c.JSON(200, web.Map{})
}

func (d *Document) edit(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	did := c.RouteValue("id").MustInt("")
	content := c.PostFormValue("content").MustString("内容不可为空")

	doc := &model.Document{ID: did}
	err := orm.Read(doc)
	web.Assert(err == nil, "读取文档信息失败")

	doc.Content = content
	doc.Modifier = uid
	doc.Time = time.Now()
	err = orm.Update(doc)
	web.Assert(err == nil, "写入数据失败")

	c.JSON(200, web.Map{})
}

func (d *Document) delete(c *web.Context) {
	id := c.RouteValue("id").MustInt("")

	doc := &model.Document{ID: id}
	err := orm.Read(doc)
	web.Assert(err == nil, "文档不存在或已被删除")

	orm.Exec("UPDATE `document` SET `parent`=? WHERE `parent`=?", doc.Parent, id)
	orm.Delete("document", id)
	c.JSON(200, web.Map{})
}
