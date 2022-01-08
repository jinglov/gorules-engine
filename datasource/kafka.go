package datasource

import (
	"github.com/Shopify/sarama"
	"github.com/jinglov/gomisc/kafka"
	"github.com/jinglov/gomisc/monitor"
	"github.com/jinglov/gomisc/pool"
	"github.com/jinglov/gorules-engine/channel"
	"github.com/jinglov/gorules-engine/decode"
	monitor2 "github.com/jinglov/gorules-engine/monitor"
	"github.com/omigo/log"
	"sync"
)

var (
	_        DataSource = &kafkaDatasource{}
	kafkaVec            = monitor.NewKafkaVec("go_rules", "engine", "kafka")
)

type kafkaDatasource struct {
	Ki         kafka.Consumer
	Name       string
	DataName   string
	Type       string
	decodeChan chan []byte
	process    int
	wg         *sync.WaitGroup
	decoder    string
	decodePool pool.StringMapPool
}

func NewKafkaDatasource(sourceName string, s *SourceCfg) (*kafkaDatasource, error) {
	kds := &kafkaDatasource{
		Type:       s.Type,
		Name:       sourceName,
		process:    s.Process,
		wg:         new(sync.WaitGroup),
		decodeChan: make(chan []byte, 65535),
		decoder:    s.Decoder,
		decodePool: pool.NewStringMapPool(),
	}
	ki, err := kafka.NewConsumerV2(s.ConsumerGroupName, s.Version, s.Topics, s.Brokers, kds.receiveMsg, kafka.WithVec(kafkaVec))
	if err != nil {
		return nil, err
	}
	kds.Ki = ki
	return kds, nil
}

func (ds *kafkaDatasource) Start(input *channel.Channel, mapping []*DataMapping) error {
	for i := 0; i < ds.process; i++ {
		ds.wg.Add(1)
		go ds.decodeProcess(input, mapping)
	}
	err := ds.Ki.Start()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (ds *kafkaDatasource) Stop() {
	err := ds.Ki.Close()
	if err != nil {
		log.Error(err)
	}
	close(ds.decodeChan)
	ds.wg.Wait()
}

func (ds *kafkaDatasource) Destory() {

}

func (ds *kafkaDatasource) receiveMsg(msg *sarama.ConsumerMessage) (err error) {
	log.Debugf("msg lenght: %d", len(msg.Value))
	if len(msg.Value) == 0 {
		return
	}
	ds.decodeChan <- msg.Value
	monitor2.ChannelVec.Inc("decode")
	return
}

func (ds *kafkaDatasource) decodeProcess(input *channel.Channel, mapping []*DataMapping) {
	defer ds.wg.Done()
	dc := decode.GetDecoder(ds.decoder)
	for {
		select {
		case msg, ok := <-ds.decodeChan:
			if !ok {
				log.Info("decode process done")
				return
			}
			monitor2.ChannelVec.Dec("decode")
			im := input.Pool.Get()
			m := ds.decodePool.Get()
			timeStamp, err := dc.DecodeReportFromByte(m, msg)
			if err != nil {
				log.Error(err)
				ds.decodePool.Put(m)
				input.Pool.Put(im)
				continue
			}
			DoMapping(timeStamp, m, mapping, im)
			ds.decodePool.Put(m)
			monitor2.DSVec.Inc(monitor2.GetDSLabel(input.Name, "input"))
			input.IChan <- im
			monitor2.ChannelVec.Inc("input")
		}
	}
}
