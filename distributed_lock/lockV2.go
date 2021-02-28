package distributed_lock

// lockV2.go - a distributed lock based on https://redis.io/commands/set

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"

	cache "go_modules/rediscache"
)

const (
	// unlock script
	unlockScript = `if redis.call("get",KEYS[1]) == ARGV[1]
then
	return redis.call("del",KEYS[1])
else
	return 0
end`
)

// LockV2
type LockV2 struct {
	key      string
	expire   int64
	rand     string
	setValue string
}

// Lock
func (t *LockV2) Lock() error {
	// make key and value
	mKey := makeKey(t.key)
	mValue := t.makeValue()

	conn := cache.GetConn()
	defer conn.Close()

	// try lock
	res, err := redis.String(conn.Do("SET", mKey, mValue, "NX", "EX", t.expire))
	if err != nil && err != redis.ErrNil {
		// redis.ErrNil has to be checked often
		return err
	}

	// lock success, set value
	if res == "OK" {
		t.setValue = mValue
		return nil
	}
	return ErrKeyConflict
}

// UnLock
func (t *LockV2) UnLock() error {
	script := redis.NewScript(1, unlockScript)
	conn := cache.GetConn()
	defer conn.Close()

	// eval script
	res, err := redis.Int(script.Do(conn, t.setValue))
	if err != nil || res == 0 {
		return err
	}
	return nil
}

// makeValue - [t.rand]+[timestamp]+[randomIntValue]
// this method may be rewrite due to different requirement
func (t *LockV2) makeValue() string {
	nt := time.Now().UnixNano()
	rand.Seed(nt)
	r := rand.Intn(100000)
	return fmt.Sprintf("%s_%d_%d", t.rand, nt, r)
}
