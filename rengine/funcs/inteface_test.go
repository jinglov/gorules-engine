package funcs

import (
	"strconv"
	"testing"
)

func getFunByName(funcs []IFun, name string) IFun {
	for _, fun := range funcs {
		if fun.Label() == name {
			return fun
		}
	}
	return nil
}

func TestGetDayInt(t *testing.T) {
	cases := []struct {
		tm string
		d  int
	}{
		{
			"191210",
			18211,
		},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			d := GetDayInt(c.tm)
			if d != c.d {
				t.Errorf("ts: %s,want:%d , res:%d", c.tm, c.d, d)
			}
		})
	}

}
