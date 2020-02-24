package gredis

import (
	"encoding/json"
	"errors"
	"fmt"
	logging "go_gin_base/hosted/logging_service"
	"go_gin_base/pkg/setting"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

func Setup() error {
	if !setting.Redis.Use {
		return errors.New("Redis Not Use")
	}
	RedisConn = &redis.Pool{
		MaxIdle:     setting.Redis.MaxIdle,
		MaxActive:   setting.Redis.MaxActive,
		IdleTimeout: setting.Redis.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.Redis.Host)
			if err != nil {
				logging.Redis.Error(fmt.Sprintf("Connect '%s' Error : %v", setting.Redis.Host, err))
				return nil, err
			}
			if setting.Redis.Password != "" {
				if _, err := c.Do("AUTH", setting.Redis.Password); err != nil {
					c.Close()
					logging.Redis.Error(fmt.Sprintf("Authorization Error : %v", err))
					return nil, err
				}
			}
			logging.Redis.Info(fmt.Sprintf("Connected '%s'", setting.Redis.Host))
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		logging.Redis.Error(fmt.Sprintf("Parse '%s' Error : %v", key, err))
		return err
	}
	_, err = conn.Do("SET", key, value)
	if err != nil {
		logging.Redis.Error(fmt.Sprintf("Set '%s' String Error : %v", key, err),
			logging.GetField("Data", string(value)))
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		logging.Redis.Error(fmt.Sprintf("Set '%s' String Expire Error : %v", key, err),
			logging.GetField("Expire", fmt.Sprintf("%ds", time)))
		return err
	}
	logging.Redis.Info(fmt.Sprintf("Set '%s' String", key),
		logging.GetField("Data", string(value)),
		logging.GetField("Expire", fmt.Sprintf("%ds", time)))
	return nil
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()
	exist, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exist
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		logging.Redis.Error(fmt.Sprintf("Get '%s' String Error : %v", key, err))
		return nil, err
	}
	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	if _, err := redis.Bool(conn.Do("DEL", key)); err != nil {
		logging.Redis.Error(fmt.Sprintf("Del '%s' String Error : %v", key, err))
		return false, err
	}
	logging.Redis.Info(fmt.Sprintf("Del '%s' String", key))
	return true, nil
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
