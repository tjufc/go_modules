package rediscache

import "testing"

func TestInitRedisCache(t *testing.T) {
	err := InitRedisCache(SERVERADDR)
	if err != nil {
		t.Errorf("error %+v", err)
		return
	}
	conn := st.Get()
	res, err := conn.Do("PING")
	if err != nil {
		t.Errorf("PING error %+v", err)
		return
	}
	t.Logf("PING res %+v", res)
	Close()
}
