package datasource

import (
	"errors"
	"github.com/jinglov/gorules-engine/channel"
	"github.com/omigo/log"
	"time"
)

func makeDataSource(name string, s *SourceCfg) (DataSource, error) {
	switch s.Type {
	case "kafka":
		return NewKafkaDatasource(name, s)
	}
	return nil, nil
}

type SourceCfg struct {
	Type              string
	ConsumerGroupName string
	Brokers           []string
	Topics            []string
	Version           string
	Process           int
	Decoder           string
}

type DataSource interface {
	Start(input *channel.Channel, mapping []*DataMapping) error
	Stop()
	Destory()
}

type DataSources struct {
	Name    string
	Source  []DataSource
	Mapping []*DataMapping
}

func NewDataSources(name string, cfgs []*SourceCfg, mapping []*DataMapping) (*DataSources, error) {
	res := &DataSources{
		Name:    name,
		Source:  make([]DataSource, 0, len(cfgs)),
		Mapping: mapping,
	}
	for _, s := range cfgs {
		kds, err := makeDataSource(name, s)
		if err != nil {
			log.Errorf("soureName: %s,err: %s", name, err)
			return res, err
		}
		if kds != nil {
			res.Source = append(res.Source, kds)
		}
	}
	return res, nil
}

func (dses *DataSources) Start(input *channel.Channel) error {
	for _, s := range dses.Source {
		err := s.Start(input, dses.Mapping)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	log.Infof("start %s datasource ok...", dses.Name)
	return nil
}

func (dses *DataSources) Stop() error {
	for _, s := range dses.Source {
		s.Stop()
	}
	log.Infof("stop %s datasource ok...", dses.Name)
	return nil
}

func (dses *DataSources) Exit() error {
	for _, s := range dses.Source {
		s.Stop()
		s.Destory()
	}
	log.Infof("exit %s datasource ok...", dses.Name)
	return nil
}

func DoMapping(timestamp int64, m map[string]string, mapping []*DataMapping, im *channel.Input) error {
	if im == nil || m == nil {
		return errors.New("channel or msg is nil")
	}
	if im.Data == nil {
		im.Data = make(map[string]string, len(mapping))
	}
	im.TimeStamp = timestamp
	// 兼容无mapping的数据源
	if len(mapping) > 0 {
		for _, f := range mapping {
			if im.Data[f.Destination] == "" && m[f.Source] != "" {
				im.Data[f.Destination] = m[f.Source]
			}
		}
	} else {
		for k, v := range m {
			if v != "" {
				im.Data[k] = v
			}
		}
	}
	im.TimeStart = time.Now()
	im.GenTheDate()
	return nil
}
