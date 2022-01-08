实时风控系统
========

本系统实时风控系统独立出来使用
主要工作流程：消费数据 -> 解码到map[string]string -> 存储数据 -> 规则过滤 -> 输出过滤 -> 数据输出

##消费数据

数据源实现以下接口，目前系统实现了kafka
```
type DataSource interface {
	DecodeChan() chan []byte
	DecodeProcess()
	Start() error
	Close()
	Destory()
}
```

###解码到map[string]string 目前实现了json上报
```
type Decoder interface {
	DecodeReportFromByte(m map[string]string, b []byte) (timestamp int64, err error)
}
```

## [规则过滤](doc/funcs.md)