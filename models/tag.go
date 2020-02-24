package models

import "github.com/jinzhu/gorm"

type Tag struct {
	Model
	Name  string `json:"name"`
	State int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag, err error) {
	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

func GetTagTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, nil
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddTag(name string, state int) (bool, error) {
	err := db.Create(&Tag{
		Name:  name,
		State: state,
	}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func EditTag(id int, data map[string]interface{}) (bool, error) {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteTag(id int) (bool, error) {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func CleanAllTag() (bool, error) {
	err := db.Unscoped().Where("deleted_at != null").Delete(&Tag{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
