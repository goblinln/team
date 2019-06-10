package controller

import (
	"crypto/md5"
	"fmt"

	"team/model"
	"team/web"
	"team/web/orm"
)

// User controller
type User int

// Register implements web.Controller interface.
func (u *User) Register(group *web.Router) {
	group.GET("", u.info)
	group.PATCH("/pswd", u.setPswd)
	group.PATCH("/avatar", u.setAvatar)
}

func (u *User) info(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	me := model.FindUser(uid)
	c.JSON(200, &web.JObject{"data": me})
}

func (u *User) setPswd(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	me := model.FindUser(uid)

	hashOld := md5.New()
	hashOld.Write([]byte(c.FormValue("oldPswd")))
	checkOld := fmt.Sprintf("%X", hashOld.Sum(nil))
	if checkOld != me.Password {
		c.JSON(200, &web.JObject{"err": "原始密码错误"})
		return
	}

	newPswd := c.FormValue("newPswd")
	cfmPswd := c.FormValue("cfmPswd")
	if newPswd != cfmPswd {
		c.JSON(200, &web.JObject{"err": "两次输入的新密码不一致"})
		return
	}

	hashNew := md5.New()
	hashNew.Write([]byte(newPswd))
	me.Password = fmt.Sprintf("%X", hashNew.Sum(nil))
	err := orm.Update(me)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "更新密码失败"})
	} else {
		c.JSON(200, &web.JObject{})
	}
}

func (u *User) setAvatar(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	me := model.FindUser(uid)

	fhs := c.MultipartForm().File["img"]
	if len(fhs) == 0 {
		c.JSON(200, &web.JObject{"err": "更新头像参数错误"})
		return
	}

	url, _, err := new(File).save(fhs[0], uid)
	if err != nil {
		c.JSON(200, &web.JObject{"err": err.Error()})
		return
	}

	me.Avatar = url
	err = orm.Update(me)
	if err != nil {
		c.JSON(200, &web.JObject{"err": "更新数据库失败"})
		return
	}

	c.JSON(200, &web.JObject{"data": me.Avatar})
}
