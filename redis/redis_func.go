package redis

import (
	"github.com/garyburd/redigo/redis"
)

func DEL(key interface{}) {
	r := redisPool.Get()
	defer r.Close()
	r.Do("DEL", key)
}

func MDEL(keys []interface{}) error {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("DEL", keys...)
	return err
}

func FLUSHDB() {
	r := redisPool.Get()
	defer r.Close()
	r.Do("FLUSHDB")
}

func KEYS(pattern string) ([]string, error) {
	r := redisPool.Get()
	defer r.Close()
	return redis.Strings(r.Do("KEYS", pattern))
}

func EXPIRE(key interface{}, ex int) {
	r := redisPool.Get()
	defer r.Close()
	r.Do("EXPIRE", key, ex)
}

func EXISTS(key interface{}) bool {
	r := redisPool.Get()
	defer r.Close()
	exists, _ := redis.Bool(r.Do("EXISTS", key))
	return exists
}

// 不推荐使用， 可以使用 PIPELINE_V2
func PIPELINE(cmds []string, args map[string][]interface{}) error {
	r := redisPool.Get()
	defer r.Close()
	r.Send("MULTI")
	for _, v := range cmds {
		r.Send(v, args[v]...)
	}
	_, err := r.Do("EXEC")
	return err
}

type OneCmd struct {
	Cmd  string
	Args []interface{}
}

func PIPELINE_V2(cmds []*OneCmd) error {
	r := redisPool.Get()
	defer r.Close()
	r.Send("MULTI")
	for _, v := range cmds {
		r.Send(v.Cmd, v.Args...)
	}
	_, err := r.Do("EXEC")
	return err
}
