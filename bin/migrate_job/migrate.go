package main

import (
	"fmt"
	redis2 "github.com/garyburd/redigo/redis"
	"klook.libs/logger"
	"migrate/common"
	"migrate/db"
	"migrate/moudle/user"
	"migrate/moudle/user_v2"
	"migrate/redis"
	"time"
)

func main() {
	for {
		max, err := user.MaxID()
		if err != nil {
			logger.Error("MaxID error", err)
			continue
		}
		maxV2, err := user_v2.MaxID()
		if err != nil {
			logger.Error("MaxID error", err)
			continue
		}
		start := maxV2
		end := maxV2 + 100
		i, err := redis.GetInt(common.MigrateEndID)
		if err != nil && err != redis2.ErrNil {
			logger.Error("redis get end id error", err)
			continue
		}
		if i > 0 && i < end {
			end = i
		}
		users, err := user.Select(user.ScopeIDRange(start, end))
		if err != nil {
			logger.Error("user Select error", err)
			continue
		}
		for _, u := range users {
			u2 := user_v2.UserV2(*u)
			if d := db.DB.Create(&u2); d.Error != nil {
				logger.Error("create user_v2 error: ", d.Error, u.ID)
				continue
			}
		}
		fmt.Printf("start: %d count: %d  proces: %.1f \n", start, end-start, float64(end)/float64(max)*100)
		time.Sleep(time.Millisecond*100)
	}
}
