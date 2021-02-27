package localcache

import (
	"time"
)

type RPCResult struct {
	ID      uint64 `json:"id"`
	Content string `json:"content"`
}

type RestfulRPC interface {
	Get(key string) (*RPCResult, error)
	Post(key string, content string) error
}

const (
	MYRPC = "MyRPC"
)

func NewRestfulRPC(rpcType string) RestfulRPC {
	if rpcType == MYRPC {
		return &MyRPC{}
	}
	return &MyRPC{}
}

// MyRPC
type MyRPC struct{}

// Get 啥也不干，干等5s
func (t *MyRPC) Get(key string) (*RPCResult, error) {
	println("rpc start")
	time.Sleep(time.Second * 5)
	println("rpc end")
	return nil, nil
}

func (t *MyRPC) Post(key string, content string) error {
	return nil
}
