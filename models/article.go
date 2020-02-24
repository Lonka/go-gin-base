package models

import ()

type AddArticleRequest struct {
	TagID   int    `json:"tag_id" validate:"required,min=1"`
	Title   string `json:"title" validate:"required,max=100"`
	Desc    string `json:"desc" validate:"required,max=255"`
	Name    string `json:"name" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=65535"`
	State   int    `json:"state"`
}

type EditArticleRequest struct {
	TagID   int    `json:"tag_id" validate:"min=1"`
	Title   string `json:"title" validate:"max=100"`
	Desc    string `json:"desc" validate:"max=255"`
	Name    string `json:"name" validate:"max=100"`
	Content string `json:"content" validate:"max=65535"`
	State   int    `json:"state"`
}

type Article struct {
	Model
	TagID   int    `json:"tag_id" gorm:"index"`
	Tag     Tag    `json:"tag"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	State   int    `json:"state"`
}

func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)
	if article.ID > 0 {
		return true
	}
	return false
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Update(data)
	return true
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:   data["tag_id"].(int),
		Title:   data["title"].(string),
		Desc:    data["desc"].(string),
		Content: data["content"].(string),
		State:   data["state"].(int),
	})
	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}
