package share

import (
	"os"
	"time"

	"team/model/user"
	"team/orm"
)

type (
	// Share schema
	Share struct {
		ID   int64     `json:"id"`
		Name string    `json:"name" orm:"type=VARCHAR(128),notnull"`
		Path string    `json:"url" orm:"type=VARCHAR(128),notnull"`
		UID  int64     `json:"uid"`
		Time time.Time `json:"time" orm:"default=CURRENT_TIMESTAMP"`
		Size int64     `json:"size"`
	}
)

// GetAll returns all shared files.
func GetAll() ([]map[string]interface{}, error) {
	rows, err := orm.Query("SELECT * FROM `share`")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := []map[string]interface{}{}
	for rows.Next() {
		one := &Share{}
		if err = orm.Scan(rows, one); err != nil {
			return nil, err
		}

		name, _ := user.FindInfo(one.UID)
		list = append(list, map[string]interface{}{
			"id":       one.ID,
			"name":     one.Name,
			"url":      one.Path,
			"uploader": name,
			"time":     one.Time.Format("2006-01-02 15:04:05"),
			"size":     one.Size,
		})
	}

	return list, nil
}

// Find shared file info by ID.
func Find(ID int64) (*Share, error) {
	file := &Share{ID: ID}
	err := orm.Read(file)
	return file, err
}

// Add a share file.
func Add(name, url string, uploader, size int64) error {
	_, err := orm.Insert(&Share{
		Name: name,
		Path: url,
		UID:  uploader,
		Time: time.Now(),
		Size: size,
	})

	return err
}

// Delete a shared file.
func Delete(ID int64) {
	file := &Share{ID: ID}
	err := orm.Read(file)
	if err == nil {
		os.Remove("." + file.Path)
		orm.Delete("share", ID)
	}
}
