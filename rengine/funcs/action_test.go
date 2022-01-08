package funcs

import (
	"strconv"
	"testing"
)

func TestAction(t *testing.T) {
	cases := []struct {
		funName string
		args    []string
		result  string
	}{
		{funName: "and", args: []string{True, True}, result: True},
		{funName: "or", args: []string{True, False}, result: True},
		{funName: "eq", args: []string{"a", "a"}, result: True},
		{funName: "neq", args: []string{"a", "a"}, result: False},
		{funName: "gt", args: []string{"2", "1"}, result: True},
		{funName: "gte", args: []string{"2", "2"}, result: True},
		{funName: "lt", args: []string{"2", "3"}, result: True},
		{funName: "lte", args: []string{"2", "3"}, result: True},
		{funName: "like", args: []string{"abcd", "a"}, result: True},
		{funName: "likeor", args: []string{"a,c", "a", "1"}, result: True},
		{funName: "in", args: []string{"a", "a", "b"}, result: True},
		{funName: "notin", args: []string{"a", "a", "b", "c"}, result: False},
		{funName: "splitin", args: []string{"a", "a,b,c", ","}, result: True},
		{funName: "splitin", args: []string{"a", "aa,b,c", ","}, result: False},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if getFunByName(ActionFuncs, c.funName) == nil {
				t.Errorf("funName: %s not register", c.funName)
				return
			}
			res := getFunByName(ActionFuncs, c.funName).(IDFun).Exec(c.args...)
			if c.result != res {
				t.Errorf("funName: %s(%v)  ===> want:%s ,res:%s ", c.funName, c.args, c.result, res)
			}
		})
	}
}
