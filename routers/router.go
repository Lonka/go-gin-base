package routers

import (
	"go_gin_base/hosted/websocket_service"
	"go_gin_base/pkg/export"
	"go_gin_base/pkg/imagi"
	"go_gin_base/pkg/qrcode"
	"go_gin_base/pkg/setting"
	"go_gin_base/pkg/upload"
	"net/http"
	"time"

	"go_gin_base/middleware/jwt"
	"go_gin_base/routers/api"
	v1 "go_gin_base/routers/api/v1"

	_ "go_gin_base/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	GinCors(r, "http://localhost:3000")

	gin.SetMode(setting.RunMode)

	r.GET("/auth", api.GetAuth)
	apiv1 := r.Group("/v1")

	apiv1.Use(jwt.JWT())
	{
		tags := apiv1.Group("/tags")
		{
			tags.GET("", v1.GetTags)
			tags.POST("", v1.AddTag)
			tags.PUT(":id", v1.EditTag)
			tags.DELETE(":id", v1.DeleteTag)

			tags.POST("export", v1.ExportTag)

			tags.GET("exportget", v1.ExportGetTag)

			tags.POST("import", v1.ImportTag)
		}

		articles := apiv1.Group("/articles")
		{
			articles.GET("", v1.GetArticles)
			articles.GET(":id", v1.GetArticle)
			articles.POST("", v1.AddArticle)
			articles.PUT(":id", v1.EditArticle)
			articles.DELETE(":id", v1.DeleteArticle)

		}

		//formdata image
		r.POST("/upload/image", api.UploadImage)
	}
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/export/excels", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/imagi", http.Dir(imagi.GetImagiFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	//url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if setting.WebSocket.Use {
		r.GET("/ws", websocket_service.Serve)
		//r.Use(GinMiddleware("http://localhost:3000"))
	}
	return r
}

func GinCors(r *gin.Engine, allowOrigin ...string) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigin,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			for _, arrow := range allowOrigin {
				if origin == arrow {
					return true
				}
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}
