package rediscache

// redis cache

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	SERVERADDR  = "127.0.0.1:6379"
	MAXIDLE     = 2
	MAXACTIVE   = 5
	IDLETIMEOUT = 240 // 240s

	// conn timeout in ms
	CONNECTTIMEOUT = 100
	READTIMEOUT    = 300
	WRITETIMEOUT   = 500
)

var (
	st *redis.Pool // pool single instance
)

// Init
func InitRedisCache(addr string) error {
	st = &redis.Pool{
		MaxIdle:     MAXIDLE,
		MaxActive:   MAXACTIVE,
		IdleTimeout: IDLETIMEOUT * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.DialTimeout("tcp", SERVERADDR, CONNECTTIMEOUT*time.Millisecond,
				READTIMEOUT*time.Millisecond, WRITETIMEOUT*time.Millisecond)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	// test connection
	testConn := st.Get()
	if _, err := testConn.Do("PING"); err != nil {
		st.Close()
		return errors.New("InitRedisCache error. redis conn fail")
	}
	testConn.Close()

	return nil
}

// Close
func Close() {
	if st != nil {
		st.Close()
	}
}

// GetConn
func GetConn() redis.Conn {
	return st.Get()
}
