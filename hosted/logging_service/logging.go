package logging_service

import (
	"go_gin_base/pkg/logging"

	"go.uber.org/zap/zapcore"
)

var (
	App       logging.LoggerBase // main, middleware, pkg
	Model     logging.LoggerBase // model
	Router    logging.LoggerBase // router
	Service   logging.LoggerBase // service, hosted
	Redis     logging.LoggerBase // pkg/gredis
	WebSocket logging.LoggerBase // pkg/wscoket, hosted/websocket_service
	Mqtt      logging.LoggerBase // pkg/mqtt
)

func Setup() {
	newLogger(&App, "App")
	newLogger(&Model, "Model")
	newLogger(&Router, "Router")
	newLogger(&Service, "Service")
	newLogger(&Redis, "Redis")
	newLogger(&WebSocket, "WebSocket")
	newLogger(&Mqtt, "MQTT")
	logging.Setup()
}

func newLogger(log *logging.LoggerBase, serviceName string) {
	log.Logger = logging.NewLogger("", zapcore.InfoLevel, 100, 3, 7, true, serviceName)
	log.WriteToAll = true
	log.ServiceName = serviceName
}

func GetField(key, value string) zapcore.Field {
	return logging.GetField(key, value)
}
