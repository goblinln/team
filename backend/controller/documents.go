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
	title := c.PostFormValue("title")
	parent := atoi(c.PostFormValue("parent"))

	rows, err := orm.Query("SELECT COUNT(*) FROM `document` WHERE `parent`=? AND `title`=?", parent, title)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "数据库连接失败"})
		return
	}

	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	if count > 0 {
		c.JSON(200, &web.JObject{"err": "同级目录下已存在同名文档"})
		return
	}

	_, err = orm.Insert(&model.Document{
		Parent:   parent,
		Title:    title,
		Author:   uid,
		Modifier: uid,
		Time:     time.Now(),
		Content:  "",
	})
	if err != nil {
		c.JSON(200, &web.JObject{"err": "写入数据失败"})
		return
	}

	c.JSON(200, &web.JObject{})
}

func (d *Document) getAll(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `document`")
	if err != nil {
		c.JSON(200, &web.JObject{"err": "数据库连接失败"})
		return
	}

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

	c.JSON(200, &web.JObject{"data": list})
}

func (d *Document) detail(c *web.Context) {
	id := atoi(c.RouteValue("id"))
	doc := &model.Document{ID: id}
	err := orm.Read(doc)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "读取文档信息失败"})
		return
	}

	creator, _ := model.FindUserInfo(doc.Author)
	modifier, _ := model.FindUserInfo(doc.Modifier)

	c.JSON(200, &web.JObject{
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
	did := atoi(c.RouteValue("id"))
	title := c.PostFormValue("title")

	doc := &model.Document{ID: did}
	err := orm.Read(doc)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "读取文档信息失败"})
		return
	}

	doc.Title = title
	doc.Modifier = uid
	doc.Time = time.Now()
	err = orm.Update(doc)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "写入数据失败"})
		return
	}

	c.JSON(200, &web.JObject{})
}

func (d *Document) edit(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	did := atoi(c.RouteValue("id"))
	content := c.PostFormValue("content")

	doc := &model.Document{ID: did}
	err := orm.Read(doc)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "读取文档信息失败"})
		return
	}

	doc.Content = content
	doc.Modifier = uid
	doc.Time = time.Now()
	err = orm.Update(doc)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "写入数据失败"})
		return
	}

	c.JSON(200, &web.JObject{})
}

func (d *Document) delete(c *web.Context) {
	id := atoi(c.RouteValue("id"))
	orm.Delete("document", id)
	c.JSON(200, &web.JObject{})
}
