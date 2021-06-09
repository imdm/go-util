package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool

func Init(dsn, password string, maxIdle, db int) {
	registerRedisPool(dsn, password, maxIdle, db)
}

func registerRedisPool(dsn string, pw string, maxIdle, db int) {
	fmt.Println("Redis Info :", dsn, pw)
	pool := &redis.Pool{
		MaxIdle: maxIdle,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", dsn, redis.DialPassword(pw))
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) > time.Minute {
				_, err := c.Do("PING")
				return err
			}
			return nil
		},
	}
	conn := pool.Get()
	defer conn.Close()

	if conn.Err() != nil {
		fmt.Printf("failed to connect redis on %s/%d, max idle conn: %d, err: %s", dsn, db, maxIdle, conn.Err().Error())
		panic(conn.Err())
	}

	fmt.Println(fmt.Sprintf("connect redis on %s/%d, max idle conn: %d", dsn, db, maxIdle))

	if rel, err := conn.Do("PING"); err == nil {
		fmt.Println("Redis PING ", rel)
	}
	redisPool = pool
}
