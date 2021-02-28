package distributed_lock

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/gomodule/redigo/redis"

	"go_modules/rediscache"
)

// lock_test

func init() {
	err := rediscache.InitRedisCache(rediscache.SERVERADDR)
	if err != nil {
		println(fmt.Sprintf("init redis cache error. %+v", err))
	}
}

func TestLockV1(t *testing.T) {
	r := rand.Intn(1000)
	lock := LockV1{
		key:    "testkey",
		expire: 300,
		rand:   strconv.Itoa(r),
	}
	// test Lock
	err := lock.Lock()
	if err != nil {
		t.Errorf("Lock error: %+v", err)
		return
	}

	// test Lock when lock is acquired
	err = lock.Lock()
	if err != ErrKeyConflict {
		t.Errorf("test Lock when lock is acquired failed. err returned: %+v", err)
		return
	}

	// print lock value and TTL
	conn := rediscache.GetConn()
	defer conn.Close()

	mKey := makeKey(lock.key)
	mValue, _ := redis.String(conn.Do("GET", mKey))
	ttl, _ := redis.Int(conn.Do("TTL", mKey))
	println(fmt.Sprintf("print lock value and TTL, value: %s, ttl: %d", mValue, ttl))

	// test UnLock
	err = lock.UnLock()
	if err != nil {
		t.Errorf("UnLock error: %+v", err)
		return
	}
}
