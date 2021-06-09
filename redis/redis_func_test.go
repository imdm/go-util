package redis_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/imdm/go-util/redis"
)

func TestPIPELINE(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)

	key := "foo"
	redis.DEL(key)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			cmds := []string{
				"INCR",
				"EXPIRE",
			}
			args := map[string][]interface{}{
				"INCR":   {key},
				"EXPIRE": {key, 3},
			}
			err := redis.PIPELINE(cmds, args)
			assert.Nil(t, err)
			t.Log("goroutine: index: ", i, "foo: ", redis.FORCEGETINT(key))
		}(i)
	}
	wg.Wait()

	t.Log(redis.FORCEGETINT(key))
	assert.Equal(t, int64(10), redis.FORCEGETINT(key))
}

func TestHMSET(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)

	key := "test_order_"
	m := map[interface{}]interface{}{1: 2, 2: 3, 3: 4, 4: 5}
	err := redis.HMSET(key, m)
	assert.Nil(t, err)
}

func TestHGETINT(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)
	value, err := redis.HGETINT("test_order_", 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), value)

}

func TestHGETALL(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)
	retM, err := redis.HGETALL("test_order_")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(retM))
	assert.Equal(t, []byte(strconv.Itoa(2)), retM["1"])
	assert.Equal(t, []byte(strconv.Itoa(3)), retM["2"])
	assert.Equal(t, []byte(strconv.Itoa(4)), retM["3"])
	assert.Equal(t, []byte(strconv.Itoa(5)), retM["4"])
}

func TestMGET(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)
	retM, err := redis.MGET([]string{"1", "2", "3", "4"})
	assert.Nil(t, err)
	assert.Equal(t, 4, len(retM))
	assert.Equal(t, []byte(strconv.Itoa(2)), retM["1"])
	assert.Equal(t, []byte(strconv.Itoa(3)), retM["2"])
	assert.Equal(t, []byte(strconv.Itoa(4)), retM["3"])
	assert.Equal(t, []byte(strconv.Itoa(5)), retM["4"])
}

func TestBit(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)
	key := "test:bit:demo"
	isOk := redis.SETBIT(key, 1001, true)
	assert.Equal(t, true, isOk)
	value, err := redis.GETBIT(key, 1001)
	assert.Nil(t, err)
	assert.Equal(t, true, value)
	isOk = redis.SETBIT(key, 1001, false)
	assert.Equal(t, true, isOk)
	value, err = redis.GETBIT(key, 1001)
	assert.Nil(t, err)
	assert.Equal(t, false, value)
}

func TestINCRBY(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)
	key := "test_incr_by_1"
	redis.DEL(key)

	var (
		res int64
		err error
	)
	res = redis.FORCEGETINT(key)
	assert.Equal(t, res, int64(0))

	res, err = redis.INCRBY(key, 2)
	assert.Nil(t, err)
	assert.Equal(t, res, int64(2))

	res, err = redis.INCRBY(key, 1)
	assert.Nil(t, err)
	assert.Equal(t, res, int64(3))

	res, err = redis.INCRBY(key, -2)
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))
}

func TestZADD(t *testing.T) {
	dsn := "127.0.0.1:6379"
	pwd := ""
	redis.Init(dsn, pwd, 10, 0)
	key := "test_incr_by_1"
	redis.DEL(key)

	var (
		res int64
		err error
	)
	err = redis.ZADD("test_zadd", 1, "v_1")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))
	err = redis.ZADD("test_zadd", 3, "v_3")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))
	err = redis.ZADD("test_zadd", 2, "v_2")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))
	err = redis.ZADD("test_zadd", 4, "v_4")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))
	err = redis.ZADD("test_zadd", 5, "v_5")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))

	i, err := redis.ZCARD("test_zadd")
	println(i)
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))

	i, err = redis.ZCOUNT("test_zadd", 2, 4)
	println(i)
	assert.Nil(t, err)
	assert.Equal(t, res, int64(1))

}
