package setting

import (
	"fmt"
	logging "go_gin_base/hosted/logging_service"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg     *ini.File
	RunMode string
)

var App = &AppSetting{}

type AppSetting struct {
	PageSize            int
	JwtSecret           string
	RuntimeRootPath     string
	LogSavePath         string
	LogFileExt          string
	PrefixUrl           string
	ExportExcelSavePath string
	QrCodeSavePath      string
	ImageSavePath       string
	ImagiSavePath       string
	FontSavePath        string
	ImageMaxSize        int
	ImageAllowExts      []string
	ClientID            string
}

var Server = &ServerSetting{}

type ServerSetting struct {
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var Database = &DatabaseSetting{}

type DatabaseSetting struct {
	DbType      string
	DbName      string
	User        string
	Password    string
	Host        string
	TablePrefix string
}

var Redis = &RedisSetting{}

type RedisSetting struct {
	Use             bool
	Host            string
	Password        string
	MaxIdle         int
	MaxActive       int
	IdleTimeout     time.Duration
	ExpireLongTime  int
	ExpireShortTime int
}

var WebSocket = &WebSocketSetting{}

type WebSocketSetting struct {
	Use bool
}

var Mqtt = &MqttSetting{}

type MqttSetting struct {
	Use  bool
	Host string
}

func Setup() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		logging.App.Fatal(fmt.Sprintf("Fail to parse 'conf/app.ini':'%v'", err))
	}
	loadBase()
	loadServer()
	loadApp()
	loadDatabase()
	loadRedis()
	loadWebSocket()
	loadMqtt()
}

func loadBase() {
	RunMode = Cfg.Section("").Key("RunMode").MustString("debug")
}

func loadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		logging.App.Fatal(fmt.Sprintf("Fail to get section 'server':'%v'", err))
	}
	Server.HttpPort = sec.Key("HttpPort").MustInt(8080)
	Server.ReadTimeout = time.Duration(sec.Key("ReadTimeout").MustInt(60)) * time.Second
	Server.WriteTimeout = time.Duration(sec.Key("WriteTimeout").MustInt(60)) * time.Second
}

func loadApp() {
	mapTo("app", App)
	App.ImageMaxSize = App.ImageMaxSize * 1024 * 1024
}

func loadDatabase() {
	mapTo("database", Database)
}

func loadRedis() {
	mapTo("redis", Redis)
	sec, err := Cfg.GetSection("server")
	if err != nil {
		logging.App.Fatal(fmt.Sprintf("Fail to get section 'server':'%v'", err))
	}
	Redis.IdleTimeout = time.Duration(sec.Key("IdleTimeout").MustInt(200)) * time.Second
}

func loadWebSocket() {
	mapTo("websocket", WebSocket)
}

func loadMqtt() {
	mapTo("mqtt", Mqtt)
}

func mapTo(section string, v interface{}) {
	err := Cfg.Section(section).MapTo(v)
	if err != nil {
		logging.App.Fatal(fmt.Sprintf("Cfg.MapTo %s err:'%v'", section, err))
	}
}
