package funcs

import "github.com/jinglov/gorules-engine/db"

func init() {
	AddSetStore(StoreSetValue{})
	AddSetStore(StoreSetValueExp{})
	AddSetStore(StoreSetNxValue{})
	AddSetStore(StoreSetNxValueExp{})
	AddSetStore(StoreAddValue{})
	AddSetStore(StoreAddValueExp{})
	AddSetStore(StoreAddValueDays{})
	AddGetStore(StoreGetValue{})
	AddGetStore(StoreGetValueExp{})
	AddGetStore(StoreGetValueDays{})
}

type StoreSetValue struct{}

func getRowDb(args []string, index int) db.DB {
	if len(args) == index {
		argDb := db.GetDb(args[index-1])
		if argDb != nil {
			return argDb
		}
	}
	return db.RowDefault
}

var _ IFun = StoreSetValue{}
var _ ISFun = StoreSetValue{}
var _ IArgc = StoreSetValue{}
var _ IDb = StoreSetValue{}

func (StoreSetValue) Label() string {
	return "storesetvalue"
}
func (StoreSetValue) Argc() int {
	return 2
}
func (StoreSetValue) DbFunc() {
	return
}
func (StoreSetValue) Exec(thedate string, args ...string) string {
	return ToString(getRowDb(args, 3).Set(thedate+args[0], args[1], DAY))
}

type StoreSetNxValue struct{}

var _ IFun = StoreSetNxValue{}
var _ ISFun = StoreSetNxValue{}
var _ IArgc = StoreSetNxValue{}
var _ IDb = StoreSetNxValue{}

func (StoreSetNxValue) Label() string {
	return "storesetnxvalue"
}
func (StoreSetNxValue) Argc() int {
	return 2
}
func (StoreSetNxValue) DbFunc() {
	return
}
func (StoreSetNxValue) Exec(thedate string, args ...string) string {
	return BToString(getRowDb(args, 3).SetNx(thedate+args[0], args[1], DAY))
}

type StoreAddValue struct{}

var _ IFun = StoreAddValue{}
var _ ISFun = StoreAddValue{}
var _ IArgc = StoreAddValue{}
var _ IDb = StoreAddValue{}

func (StoreAddValue) Label() string {
	return "storeaddvalue"
}
func (StoreAddValue) Argc() int {
	return 2
}
func (StoreAddValue) DbFunc() {
	return
}
func (StoreAddValue) Exec(thedate string, args ...string) string {
	return ToString(getRowDb(args, 3).Add(thedate+args[0], ToInt(args[1]), DAY))
}

type StoreGetValue struct{}

var _ IFun = StoreGetValue{}
var _ ISFun = StoreGetValue{}
var _ IArgc = StoreGetValue{}
var _ IDb = StoreGetValue{}

func (StoreGetValue) Label() string {
	return "storegetvalue"
}
func (StoreGetValue) Argc() int {
	return 1
}
func (StoreGetValue) DbFunc() {
	return
}
func (StoreGetValue) Exec(thedate string, args ...string) string {
	return getRowDb(args, 2).Get(thedate + args[0])
}

type StoreSetValueExp struct{}

var _ IFun = StoreSetValueExp{}
var _ IDFun = StoreSetValueExp{}
var _ IArgc = StoreSetValueExp{}
var _ IDb = StoreSetValueExp{}

func (StoreSetValueExp) Label() string {
	return "storesetvalueexp"
}
func (StoreSetValueExp) Argc() int {
	return 3
}
func (StoreSetValueExp) DbFunc() {
	return
}
func (StoreSetValueExp) Exec(args ...string) string {
	expire := ParseExp(args[0])
	return ToString(getRowDb(args, 4).Set(args[1], args[2], expire))
}

type StoreSetNxValueExp struct{}

var _ IFun = StoreSetNxValueExp{}
var _ IDFun = StoreSetNxValueExp{}
var _ IArgc = StoreSetNxValueExp{}
var _ IDb = StoreSetNxValueExp{}

func (StoreSetNxValueExp) Label() string {
	return "storesetnxvalueexp"
}
func (StoreSetNxValueExp) Argc() int {
	return 3
}
func (StoreSetNxValueExp) DbFunc() {
	return
}
func (StoreSetNxValueExp) Exec(args ...string) string {
	expire := ParseExp(args[0])
	return BToString(getRowDb(args, 4).SetNx(args[1], args[2], expire))
}

type StoreAddValueExp struct{}

var _ IFun = StoreAddValueExp{}
var _ IDFun = StoreAddValueExp{}
var _ IArgc = StoreAddValueExp{}
var _ IDb = StoreAddValueExp{}

func (StoreAddValueExp) Label() string {
	return "storeaddvalueexp"
}
func (StoreAddValueExp) Argc() int {
	return 3
}
func (StoreAddValueExp) DbFunc() {
	return
}
func (StoreAddValueExp) Exec(args ...string) string {
	expire := ParseExp(args[0])
	return ToString(getRowDb(args, 4).Add(args[1], ToInt(args[2]), expire))
}

type StoreGetValueExp struct{}

var _ IFun = StoreGetValueExp{}
var _ IDFun = StoreGetValueExp{}
var _ IArgc = StoreGetValueExp{}
var _ IDb = StoreGetValueExp{}

func (StoreGetValueExp) Label() string {
	return "storegetvalueexp"
}
func (StoreGetValueExp) Argc() int {
	return 1
}
func (StoreGetValueExp) DbFunc() {
	return
}
func (StoreGetValueExp) Exec(args ...string) string {
	return getRowDb(args, 2).Get(args[0])
}

type StoreAddValueDays struct{}

var _ IFun = StoreAddValueDays{}
var _ ISFun = StoreAddValueDays{}
var _ IArgc = StoreAddValueDays{}
var _ IDb = StoreAddValueDays{}
var storeGetValueDays = StoreGetValueDays{}

func (StoreAddValueDays) Label() string {
	return "storeaddvaluedays"
}
func (StoreAddValueDays) Argc() int {
	return 3
}
func (StoreAddValueDays) DbFunc() {
	return
}
func (StoreAddValueDays) Exec(thedate string, args ...string) string {
	dayInt := ToString(GetDayInt(thedate))
	maxDay := ToInt(args[0])
	getRowDb(args, 4).Add(dayInt+args[1], ToInt(args[2]), maxDay*DAY)
	if len(args) == 4 {
		return storeGetValueDays.Exec(thedate, args[0], args[1], args[3])
	} else {
		return storeGetValueDays.Exec(thedate, args[0], args[1])
	}
}

type StoreGetValueDays struct{}

var _ IFun = StoreGetValueDays{}
var _ ISFun = StoreGetValueDays{}
var _ IArgc = StoreGetValueDays{}
var _ IDb = StoreGetValueDays{}

func (StoreGetValueDays) Label() string {
	return "storegetvaluedays"
}
func (StoreGetValueDays) Argc() int {
	return 2
}
func (StoreGetValueDays) DbFunc() {
	return
}
func (StoreGetValueDays) Exec(thedate string, args ...string) string {
	dayInt := GetDayInt(thedate)
	maxDay := ToInt(args[0])
	if maxDay <= 0 {
		return False
	}
	var count int
	for i := dayInt; i > dayInt-maxDay; i-- {
		res := getRowDb(args, 3).Get(ToString(i) + args[1])
		count += ToInt(res)
	}
	return ToString(count)
}
