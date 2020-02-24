package tag_service

import (
	"encoding/json"
	"errors"
	"go_gin_base/models"
	"go_gin_base/pkg/export"
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/gredis"
	"go_gin_base/pkg/setting"
	"go_gin_base/service/cache_service"
	"io"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/unknwon/com"
)

type Tag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	State    int    `json:"state"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}

type AddTagRequest struct {
	Name  string `json:"name" validate:"required,max=100" example:"name"`
	State int    `json:"state"`
}

type EditTagRequest struct {
	Name  string `json:"name" validate:"max=100" example:"name"`
	State int    `json:"state"`
}

type ExportTagRequest struct {
	Name  string `json:"name" validate:"max=100" example:" "`
	State int    `json:"state" example:"-1"`
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)
	key, err := t.getCache("List", &cacheTags)
	if err == nil {
		return cacheTags, nil
	}

	tags, err = models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	if setting.Redis.Use {
		gredis.Set(key, tags, setting.Redis.ExpireShortTime)
	}
	return tags, nil
}

func (t *Tag) Count() (int, error) {
	var (
		count, cacheCount int
	)
	key, err := t.getCache("Count", &cacheCount)
	if err == nil {
		return cacheCount, nil
	}

	count, err = models.GetTagTotal(t.getMaps())
	if err != nil {
		return -1, err
	}

	if setting.Redis.Use {
		gredis.Set(key, count, setting.Redis.ExpireShortTime)
	}
	return count, err
}

func (t *Tag) Add() (bool, error) {
	ok, err := models.AddTag(t.Name, t.State)
	if ok {
		updateCache()
	}
	return ok, err
}

func (t *Tag) Edit() (bool, error) {
	ok, err := models.EditTag(t.ID, t.getMaps())
	if ok {
		updateCache()
	}
	return ok, err
}

func (t *Tag) Delete() (bool, error) {
	ok, err := models.DeleteTag(t.ID)
	if ok {
		updateCache()
	}
	return ok, err
}

func (t *Tag) CleanAll() (bool, error) {
	ok, err := models.CleanAllTag()
	if ok {
		updateCache()
	}
	return ok, err
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}
	xlsFile := xlsx.NewFile()
	sheet, err := xlsFile.AddSheet("Tags")
	if err != nil {
		return "", err
	}
	titles := []string{"ID", "Name", "Created By", "Created Time", "Updated By", "Updated Time"}
	row := sheet.AddRow()
	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}
	for _, v := range tags {
		values := []string{
			com.ToStr(v.ID),
			v.Name,
			v.CreatedBy,
			v.CreatedAt.Format("2006-01-02 03:04:05"),
			v.UpdatedBy,
			v.UpdatedAt.Format("2006-01-02 03:04:05"),
		}
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	time := time.Now().Format("20060102030405")
	fileName := "tags-" + time + export.ExcelEXT
	exportPath := export.GetExcelFullPath()
	err = file.CheckSrc(exportPath)
	if err != nil {
		return "", err
	}
	err = xlsFile.Save(exportPath + fileName)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (t *Tag) Import(r io.ReaderAt, size int64) error {
	xlsFile, err := xlsx.OpenReaderAt(r, size)
	if err != nil {
		return nil
	}
	sheet := xlsFile.Sheet["Tags"]
	for i, row := range sheet.Rows {
		if i > 0 {
			var data []string
			for _, cell := range row.Cells {
				data = append(data, cell.Value)
			}
			models.AddTag(data[0], 1)
		}
	}

	return nil
}

func updateCache() {
	if setting.Redis.Use {
		if cache_service.UpdateTagImmediatly {
			key := cache_service.GetTagsKey()
			gredis.LikeDeletes(key)
		}
	}
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	name := strings.TrimSpace(t.Name)
	if name != "" {
		maps["name"] = name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}

func (t *Tag) getCache(category string, result interface{}) (string, error) {
	var (
		key string
	)
	if !setting.Redis.Use {
		return key, errors.New("Not Exist")
	}

	cache := cache_service.Tag{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	switch category {
	case "List":
		key = cache.GetTagsKey()
		break
	case "Count":
		key = cache.CountTagsKey()
		break
	}

	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			return "", err
		} else {
			json.Unmarshal(data, &result)
			return key, nil
		}
	}
	return key, errors.New("Not Exist")
}
