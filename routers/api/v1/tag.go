package v1

import (
	"fmt"
	"go_gin_base/models"
	"go_gin_base/pkg/app"
	"go_gin_base/pkg/e"
	"go_gin_base/pkg/export"
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/setting"
	"go_gin_base/pkg/util"
	"go_gin_base/service/tag_service"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unknwon/com"
)

// @Summary Get multiple Tags
// @Tags Tags
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param name query string false "Name"
// @Param state query int false "State" Enums(0,1)
// @Success 200 {object} app.SwagResponse
// @Failure 500 {object} app.SwagErrorResponse
// @Router /v1/tags [get]
func GetTags(c *gin.Context) {
	//GET("/tags")
	appGin := app.Gin{C: c}
	name := c.Query("name")
	data := make(map[string]interface{})

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	tagService := tag_service.Tag{Name: name, State: state, PageNum: util.GetPage(c), PageSize: setting.App.PageSize}

	tags, err := tagService.GetAll()
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, data)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, data)
		return
	}
	data["lists"] = tags
	data["total"] = count
	appGin.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary Add Tag
// @Tags Tags
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param Params body tag_service.AddTagRequest true "Params"
// @Success 200 {object} app.SwagResponse
// @Failure 400 {object} app.SwagErrorResponse
// @Failure 500 {object} app.SwagErrorResponse
// @Router /v1/tags [post]
func AddTag(c *gin.Context) {
	//POST("/tags")
	appGin := app.Gin{C: c}
	var vfs []models.ValidField
	data := make(map[string]interface{})
	var body tag_service.AddTagRequest

	c.ShouldBind(&body)

	if ok, vf := util.Contains(com.ToStr(body.State), "state", "0,1"); !ok {
		vfs = append(vfs, *vf)
	}
	if ok, vf := util.ValidStruct(body); !ok {
		vfs = append(vfs, *vf...)
	}

	if len(vfs) > 0 {
		data["errors"] = vfs
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, data)
		return
	}

	tagService := tag_service.Tag{Name: body.Name, State: body.State}

	exist, err := tagService.ExistByName()

	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if exist {
		appGin.Response(http.StatusOK, e.ERROR_EXIST_TAG, data)
		return
	}

	_, err = tagService.Add()
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary Edit Tag
// @Tags Tags
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param ID path int true "Tag ID"
// @Param Params body tag_service.EditTagRequest true "Params"
// @Success 200 {object} app.SwagResponse
// @Failure 400 {object} app.SwagErrorResponse
// @Failure 500 {object} app.SwagErrorResponse
// @Router /v1/tags/{ID} [put]
func EditTag(c *gin.Context) {
	//PUT("/tags/:id")
	appGin := app.Gin{C: c}

	var vfs []models.ValidField
	data := make(map[string]interface{})
	var body tag_service.EditTagRequest

	c.ShouldBind(&body)

	id := com.StrTo(c.Param("id")).MustInt()
	if ok, vf := util.HasMinOf(id, "id", "1"); !ok {
		vfs = append(vfs, *vf)
	}
	if ok, vf := util.Contains(com.ToStr(body.State), "state", "0,1"); !ok {
		vfs = append(vfs, *vf)
	}
	if ok, vf := util.ValidStruct(body); !ok {
		vfs = append(vfs, *vf...)
	}

	if len(vfs) > 0 {
		data["errors"] = vfs
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, data)
		return
	}

	tagService := tag_service.Tag{ID: id, Name: body.Name, State: body.State}

	exist, err := tagService.ExistByID()

	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if !exist {
		appGin.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, data)
		return
	}

	_, err = tagService.Edit()
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary Delete Tag
// @Tags Tags
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param ID path int true "Tag ID"
// @Success 200 {object} app.SwagResponse
// @Failure 400 {object} app.SwagErrorResponse
// @Failure 500 {object} app.SwagErrorResponse
// @Router /v1/tags/{ID} [delete]
func DeleteTag(c *gin.Context) {
	//DELETE("/tags/:id")
	appGin := app.Gin{C: c}

	var vfs []models.ValidField
	data := make(map[string]interface{})

	id := com.StrTo(c.Param("id")).MustInt()
	if ok, vf := util.HasMinOf(id, "id", "1"); !ok {
		vfs = append(vfs, *vf)
	}

	if len(vfs) > 0 {
		data["errors"] = vfs
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, data)
		return
	}

	tagService := tag_service.Tag{ID: id}

	exist, err := tagService.ExistByID()

	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if !exist {
		appGin.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, data)
		return
	}
	_, err = tagService.Delete()
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary Export Tag
// @Tags Tags
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param Params body tag_service.ExportTagRequest true "Params"
// @Success 200 {object} app.SwagResponse
// @Failure 400 {object} app.SwagErrorResponse
// @Router /v1/tags/export [post]
func ExportTag(c *gin.Context) {
	appGin := app.Gin{C: c}
	var vfs []models.ValidField
	data := make(map[string]interface{})
	var body tag_service.ExportTagRequest

	c.ShouldBind(&body)
	if ok, vf := util.ValidStruct(body); !ok {
		vfs = append(vfs, *vf...)
	}
	if len(vfs) > 0 {
		data["errors"] = vfs
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, data)
		return
	}
	tagService := tag_service.Tag{Name: body.Name, State: body.State}
	fileName, err := tagService.Export()
	if err != nil {
		appGin.Response(http.StatusOK, e.ERROR_EXPORT_FAIL, nil)
		return
	}
	data["export_url"] = export.GetExcelFullUrl(fileName)
	data["export_save_url"] = export.GetExcelPath() + fileName
	appGin.Response(http.StatusOK, e.SUCCESS, data)

}

// @Summary Export Tag
// @Tags Tags
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param name query string false "Name"
// @Param state query int false "State" Enums(-1,0,1)
// @Success 200 {object} app.SwagResponse
// @Failure 400 {object} app.SwagErrorResponse
// @Router /v1/tags/exportget [get]
func ExportGetTag(c *gin.Context) {
	appGin := app.Gin{C: c}

	name := c.Query("name")
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{Name: name, State: state, PageNum: util.GetPage(c), PageSize: setting.App.PageSize}

	fileName, err := tagService.Export()
	if err != nil {

		appGin.Response(http.StatusOK, e.ERROR_EXPORT_FAIL, nil)
		return
	}
	appGin.C.Header("Content-Description", "File Transfer")
	appGin.C.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	appGin.C.Header("Content-Transfer-Encoding", "binary")
	appGin.C.Header("Content-Type", "application/octet-stream")
	appGin.C.Header("Expires", "0")
	appGin.C.Header("Cache-Control", "must-revalidate")
	appGin.C.Header("Pragma", "public")
	filePath := export.GetExcelPath() + fileName
	if isExist := file.CheckExist(filePath); isExist == false {
		// if the path not found
		http.NotFound(appGin.C.Writer, appGin.C.Request)
		return
	}
	appGin.C.File(setting.App.RuntimeRootPath + filePath)
}

// @Summary Import Tag
// @Tags Tags
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param file formData file true "File"
// @Success 200 {object} app.SwagUploadResponse
// @Router /v1/tags/import [post]
func ImportTag(c *gin.Context) {
	appGin := app.Gin{C: c}
	file, _, err := appGin.C.Request.FormFile("file")
	if err != nil {
		appGin.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	tagService := tag_service.Tag{}
	size, err := file.Seek(0, 2)
	err = tagService.Import(file, size)
	if err != nil {
		appGin.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appGin.Response(http.StatusOK, e.SUCCESS, nil)
}
