package db

import (
	"github.com/jinglov/gomisc/drivers"
	"strconv"
	"testing"
	"time"
)

var cfg DBConf

func init() {
	cfg.Redis, _ = drivers.NewRedis(nil, "test", "192.168.57.3", "", 6379, 1)
	cfg.PipeBatch = 10
	cfg.PipeProcess = 1
	GetDb("redis").(*Redis).conn = cfg.Redis
}

func TestRedisP_AppendList(t *testing.T) {
	cases := []struct {
		key   string
		value string
		exp   int
		res   int
	}{
		{"listkey1", "val1", 2, 1},
		{"listkey2", "val2", 3, 1},
		{"listkey3", "val31", 3, 1},
		{"listkey3", "val32", 3, 2},
		{"listkey3", "val32", 3, 2},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			GetDb("redisp").Init(cfg)
			GetDb("redisp").AppendList(c.key, c.value, c.exp)
			GetDb("redisp").Stop()
			v := GetDb("redis").GetList(c.key)
			if len(v) != c.res {
				t.Errorf("key:%s , want:%d , res: %s", c.key, c.res, v)
			}
		})
	}
}

func TestRedisP_Set(t *testing.T) {
	cases := []struct {
		key   string
		value string
		exp   int
		sleep int
		res   string
	}{
		{"key1", "val1", 10, 0, "val1"},
		{"key2", "val2", 3, 4, ""},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			GetDb("redisp").Init(cfg)
			GetDb("redisp").Set(c.key, c.value, c.exp)
			GetDb("redisp").Stop()
			time.Sleep(time.Duration(c.sleep) * time.Second)
			v := GetDb("redis").Get(c.key)
			if v != c.res {
				t.Errorf("key:%s , want:%s , res: %s", c.key, c.res, v)
			}
		})
	}
}

func TestRedisP_Add(t *testing.T) {
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
			GetDb("redisp").Init(cfg)
			GetDb("redisp").Add(c.key, c.value, c.exp)
			GetDb("redisp").Stop()
			vv := GetDb("redis").Get(c.key)
			v, _ := strconv.Atoi(vv)
			if v != c.res {
				t.Errorf("key:%s , want:%d , res: %d", c.key, c.res, v)
			}
		})
	}
}
