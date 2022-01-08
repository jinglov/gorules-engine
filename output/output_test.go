package output

import (
	"strconv"
	"testing"
	"time"

	"github.com/jinglov/gorules-engine/channel"
)

var pushMapping = map[string]int{
	"tcid-T001": 1,
	"tcid-T002": 2,
}

func TestExtractOut(t *testing.T) {
	cases := []struct {
		output *channel.Output
	}{
		{
			output: &channel.Output{
				Cid:        "tcid",
				Timestamp:  int64(time.Now().Second()) * 1000,
				Keys:       []string{"fireeyes_device", "fireeyes_member_id"},
				ReasonCode: "T001",
				SendTo:     "default",
				Data: map[string]string{
					"cid":            "tcid",
					"openid":         "tk2",
					"_isr":           "1",
					"_db_device_cnt": "3",
				},
			},
		},
	}

	oo := NewOutputs("test", nil, nil, "", pushMapping)
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			oo.ExtractOut(c.output)
			t.Log(c.output.Data)
		})
	}
}

func TestStartTop(t *testing.T) {
	o := NewOutputs("test", nil, nil, "", pushMapping)
	input := channel.NewChannel("test", 10, 10, 10)
	o.Start(input, 2)
	o.Stop()
}
