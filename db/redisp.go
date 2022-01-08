package db

import (
	"github.com/go-redis/redis"
	"github.com/jinglov/gomisc/drivers"
	"github.com/omigo/log"
	"sync"
)

func init() {
	RegisterDb("redisp", &RedisP{})
}

type RedisP struct {
	conn     *drivers.Redis
	pch      chan redis.Cmder
	batchCmd int
	wg       sync.WaitGroup
}

var _ DB = &RedisP{}

func (as *RedisP) Init(cfg DBConf) bool {
	if cfg.Redis == nil || cfg.PipeProcess <= 0 || cfg.PipeBatch <= 0 {
		return false
	}
	as.batchCmd = cfg.PipeBatch
	cmd := cfg.Redis.Client.Ping()
	if cmd.Err() != nil {
		return false
	}
	as.conn = cfg.Redis
	as.pch = make(chan redis.Cmder, cfg.PipeBatch*cfg.PipeProcess)
	for i := 0; i < cfg.PipeProcess; i++ {
		as.wg.Add(1)
		go as.process(&as.wg)
	}
	return true
}

func (as *RedisP) process(wg *sync.WaitGroup) {
	defer wg.Done()
	if !as.check() {
		log.Error("start redis process check error")
		return
	}
	pconn := as.conn.Client.Pipeline()
	index := 0
	log.Info("start redis pipe process...")
	for {
		select {
		case cmd, ok := <-as.pch:
			if !ok {
				pconn.Exec()
				log.Info("redis pipe exit...")
				pconn.Close()
				return
			}
			log.Debug(cmd)
			pconn.Process(cmd)
			index++
			if index%as.batchCmd == 0 {
				pconn.Exec()
				index = 0
			}
		}
	}
}

func (as *RedisP) Stop() {
	if !as.check() {
		return
	}
	close(as.pch)
	as.wg.Wait()
}

func (as *RedisP) expire(key string, exp int) {
	as.pch <- redis.NewBoolCmd("expire", key, exp)
}

func (as *RedisP) AppendList(key, value string, exp int) int {
	if !as.check() {
		return 0
	}
	as.pch <- redis.NewIntCmd("hset", key, value, 1)
	as.expire(key, exp)
	return 0
}
func (as *RedisP) GetList(key string) []string {
	return nil
}
func (as *RedisP) GetListLen(key string) int {
	return 0
}

func (as *RedisP) Set(key, value string, exp int) int {
	if !as.check() {
		return 0
	}
	as.pch <- redis.NewStatusCmd("set", key, value, "ex", exp)
	return 0
}

func (as *RedisP) SetNx(key, value string, exp int) bool {
	if !as.check() {
		return false
	}
	as.pch <- redis.NewStatusCmd("setnx", key, value, "ex", exp)
	return false
}
func (as *RedisP) Add(key string, value, exp int) int {
	if !as.check() {
		return 0
	}
	as.pch <- redis.NewIntCmd("incrby", key, value)
	as.expire(key, exp)
	return 0
}
func (as *RedisP) Get(key string) string {
	return ""
}

func (as *RedisP) check() bool {
	if as.conn == nil {
		return false
	}
	if as.conn.Client == nil {
		return false
	}
	return true
}

func (as *RedisP) UniqueOutput(exp int, theDate, cid, keyType, key, reason string) bool {
	return as.SetNx("data_"+keyType+theDate+cid+key+reason, "", exp)
}

func (as *RedisP) GetReasonList(thedate, cid, keyType, key string) []string {
	return []string{}
}
