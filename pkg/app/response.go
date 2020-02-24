package app

import (
	"go_gin_base/models"
	"go_gin_base/pkg/e"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type SwagResponse struct {
	Code int          `json:"code" example:"200"`
	Msg  string       `json:"msg" example:"ok"`
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Lists     interface{}         `json:"lists"`
	Total     int                 `json:"total"`
	Errors    []models.ValidField `json:"error"`
	FieldName interface{}         `json:"field_name"`
}

type SwagErrorResponse struct {
	Code int               `json:"code" example:"500"`
	Msg  string            `json:"msg" example:"server error"`
	Data ErrorResponseData `json:"data"`
}

type ErrorResponseData struct {
	Errors []models.ValidField `json:"error"`
}

type SwagUploadResponse struct {
	Code int                `json:"code" example:"200"`
	Msg  string             `json:"msg" example:"ok"`
	Data UploadResponseData `json:"data"`
}

type UploadResponseData struct {
	ImageUrl     string `json:"image_url"`
	ImageSaveUrl string `json:"image_save_url"`
}

type SwagAuthResponse struct {
	Code int                `json:"code" example:"200"`
	Msg  string             `json:"msg" example:"ok"`
	Data UploadResponseData `json:"data"`
}

type AuthResponseData struct {
	Token string `json:"token"`
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}
