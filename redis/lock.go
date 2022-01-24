package redis

import "github.com/garyburd/redigo/redis"

const (
	NX = "NX" // 不存在则执行
	EX = "EX" // 过期时间(秒)  PX 毫秒
	OK = "OK" // 操作成功
)

func Lock(key, requestId string, expire int64) bool {
	r := redisPool.Get()
	defer r.Close()
	msg, _ := redis.String(r.Do("SET", key, requestId, NX, EX, expire))
	if msg == OK {
		return true
	}
	return false
}

func UnLock(key, requestId string) bool {
	r := redisPool.Get()
	defer r.Close()
	msg, _ := redis.String(r.Do("GET", key))
	if msg == requestId {
		msg, _ := redis.Int64(r.Do("DEL", key))
		// 避免操作时间过长,自动过期时再删除返回结果为0
		if msg == 1 || msg == 0 {
			return true
		}
		return false
	}
	return false
}
