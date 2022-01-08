package saver

import (
	"github.com/jinglov/gorules-engine/channel"
	"github.com/jinglov/gorules-engine/monitor"
	"github.com/jinglov/gorules-engine/rengine"
	"github.com/jinglov/gorules-engine/rengine/funcs"
	"github.com/omigo/log"
	"sync"
)

type Savers struct {
	Name         string
	Savers       []*Saver
	filterEngine *rengine.Engine
	execEngine   *rengine.Engine
	wg           sync.WaitGroup
	exit         chan struct{}
}

type Saver struct {
	Filter       string
	PFilter      *rengine.Prule
	Exec         string
	PExec        *rengine.Prule
	ResDataKey   string
	StartLabel   *monitor.SaverLabels
	TrueLabel    *monitor.SaverLabels
	FalseLabel   *monitor.SaverLabels
	actCalLabels *monitor.ActionCalcuateLabels
}

func NewSavers(dataName string, ss []*Saver, fe *rengine.Engine, ee *rengine.Engine) *Savers {
	res := &Savers{
		Name:         dataName,
		Savers:       make([]*Saver, 0, len(ss)),
		filterEngine: fe,
		execEngine:   ee,
		wg:           sync.WaitGroup{},
		exit:         make(chan struct{}),
	}
	var err error
	for _, r := range ss {
		r.PFilter, err = rengine.NewPrule(fe, []byte(r.Filter))
		if err != nil {
			log.Errorf("dataSource: %s , err: %s", dataName, err)
			continue
		}
		r.PExec, err = rengine.NewPrule(ee, []byte(r.Exec))
		if err != nil {
			log.Errorf("dataSource: %s , err: %s", dataName, err)
			continue
		}
		r.StartLabel = &monitor.SaverLabels{DataSource: dataName, ResDataKey: r.ResDataKey, Status: "start"}
		r.TrueLabel = &monitor.SaverLabels{DataSource: dataName, ResDataKey: r.ResDataKey, Status: funcs.True}
		r.FalseLabel = &monitor.SaverLabels{DataSource: dataName, ResDataKey: r.ResDataKey, Status: funcs.False}
		r.actCalLabels = monitor.GetActionCalcuateLabels(r.ResDataKey)
		res.Savers = append(res.Savers, r)
	}
	return res
}

func (ss *Savers) Start(input *channel.Channel, p int) {
	if p <= 0 {
		p = 1
	}
	for i := 0; i < p; i++ {
		ss.wg.Add(1)
		go ss.run(input)
	}
	log.Infof("start %s saver ok...", ss.Name)
	return
}

func (ss *Savers) Stop() {
	close(ss.exit)
	ss.wg.Wait()
	log.Infof("stop %s saver ok...", ss.Name)

}

func (ss *Savers) Exit(input *channel.Channel) {
	close(input.IChan)
	ss.wg.Wait()
	log.Infof("exit %s saver ok...", ss.Name)
}

func (ss *Savers) run(input *channel.Channel) {
	defer ss.wg.Done()
	for {
		select {
		case msg, ok := <-input.IChan:
			if !ok {
				log.Info("kafka 'run' process exit")
				return
			}
			monitor.ChannelVec.Dec("input")
			ss.DoSaver(msg)
			monitor.ChannelVec.Inc("saved")
			input.SChan <- msg
			monitor.DSVec.Inc(monitor.GetDSLabel(ss.Name, "saver"))
		case <-ss.exit:
			log.Info("saver receive exit chan")
			return
		}
	}
}

// 处理保存的逻辑
func (ss *Savers) DoSaver(msg *channel.Input) {
	if ss == nil {
		return
	}
	// var sok string
	for _, s := range ss.Savers {

		monitor.SaverVec.Inc(s.StartLabel)
		if rengine.ParserBool(ss.filterEngine, msg.TheDate, msg.Data, s.PFilter, monitor.SaverCalVec, s.actCalLabels) {
			res := rengine.ParserString(ss.execEngine, msg.TheDate, msg.Data, s.PExec, monitor.SaverCalVec, s.actCalLabels)
			if res != "" && res != funcs.False && s.ResDataKey != "" {
				msg.Data[s.ResDataKey] = res
				monitor.SaverVec.Inc(s.TrueLabel)
				// sok += fmt.Sprintf("k:%s,v:%s |", s.ResDataKey, res)
			} else {
				monitor.SaverVec.Inc(s.FalseLabel)
			}
		}
	}
	// if len(sok) > 0 {
	// 	log.Debug(sok)
	// 	log.Debug(input.Data)
	// }
}
