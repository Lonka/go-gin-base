package api

import (
	logging "go_gin_base/hosted/logging_service"
	"go_gin_base/pkg/app"
	"go_gin_base/pkg/e"
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Upload Image
// @Tags Upload
// @Produce  json
// @param Authorization header string true "Authorization" default(bearer token)
// @Param image formData file true "Image"
// @Success 200 {object} app.SwagUploadResponse
// @Success 400 {object} app.SwagErrorResponse
// @Success 500 {object} app.SwagErrorResponse
// @Router /upload/image [post]
func UploadImage(c *gin.Context) {
	appGin := app.Gin{C: c}

	code := e.SUCCESS
	data := make(map[string]string)
	mfile, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Router.Warn(err.Error())
		code = e.ERROR
		appGin.Response(http.StatusOK, code, data)
		return
	}

	if image == nil {
		code = e.INVALID_PARAMS
		appGin.Response(http.StatusBadRequest, code, data)
		return
	}

	imageName := upload.GetImageMD5Name(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(mfile) {
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		appGin.Response(http.StatusBadRequest, code, data)
		return
	}

	err = file.CheckSrc(fullPath)
	if err != nil {
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
		logging.Router.Warn(err.Error())
		appGin.Response(http.StatusInternalServerError, code, data)
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
		logging.Router.Warn(err.Error())
		appGin.Response(http.StatusInternalServerError, code, data)
		return
	}

	data["image_url"] = upload.GetImageFullUrl(imageName)
	data["image_save_url"] = savePath + imageName
	appGin.Response(http.StatusOK, code, data)
}
