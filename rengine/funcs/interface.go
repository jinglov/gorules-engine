package funcs

import (
	"reflect"
	"strconv"
	"time"
)

// 基本函数需要实现此interface
type IFun interface {
	// 函数名，返回结果全小写
	Label() string
}

// 普通函数需要实现此interface
type IDFun interface {
	// 函数执行体
	Exec(args ...string) string
}

// 需要把map[string]string传进去的函数需要实现此interface
type IMFun interface {
	// 函数执行体
	Exec(m map[string]string, args ...string) string
}

// 需要把thedate传进去的函数，需要实现此interface。一般在数据库操作中会用到
type ISFun interface {
	// 函数执行体
	Exec(thedate string, args ...string) string
}

// 高速方法需要实现此interface，高速方法指没有IO操作的各种方法
type IHeight interface {
	Height()
}

// 参数顺序可变，并且可以优先执行高效率的方法
type IOrder interface {
	// 如果当前结果和这个结果相等，则没必要执行后续的参数了。
	OrderRes() string
}

// 指定参数个数需要实现此interface
type IArgc interface {
	// 返回参数个数
	Argc() int
}

// 指定参数个数为几种
type IArgcs interface {
	Argcs() []int
}

// 最少参数个数
type IArgcMin interface {
	ArgcMin() int
}

// 哪些参数必须为固定值
type IMustValue interface {
	MustValue() []int
}

type ICache interface {
	Cache()
}

type IDb interface {
	DbFunc()
}

var HeightType = reflect.TypeOf((*IHeight)(nil)).Elem()
var OrderType = reflect.TypeOf((*IOrder)(nil)).Elem()
var DFunType = reflect.TypeOf((*IDFun)(nil)).Elem()
var MFunType = reflect.TypeOf((*IMFun)(nil)).Elem()
var SFunType = reflect.TypeOf((*ISFun)(nil)).Elem()
var ArgcType = reflect.TypeOf((*IArgc)(nil)).Elem()
var ArgcMinType = reflect.TypeOf((*IArgcMin)(nil)).Elem()
var ArgcsType = reflect.TypeOf((*IArgcs)(nil)).Elem()
var MustValueType = reflect.TypeOf((*IMustValue)(nil)).Elem()
var CacheType = reflect.TypeOf((*ICache)(nil)).Elem()
var DbType = reflect.TypeOf((*IDb)(nil)).Elem()

const (
	True       = "1"
	False      = string(0)
	DELIMIT    = string(30)
	DayType    = 'd'
	SecondType = 's'
	DAY        = 86400
)

func ToInt(s string) int {
	if s == False {
		return 0
	}
	num, _ := strconv.Atoi(s)
	return num
}
func ToString(i int) string {
	return strconv.Itoa(i)
}

func BToString(b bool) string {
	if b {
		return True
	}
	return False
}

// 支持 2 种
// 123s， 123秒
// 123d, 123天
func ParseExp(arg string) int {
	if len(arg) == 0 {
		return 0
	}
	switch arg[len(arg)-1] {
	case DayType:
		return ToInt(arg[:len(arg)-1]) * DAY
	case SecondType:
		return ToInt(arg[:len(arg)-1])
	default:
		return 0
	}
}

func GetDayInt(thedate string) int {
	var tm time.Time
	if thedate == "" {
		tm = time.Now()
	} else {
		var err error
		tm, err = time.ParseInLocation("060102", thedate, time.Local)
		if err != nil {
			tm = time.Now()
		}
	}
	_, offset := tm.Zone()
	return int((tm.Unix() + int64(offset)) / 86400)
}
