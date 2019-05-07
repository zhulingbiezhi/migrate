package main

import (
	"context"
	"fmt"
	redis2 "github.com/garyburd/redigo/redis"
	"klook.libs/logger"
	"math/rand"
	"migrate/common"
	"migrate/moudle/user"
	"migrate/moudle/user_v2"
	"migrate/redis"
	"runtime"
	"time"
)

const (
	NewWriteMask = iota
	NewReadMask  = iota << 1
	OldWriteMask = iota << 1
	OldReadMask  = iota << 1
)

var defaultMask = int64(OldWriteMask | OldReadMask)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	go func(ctx context.Context) {
		for {
			Write()
			select {
			case <-ctx.Done():
			default:
				//time.Sleep(time.Millisecond * 10)
			}
		}
	}(ctx)
	go func(ctx context.Context) {
		for {
			Write()
			select {
			case <-ctx.Done():
			default:
				//time.Sleep(time.Millisecond * 10)
			}
		}
	}(ctx)
	for {
		runtime.Gosched()
	}
}

func Write() {
	var err error
	flag, err := GetMigrateMask()
	if err != nil {
		logger.Error("redis error ", err)
		time.Sleep(time.Second * 2)
		return
	}
	var u *user.User
	if (OldWriteMask & flag) > 0 {
		if u, err = WriteOldUser(); err != nil {
			logger.Error("WriteOldUser error ", err)
			time.Sleep(time.Second * 2)
			return
		}
	}
	if (NewWriteMask & flag) > 0 {
		u2 := user_v2.UserV2(*u)
		if err = WriteNewUser(&u2); err != nil {
			logger.Error("WriteNewUser error ", err)
			time.Sleep(time.Second * 2)
			return
		}
	}
	if u.ID%100 == 0 {
		fmt.Println(u.ID)
	}
}

func GetMigrateMask() (int64, error) {
	i, err := redis.GetInt(common.MigrateKey)
	if err != nil {
		if err == redis2.ErrNil {
			return defaultMask, nil
		}
		return 0, err
	}
	return i, nil
}

func WriteOldUser() (*user.User, error) {
	var err error
	var id = fmt.Sprint(rand.Int() % 10)
	u := &user.User{
		Name:       "damon.hu -- " + id,
		FirstName:  "damon",
		LastName:   "hu",
		Email:      "damon.hu" + id + "@qq.com",
		Sex:        "男",
		Flag:       99,
		CreateTime: time.Now(),
	}
	u.ID, err = u.Insert()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return u, nil
}

func ReadOldUser() ([]*user.User, error) {
	t := time.Now()
	at := t.Add(time.Second * 10)
	users, err := user.Select(user.ScopeStartTime(&t), user.ScopeEndTime(&at))
	if err != nil {
		return nil, err
	}
	return users, nil
}

func WriteNewUser(u *user_v2.UserV2) error {
	var err error
	id := u.ID
	u.Flag = 0
	u.ID, err = u.Insert()
	if err != nil {
		logger.Error(err)
		return err
	}
	if u.ID != id {
		logger.Error("new and old id not match !", id, u.ID)
	}
	if err = redis.SetIntNX(common.MigrateEndID, u.ID); err != nil {
		return err
	}
	return nil
}

func ReadNewUser() ([]*user_v2.UserV2, error) {
	t := time.Now()
	at := t.Add(time.Second * 10)
	users, err := user_v2.Select(user.ScopeStartTime(&t), user.ScopeEndTime(&at))
	if err != nil {
		return nil, err
	}
	return users, nil
}
