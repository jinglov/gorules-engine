package db

import (
	"github.com/jinglov/gomisc/drivers"
	"github.com/omigo/log"
	"time"
)

func init() {
	RegisterDb("redis", &Redis{})
}

type Redis struct {
	conn *drivers.Redis
}

var _ DB = &Redis{}

func (as *Redis) Init(cfg DBConf) bool {
	if cfg.Redis == nil {
		return false
	}
	as.conn = cfg.Redis
	res := as.conn.Client.Ping()
	return res.Err() == nil
}

func (as *Redis) Stop() {

}
func (as *Redis) AppendList(key, value string, exp int) int {
	if !as.check() {
		return 0
	}
	uniqKey := key + value
	if BloomTestAndAdd(uniqKey) {
		log.Debug("bloom test true")
		return 0
	}
	res := as.conn.Client.HSetNX(key, value, 1)
	if !res.Val() {
		return 0
	}
	len := as.GetListLen(key)
	if len == 1 {
		as.conn.Expire(key, time.Duration(exp)*time.Second)
	}
	return len
}
func (as *Redis) GetList(key string) []string {
	if !as.check() {
		return nil
	}
	res := as.conn.Client.HKeys(key)
	return res.Val()
}
func (as *Redis) GetListLen(key string) int {
	if !as.check() {
		return 0
	}
	res := as.conn.Client.HLen(key)
	return int(res.Val())
}

func (as *Redis) Set(key, value string, exp int) int {
	if !as.check() {
		return 0
	}
	_, err := as.conn.Set(key, value, time.Duration(exp)*time.Second)
	if err != nil {
		return 0
	}
	return 1
}

func (as *Redis) SetNx(key, value string, exp int) bool {
	if !as.check() {
		return false
	}
	uniqKey := key + value
	if BloomTestAndAdd(uniqKey) {
		return false
	}
	res, err := as.conn.Setnx(exp, key, value)
	if err != nil {
		return false
	}
	return res
}

func (as *Redis) Add(key string, value, exp int) int {
	if !as.check() {
		return 0
	}
	res := as.conn.Client.IncrBy(key, int64(value))
	if res.Val() == int64(value) {
		as.conn.Expire(key, time.Duration(exp)*time.Second)
	}
	return int(res.Val())
}
func (as *Redis) Get(key string) string {
	if !as.check() {
		return ""
	}
	res, err := as.conn.Get(key)
	if err != nil {
		return ""
	}
	return res
}

func (as *Redis) check() bool {
	if as.conn == nil {
		log.Error("a")
		return false
	}
	if as.conn.Client == nil {
		log.Error("b")
		return false
	}
	return true
}

func (as *Redis) UniqueOutput(exp int, theDate, cid, keyType, key, reason string) bool {
	len := as.AppendList("reasoncode_"+keyType+theDate+cid+key, reason, exp)
	if len > 0 {
		return true
	}
	return false
}

func (as *Redis) GetReasonList(theDate, cid, keyType, key string) []string {
	return []string{}
}
