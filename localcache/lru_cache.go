package localcache

import (
	"errors"

	glru "github.com/hashicorp/golang-lru"
)

// 业务的本地缓存，使用lru策略
// 更新策略：如果本地缓存没有，则请求rpc获取后更新
// 缓存策略不一定是lru，仅仅是把代码拿来复用而已。

const (
	lruSize = 10
)

type rpcFunc func(key string) (string, error)

type LRUCache struct {
	lru *glru.Cache
	rpc rpcFunc
}

func NewLRUCache() *LRUCache {
	cache, _ := glru.New(lruSize)
	r := func(key string) (string, error) {
		res, err := NewRestfulRPC(MYRPC).Get(key)
		if err != nil {
			return "", err
		}
		return res.Content, nil
	}
	return &LRUCache{lru: cache, rpc: r}
}

// Get
func (t *LRUCache) Get(key string) (string, error) {
	// 一般策略：先get，如果没有请求rpc后，没有就set
	if res, ok := t.lru.Get(key); ok {
		content, typeOK := res.(string)
		if !typeOK {
			return "", errors.New("ret value type is not string")
		}
		return content, nil
	}

	// 如果有rpc耗时较大，且key是动态变化的，需要减少rpc的调用。怎么办？
	// 1. 把rpc调用放到锁里，只调用1次rpc
	// 2. 《go程序设计》中的方法，使用chan广播消息。但有个问题是如果rpc调用失败了，其余的goroutine也会获取到失败的信息。
	res, err := t.rpc(key)
	if err != nil {
		return "", errors.New("rpc Get error")
	}
	t.lru.Add(key, res)
	return res, nil
}

// Set
func (t *LRUCache) Set(key string, content string) error {
	t.lru.Add(key, content)
	return nil
}
