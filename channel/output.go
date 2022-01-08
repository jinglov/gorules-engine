package channel

import "time"

type Output struct {
	TheDate    string
	Timestamp  int64
	TimeStart  time.Time
	Cid        string
	Data       map[string]string
	Keys       []string
	ReasonCode string
	SendTo     string

	Expire     int
	UniqueType string
}
