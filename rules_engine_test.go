package gorulles_engine

import (
	"github.com/jinglov/gorules-engine/datasource"
	"github.com/jinglov/gorules-engine/db"
	"github.com/jinglov/gorules-engine/output"
	"github.com/jinglov/gorules-engine/rengine"
	"github.com/jinglov/gorules-engine/ruler"
	"github.com/jinglov/gorules-engine/saver"
	"strconv"
	"testing"
	"time"
)

var (
	saveEngine = rengine.NewFilterEngine("saver", false, true)
	execEngine = rengine.NewExecEngine("exec", true)
	ruleEngine = rengine.NewFilterEngine("rule", true, true)
	dms        []*datasource.DataMapping
	ss         []*saver.Saver
	rs         []*ruler.Rule
	oo         map[string]*output.OutputRule
	name       = "test"
)

func init() {
	dbcfg := db.DBConf{
		ListDefault:      "redis",
		RowDefault:       "redis",
		DuplicateDefault: "redis",
	}
	db.InitDB(dbcfg)
	dms = []*datasource.DataMapping{
		{Source: "k1", Destination: "dk1"},
		{Source: "k2", Destination: "dk2"},
		{Source: "k3", Destination: "dk3"},
	}
	rs = []*ruler.Rule{
		{
			SendTo:     "default",
			ReasonCode: "TEST01",
			ReasonExp:  "",
			Filter:     "and(neq(getSrc('dk1'),'') , neq(getSrc('dk2') , ''))",
			Cid:        "'test'",
			Data:       []string{"dk1", "dk2"},
			Output:     []string{"dk1", "dk2"},
			Expire:     0,
		},
		{
			SendTo:     "default",
			ReasonCode: "TEST02",
			ReasonExp:  "",
			Filter:     "and(neq(getSrc('dk1'),'') , neq(getSrc('_saver1') , ''))",
			Cid:        "'test'",
			Data:       []string{"dk1", "dk2", "_saver1"},
			Output:     []string{"dk1"},
			Expire:     0,
		},
	}
	ss = []*saver.Saver{
		{
			Filter:     "neq(getSrc('dk1') , '')",
			Exec:       "'testsave1'",
			ResDataKey: "_saver1",
		},
	}
	oo = map[string]*output.OutputRule{
		"dk1": {KeyType: "test", Exec: "getSrc('dk1')"},
		"dk2": {KeyType: "test", Exec: "getSrc('dk2')"},
	}
}

func TestGoRulesEngineSync(t *testing.T) {
	xw := NewRulesEngine("teset", saveEngine, execEngine, ruleEngine, 1024, 1024, 64)
	var err error
	xw.sources, err = datasource.NewDataSources(name, nil, dms)
	if err != nil {
		t.Error(err)
		return
	}
	xw.savers = saver.NewSavers(name, ss, saveEngine, execEngine)
	xw.rules = ruler.NewRules(name, rs, ruleEngine)
	xw.outputs = output.NewOutputs(name, oo, ruleEngine, "redis", nil)
	cases := []struct {
		timestamp int64
		m         map[string]string
		res       int
	}{
		{
			timestamp: time.Now().Unix(),
			m:         map[string]string{},
			res:       0,
		},
		{
			timestamp: time.Now().Unix(),
			m:         map[string]string{"k1": "v1"},
			res:       1,
		},
		{
			timestamp: time.Now().Unix(),
			m:         map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
			res:       2,
		},
		{
			timestamp: time.Now().Unix(),
			m:         map[string]string{"k1": "v21", "k2": "v22", "k3": "v23"},
			res:       3,
		},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			o, e := xw.SyncCalculate(c.timestamp, c.m)
			if e != nil {
				t.Error(e)
				return
			}
			if len(o) != c.res {
				t.Errorf("want len:%d , res:%d", c.res, len(o))
				t.Log(o)
			}
		})
	}
}
