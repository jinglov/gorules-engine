package db

import (
	"github.com/jinglov/gomisc/drivers"
	"strconv"
	"testing"
	"time"
)

func connectRedisTest(t *testing.T) {
	var err error
	GetDb("redis").(*Redis).conn, err = drivers.NewRedisV2("172.25.23.57", 6379, 12,
		drivers.WithRedisPassword("actRedis"))
	if err != nil {
		t.Error(err)
	}
}

func TestRedis_AppendList(t *testing.T) {
	InitBloom(10000)
	connectRedisTest(t)
	cases := []struct {
		key   string
		value string
		exp   int
		res   int
	}{
		{"listkey1", "val1", 10, 1},
		{"listkey2", "val2", 3, 1},
		{"listkey3", "val31", 10, 1},
		{"listkey3", "val32", 10, 2},
		{"listkey3", "val32", 10, 0},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			v := GetDb("redis").AppendList(c.key, c.value, c.exp)
			if v != c.res {
				t.Errorf("key:%s , want:%d , res: %d", c.key, c.res, v)
			}
		})
	}
}

func TestRedis_GetList(t *testing.T) {
	connectRedisTest(t)
	cases := []struct {
		key   string
		value string
		exp   int
		res   int
	}{
		{"listkey1", "val1", 10, 1},
		{"listkey2", "val2", 3, 1},
		{"listkey3", "val31", 10, 1},
		{"listkey3", "val32", 10, 2},
		{"listkey3", "val32", 10, 2},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			GetDb("redis").AppendList(c.key, c.value, c.exp)
			v := GetDb("redis").GetList(c.key)
			t.Log(v)
			if len(v) != c.res {
				t.Errorf("key:%s , want:%d , res: %s", c.key, c.res, v)
			}
		})
	}
}

func TestRedis_GetListLen(t *testing.T) {
	connectRedisTest(t)

	cases := []struct {
		key   string
		value string
		exp   int
		res   int
	}{
		{"listkey1", "val1", 10, 1},
		{"listkey2", "val2", 3, 1},
		{"listkey3", "val31", 10, 1},
		{"listkey3", "val32", 10, 2},
		{"listkey3", "val32", 10, 2},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			GetDb("redis").AppendList(c.key, c.value, c.exp)
			v := GetDb("redis").GetListLen(c.key)
			if v != c.res {
				t.Errorf("key:%s , want:%d , res: %d", c.key, c.res, v)
			}
		})
	}
}

func TestRedis_Set(t *testing.T) {
	connectRedisTest(t)

	cases := []struct {
		key   string
		value string
		exp   int
		sleep int
		res   string
	}{
		{"key1", "val1", 10, 0, "val1"},
		{"key1", "val1", 10, 0, "val1"},
		{"key2", "val2", 3, 4, ""},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			setRes := GetDb("redis").Set(c.key, c.value, c.exp)
			t.Log("keyval", c.key, c.value)
			t.Log("setRes", setRes)
			time.Sleep(time.Duration(c.sleep) * time.Second)
			v := GetDb("redis").Get(c.key)
			if v != c.res {
				t.Errorf("key:%s , want:%s , res: %s", c.key, c.res, v)
			}
		})
	}
}

func TestRedis_Add(t *testing.T) {
	connectRedisTest(t)

	cases := []struct {
		key   string
		value int
		exp   int
		sleep int
		res   int
	}{
		{"key1", 1, 10, 0, 1},
		{"key2", 2, 3, 4, 2},
		{"key2", 3, 3, 4, 5},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			v := GetDb("redis").Add(c.key, c.value, c.exp)
			if v != c.res {
				t.Errorf("key:%s , want:%d , res: %d", c.key, c.res, v)
			}
		})
	}
}
