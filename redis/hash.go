package redis

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

func HGETALL(name string) (map[string]interface{}, error) {
	r := redisPool.Get()
	defer r.Close()
	values, err := redis.Values(r.Do("HGETALL", name))
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("expects even number of values result")
	}
	ret := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		ret[string(values[i].([]byte))] = values[i+1]
	}
	return ret, nil
}

func HGET(name string, key interface{}) (interface{}, error) {
	r := redisPool.Get()
	defer r.Close()
	return r.Do("HGET", name, key)
}

func HGETINT(name string, key interface{}) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Int64(r.Do("HGET", name, key))
}

func FORCEHGETINT(name string, key interface{}) int64 {
	r := redisPool.Get()
	defer r.Close()
	v, err := redis.Int64(r.Do("HGET", name, key))
	if err == redis.ErrNil {
		return 0
	}
	return v
}

func HSET(name string, key, value interface{}) bool {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("HSET", name, key, value)
	return err == nil
}

func HMSET(name string, m map[interface{}]interface{}) error {
	r := redisPool.Get()
	defer r.Close()
	var args []interface{}
	args = append(args, name)
	for k, v := range m {
		args = append(args, k, v)
	}
	_, err := r.Do("HMSET", args...)
	return err
}

func HMGET(key string, fields []interface{}) (map[interface{}]interface{}, error) {
	r := redisPool.Get()
	defer r.Close()
	var args []interface{}
	args = append(args, key)
	for _, v := range fields {
		args = append(args, v)
	}
	values, err := redis.Values(r.Do("HMGET", args...))
	if err != nil {
		return nil, err
	}

	res := make(map[interface{}]interface{})
	for i, field := range fields {
		res[field] = values[i]
	}
	return res, nil
}

func HSETERR(name string, key, value interface{}) error {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("HSET", name, key, value)
	return err
}

func HINCRBY(name string, key interface{}, amount int64) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Int64(r.Do("HINCRBY", name, key, amount))
}

func HINCRBYFLOAT(name string, key interface{}, amount float64) (float64, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Float64(r.Do("HINCRBYFLOAT", name, key, amount))
}
