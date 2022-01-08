package ruler

import (
	"github.com/jinglov/gorules-engine/channel"
	"testing"
)

func TestStartTop(t *testing.T) {
	r := NewRules("test", make([]*Rule, 0), nil)

	input := channel.NewChannel("test", 10, 10, 10)
	r.Start(input, 3)

	r.Stop()
}
