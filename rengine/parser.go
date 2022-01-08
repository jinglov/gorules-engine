package rengine

import (
	"strings"
	"sync"

	"github.com/jinglov/gorules-engine/monitor"
	"github.com/jinglov/gorules-engine/rengine/funcs"
	"github.com/omigo/log"
)

type Prule struct {
	Type    string
	Name    string
	LowRule bool
	Index   int
	HArgs   []*Prule
	Args    []*Prule
}

type parserV3 struct {
	E              *Engine
	HfuncTimes     int
	LfuncTimes     int
	CachefuncTimes int
	items          []string
}

func (p *parserV3) pushItem(item string) {
	if item == "" {
		item = funcs.False
	}
	p.items = append(p.items, item)
}

func (p *parserV3) pushResult(item string) {
	if item == "" {
		item = funcs.False
	}
	// if item != False {
	p.items = append(p.items, item)
	// }
}

func (p *parserV3) popItem() string {
	l := len(p.items)
	if l == 0 {
		return ""
	}
	item := p.items[l-1]
	p.items = p.items[:l-1]
	return item
}

func (p *parserV3) lastItem() string {
	l := len(p.items)
	if l == 0 {
		return ""
	}
	item := p.items[l-1]
	return item
}

func (p *parserV3) calcuateFuncTimes(actVec *monitor.CalcuateVec, actCallabels *monitor.ActionCalcuateLabels) {
	log.Debugf("HfuncTimes:%d, LfuncTimes:%d, CachefuncTimes:%d", p.HfuncTimes, p.LfuncTimes, p.CachefuncTimes)
	if actVec == nil || actCallabels == nil {
		return
	}

	if p.HfuncTimes > 0 {
		actVec.Add(actCallabels.HFunc, p.HfuncTimes)
	}

	if p.LfuncTimes > 0 {
		actVec.Add(actCallabels.LFunc, p.LfuncTimes)
	}

	if p.CachefuncTimes > 0 {
		actVec.Add(actCallabels.CacheFunc, p.CachefuncTimes)
	}
}

func ParserBool(e *Engine, thedate string, m map[string]string, exp *Prule, actVec *monitor.CalcuateVec, actCallabels *monitor.ActionCalcuateLabels) bool {
	p := newParser(e)
	p.parser(thedate, m, exp)
	res := p.Bool()
	p.calcuateFuncTimes(actVec, actCallabels)
	releaseParser(p)
	return res
}

func ParserString(e *Engine, theDate string, m map[string]string, exp *Prule, actVec *monitor.CalcuateVec, actCallabels *monitor.ActionCalcuateLabels) string {
	p := newParser(e)
	p.parser(theDate, m, exp)
	res := p.String()
	p.calcuateFuncTimes(actVec, actCallabels)
	releaseParser(p)
	return res
}

func ParserSlice(e *Engine, thedate string, m map[string]string, exp *Prule, actVec *monitor.CalcuateVec, actCallabels *monitor.ActionCalcuateLabels) []string {
	p := newParser(e)
	p.parser(thedate, m, exp)
	res := p.StringSlice()
	p.calcuateFuncTimes(actVec, actCallabels)
	releaseParser(p)
	return res
}

func (p *parserV3) parser(thedate string, m map[string]string, exp *Prule) {
	args := strSlicePool.Get()
	defer strSlicePool.Put(args)

	for _, arg := range exp.HArgs {
		switch arg.Type {
		case ItemTypeFun:
			p.parser(thedate, m, arg)
		case ItemTypeValue:
			p.pushItem(arg.Name)
		}
	}

	var expFun funcs.IFun
	if exp.Type == ItemTypeFun {
		expFun = p.E.funcs[exp.Name]
	}

	if exp.Type == ItemTypeFun && len(exp.HArgs) > 0 {
		for range exp.HArgs {
			args = append(args, p.popItem())
		}
		p.exec(thedate, m, expFun, args, exp.LowRule)
		res := p.lastItem()
		if fun, ok := p.E.IOrderMap[expFun.Label()]; ok {
			if res == fun.OrderRes() {
				return
			}
		}
		if len(exp.Args) > 0 {
			args = args[:0]
			args = append(args, p.popItem())
		}
	}
	// args := make([]string, 0, len(exp.Args))
	for _, arg := range exp.Args {
		switch arg.Type {
		case ItemTypeFun:
			p.parser(thedate, m, arg)
		case ItemTypeValue:
			p.pushItem(arg.Name)
		}
	}
	// log.Debugf("funname:%s , args:%v , itme:%v", exp.Name, args, p.items)
	if exp.Type == ItemTypeFun && len(exp.Args) > 0 {
		for range exp.Args {
			args = append(args, p.popItem())
		}

		p.exec(thedate, m, expFun, args, exp.LowRule)
	}
}

