package controller

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"team/model"
	"team/orm"
	"team/web"
)

// File controller
type File int

// Register implements web.Controller interface.
func (f *File) Register(group *web.Router) {
	group.POST("/upload", f.upload)
	group.POST("/share", f.share)
	group.GET("/share/list", f.getShareList)
	group.GET("/share/:id", f.download)
	group.DELETE("/share/:id", f.deleteShare)
}

func (f *File) upload(c *web.Context) {
	formFiles := c.MultipartForm().File
	uid := c.Session.Get("uid").(int64)

	for _, files := range formFiles {
		for _, fh := range files {
			url, size := f.save(fh, uid)

			c.JSON(200, web.Map{
				"data": web.Map{
					"url":  url,
					"size": size,
				},
			})
			return
		}
	}

	c.JSON(200, web.Map{})
}

func (f *File) share(c *web.Context) {
	formFiles := c.MultipartForm().File
	uid := c.Session.Get("uid").(int64)

	for _, files := range formFiles {
		for _, fh := range files {
			url, size := f.save(fh, uid)

			orm.Insert(&model.Share{
				Name: fh.Filename,
				Path: url,
				UID:  uid,
				Time: time.Now(),
				Size: size,
			})

			c.JSON(200, web.Map{})
			return
		}
	}

	c.JSON(200, web.Map{})
}

func (f *File) getShareList(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `share`")
	assert(err == nil, "拉取文件列表失败")
	defer rows.Close()

	list := []map[string]interface{}{}
	for rows.Next() {
		one := &model.Share{}
		err = orm.Scan(rows, one)
		if err == nil {
			name, _ := model.FindUserInfo(one.UID)
			list = append(list, map[string]interface{}{
				"id":       one.ID,
				"name":     one.Name,
				"url":      one.Path,
				"uploader": name,
				"time":     one.Time.Format("2006-01-02 15:04:05"),
				"size":     one.Size,
			})
		}
	}

	c.JSON(200, web.Map{"data": list})
}

func (f *File) download(c *web.Context) {
	id := atoi(c.RouteValue("id"))
	share := &model.Share{ID: id}
	err := orm.Read(share)
	assert(err == nil, "文件不存在或已被删除")

	c.FileWithName(200, "."+share.Path, share.Name)
}

func (f *File) deleteShare(c *web.Context) {
	id := atoi(c.RouteValue("id"))
	orm.Delete("share", id)
	c.JSON(200, web.Map{})
}

func (f *File) save(fh *multipart.FileHeader, uploader int64) (string, int64) {
	reader, err := fh.Open()
	assert(err == nil, "打开上传文件失败")
	defer reader.Close()

	dir := fmt.Sprintf("uploads/%d", uploader)
	err = os.MkdirAll(dir, 0777)
	assert(err == nil, "创建上传目录失败")

	now := time.Now()
	path := fmt.Sprintf("%s/%d_%s", dir, now.Unix(), fh.Filename)
	writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	assert(err == nil, "保存上传文件失败")
	defer writer.Close()

	size, err := io.Copy(writer, reader)
	assert(err == nil, "写入上传文件失败")

	return ("/" + path), size
}
