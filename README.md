# Go Gin Base

An example of gin base on [eddycjy/go-gin-example](https://github.com/eddycjy/go-gin-example)


## How to run
### Required
 - Install MySQL and create a **db database** and import [SQL](https://github.com/Lonka/go-gin-base/db.sql)
 - [Swagger](https://github.com/go-swagger/go-swagger)

### Conf
Setting configuration via `conf/app.ini`.

### Run
```
$ cd ./go-gin-base
$ swag init
$ go run main.go
```

## Features
 - Go Mod
 - gin
 - gin cors
 - Jwt-go
 - go-playground/validator
 - gorm (callback)
 - Redis
 - MQTT Client (mosquitto ca)
 - WebSocket (broadcast, namesapce, event)
 - logging
 - export excel
 - minimum docker image


