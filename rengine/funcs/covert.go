package funcs

import (
	"bytes"
	"encoding/hex"
	"github.com/spaolacci/murmur3"
	"strconv"
	"strings"
	"time"
)

func init() {
	AddCovert(Len{})
	AddCovert(Sum{})
	AddCovert(Subtract{})
	AddCovert(Div{})
	AddCovert(Multiply{})
	AddCovert(Concat{})
	AddCovert(Split{})
	AddCovert(Upper{})
	AddCovert(Lower{})
	AddCovert(Left{})
	AddCovert(Right{})
	AddCovert(DtToUnix{})
	AddCovert(ToDate{})
	AddCovert(TrimSpace{})
	AddCovert(Hash{})
}

var CovertFuncs = []IFun{}

func AddCovert(fun IFun) {
	CovertFuncs = append(CovertFuncs, fun)
}

// Len
// 计算字符串长度
// len('abc') == '3'
// 这是一个高性能的方法
// 限制只能有1个参数
type Len struct{}

var _ IFun = Len{}
var _ IDFun = Len{}
var _ IHeight = Len{}
var _ IArgc = Len{}

func (Len) Height() {}
func (Len) Label() string {
	return "len"
}
func (Len) Argc() int {
	return 1
}
func (Len) Exec(args ...string) string {
	return ToString(len(args[0]))
}

// Sum
// 求和
// sum('1','2','3') == '6'
// 这是一个高性能的方法
type Sum struct{}

var _ IFun = Sum{}
var _ IDFun = Sum{}
var _ IHeight = Sum{}
var _ IArgcMin = Sum{}

func (Sum) Height() {}
func (Sum) Label() string {
	return "sum"
}
func (Sum) ArgcMin() int {
	return 2
}
func (Sum) Exec(args ...string) string {
	c := 0
	for _, i := range args {
		c += ToInt(i)
	}
	return ToString(c)
}

// Subtract
// 减法
// subtract('8','5') == '3'
// 这是一个高性能的方法
// 限制只能有2个参数
type Subtract struct{}

var _ IFun = Subtract{}
var _ IDFun = Subtract{}
var _ IHeight = Subtract{}
var _ IArgc = Subtract{}

func (Subtract) Height() {}
func (Subtract) Label() string {
	return "subtract"
}
func (Subtract) Argc() int {
	return 2
}
func (Subtract) Exec(args ...string) string {
	return ToString(ToInt(args[0]) - ToInt(args[1]))
}

// Div
// 除法（只取商的整数部份）
// div('8','2') == '4'
// div('5' , '2') == '2'
// 这是一个高性能的方法
// 限制只能有2个参数

type Div struct{}

var _ IFun = Div{}
var _ IDFun = Div{}
var _ IHeight = Div{}
var _ IArgc = Div{}

func (Div) Height() {}
func (Div) Label() string {
	return "div"
}
func (Div) Argc() int {
	return 2
}
func (Div) Exec(args ...string) string {
	x := ToInt(args[0])
	y := ToInt(args[1])
	if y != 0 {
		return ToString(x / y)
	}
	return False
}

// Multipy
// 乘法
// multiply('2','3','4') == '24'
// multiply('10','20','0') == '0'
// 这是一个高性能的方法

type Multiply struct{}

var _ IFun = Multiply{}
var _ IDFun = Multiply{}
var _ IHeight = Multiply{}
var _ IArgcMin = Multiply{}

func (Multiply) Height() {}
func (Multiply) Label() string {
	return "multiply"
}
func (Multiply) ArgcMin() int {
	return 2
}
func (Multiply) Exec(args ...string) string {
	c := 1
	for _, i := range args {
		c *= ToInt(i)
	}
	return ToString(c)
}

// Concat
// 字符串连接
// concat('abc' , 'de' , 'f') == 'abcdef'
// 这是一个高性能的方法
type Concat struct{}

var _ IFun = Concat{}
var _ IDFun = Concat{}
var _ IHeight = Concat{}
var _ IArgcMin = Concat{}

func (Concat) Height() {}
func (Concat) Label() string {
	return "concat"
}
func (Concat) ArgcMin() int {
	return 2
}
func (Concat) Exec(args ...string) string {
	switch len(args) {
	case 0:
		return ""
	case 1:
		return args[0]
	}
	n := 0
	for _, i := range args {
		n += len(i)
	}

	b := new(bytes.Buffer)
	b.Grow(n)
	for _, s := range args {
		b.WriteString(s)
	}
	return b.String()
}

// Split
// 取一个数组中，第N个元素下标从0开始
// split('a,b,c,d' , ',' , 2) == 'c'
// 这是一个高性能的方法
// 限制只能有3个参数
type Split struct{}

var _ IFun = Split{}
var _ IDFun = Split{}
var _ IHeight = Split{}
var _ IArgc = Split{}

func (Split) Height() {}
func (Split) Label() string {
	return "split"
}
func (Split) Argc() int {
	return 3
}
func (Split) Exec(args ...string) string {
	ss := strings.Split(args[0], args[1])
	index := ToInt(args[2])

	if len(ss) < index {
		return ""
	}
	if index < 0 {
		index += len(ss)
	}
	if index >= len(ss) || index < 0 {
		return ""
	}
	return ss[index]
}

// Upper
// 字符串转大写
// upper('abc') == 'ABC'
// 这是一个高性能的方法
// 限制只能有1个参数
type Upper struct{}

var _ IFun = Upper{}
var _ IDFun = Upper{}
var _ IHeight = Upper{}
var _ IArgc = Upper{}

