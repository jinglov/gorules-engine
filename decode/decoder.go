package decode

import (
	"log"
)

type Decoder interface {
	DecodeReportFromByte(m map[string]string, b []byte) (timestamp int64, err error)
}

var decoders = make(map[string]Decoder, 1)
var DefaultDecoder Decoder

func RegisterDecoder(name string, d Decoder) {
	if _, ok := decoders[name]; ok {
		log.Fatalln(name + " decoder duplic register")
	}
	decoders[name] = d
}

func GetDecoder(name string) Decoder {
	if d, ok := decoders[name]; ok {
		return d
	}
	return DefaultDecoder
}

func SetDefaultDecoder(d Decoder) {
	DefaultDecoder = d
}