func (p *parserV3) exec(thedate string, m map[string]string, expFun funcs.IFun, args []string, lowRule bool) {
	if p.E == nil {
		log.Errorf("engine is nil,funName:%s, args: %v", expFun.Label(), args)
		return
	}
	funName := expFun.Label()
	if fun, ok := p.E.IArgcMap[funName]; ok {
		if _, ok := p.E.IDbMap[funName]; ok && p.E.AllowDbParam {
			if len(args) != fun.Argc() && len(args) != fun.Argc()+1 {
				log.Errorf("%s args must (%d,%d) ,not: %d, args: %v", funName, fun.Argc(), fun.Argc()+1, len(args), args)
				return
			}
		} else {
			if len(args) != fun.Argc() {
				log.Errorf("%s args must %d ,not: %d, args: %v", funName, fun.Argc(), len(args), args)
				return
			}
		}
	}

	if label, ok := p.E.vecLabels[expFun.Label()]; ok {
		if p.E.vec != nil {
			p.E.vec.Inc(label)
		}
	}

	// 参数翻转
	l := len(args)
	for i := 0; i < l/2; i++ {
		args[i], args[l-i-1] = args[l-i-1], args[i]
	}
	_, cacheOK := p.E.ICacheMap[funName]
	var cacheKey string
	if cacheOK {
		cacheKey = "_gorules_c_" + funcs.DELIMIT + funName + funcs.DELIMIT + strings.Join(args, funcs.DELIMIT)
		if v, ok := m[cacheKey]; ok {
			p.pushResult(v)
			p.CachefuncTimes++
			return
		}
	}

	if lowRule {
		p.LfuncTimes++
	} else {
		p.HfuncTimes++
	}

	if fun, ok := p.E.IDFunMap[funName]; ok {
		p.pushResult(fun.Exec(args...))
	} else if fun, ok := p.E.IMFunMap[funName]; ok {
		p.pushResult(fun.Exec(m, args...))
	} else if fun, ok := p.E.ISFunMap[funName]; ok {
		p.pushResult(fun.Exec(thedate, args...))
	}
	if cacheOK {
		m[cacheKey] = p.lastItem()
	}
}

func (p *parserV3) Bool() bool {
	if len(p.items) != 1 {
		return false
	}
	if p.popItem() == funcs.True {
		return true
	}
	return false
}

func (p *parserV3) String() string {
	if len(p.items) != 1 {
		return ""
	}
	return p.popItem()
}

func (p *parserV3) StringSlice() []string {
	if len(p.items) != 1 {
		return []string{}
	}
	item := p.popItem()
	if item == funcs.False || item == "" {
		return []string{}
	}
	return strings.Split(item, funcs.DELIMIT)
}

type parserPool struct {
	sync.Pool
}

var (
	parserp *parserPool
)

func init() {
	parserp = newParserPool()
}

func newParserPool() *parserPool {
	return &parserPool{
		sync.Pool{New: func() interface{} {
			return &parserV3{
				items: make([]string, 0, 8),
			}
		}}}
}

func (mp *parserPool) Get() *parserV3 {
	i := mp.Pool.Get().(*parserV3)
	return i
}

func (mp *parserPool) Put(m *parserV3) {
	m.items = m.items[:0]
	m.HfuncTimes = 0
	m.LfuncTimes = 0
	m.CachefuncTimes = 0
	mp.Pool.Put(m)
}

func releaseParser(m *parserV3) {
	parserp.Put(m)
}

func newParser(e *Engine) *parserV3 {
	p := parserp.Get()
	p.E = e
	return p
}

var strSlicePool *stringSlicePool

func init() {
	strSlicePool = newStringSlicePool()
}

type stringSlicePool struct {
	sync.Pool
}

func newStringSlicePool() *stringSlicePool {
	return &stringSlicePool{
		sync.Pool{
			New: func() interface{} {
				return []string{}
			},
		},
	}
}

func (f *stringSlicePool) Get() []string {
	return f.Pool.Get().([]string)[:0]
}
func (f *stringSlicePool) Put(s []string) {
	f.Pool.Put(s)
}
