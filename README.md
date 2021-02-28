# go_modules

+ 常见功能的go实现练习
+ 学习

## localcache

+ 依赖[一个开源lru实现](github.com/hashicorp/golang-lru)
+ `Get`策略：如果缓存没有，则调用rpc获取
+ 优化：根据使用场景优化rpc的调用

## distributed_lock1

分布式锁1: 基于redis 官方文档`SETNX`内的方法实现

## redis_cache

redis缓存