func (Upper) Height() {}
func (Upper) Label() string {
	return "upper"
}
func (Upper) Argc() int {
	return 1
}
func (Upper) Exec(args ...string) string {
	return strings.ToUpper(args[0])

}

// Lower
// 字符串转小写
// lower('ABC') == 'abc'
// lower('Abc') == 'abc'
// 这是一个高性能的方法
// 限制只能有1个参数
type Lower struct{}

var _ IFun = Lower{}
var _ IDFun = Lower{}
var _ IHeight = Lower{}
var _ IArgc = Lower{}

func (Lower) Height() {}
func (Lower) Label() string {
	return "lower"
}
func (Lower) Argc() int {
	return 1
}
func (Lower) Exec(args ...string) string {
	return strings.ToLower(args[0])
}

// Left
// 字符串从左边取N个字符，负值为从右边开始数
// left('abcdefg' , 2) == 'ab'
// left('abcdefg' , -2) == 'abcde'
// left('abcdefg' , 8) == 'abcdefg'
// left('abcdefg' , -8) == ''
// 这是一个高性能的方法
// 限制只能有2个参数
type Left struct{}

var _ IFun = Left{}
var _ IDFun = Left{}
var _ IHeight = Left{}
var _ IArgc = Left{}

func (Left) Height() {}
func (Left) Label() string {
	return "left"
}
func (Left) Argc() int {
	return 2
}
func (Left) Exec(args ...string) string {
	lens := ToInt(args[1])
	if len(args[0]) < lens {
		return args[0]
	}
	if lens < 0 {
		lens += len(args[0])
		if lens < 0 {
			lens = 0
		}
	}
	return args[0][:lens]
}

// Right
// 字符串从右边取N个字符，负值为从左边开始数
// right('abcdefg' , 2) == 'fg'
// right('abcdefg' , -2) == 'cdefg'
// right('abcdefg' , 8) == 'abcdefg'
// right('abcdefg' , -8) == ''
// 这是一个高性能的方法
// 限制只能有2个参数
type Right struct{}

var _ IFun = Right{}
var _ IDFun = Right{}
var _ IHeight = Right{}
var _ IArgc = Right{}

func (Right) Height() {}
func (Right) Label() string {
	return "right"
}
func (Right) Argc() int {
	return 2
}
func (Right) Exec(args ...string) string {
	lens := ToInt(args[1])
	strLen := len(args[0])
	if lens < 0 {
		lens += strLen
		if lens < 0 {
			lens = 0
		}
	}
	if len(args[0]) < lens {
		return args[0]
	}
	return args[0][strLen-lens : strLen]
}

// DtToUnix
// 把日期时间转换成时间戳，传空时返回当前时间戳
// dttounix('2019-11-01 12:05:00') == '1572581100'
// 这是一个高性能的方法
type DtToUnix struct{}

var _ IFun = DtToUnix{}
var _ IDFun = DtToUnix{}
var _ IHeight = DtToUnix{}
var _ IArgc = DtToUnix{}

// todo 支持更多时间日期格式
var TimeFormats = []string{
	"2006-01-02T15:04:05-07:00",
	"2006-01-02 15:04:05",
}

func (DtToUnix) Height() {}
func (DtToUnix) Label() string {
	return "dttounix"
}
func (DtToUnix) Argc() int {
	return 1
}
func (DtToUnix) Exec(args ...string) string {
	if len(args) < 1 || args[0] == False || args[0] == "0" || args[0] == "" {
		return ToString(int(time.Now().Unix()))
	}
	for _, f := range TimeFormats {
		t, err := time.ParseInLocation(f, args[0], time.Local)
		if err == nil {
			return ToString(int(t.Unix()))
		}
	}
	return False
}

// ToDate
// 把时间戳转成日期
// todate('1572581100') == '2019-11-01'
// 这是一个高性能的方法
// 限制只能有1个参数
type ToDate struct{}

var _ IFun = ToDate{}
var _ IDFun = ToDate{}
var _ IHeight = ToDate{}
var _ IArgc = ToDate{}

func (ToDate) Height() {}
func (ToDate) Label() string {
	return "todate"
}
func (ToDate) Argc() int {
	return 1
}
func (ToDate) Exec(args ...string) string {
	tm, _ := strconv.Atoi(args[0])
	tt := time.Unix(int64(tm), 0)
	return tt.Format("2006-01-02")
}

// TrimSpace
// 字符串两头去空
// TrimSpace('   abcd    ') == 'abc d'
// 这是一个高性能的方法
// 限制只能有1个参数
type TrimSpace struct{}

var _ IFun = TrimSpace{}
var _ IDFun = TrimSpace{}
var _ IHeight = TrimSpace{}
var _ IArgc = TrimSpace{}

func (TrimSpace) Height() {}
func (TrimSpace) Label() string {
	return "trimspace"
}
func (TrimSpace) Argc() int {
	return 1
}
func (TrimSpace) Exec(args ...string) string {
	return strings.TrimSpace(args[0])
}

// Hash
// 字符串两头去空
// Hash('12355abefdass') == '0123456789abcdef'
// 这是一个高性能的方法
// 限制只能有1个参数
type Hash struct{}

var _ IFun = Hash{}
var _ IDFun = Hash{}
var _ IHeight = Hash{}
var _ IArgc = Hash{}

func (Hash) Height() {}
func (Hash) Label() string {
	return "hash"
}
func (Hash) Argc() int {
	return 1
}
func (Hash) Exec(args ...string) string {
	h := murmur3.New128()
	h.Write([]byte(args[0]))
	res := h.Sum(nil)
	return hex.EncodeToString(res)
}
