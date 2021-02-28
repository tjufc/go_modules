# go_modules

+ 常见功能的go实现练习
+ 学习

## localcache

+ 依赖[一个开源lru实现](github.com/hashicorp/golang-lru)
+ `Get`策略：如果缓存没有，则调用rpc获取
+ 优化：根据使用场景优化rpc的调用

## distributed_lock

分布式锁: 基于redis实现

+ lockV1.go: redis官方文档`SETNX`内的方法实现
+ lockV2.go: 基于`SET NX EX`和释放锁脚本

以上实现都只能保证单机实例服务ok，在分布式集群服务的情况下可能存在问题。

+ TODO [Redlock](https://redis.io/topics/distlock)

## redis_cache

redis缓存