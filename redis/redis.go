package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379",
				redis.DialConnectTimeout(time.Second*10),
				redis.DialReadTimeout(time.Second*30),
				redis.DialWriteTimeout(time.Second*30))
			if err != nil {
				return nil, err
			}
			//if _, err := c.Do("AUTH", password); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Get(key string) (string, error) {
	return redis.String(pool.Get().Do("GET", key))
}

func GetInt(key string) (int64, error) {
	return redis.Int64(pool.Get().Do("GET", key))
}

func SetIntNX(key string, val int64) error {
	_, err := pool.Get().Do("SET", key, val, "NX")
	if err != nil {
		return err
	}
	return nil
}

func SetInt(key string, val int64) error {
	_, err := pool.Get().Do("SET", key, val)
	if err != nil {
		return err
	}
	return nil
}
