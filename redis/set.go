package redis

import "github.com/garyburd/redigo/redis"

func SADD(set string, item interface{}) bool {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("SADD", set, item)
	return err == nil
}

func SISMEMBER(set string, item interface{}) bool {
	r := redisPool.Get()
	defer r.Close()
	in, _ := redis.Bool(r.Do("SISMEMBER", set, item))
	return in
}

func SREM(set string, item interface{}) bool {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("SREM", set, item)
	return err == nil
}

func SMEMBERS(key string) ([]interface{}, error) {
	r := redisPool.Get()
	defer r.Close()
	values, err := redis.Values(r.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	return values, nil
}
