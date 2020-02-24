package v1

import (
	"go_gin_base/models"
	"go_gin_base/pkg/app"
	"go_gin_base/pkg/e"
	"go_gin_base/pkg/setting"
	"go_gin_base/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// @Summary Get Article
// @Tags Articles
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param ID path int true "Article ID"
// @Success 200 {object} app.SwagResponse
// @Failure 500 {object} app.SwagResponse
// @Router /v1/articles/{ID} [get]
func GetArticle(c *gin.Context) {
	appGin := app.Gin{C: c}

	var vfs []models.ValidField
	data := make(map[string]interface{})

	id := com.StrTo(c.Param("id")).MustInt()
	if ok, vf := util.HasMinOf(id, "id", "1"); !ok {
		vfs = append(vfs, *vf)
	}

	code := e.SUCCESS
	if len(vfs) > 0 {
		code = e.INVALID_PARAMS
		data["errors"] = vfs
	} else if !models.ExistArticleByID(id) {
		//code = e
	} else {
		data["data"] = models.GetArticle(id)
	}
	appGin.Response(http.StatusOK, code, data)
}

// @Summary Get multiple Articles
// @Tags Articles
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param tag_id query string false "Tag ID"
// @Param state query int false "State" Enums(0,1)
// @Success 200 {object} app.SwagResponse
// @Failure 500 {object} app.SwagResponse
// @Router /v1/articles [get]
func GetArticles(c *gin.Context) {
	appGin := app.Gin{C: c}

	var vfs []models.ValidField
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		if ok, vf := util.Contains(com.ToStr(state), "state", "0,1"); !ok {
			vfs = append(vfs, *vf)
		}
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		if ok, vf := util.HasMinOf(tagId, "tag_id", "1"); !ok {
			vfs = append(vfs, *vf)
		}
	}
	code := e.SUCCESS
	if len(vfs) > 0 {
		code = e.INVALID_PARAMS
		data["errors"] = vfs
		// } else if !models.ExistTagByID(tagId) {
		// 	code = e.ERROR_NOT_EXIST_TAG
	} else {
		data["lists"] = models.GetArticles(util.GetPage(c), setting.App.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)

	}
	appGin.Response(http.StatusOK, code, data)
}

// @Summary Add Article
// @Tags Articles
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param Params body models.AddArticleRequest true "Params"
// @Success 200 {object} app.SwagResponse
// @Router /v1/articles [post]
func AddArticle(c *gin.Context) {

	appGin := app.Gin{C: c}
	var vfs []models.ValidField
	data := make(map[string]interface{})
	var body models.AddArticleRequest

	c.ShouldBind(&body)

	if ok, vf := util.Contains(com.ToStr(body.State), "state", "0,1"); !ok {
		vfs = append(vfs, *vf)
	}
	if ok, vf := util.ValidStruct(body); !ok {
		vfs = append(vfs, *vf...)
	}

	code := e.SUCCESS
	if len(vfs) > 0 {
		code = e.INVALID_PARAMS
		data["errors"] = vfs
	} else if exist, _ := models.ExistTagByID(body.TagID); !exist {
		code = e.ERROR_NOT_EXIST_TAG
	} else {
		data["tag_id"] = body.TagID
		data["title"] = body.Title
		data["desc"] = body.Desc
		data["content"] = body.Content
		data["state"] = body.State
		models.AddArticle(data)
	}

	appGin.Response(http.StatusOK, code, data)
}

// @Summary Edit Article
// @Tags Articles
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param ID path int true "Article ID"
// @Param Params body models.EditArticleRequest true "Params"
// @Success 200 {object} app.SwagResponse
// @Router /v1/articles/{ID} [put]
func EditArticle(c *gin.Context) {
	//TODO : 該怎麼judy int 值是空的不更新
	appGin := app.Gin{C: c}
	var vfs []models.ValidField
	data := make(map[string]interface{})
	var body models.AddArticleRequest

	c.ShouldBind(&body)

	// var state int = -1
	// if arg := c.Query("state"); arg != "" {
	// 	state = com.StrTo(arg).MustInt()
	// 	if ok, vf := util.Contains(com.ToStr(state), "state", "0,1"); !ok {
	// 		vfs = append(vfs, *vf)
	// 	}
	// }

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

	code := e.SUCCESS
	if len(vfs) > 0 {
		code = e.INVALID_PARAMS
		data["errors"] = vfs
	} else if !models.ExistArticleByID(id) {
		//code = e.ERROR_NOT_EXIST_TAG
	} else if exist, _ := models.ExistTagByID(body.TagID); !exist {
		code = e.ERROR_NOT_EXIST_TAG
	} else {
		if body.TagID > 0 {
			data["tag_id"] = body.TagID
		}
		if body.Title != "" {
			data["title"] = body.Title
		}
		if body.Desc != "" {
			data["desc"] = body.Desc
		}
		if body.Content != "" {
			data["content"] = body.Content
		}
		data["state"] = body.State

		models.EditArticle(id, data)
	}

	appGin.Response(http.StatusOK, code, data)
}

// @Summary Delete Articles
// @Tags Articles
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param ID path int true "Article ID"
// @Success 200 {object} app.SwagResponse
// @Router /v1/articles/{ID} [delete]
func DeleteArticle(c *gin.Context) {
	appGin := app.Gin{C: c}

	var vfs []models.ValidField
	data := make(map[string]interface{})

	id := com.StrTo(c.Param("id")).MustInt()
	if ok, vf := util.HasMinOf(id, "id", "1"); !ok {
		vfs = append(vfs, *vf)
	}

	code := e.SUCCESS
	if len(vfs) > 0 {
		code = e.INVALID_PARAMS
		data["errors"] = vfs
	} else if !models.ExistArticleByID(id) {
		//code = e.ERROR_NOT_EXIST_TAG
	} else {
		models.DeleteArticle(id)
	}

	appGin.Response(http.StatusOK, code, data)
}
