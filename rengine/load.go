package rengine

import (
	"bytes"
	"fmt"
	"github.com/jinglov/gorules-engine/rengine/funcs"
	"github.com/omigo/log"
	"reflect"
	"strings"
	"unicode"
)

const (
	FunStartToken      = '('
	MustFunStart  byte = 1
	FunEndToken        = ')'
	MustFunEnd    byte = 2
	ValueToken         = '\''
	MustValue     byte = 4
	NewArgToken        = ','
	MustNewArg    byte = 8
	ItemTypeFun        = "fun"
	ItemTypeValue      = "value"
)

func NewPrule(e *Engine, b []byte) (*Prule, error) {
	p := &Prule{}
	var n = 0
	var check byte = 0
	var funToken = 0
	check |= MustFunStart
	check |= MustValue
	p.Index = n
	err := parserToStruct(p, &n, b, &check, &funToken, e)
	if err != nil {
		log.Error("NewPrule err", err)
		return nil, err
	}
	err = checkPrule(p, e)
	return p, err
}

func checkPrule(p *Prule, e *Engine) error {
	for _, a := range append(p.Args, p.HArgs...) {
		if a.Type == ItemTypeFun {
			expFun, ok := e.funcs[a.Name]
			if !ok {
				return fmt.Errorf("Function: %s not register", a.Name)
			}
			if reflect.TypeOf(expFun).Implements(funcs.ArgcType) {
				argc := len(a.Args) + len(a.HArgs)
				argcAllow := expFun.(funcs.IArgc).Argc()
				if e.AllowDbParam && reflect.TypeOf(expFun).Implements(funcs.DbType) {
					if argc != argcAllow && argc != argcAllow+1 {
						return fmt.Errorf("Index:%d   funName: %s want argc: [%d,%d] , than: %d", a.Index, a.Name, argcAllow, argcAllow+1, argc)
					}
				} else {
					if argcAllow != argc {
						return fmt.Errorf("Index:%d   funName: %s want argc: %d , than: %d", a.Index, a.Name, argcAllow, argc)
					}
				}
			}

			if reflect.TypeOf(expFun).Implements(funcs.ArgcsType) {
				argcs := expFun.(funcs.IArgcs).Argcs()
				argc := len(a.Args) + len(a.HArgs)
				hasErr := true
				for _, argN := range argcs {
					if argN == argc {
						hasErr = false
						break
					}
				}
				if hasErr {
					return fmt.Errorf("Index:%d   funName: %s want argcs: %d , than: %d", a.Index, a.Name, expFun.(funcs.IArgcs).Argcs(), len(a.Args)+len(a.HArgs))
				}
			}
			if reflect.TypeOf(expFun).Implements(funcs.ArgcMinType) {
				if expFun.(funcs.IArgcMin).ArgcMin() > len(a.Args)+len(a.HArgs) {
					return fmt.Errorf("Index: %d funName: %s want argcmin: %d , than: %d", a.Index, a.Name, expFun.(funcs.IArgcMin).ArgcMin(), len(a.Args)+len(a.HArgs))
				}
			}
			if reflect.TypeOf(expFun).Implements(funcs.MustValueType) {
				for _, n := range expFun.(funcs.IMustValue).MustValue() {
					p := a.Args[n]
					if p.Type != ItemTypeValue {
						return fmt.Errorf("Index: %d funName: %s at args:%d must value but is fun (name: %s)", a.Index, a.Name, n, p.Name)
					}
				}
			}
			err := checkPrule(a, e)
			if err != nil {
				log.Error("checkPrule err", err)
				return err
			}
		}
	}
	return nil
}

func parserToStruct(p *Prule, n *int, exp []byte, check *byte, funToken *int, engine *Engine) error {
	var valueFlag, cmdFlag bool
	cmdFlag = true
	subp := &Prule{Index: *n}
	var buf bytes.Buffer
	l := len(exp)
	for {
		if l == *n {
			if *funToken > 0 {
				log.Info(string(exp[:*n]))
				return fmt.Errorf("must )x%d at index: %d", *funToken, *n-1)
			}
			return nil
		}
		b := exp[*n]
		*n++
		switch b {
		case FunStartToken:
			if valueFlag {
				buf.WriteByte(b)
				continue
			}
			if MustFunStart&*check == 0 {
				log.Info(string(exp[:*n]))
				return fmt.Errorf("must not ( at index: %d", *n-1)
			}
			*funToken++
			subp.Name = strings.ToLower(buf.String())
			if expFun, ok := engine.funcs[subp.Name]; ok {
				if !reflect.TypeOf(expFun).Implements(funcs.DFunType) &&
					!reflect.TypeOf(expFun).Implements(funcs.MFunType) &&
					!reflect.TypeOf(expFun).Implements(funcs.SFunType) {
					log.Fatalf("fun: %s must implement IDFun OR IMFun OR ISFun", subp.Name)
				}
				if !reflect.TypeOf(expFun).Implements(funcs.HeightType) {
					subp.LowRule = true
				}
				subp.Type = ItemTypeFun
				*check &= ^MustNewArg
				*check |= MustValue
				*check &= ^MustFunEnd
				err := parserToStruct(subp, n, exp, check, funToken, engine)
				if err != nil {
					return err
				}
				if subp.LowRule {
					p.LowRule = true
				}
				pFun := engine.funcs[p.Name]
				if !subp.LowRule && pFun != nil && reflect.TypeOf(pFun).Implements(funcs.OrderType) {
					p.HArgs = append(p.HArgs, subp)
				} else {
					p.Args = append(p.Args, subp)
				}
				buf.Reset()
			} else {
				return fmt.Errorf("Function: %s not register", subp.Name)
			}
		case FunEndToken:
			if valueFlag {
				buf.WriteByte(b)
				continue
			}
			if MustFunEnd&*check == 0 {
				log.Info(string(exp[:*n]))
				return fmt.Errorf("must not ) at index: %d", *n-1)
			}
			*funToken--
			if *funToken < 0 {
				log.Info(string(exp[:*n]))
				return fmt.Errorf("must not ) at index: %d", *n-1)
			}
			*check &= ^MustFunStart
			*check &= ^MustValue
			*check |= MustNewArg
			*check |= MustFunEnd
			return nil
		case ValueToken:
			if MustValue&*check == 0 {
				log.Info(string(exp[:*n]))
				return fmt.Errorf("must not \"'\" at index: %d", *n-1)
			}
			if valueFlag {
				// 把值压入栈
				subp.Type = ItemTypeValue
				subp.Name = buf.String()
				p.Args = append(p.Args, subp)
				buf.Reset()
			}
			valueFlag = !valueFlag
			*check &= ^MustValue
			*check &= ^MustFunStart
			*check |= MustValue
			*check |= MustFunEnd
		case NewArgToken:
			if valueFlag {
				buf.WriteByte(b)
				continue
			}
			if NewArgToken&*check == 0 {
				log.Info(string(exp[:*n]))
				return fmt.Errorf("must not , at index: %d", *n-1)
			}
			subp = &Prule{Index: *n}
			cmdFlag = true
			*check &= ^MustNewArg
			*check |= MustFunStart
			*check &= ^MustFunEnd
			*check |= MustValue
		default:
			if valueFlag {
				*check |= MustValue
				buf.WriteByte(b)
				continue
			}
			if cmdFlag {
				if !unicode.IsSpace(rune(b)) {
					// 在命令状态下，后面的关键字符只能是(
					// *check |= MustFunStart
					*check &= ^MustValue
					*check &= ^MustNewArg
					buf.WriteByte(b)
				}
			}
		}
	}
	return nil
}
