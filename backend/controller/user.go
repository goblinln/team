package controller

import (
	"crypto/md5"
	"fmt"

	"team/model"
	"team/orm"
	"team/web"
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
	c.JSON(200, web.Map{"data": me})
}

func (u *User) setPswd(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	me := model.FindUser(uid)

	hashOld := md5.New()
	hashOld.Write([]byte(c.FormValue("oldPswd")))
	checkOld := fmt.Sprintf("%X", hashOld.Sum(nil))
	assert(checkOld == me.Password, "原始密码错误")

	newPswd := c.FormValue("newPswd")
	cfmPswd := c.FormValue("cfmPswd")
	assert(newPswd == cfmPswd, "两次输入的新密码不一致")

	hashNew := md5.New()
	hashNew.Write([]byte(newPswd))
	me.Password = fmt.Sprintf("%X", hashNew.Sum(nil))
	err := orm.Update(me)
	assert(err == nil, "更新密码失败")

	c.JSON(200, web.Map{})
}

func (u *User) setAvatar(c *web.Context) {
	uid := c.Session.Get("uid").(int64)
	me := model.FindUser(uid)

	fhs := c.MultipartForm().File["img"]
	assert(len(fhs) > 0, "更新头像参数错误")

	me.Avatar, _ = new(File).save(fhs[0], uid)
	err := orm.Update(me)
	assert(err == nil, "更新数据库失败")

	c.JSON(200, web.Map{"data": me.Avatar})
}
