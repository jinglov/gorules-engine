package output

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/jinglov/gorules-engine/channel"
	"github.com/jinglov/gorules-engine/db"
	"github.com/jinglov/gorules-engine/monitor"
	"github.com/jinglov/gorules-engine/rengine"
	"github.com/omigo/log"
)

var (
	GlobalTagPrefixKey        = "globalcid-"
	GlobalTagMonitorPrefixKey = "globalcid-"
)

type Outputs struct {
	name           string
	Rules          map[string]*OutputRule
	PushConfig     *PushConfig
	Engine         *rengine.Engine
	wg             sync.WaitGroup
	DuplicateStore string
	exit           chan struct{}
}

type PushConfig struct {
	PushMapping map[string]int
	PushGlobal  map[string]int
}

type OutputRule struct {
	KeyType string
	Exec    string
	PExec   *rengine.Prule
}

type OutputData struct {
	DataSource     string            `json:"data_source"`
	Cid            string            `json:"cid"`
	ReasonCode     string            `json:"reason_code"`
	TagCode        int               `json:"tag_code"`
	IsGlobal       bool              `json:"is_global"`
	Key            string            `json:"key"`
	KeyType        string            `json:"key_type"`
	Timestamp      int64             `json:"timestamp"`
	DataTimestamp  int64             `json:"data_timestamp"`
	StartTimestamp int64             `json:"start_timestamp"`
	EndTimestamp   int64             `json:"end_timestamp"`
	Ext            map[string]string `json:"ext"`
}

func (od *OutputData) toByte() ([]byte, error) {
	return json.Marshal(od)
}

func NewOutputs(name string, rules map[string]*OutputRule, e *rengine.Engine, dstore string, pushConfig *PushConfig) *Outputs {
	res := &Outputs{
		name:       name,
		Rules:      make(map[string]*OutputRule, len(rules)),
		PushConfig: pushConfig,
		Engine:     e,
		wg:         sync.WaitGroup{},
		exit:       make(chan struct{}),
	}
	if dstore != "" && db.GetDb(dstore) != nil {
		res.DuplicateStore = dstore
	}
	var err error
	for k, r := range rules {
		r.PExec, err = rengine.NewPrule(e, []byte(r.Exec))
		if err != nil {
			log.Error(err)
			continue
		}
		res.Rules[k] = r
	}
	return res
}

func (os *Outputs) Start(input *channel.Channel, p int) {
	if p <= 0 {
		p = 1
	}
	for i := 0; i < p; i++ {
		os.wg.Add(1)
		go os.run(input)
	}
	log.Infof("start %s output ok...", os.name)
	return
}

func (os *Outputs) Stop() {
	close(os.exit)
	os.wg.Wait()
	log.Infof("stop %s output ok...", os.name)
}

func (os *Outputs) Exit(input *channel.Channel) {
	close(input.OChan)
	os.wg.Wait()
	log.Infof("exit %s output ok...", os.name)
}

func (os *Outputs) run(input *channel.Channel) {
	defer os.wg.Done()
	for {
		select {
		case msg, ok := <-input.OChan:
			monitor.ChannelVec.Dec("output")
			if !ok {
				log.Info("output process done")
				return
			}
			outMsgs := os.ExtractOut(msg)
			for _, m := range outMsgs {
				b, e := m.toByte()
				if e != nil {
					log.Error(e)
					continue
				}
				log.Debugf("sendData: %s", string(b))
				// sink.SendByte(msg.SendTo, b)
			}
		case <-os.exit:
			log.Info("output receive exit chan")
			return
		}
	}
}

// 处理输出消息，是否扩展成多个，扩展之后把消息发到到sink中
func (os *Outputs) ExtractOut(msg *channel.Output) (outMsgs []*OutputData) {
	if os == nil {
		return
	}
	calLabels := monitor.GetActionCalcuateLabels(msg.ReasonCode)
	for _, k := range msg.Keys {
		if v, ok := os.Rules[k]; ok {
			// keys := parser.OutputSlice(msg.TheDate, msg.Data, v.PExec)
			keys := rengine.ParserSlice(os.Engine, msg.TheDate, msg.Data, v.PExec, monitor.OutputCodeCalVec, calLabels)
			// log.Debug(msg.Data)
			// log.Debug(keys)
			for _, key := range keys {
				if key == "" {
					continue
				}
				if !os.checkReasonCode(msg.TheDate, msg.Cid, v.KeyType, key, msg.ReasonCode, msg.UniqueType, msg.Expire) {
					continue
				}
				tagCode := os.PushConfig.PushMapping[msg.Cid+"-"+msg.ReasonCode]
				var isGlobal bool
				if tagCode > 0 && os.PushConfig.PushGlobal[GlobalTagPrefixKey+v.KeyType+"-"+msg.ReasonCode] == tagCode {
					isGlobal = true
				}
				outMsgs = append(outMsgs, &OutputData{
					DataSource:     os.name,
					Cid:            msg.Cid,
					ReasonCode:     msg.ReasonCode,
					TagCode:        tagCode,
					IsGlobal:       isGlobal,
					Key:            key,
					KeyType:        v.KeyType,
					Timestamp:      time.Now().UnixNano() / int64(time.Millisecond),
					DataTimestamp:  msg.Timestamp,
					StartTimestamp: msg.TimeStart.UnixNano() / int64(time.Millisecond),
					EndTimestamp:   time.Now().UnixNano() / int64(time.Millisecond),
					Ext:            msg.Data,
				})

				monitor.OutputCidVec.Inc(monitor.GetOutputRuleLabel(msg.Cid))
				monitor.OutputCodeVec.Inc(monitor.GetOutputRuleLabel(msg.ReasonCode))

				monitor.OutputTagVec.Inc(monitor.GetOutputRuleLabel(strconv.Itoa(tagCode)))
				if isGlobal {
					monitor.OutputTagVec.Inc(monitor.GetOutputRuleLabel(GlobalTagMonitorPrefixKey + strconv.Itoa(tagCode)))
				}

				monitor.OutputKeyVec.Inc(monitor.GetOutputRuleLabel(v.KeyType))
				monitor.HitVec.Observe(monitor.GetHitLabel(os.name, msg.ReasonCode),
					float64(time.Since(msg.TimeStart))/float64(time.Millisecond))
			}
		}
	}
	return
}

// 其中aerospike不能完美支持指定过期时间，redis不能支持根据key找到作弊类型
func (os *Outputs) checkReasonCode(theDate, cid, keyType, key, reason, uniqueType string, exp int) bool {
	// 不去重
	if exp == -1 {
		return true
	}
	if exp == 0 {
		exp = 86400
	}
	var udb db.DB
	if uniqueType != "" {
		udb = db.GetDb(uniqueType)
	} else {
		udb = db.GetDb(os.DuplicateStore)
	}
	if udb == nil {
		return true
	}
	return udb.UniqueOutput(exp, theDate, cid, keyType, key, reason)
}
