package redis

import (
	"github.com/garyburd/redigo/redis"
)

func LPUSH(name string, v interface{}) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Int64(r.Do("LPUSH", name, v))
}

func RPUSH(name string, v interface{}) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Int64(r.Do("RPUSH", name, v))
}

func LPOP(name string) (interface{}, error) {
	r := redisPool.Get()
	defer r.Close()
	return r.Do("LPOP", name)
}

func LLEN(name string) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Int64(r.Do("LLEN", name))
}
