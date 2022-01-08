package funcs

import (
	"strings"
)

func init() {
	AddAction(And{})
	AddAction(Or{})
	AddAction(Eq{})
	AddAction(Neq{})
	AddAction(Gt{})
	AddAction(Gte{})
	AddAction(Lt{})
	AddAction(Lte{})
	AddAction(Like{})
	AddAction(LikeOr{})
	AddAction(In{})
	AddAction(NotIn{})
	AddAction(SplitIn{})
	AddAction(LtVersion{})
	AddAction(EqVersion{})
	AddAction(GtVersion{})
}

var ActionFuncs = []IFun{}

func AddAction(fun IFun) {
	ActionFuncs = append(ActionFuncs, fun)
}

// And
// 逻辑与
// and('1','1') == '1'
// and('1','0') == '0'
// 所有参数只要有一个非，则结果为非
// 这是一个高性能的方法
// 这个方法的参数顺序不重要，可以根据参数性能优化顺序
type And struct{}

var _ IFun = And{}
var _ IDFun = And{}
var _ IHeight = And{}
var _ IOrder = And{}
var _ IArgcMin = And{}

func (And) Height() {}
func (And) OrderRes() string {
	return False
}
func (And) Label() string {
	return "and"
}
func (And) ArgcMin() int {
	return 2
}
func (And) Exec(args ...string) string {
	for _, arg := range args {
		if arg == False {
			return False
		}
	}
	return True
}

// Or
// 逻辑或
// or('1','0') == '1'
// or('0','0') == '0'
// 所有参数只要求一个真，则结果为真
// 这是一个高性能的方法
// 这个方法的参数顺序不重要，可以根据参数性能优化顺序
type Or struct{}

var _ IFun = Or{}
var _ IDFun = Or{}
var _ IHeight = Or{}
var _ IOrder = Or{}
var _ IArgcMin = Or{}

func (Or) Height() {}
func (Or) OrderRes() string {
	return True
}
func (Or) Label() string {
	return "or"
}
func (Or) ArgcMin() int {
	return 2
}
func (Or) Exec(args ...string) string {
	for _, arg := range args {
		if arg == True {
			return True
		}
	}
	return False
}

// Eq
// 相等
// eq('a','b') == '0'
// eq('a','a') == '1'
// 判断所有参数是否相等
// 这是一个高性能的方法
// 这个方法的参数顺序不重要，可以根据参数性能优化顺序
type Eq struct{}

var _ IFun = Eq{}
var _ IDFun = Eq{}
var _ IHeight = Eq{}
var _ IArgcMin = Eq{}

func (Eq) Height() {}
func (Eq) Label() string {
	return "eq"
}
func (Eq) ArgcMin() int {
	return 2
}
func (Eq) Exec(args ...string) string {
	for _, i := range args[1:] {
		if i != args[0] {
			return False
		}
	}
	return True
}

// Neq
// 不相等
// neq('a','a') == '0'
// neq('a','b') == '1'
// 所有参数有任何两个不相等
// 这是一个高性能的方法
// 这个方法的参数顺序不重要，可以根据参数性能优化顺序
type Neq struct{}

var _ IFun = Neq{}
var _ IDFun = Neq{}
var _ IHeight = Neq{}
var _ IArgc = Neq{}

func (Neq) Height() {}
func (Neq) Label() string {
	return "neq"
}
func (Neq) Argc() int {
	return 2
}
func (Neq) Exec(args ...string) string {
	for _, i := range args[1:] {
		if i != args[0] {
			return True
		}
	}
	return False
}

// Gt
// 参数1数值大于参数2数值
// gt('1','2') == '0'
// gt('2','1') == '1'
// gt('1','1') == '0'
// 这是一个高性能的方法
// 限制只能有2个参数
type Gt struct{}

var _ IFun = Gt{}
var _ IDFun = Gt{}
var _ IHeight = Gt{}
var _ IArgc = Gt{}

func (Gt) Height() {}
func (Gt) Label() string {
	return "gt"
}
func (Gt) Argc() int {
	return 2
}
func (Gt) Exec(args ...string) string {
	if ToInt(args[0]) > ToInt(args[1]) {
		return True
	}
	return False
}

// Gte
// 参数1数值大于等于参数2数值
// gte('1','2') == '0'
// gte('2','1') == '1'
// gte('1','1') == '1'
// 这是一个高性能的方法
// 限制只能有2个参数
type Gte struct{}

