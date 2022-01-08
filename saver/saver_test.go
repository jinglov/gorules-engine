package saver

import (
	"github.com/jinglov/gorules-engine/channel"
	"testing"
)

func TestStartStop(t *testing.T) {
	s := NewSavers("test", make([]*Saver, 0), nil, nil)
	input := channel.NewChannel("test", 10, 10, 10)
	s.Start(input, 2)

	s.Stop()
}
