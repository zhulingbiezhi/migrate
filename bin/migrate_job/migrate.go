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
	var endID = int64(0)
	var curID = int64(0)
	var max = int64(0)
	var process = float64(0.0)
	var err error
	for {
		curID, err = redis.GetInt(common.MigrateUserCurID)
		if err != nil && err != redis2.ErrNil {
			logger.Error("redis get current id error", err)
			continue
		}

		start := curID
		end := curID + 100

		if endID <= 0 {
			endID, err = redis.GetInt(common.MigrateUserEndID)
			if err != nil && err != redis2.ErrNil {
				logger.Error("redis get end id error", err)
				continue
			}
			max, err = user.MaxID()
			if err != nil {
				logger.Error("MaxID error", err)
				continue
			}
		} else {
			max = endID
			if endID < end {
				end = endID - 1
			}
		}

		users, err := user.Select(user.ScopeIDRange(start, end))
		if err != nil {
			logger.Error("user Select error", err)
			continue
		}
		if len(users) == 0{
			fmt.Println("没有数据")
			if endID >0 {
				fmt.Println("迁移结束")
				break
			}
			time.Sleep(time.Second)
			continue
		}

		var maxID int64
		for _, u := range users {
			u2 := user_v2.UserV2(*u)
			if d := db.DB.Create(&u2); d.Error != nil {
				logger.Error("create user_v2 error: ", d.Error, u.ID)
				continue
			}
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		if maxID > curID{
			if err = redis.SetInt(common.MigrateUserCurID, maxID);err != nil {
				logger.Error("redis set cur id error", err)
				continue
			}
		}
		process = float64(end)/float64(max)*100
		fmt.Printf("start: %d end: %d , end_id: %d  proces: %.2f \n", start, maxID, endID, process)
		time.Sleep(time.Millisecond * 100)
	}
}
