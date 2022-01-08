package rengine

import (
	"github.com/jinglov/gorules-engine/monitor"
	"github.com/jinglov/gorules-engine/rengine/funcs"
	"github.com/omigo/log"
	"reflect"
	"sync"
)

type Engine struct {
	Name         string
	AllowDbParam bool
	funcs        map[string]funcs.IFun
	vec          *monitor.FuncVec
	vecLabels    map[string]*monitor.FuncLabel
	mu           sync.RWMutex
	storeCache   bool
	IOrderMap    map[string]funcs.IOrder
	IDFunMap     map[string]funcs.IDFun
	IMFunMap     map[string]funcs.IMFun
	ISFunMap     map[string]funcs.ISFun
	IArgcMap     map[string]funcs.IArgc
	IDbMap       map[string]funcs.IDb
	ICacheMap    map[string]funcs.ICache
}

func (e *Engine) Register(exec funcs.IFun) {
	e.mu.Lock()
	defer e.mu.Unlock()
	label := exec.Label()
	if _, ok := e.funcs[label]; ok {
		log.Panicf("func duplic register:%s", label)
	}
	e.funcs[label] = exec
	e.vecLabels[label] = &monitor.FuncLabel{Name: label}
	if reflect.TypeOf(exec).Implements(funcs.OrderType) {
		e.IOrderMap[label] = exec.(funcs.IOrder)
	}
	if reflect.TypeOf(exec).Implements(funcs.DFunType) {
		e.IDFunMap[label] = exec.(funcs.IDFun)
	}
	if reflect.TypeOf(exec).Implements(funcs.MFunType) {
		e.IMFunMap[label] = exec.(funcs.IMFun)
	}
	if reflect.TypeOf(exec).Implements(funcs.SFunType) {
		e.ISFunMap[label] = exec.(funcs.ISFun)
	}
	if reflect.TypeOf(exec).Implements(funcs.ArgcType) {
		e.IArgcMap[label] = exec.(funcs.IArgc)
	}
	if reflect.TypeOf(exec).Implements(funcs.DbType) {
		e.IDbMap[label] = exec.(funcs.IDb)
	}
	if e.storeCache && reflect.TypeOf(exec).Implements(funcs.CacheType) {
		e.ICacheMap[label] = exec.(funcs.ICache)
	}
}

func (e *Engine) MuliteRegister(addFuncs []funcs.IFun) {
	for _, fun := range addFuncs {
		e.Register(fun)
	}
}

func NewEngine(name string) *Engine {
	e := &Engine{
		Name:      name,
		funcs:     make(map[string]funcs.IFun),
		vec:       monitor.NewFuncVec("go_rules", "run", name),
		vecLabels: make(map[string]*monitor.FuncLabel),
		IOrderMap: make(map[string]funcs.IOrder),
		IDFunMap:  make(map[string]funcs.IDFun),
		IMFunMap:  make(map[string]funcs.IMFun),
		ISFunMap:  make(map[string]funcs.ISFun),
		IArgcMap:  make(map[string]funcs.IArgc),
		IDbMap:    make(map[string]funcs.IDb),
		ICacheMap: make(map[string]funcs.ICache),
	}
	return e
}

func NewFilterEngine(name string, storeCache bool, allowDbParam bool) *Engine {
	e := NewEngine(name)
	e.storeCache = storeCache
	e.AllowDbParam = allowDbParam
	e.MuliteRegister(funcs.ActionFuncs)
	e.MuliteRegister(funcs.CovertFuncs)
	e.MuliteRegister(funcs.DataFuncs)
	e.MuliteRegister(funcs.StoreGetFuncs)
	return e
}

func NewExecEngine(name string, allowDbParam bool) *Engine {
	e := NewEngine(name)
	e.AllowDbParam = allowDbParam
	e.MuliteRegister(funcs.ActionFuncs)
	e.MuliteRegister(funcs.CovertFuncs)
	e.MuliteRegister(funcs.DataFuncs)
	e.MuliteRegister(funcs.StoreGetFuncs)
	e.MuliteRegister(funcs.StoreSetFuncs)
	return e
}
