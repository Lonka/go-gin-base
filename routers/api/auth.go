package api

import (
	"go_gin_base/models"
	"go_gin_base/pkg/app"
	"go_gin_base/pkg/e"
	"go_gin_base/pkg/util"
	"go_gin_base/service/auth_service"

	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Authorization
// @Tags Authorization
// @Produce  json
// @Param username query string true "User Name"
// @Param password query string true "Password"
// @Success 200 {object} app.SwagAuthResponse
// @Success 400 {object} app.SwagErrorResponse
// @Success 401 {object} app.SwagErrorResponse
// @Success 500 {object} app.SwagErrorResponse
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appGin := app.Gin{C: c}

	var vfs []models.ValidField
	data := make(map[string]interface{})

	username := c.Query("username")
	password := c.Query("password")

	authService := auth_service.Auth{Username: username, Password: password}

	if ok, vf := util.ValidStruct(authService); !ok {
		vfs = append(vfs, *vf...)
	}

	if len(vfs) > 0 {
		data["errors"] = vfs
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, data)
		return
	}

	isExist, err := authService.Check()

	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, data)
		return
	}

	if !isExist {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH, data)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, data)
		return
	} else {
		data["token"] = token
	}

	appGin.Response(http.StatusOK, e.SUCCESS, data)
}