var _ IFun = Gte{}
var _ IDFun = Gte{}
var _ IHeight = Gte{}
var _ IArgc = Gte{}

func (Gte) Height() {}
func (Gte) Label() string {
	return "gte"
}
func (Gte) Argc() int {
	return 2
}
func (Gte) Exec(args ...string) string {
	if ToInt(args[0]) >= ToInt(args[1]) {
		return True
	}
	return False
}

// Lt
// 参数1数值小于参数2数值
// lt('1','2') == '1'
// lt('2','1') == '0'
// lt('1','1') == '0'
// 这是一个高性能的方法
// 限制只能有2个参数
type Lt struct{}

var _ IFun = Lt{}
var _ IDFun = Lt{}
var _ IHeight = Lt{}
var _ IArgc = Lt{}

func (Lt) Height() {}
func (Lt) Label() string {
	return "lt"
}
func (Lt) Argc() int {
	return 2
}
func (Lt) Exec(args ...string) string {
	if ToInt(args[0]) < ToInt(args[1]) {
		return True
	}
	return False
}

// Lte
// 参数1数值小于等于参数2数值
// lte('1','2') == '1'
// lte('2','1') == '0'
// lte('1','1') == '1'
// 这是一个高性能的方法
// 限制只能有2个参数
type Lte struct{}

var _ IFun = Lte{}
var _ IDFun = Lte{}
var _ IHeight = Lte{}
var _ IArgc = Lte{}

func (Lte) Height() {}
func (Lte) Label() string {
	return "lte"
}
func (Lte) Argc() int {
	return 2
}
func (Lte) Exec(args ...string) string {
	if ToInt(args[0]) <= ToInt(args[1]) {
		return True
	}
	return False
}

// Like
// 参数1字符串中是否包含参数2
// like('abc' , 'ab') == '1'
// like('abc' ,'cd') == '0'
// 这是一个高性能的方法
// 限制只能有2个参数
type Like struct{}

var _ IFun = Like{}
var _ IDFun = Like{}
var _ IHeight = Like{}
var _ IArgc = Like{}

func (Like) Height() {}
func (Like) Label() string {
	return "like"
}
func (Like) Argc() int {
	return 2
}
func (Like) Exec(args ...string) string {
	if strings.Index(args[0], args[1]) != -1 {
		return True
	}
	return False
}

// LikeOr
// 参数1字符串中是否包含其它参数中的其中一个字符串
// likeor('abc' , 'a','d','e') == '1'
// likeor('abc' , 'd','e') == '0'
// 这是一个高性能的方法
type LikeOr struct{}

var _ IFun = LikeOr{}
var _ IDFun = LikeOr{}
var _ IHeight = LikeOr{}
var _ IArgcMin = LikeOr{}

func (LikeOr) Height() {}
func (LikeOr) Label() string {
	return "likeor"
}
func (LikeOr) ArgcMin() int {
	return 2
}
func (LikeOr) Exec(args ...string) string {
	if args[0] == "" || args[0] == False {
		return False
	}
	for _, i := range args[1:] {
		if strings.Index(args[0], i) != -1 {
			return True
		}
	}

	return False
}

// In
// 参数1，是否包含在其它参数中
// in('a' , 'a' , 'b' , 'c') == '1'
// in('a' , 'b' , 'c' , 'e') == '0'
// 这是一个高性能的方法
type In struct{}

var _ IFun = In{}
var _ IDFun = In{}
var _ IHeight = In{}
var _ IArgcMin = In{}

func (In) Height() {}
func (In) Label() string {
	return "in"
}
func (In) ArgcMin() int {
	return 2
}
func (In) Exec(args ...string) string {
	for _, item := range args[1:] {
		if args[0] == item {
			return True
		}
	}
	return False
}

// NotIn
// 参数1，是否不在其它参数中
// notin('a' , 'a' , 'b' , 'c') == '0'
// notin('a' , 'b' , 'c' , 'e') == '1'
// 这是一个高性能的方法
type NotIn struct{}

var _ IFun = NotIn{}
var _ IDFun = NotIn{}
var _ IHeight = NotIn{}
var _ IArgcMin = NotIn{}

func (NotIn) Height() {}
func (NotIn) Label() string {
	return "notin"
}
func (NotIn) ArgcMin() int {
	return 2
}
func (NotIn) Exec(args ...string) string {
	for _, item := range args[1:] {
		if args[0] == item {
			return False
		}
	}
	return True
}

