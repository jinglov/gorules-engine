package funcs

import (
	"strconv"
	"testing"
)

func TestCovert(t *testing.T) {
	cases := []struct {
		funName string
		args    []string
		result  string
	}{
		{funName: "sum", args: []string{"1", "2", "3"}, result: "6"},
		{funName: "subtract", args: []string{"5", "2"}, result: "3"},
		{funName: "div", args: []string{"6", "3"}, result: "2"},
		{funName: "multiply", args: []string{"2", "3", "4"}, result: "24"},
		{funName: "concat", args: []string{"a", "b", "cd"}, result: "abcd"},
		{funName: "left", args: []string{"abc", "4"}, result: "abc"},
		{funName: "left", args: []string{"abc", "2"}, result: "ab"},
		{funName: "left", args: []string{"abc", "-1"}, result: "ab"},
		{funName: "left", args: []string{"abc", "-4"}, result: ""},
		{funName: "right", args: []string{"abc", "2"}, result: "bc"},
		{funName: "right", args: []string{"abc", "-1"}, result: "bc"},
		{funName: "right", args: []string{"abc", "-4"}, result: ""},
		{funName: "split", args: []string{"a,c", ",", "0"}, result: "a"},
		{funName: "split", args: []string{"a,c", ",", "1"}, result: "c"},
		{funName: "split", args: []string{"a,c", ",", "2"}, result: ""},
		{funName: "split", args: []string{"a,c", ",", "-1"}, result: "c"},
		{funName: "hash", args: []string{"111"}, result: "fd8f478a42a33d551f5db95b90d4da89"},
		{funName: "gettuid", args: []string{"ACDirMPfJvxOmqXJDsAat1Sfdv9mpiYAN-NxdHQ"}, result: "4qzD3yb8TpqlyQ7AGrdUnw"},
		{funName: "gettuid", args: []string{"xBSLpCf-z0Iymmqx7D5H0inR"}, result: ""},
		{funName: "todate", args: []string{""}, result: "1970-01-01"},
		{funName: "todate", args: []string{"1572581100"}, result: "2019-11-01"},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if getFunByName(CovertFuncs, c.funName) == nil {
				t.Errorf("funName%s not register", c.funName)
			}
			res := getFunByName(CovertFuncs, c.funName).(IDFun).Exec(c.args...)
			if c.result != res {
				t.Errorf("funName:%s(%v)  ===> want:%s ,res:%s ", c.funName, c.args, c.result, res)
			}
		})
	}
}
