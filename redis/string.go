package redis

import (
	"github.com/garyburd/redigo/redis"
)

func GET(key string) (string, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.String(r.Do("GET", key))
}

func GETINT(key string) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Int64(r.Do("GET", key))
}

func FORCEGETINT(key string) int64 {
	r := redisPool.Get()
	defer r.Close()
	v, err := redis.Int64(r.Do("GET", key))
	if err == redis.ErrNil {
		return 0
	}
	return v
}

func SET(key, value string) bool {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("SET", key, value)
	return err == nil
}

func INCR(key string) error {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("INCR", key)
	if err != nil {
		return err
	}
	return nil
}

func INCRBY(key string, amount int64) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Int64(r.Do("INCRBY", key, amount))
}

func MGET(keys []string) (map[string]interface{}, error) {
	//r := redisPool.Get()
	//defer r.Close()
	//return redis.String(r.Do("MGET", key))

	r := redisPool.Get()
	defer r.Close()
	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	values, err := redis.Values(r.Do("MGET", args...))
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{}, 0)
	for i, k := range keys {
		if values[i] != nil {
			res[k] = values[i]
		}
	}
	return res, nil
}

func SETBIT(key string, number int32, value bool) bool {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("SETBIT", key, number, value)
	return err == nil
}

func GETBIT(key string, number int32) (bool, error) {
	r := redisPool.Get()
	defer r.Close()

	return redis.Bool(r.Do("GETBIT", key, number))
}

func SETEX(ex int, key, value interface{}) bool {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("SETEX", key, ex, value)
	return err == nil
}

func SETERR(key, value string) error {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("SET", key, value)
	return err
}

func SETNX(key, value interface{}) (bool, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Bool(r.Do("SETNX", key, value))
}
