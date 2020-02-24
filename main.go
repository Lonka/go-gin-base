package main

import (
	"go_gin_base/models"
	"go_gin_base/pkg/gredis"
	"go_gin_base/pkg/setting"
	"go_gin_base/routers"

	// "context"
	// "fmt"
	// "net/http"
	// "os"
	// "os/signal"
	// "time"

	logging "go_gin_base/hosted/logging_service"
	mqtt "go_gin_base/hosted/mqtt_service"
	websocket "go_gin_base/hosted/websocket_service"
)

// @title BaBao Wealthy API
// @version v1

// @contact.name Lonka.Liu
// @contact.url https://#
// @contact.email lonka.liu@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	logging.Setup()
	setting.Setup()
	models.Setup()
	if setting.Redis.Use {
		gredis.Setup()
	}
	if setting.WebSocket.Use {
		websocket.Setup()
	}
	if setting.Mqtt.Use {
		mqtt.Setup()
	}

	router := routers.InitRouter()

	router.Run()

}
