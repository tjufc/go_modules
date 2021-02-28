package distributed_lock

// lockV1.go 分布式锁的一种实现
// redis命令SETNX官方文档中介绍的方法实现

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	cache "go_modules/rediscache"
)

const (
	// key expire time 300s
	KEYEXIRE = 300
)

type LockV1 struct {
	key    string
	expire int64
	rand   string
}

// Lock
func (t *LockV1) Lock() error {
	// make key and value
	mKey := makeKey(t.key)
	nt := time.Now().Unix()
	ext := nt + t.expire
	mValue := t.makeValue(nt, ext)

	// get redis connection
	conn := cache.GetConn()
	defer conn.Close()

	// try lock by SETNX
	res, err := redis.Int(conn.Do("SETNX", mKey, mValue))
	if err != nil {
		return err
	}
	// lock success. set expire and return
	if res == 1 {
		conn.Do("EXPIRE", mKey, KEYEXIRE)
		return nil
	}

	// lock has been acquired
	curValue, _ := redis.String(conn.Do("GET", mKey))
	_, curExt := t.parseValue(curValue)
	if curExt > nt {
		return ErrKeyConflict
	}

	// lock expired or deleted. try lock
	ret, err := redis.String(conn.Do("GETSET", mKey, mValue))
	if err != nil && err != redis.ErrNil {
		return err
	}
	// lock has been acquired
	if err == nil && ret != mValue {
		return ErrKeyConflict
	}

	// lock success. set expire and return
	conn.Do("EXPIRE", mKey, KEYEXIRE)
	return nil
}

// UnLock
func (t *LockV1) UnLock() error {
	mKey := makeKey(t.key)
	nt := time.Now().Unix()

	conn := cache.GetConn()
	defer conn.Close()

	// DEL key only if key not expired
	res, err := redis.String(conn.Do("GET", mKey))
	if err != nil {
		if err == redis.ErrNil {
			return nil
		}
		return err
	}
	_, ext := t.parseValue(res)
	if ext > nt {
		_, err := conn.Do("DEL", mKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// makeValue - rand_curTime_expireTime
func (t *LockV1) makeValue(curTime int64, expireTime int64) string {
	return fmt.Sprintf("%s_%d_%d", t.rand, curTime, expireTime)
}

// parseValue - return (0, 0) if value format error
func (t *LockV1) parseValue(value string) (int64, int64) {
	res := strings.Split(value, "_")
	if len(res) != 3 {
		return 0, 0
	}
	cur, err := strconv.ParseInt(res[1], 10, 64)
	if err != nil {
		return 0, 0
	}
	curEx, err := strconv.ParseInt(res[2], 10, 64)
	if err != nil {
		return 0, 0
	}
	return cur, curEx
}