// SplitIn
// 参数1，是否不在其它参数中
// SplitIn('a' , 'a,b,c') == '0'
// SplitIn('a' , 'b,c,e') == '1'
// 这是一个高性能的方法
type SplitIn struct{}

var _ IFun = SplitIn{}
var _ IDFun = SplitIn{}
var _ IHeight = SplitIn{}
var _ IArgc = SplitIn{}

func (SplitIn) Height() {}
func (SplitIn) Label() string {
	return "splitin"
}
func (SplitIn) Argc() int {
	return 3
}
func (SplitIn) Exec(args ...string) string {
	items := strings.Split(args[1], args[2])
	for _, item := range items {
		if args[0] == item {
			return True
		}
	}
	return False
}

// LtVersion
// 参数1的版本是否比参数2的版本低(只比较前3位)
// ltversion('2.8.0' , '2.8.1') == '1'
// ltversion('2.8.0' , '2.8.0') == '0'
// ltversion('2.8.0.1' , '2.8.0.2') == '0'
// 这是一个高性能的方法
// 限制只能有2个参数
type LtVersion struct{}

var _ IFun = LtVersion{}
var _ IDFun = LtVersion{}
var _ IHeight = LtVersion{}
var _ IArgc = LtVersion{}

func (LtVersion) Height() {}
func (LtVersion) Label() string {
	return "ltversion"
}

func (LtVersion) Argc() int {
	return 2
}

func (LtVersion) Exec(args ...string) string {
	if versionCompare(args[0], args[1]) == -1 {
		return True
	}
	return False
}

// EqVersion
// 参数1的版本和参数2的版本是否相等（只比较前3位)
// eqversion('2.8.0' , '2.8.1') == '0'
// eqversion('2.8.0' , '2.8.0') == '1'
// eqversion('2.8.0.1' , '2.8.0.2') == '1'
// 这是一个高性能的方法
// 限制只能有2个参数
type EqVersion struct{}

var _ IFun = EqVersion{}
var _ IDFun = EqVersion{}
var _ IHeight = EqVersion{}
var _ IArgc = EqVersion{}

func (EqVersion) Height() {}
func (EqVersion) Label() string {
	return "eqversion"
}

func (EqVersion) Argc() int {
	return 2
}

func (EqVersion) Exec(args ...string) string {
	if versionCompare(args[0], args[1]) == 0 {
		return True
	}
	return False
}

// GtVersion
// 参数1的版本是否比参数2的版本高（只比较前3位）
// gtversion('2.8.0' , '2.8.1') == '0'
// gtversion('2.8.0' , '2.8.0') == '1'
// gtversion('2.8.0.1' , '2.8.0.2') == '0'
// 这是一个高性能的方法
// 限制只能有2个参数
type GtVersion struct{}

var _ IFun = GtVersion{}
var _ IDFun = GtVersion{}
var _ IHeight = GtVersion{}
var _ IArgc = GtVersion{}

func (GtVersion) Height() {}
func (GtVersion) Label() string {
	return "gtversion"
}

func (GtVersion) Argc() int {
	return 2
}

func (GtVersion) Exec(args ...string) string {
	if versionCompare(args[0], args[1]) == 1 {
		return True
	}
	return False
}

// 版本对比
// a == b return 0
// a > b return 1
// a < b return -1
func versionCompare(a, b string) int8 {
	// 如果 a b 字面值相同，版本号一致
	if a == b {
		return 0
	}

	// 分隔开比较 2.7.4.1228.1800  2.7.10.1228.1800
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")

	for i := 0; i < len(as); i++ {
		// 只比较前三位版本，后面的小版本不做比较 2018-07-24 14:25
		if i == 3 {
			return 0
		}
		if i == len(bs) {
			// i=2 2.7.4 > 2.7 or 2.7.0 > 2.7
			return 1
		}
		if len(as[i])-len(bs[i]) > 0 { // 2.7.100 > 2.7.4
			return 1
		} else if len(as[i])-len(bs[i]) < 0 { // 2.7.4 > 2.7.100
			return -1
		}
		if as[i] > bs[i] { // 2.7.9 > 2.7.4
			return 1
		} else if as[i] < bs[i] {
			return -1
		}
	}
	if len(as) == 3 {
		return 0
	}
	// 2.7 < 2.7.4 or 2.7 < 2.7.0
	return -1
}
