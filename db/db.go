package db

import (
	"github.com/jinglov/gomisc/drivers"
	"github.com/omigo/log"
)

var (
	//
	// RowDefault     = "redis"
	// ListDefault    = "aerospike"
	// DuplicateStore = "aerospike"
	RowDefault, ListDefault, DuplicateStore DB
)

type DBConf struct {
	ListDefault      string
	RowDefault       string
	DuplicateDefault string
	BloomSize        uint
	PipeBatch        int
	PipeProcess      int
	Aerospike        *drivers.AsStore
	Redis            *drivers.Redis
}

type DB interface {
	Init(cfg DBConf) bool
	Stop()

	AppendList(key, value string, exp int) int
	GetList(key string) []string
	GetListLen(key string) int

	Set(key, value string, exp int) int
	SetNx(key, value string, exp int) bool
	Add(key string, value, exp int) int
	Get(key string) string

	UniqueOutput(exp int, theDate, cid, keyType, key, reason string) bool
	GetReasonList(theDate, cid, keyType, key string) []string
}

var dbs = make(map[string]DB, 2)

func GetDb(name string) DB {
	if db, ok := dbs[name]; ok {
		return db
	}
	log.Error("db " + name + " not register.")
	return nil
}

func RegisterDb(name string, db DB) {
	if _, ok := dbs[name]; ok {
		log.Errorf("db:% is duplicate register.", name)
	}
	dbs[name] = db
}

func InitDB(cfgs DBConf) {
	for name, db := range dbs {
		if ok := db.Init(cfgs); ok {
			if name == cfgs.RowDefault {
				RowDefault = GetDb(name)
			}
			if name == cfgs.ListDefault {
				ListDefault = GetDb(name)
			}
			if name == cfgs.DuplicateDefault {
				DuplicateStore = GetDb(name)
			}
		}
	}
	if RowDefault == nil {
		RowDefault = GetDb("redis")
	}
	if ListDefault == nil {
		ListDefault = GetDb("aerospike")
	}
	if DuplicateStore == nil {
		DuplicateStore = GetDb("aerospike")
	}
	if RowDefault == nil {
		RowDefault = GetDb("aerospike")
	}
	if ListDefault == nil {
		ListDefault = GetDb("redis")
	}
	if DuplicateStore == nil {
		DuplicateStore = GetDb("redis")
	}
}

func Stop() {
	for _, db := range dbs {
		db.Stop()
	}
}
