package mqtt_service

import (
	logging "go_gin_base/hosted/logging_service"
	mqttPkg "go_gin_base/pkg/mqtt"

	websocket "go_gin_base/hosted/websocket_service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Setup() {
	mqttPkg.OnSubscribe(OnSubscribe)
	mqttPkg.Serve()
}

func OnSubscribe(client mqtt.Client) {
	mqttPkg.Subscribe(client, "test/#", 0, test)
}

func test(client mqtt.Client, message mqtt.Message) {
	data := string(message.Payload())
	websocket.Broadcast("bbl", "Heartbeat", data)
	logging.Mqtt.Info(data)
}
