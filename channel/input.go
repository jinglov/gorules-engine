package channel

import (
	"time"
)

type Input struct {
	Data      map[string]string
	TimeStamp int64
	TimeStart time.Time
	TheDate   string
}

var (
	MinTimestampS = int64(time.Now().Second())
	MinTimestampM = int64(1_000_000_000_000)
)

func (ipt *Input) GenTheDate() {
	// 如果时间戳是毫秒
	if ipt.TimeStamp > MinTimestampM {
		ipt.TheDate = time.Unix(ipt.TimeStamp/1000, 0).Format("060102")
		return
	}
	// 如果时间戳比较早，则用当前时间
	if ipt.TimeStamp <= MinTimestampS {
		ipt.TheDate = time.Now().Format("060102")
		ipt.TimeStamp = time.Now().Unix() * 1000
		return
	}
	// 其它时间戳
	ipt.TheDate = time.Unix(ipt.TimeStamp, 0).Format("060102")
	ipt.TimeStamp *= 1000
}

func (ipt *Input) reset() *Input {
	ipt.TimeStamp = 0
	ipt.TheDate = ""
	for k := range ipt.Data {
		delete(ipt.Data, k)
	}
	return ipt
}
