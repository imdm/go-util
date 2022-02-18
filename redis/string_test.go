package redis_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/imdm/go-util/redis"
)

func TestSETNX(t *testing.T) {
	var (
		dsn     = "127.0.0.1:6379"
		pwd     = ""
		key     = fmt.Sprintf("test_setnx_%d", time.Now().Unix())
		okCount int
	)
	redis.Init(dsn, pwd, 10, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, err := redis.SETNX(key, 1)
			assert.Nil(t, err)
			if ok {
				okCount++
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, 1, okCount)
}

func TestINCRBYFLOAT(t *testing.T) {
	redis.Init("127.0.0.1:6379", "", 10, 15)
	key := "test_incrbyfloat"
	setRe := redis.SET(key, "0.32")
	assert.Equal(t, true, setRe)
	value, err := redis.INCRBYFLOAT(key, 0.5)
	assert.Empty(t, err)
	assert.Equal(t, 0.82, value)
	redis.DEL(key)
}
