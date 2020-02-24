package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	logging "go_gin_base/hosted/logging_service"
	"go_gin_base/pkg/setting"
	"io/ioutil"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	Client mqtt.Client
	opts   = mqtt.NewClientOptions()
)

func Serve() {
	tlsConfig := newTLSConfig()
	opts.AddBroker("tls://" + setting.Mqtt.Host)
	opts.SetClientID(setting.App.ClientID)
	opts.SetTLSConfig(tlsConfig)
	opts.SetKeepAlive(20 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Second * time.Duration(5))
	Client = mqtt.NewClient(opts)
	go connect()
}

func OnSubscribe(onSubscribeEvent mqtt.OnConnectHandler) {
	opts.OnConnect = onSubscribeEvent
}

func Subscribe(client mqtt.Client, topic string, qos byte, subscribeEvent mqtt.MessageHandler) {
	if token := client.Subscribe(topic, qos, subscribeEvent); token.Wait() && token.Error() != nil {
		logging.Mqtt.Error(fmt.Sprintf("Subscribe Error:%v", token.Error()))
	}
}

func UnSubscribe(topic string) {
	if token := Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		logging.Mqtt.Error(fmt.Sprintf("UnSubscribe Error:%v", token.Error()))
	}
}

func Publish(topic string, qos byte, retained bool, payload interface{}) {
	if token := Client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		logging.Mqtt.Error(fmt.Sprintf("Publish Error:%v", token.Error()))
	}
}

func connect() {
	d := time.Duration(time.Second * 5)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		if !Client.IsConnected() {
			if token := Client.Connect(); token.Wait() && token.Error() != nil {
				logging.Mqtt.Error(fmt.Sprintf("Connect Error: %v", token.Error()))
			} else {
				logging.Mqtt.Info(fmt.Sprintf("Host: %s is connected", setting.Mqtt.Host))
				return
			}

		}
	}
}

func newTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	caPem, err := ioutil.ReadFile("cert/ca.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(caPem)
	}
	return &tls.Config{
		RootCAs:            certpool,
		InsecureSkipVerify: true,
	}
}
