package redis

func ZADD(key string, score int64, value string) error {
	r := redisPool.Get()
	defer r.Close()
	_, err := r.Do("ZADD", key, score, value)
	return err
}

func ZCARD(key string) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	i, err := r.Do("ZCARD", key)
	return i.(int64), err
}

func ZCOUNT(key string, min, max int64) (int64, error) {
	r := redisPool.Get()
	defer r.Close()
	i, err := r.Do("ZCOUNT", key, min, max)
	return i.(int64), err
}
