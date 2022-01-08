package gorulles_engine

import (
	"errors"

	"github.com/jinglov/gorules-engine/channel"
	"github.com/jinglov/gorules-engine/datasource"
	"github.com/jinglov/gorules-engine/monitor"
	"github.com/jinglov/gorules-engine/output"
	"github.com/jinglov/gorules-engine/rengine"
	"github.com/jinglov/gorules-engine/ruler"
	"github.com/jinglov/gorules-engine/saver"
)

type RulesEngine struct {
	name         string
	channel      *channel.Channel
	filterEngine *rengine.Engine
	execEngine   *rengine.Engine
	ruleEngine   *rengine.Engine
	sources      *datasource.DataSources
	savers       *saver.Savers
	rules        *ruler.Rules
	outputs      *output.Outputs
}

func (re *RulesEngine) Exit() {
	if re.sources != nil {
		re.sources.Exit()
	}
	if re.savers != nil {
		re.savers.Exit(re.channel)
	}
	if re.rules != nil {
		re.rules.Exit(re.channel)
	}
	if re.outputs != nil {
		re.outputs.Exit(re.channel)
	}
}

func NewRulesEngine(name string, saverFilterEngine *rengine.Engine, saverExecEngine *rengine.Engine, ruleEngine *rengine.Engine, ilen, slen, olen int) *RulesEngine {
	return &RulesEngine{
		name:         name,
		channel:      channel.NewChannel(name, ilen, slen, olen),
		filterEngine: saverFilterEngine,
		execEngine:   saverExecEngine,
		ruleEngine:   ruleEngine,
	}
}

func (re *RulesEngine) StartSource(dsCfg []*datasource.SourceCfg, mapping []*datasource.DataMapping) (err error) {
	if re.sources != nil {
		re.sources.Exit()
	}
	re.sources, err = datasource.NewDataSources(re.name, dsCfg, mapping)
	if err != nil {
		return err
	}
	return re.sources.Start(re.channel)
}

func (re *RulesEngine) StartSaver(ss []*saver.Saver, pcnt int) int {
	if re.savers != nil {
		re.savers.Stop()
	}
	re.savers = saver.NewSavers(re.name, ss, re.filterEngine, re.execEngine)
	re.savers.Start(re.channel, pcnt)
	return len(re.savers.Savers)
}

func (re *RulesEngine) StartRuler(rs []*ruler.Rule, pcnt int) int {
	if re.rules != nil {
		re.rules.Stop()
	}
	re.rules = ruler.NewRules(re.name, rs, re.ruleEngine)

	re.rules.Start(re.channel, pcnt)

	return len(re.rules.Rules)
}

func (re *RulesEngine) StartOutput(rules map[string]*output.OutputRule, dstore string, pcnt int, pushConfig *output.PushConfig) int {
	if re.outputs != nil {
		re.outputs.Stop()
	}
	re.outputs = output.NewOutputs(re.name, rules, re.ruleEngine, dstore, pushConfig)
	re.outputs.Start(re.channel, pcnt)
	return len(re.outputs.Rules)
}

// SyncCalculate 同步计算
// @Params timestamp 数据时间戳
// @Params m 数据内容
// @Return 命中结果和错误
func (re *RulesEngine) SyncCalculate(timestamp int64, m map[string]string) ([]*output.OutputData, error) {
	if re == nil || re.sources == nil || re.sources.Mapping == nil {
		return nil, errors.New("source or mapping is nil")
	}
	msg := re.channel.Pool.Get()
	err := datasource.DoMapping(timestamp, m, re.sources.Mapping, msg)
	if err != nil {
		return nil, err
	}
	// log.Json(msg)
	re.savers.DoSaver(msg)
	// log.Json(msg)
	outs := re.rules.DoRuler(msg)
	// log.Json(msg)
	re.channel.Pool.Put(msg)

	res := make([]*output.OutputData, 0)
	for _, o := range outs {
		extOut := re.outputs.ExtractOut(o)
		res = append(res, extOut...)
	}
	// log.JSON(res)
	return res, nil
}

// SyncCalculate 异步计算
// @Params timestamp 数据时间戳
// @Params m 数据内容
func (re *RulesEngine) ASyncCalculate(timestamp int64, m map[string]string) error {
	if re == nil || re.sources == nil || re.sources.Mapping == nil {
		return errors.New("source or mapping is nil")
	}
	msg := re.channel.Pool.Get()
	err := datasource.DoMapping(timestamp, m, re.sources.Mapping, msg)
	if err != nil {
		return err
	}
	monitor.DSVec.Inc(monitor.GetDSLabel(re.channel.Name, "input"))
	re.channel.IChan <- msg
	monitor.ChannelVec.Inc("input")
	return nil
}
