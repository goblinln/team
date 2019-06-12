package controller

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"team/model"
	"team/web"
	"team/web/orm"
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
			url, size, err := f.save(fh, uid)
			if err != nil {
				c.JSON(200, &web.JObject{"err": err.Error()})
				return
			}

			c.JSON(200, &web.JObject{
				"data": &web.JObject{
					"url":  url,
					"size": size,
				},
			})
			return
		}
	}

	c.JSON(200, &web.JObject{})
}

func (f *File) share(c *web.Context) {
	formFiles := c.MultipartForm().File
	uid := c.Session.Get("uid").(int64)

	for _, files := range formFiles {
		for _, fh := range files {
			url, size, err := f.save(fh, uid)
			if err != nil {
				c.JSON(200, &web.JObject{"err": err.Error()})
				return
			}

			orm.Insert(&model.Share{
				Name: fh.Filename,
				Path: url,
				UID:  uid,
				Time: time.Now(),
				Size: size,
			})

			c.JSON(200, &web.JObject{})
			return
		}
	}

	c.JSON(200, &web.JObject{})
}

func (f *File) getShareList(c *web.Context) {
	rows, err := orm.Query("SELECT * FROM `share`")
	if err != nil {
		c.JSON(200, &web.JObject{"err": "拉取文件列表失败"})
		return
	}

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

	c.JSON(200, &web.JObject{"data": list})
}

func (f *File) download(c *web.Context) {
	id := atoi(c.RouteValue("id"))
	share := &model.Share{ID: id}
	err := orm.Read(share)
	if err != nil {
		c.HTML(404, "文件不存在或已被删除")
		return
	}

	c.FileWithName(200, "."+share.Path, share.Name)
}

func (f *File) deleteShare(c *web.Context) {
	id := atoi(c.RouteValue("id"))
	orm.Delete("share", id)
	c.JSON(200, &web.JObject{})
}

func (f *File) save(fh *multipart.FileHeader, uploader int64) (string, int64, error) {
	reader, err := fh.Open()
	if err != nil {
		return "", 0, errors.New("打开上传文件失败")
	}

	defer reader.Close()

	dir := fmt.Sprintf("uploads/%d", uploader)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return "", 0, errors.New("创建上传目录失败")
	}

	now := time.Now()
	path := fmt.Sprintf("%s/%d_%s", dir, now.Unix(), fh.Filename)
	writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", 0, errors.New("保存上传文件失败")
	}

	defer writer.Close()

	size, err := io.Copy(writer, reader)
	if err != nil {
		return "", 0, errors.New("写入上传文件失败")
	}

	return ("/" + path), size, nil
}
