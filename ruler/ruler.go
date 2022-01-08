package ruler

import (
	"sync"
	"time"

	"github.com/jinglov/gorules-engine/channel"
	"github.com/jinglov/gorules-engine/monitor"
	"github.com/jinglov/gorules-engine/rengine"
	"github.com/omigo/log"
)

type Rules struct {
	Name   string
	Rules  []*Rule
	engine *rengine.Engine
	wg     sync.WaitGroup
	exit   chan struct{}
}

type Rule struct {
	SendTo       string
	ReasonCode   string
	ReasonExp    string
	PReasonExp   *rengine.Prule
	Filter       string
	PFilter      *rengine.Prule
	Cid          string
	PCid         *rengine.Prule
	Data         []string
	Output       []string
	UniqueType   string
	Expire       int
	actCalLabels *monitor.ActionCalcuateLabels
}

func NewRules(dataName string, rs []*Rule, e *rengine.Engine) *Rules {
	res := &Rules{
		Name:   dataName,
		Rules:  make([]*Rule, 0, len(rs)),
		engine: e,
		wg:     sync.WaitGroup{},
		exit:   make(chan struct{}),
	}
	var err error
	for _, r := range rs {
		r.PFilter, err = rengine.NewPrule(e, []byte(r.Filter))
		if err != nil {
			log.Errorf("dataSource: %s , err: %s", dataName, err)
			continue
		}
		r.PCid, err = rengine.NewPrule(e, []byte(r.Cid))
		if err != nil {
			log.Errorf("dataSource: %s , err: %s", dataName, err)
			continue
		}
		if len(r.ReasonExp) > 0 {
			r.PReasonExp, err = rengine.NewPrule(e, []byte(r.ReasonExp))
			if err != nil {
				log.Errorf("dataSource: %s , err: %s", dataName, err)
				continue
			}
		}
		r.actCalLabels = monitor.GetActionCalcuateLabels(r.ReasonCode)
		res.Rules = append(res.Rules, r)
	}
	return res
}

func (rs *Rules) Start(input *channel.Channel, p int) {
	if p <= 0 {
		p = 1
	}
	for i := 0; i < p; i++ {
		rs.wg.Add(1)
		go rs.run(input)
	}
	log.Infof("start %s rule ok...", rs.Name)
	return
}
func (rs *Rules) Stop() {
	close(rs.exit)
	rs.wg.Wait()
	log.Infof("stop %s rule ok...", rs.Name)

}

func (rs *Rules) Exit(input *channel.Channel) {
	close(input.SChan)
	rs.wg.Wait()
	log.Infof("exit %s rule ok...", rs.Name)
}

func (rs *Rules) run(input *channel.Channel) {
	defer rs.wg.Done()
	for {
		select {
		case msg, ok := <-input.SChan:
			if !ok {
				log.Info("rule checker done")
				return
			}
			monitor.ChannelVec.Dec("saved")
			outs := rs.DoRuler(msg)
			for _, out := range outs {
				input.OChan <- out
				monitor.ChannelVec.Inc("output")
			}
			// 处理完所有逻辑，把Input放回对象池复用
			input.Pool.Put(msg)
		case <-rs.exit:
			log.Info("rule receive exit chan")
			return
		}
	}
}

// 处理保存的逻辑
func (rs *Rules) DoRuler(msg *channel.Input) (outs []*channel.Output) {
	if rs == nil {
		return
	}
	for _, r := range rs.Rules {
		// if parser.FilterBool(input.TheDate, input.Data, r.PFilter) {
		if rengine.ParserBool(rs.engine, msg.TheDate, msg.Data, r.PFilter, monitor.RulerCalVec, r.actCalLabels) {
			// 命中规则
			m := make(map[string]string, len(r.Data))
			for _, d := range r.Data {
				m[d] = msg.Data[d]
			}
			log.Debugf("hit:%s", r.ReasonCode)
			cid := rengine.ParserString(rs.engine, msg.TheDate, msg.Data, r.PCid, monitor.RulerCalVec, r.actCalLabels)

			if len(r.ReasonCode) > 0 {
				outs = append(outs, &channel.Output{
					TheDate:    msg.TheDate,
					Timestamp:  msg.TimeStamp,
					TimeStart:  msg.TimeStart,
					Cid:        cid,
					SendTo:     r.SendTo,
					ReasonCode: r.ReasonCode,
					Keys:       r.Output,
					Data:       m,
					Expire:     r.Expire,
					UniqueType: r.UniqueType,
				})
			}
			if r.PReasonExp != nil {
				reasonCodes := rengine.ParserSlice(rs.engine, msg.TheDate, msg.Data, r.PReasonExp, monitor.RulerCalVec, r.actCalLabels)
				for _, code := range reasonCodes {
					outs = append(outs, &channel.Output{
						TheDate:    msg.TheDate,
						Timestamp:  msg.TimeStamp,
						TimeStart:  msg.TimeStart,
						Cid:        cid,
						SendTo:     r.SendTo,
						ReasonCode: code,
						Keys:       r.Output,
						Data:       m,
						Expire:     r.Expire,
						UniqueType: r.UniqueType,
					})
				}
			}
			monitor.RuleCidVec.Inc(monitor.GetRuleLabels(cid))
			monitor.RuleCodeVec.Inc(monitor.GetRuleLabels(r.ReasonCode))
		}
	}
	monitor.TimeVec.Observe(monitor.GetTimeLabel(rs.Name),
		float64(time.Since(msg.TimeStart))/float64(time.Millisecond))

	monitor.DSVec.Inc(monitor.GetDSLabel(rs.Name, "checker"))
	return
}
