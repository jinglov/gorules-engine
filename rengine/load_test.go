package rengine

import (
	"fmt"
	"strconv"
	"testing"
)

func TestInit(t *testing.T) {
	e := NewExecEngine("test", false)
	rules := `and(
        eq(getSrc('isnew') , '1') ,
        neq(getSrc('openid') , '') ,
        eq(getSrc('source') , 'android') ,
        neq(getSrc('cbinlm') , '') ,
        neq(getSrc('crootmod') , '') ,
        aa(getSrc('cmf') , '')
    )
`
	pr, err := NewPrule(e, []byte(rules))
	if err != nil {
		t.Error(err)
	} else {
		printp(0, pr)
	}
}
func TestNewPrule(t *testing.T) {
	e := NewExecEngine("test", false)
	cases := []struct {
		rule     string
		hasError bool
	}{
		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf')   '')`, hasError: true},
		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf'),'')`, hasError: false},
		{rule: `and(
					eq(left(getSrc('timestamp'), '10'), todate(dtToUnix(''))),
					gte(getSrc('level') , '6'),
					lte(div(getSrc('period'), '3600'), '4')
                    )`, hasError: false},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := NewPrule(e, []byte(c.rule))
			if (err != nil) != c.hasError {
				t.Errorf("want error %v , res: %v", c.hasError, err != nil)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

func printp(i int, p *Prule) {
	for _, a := range p.Args {
		fmt.Printf("level: %d     == > name:%s \n", i, a.Name)
		if len(a.Args) > 0 || len(a.HArgs) > 0 {
			printp(i+1, a)
		}
	}
	for _, a := range p.HArgs {
		fmt.Printf("hilevel: %d     == > name:%s \n", i, a.Name)
		if len(a.HArgs) > 0 || len(a.Args) > 0 {
			printp(i+1, a)
		}
	}
}

func TestNewPrule2(t *testing.T) {
	e1 := NewExecEngine("test1", false)
	e2 := NewExecEngine("test2", true)

	cases := []struct {
		rule     string
		e        *Engine
		hasError bool
	}{
		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf')   '')`, e: e1, hasError: true},
		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf'))`, e: e1, hasError: false},
		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf') , 'redis')`, e: e1, hasError: true},
		{rule: `and(
					eq(left(getSrc('timestamp'), '10'), todate(dtToUnix(''))),
					gte(getSrc('level') , '6'),
					lte(div(getSrc('period'), '3600'), '4')
                    )`, e: e1, hasError: false},

		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')))`, e: e2, hasError: true},
		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf'),'redis','cid')`, e: e2, hasError: true},

		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf'))`, e: e2, hasError: false},
		{rule: `StoreAppendList(concat('id_key' , getSrc('cid')),getSrc('cmf'),'redis')`, e: e2, hasError: false},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := NewPrule(c.e, []byte(c.rule))
			if (err != nil) != c.hasError {
				t.Errorf("want error %v , res: %v", c.hasError, err != nil)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}
