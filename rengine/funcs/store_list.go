package funcs

import (
	"strings"

	"github.com/jinglov/gorules-engine/db"
)

func init() {
	AddSetStore(StoreAppendList{})
	AddSetStore(StoreAppendListExp{})
	AddSetStore(StoreAppendListDays{})
	AddGetStore(StoreGetList{})
	AddGetStore(StoreGetListExp{})
	AddGetStore(StoreGetListDays{})
	AddGetStore(StoreGetListLen{})
	AddGetStore(StoreGetListLenExp{})
	AddGetStore(StoreGetListLenDays{})
	AddGetStore(StoreGetListDelimit{})
	AddGetStore(StoreGetListDelimitExp{})
	AddGetStore(StoreGetListDelimitDays{})
	AddGetStore(StoreGetReasonList{})
}

type StoreAppendList struct{}

func getListDb(args []string, index int) db.DB {
	if len(args) == index {
		argDb := db.GetDb(args[index-1])
		if argDb != nil {
			return argDb
		}
	}
	return db.ListDefault
}

var _ IFun = StoreAppendList{}
var _ ISFun = StoreAppendList{}
var _ IArgc = StoreAppendList{}
var _ IDb = StoreAppendList{}

func (StoreAppendList) Label() string {
	return "storeappendlist"
}
func (StoreAppendList) Argc() int {
	return 2
}
func (StoreAppendList) DbFunc() {
	return
}
func (StoreAppendList) Exec(thedate string, args ...string) string {
	return ToString(getListDb(args, 3).AppendList(thedate+args[0], args[1], DAY))
}

type StoreGetList struct{}

var _ IFun = StoreGetList{}
var _ ISFun = StoreGetList{}
var _ IArgc = StoreGetList{}
var _ IDb = StoreGetList{}

func (StoreGetList) Label() string {
	return "storegetlist"
}
func (StoreGetList) Argc() int {
	return 1
}
func (StoreGetList) DbFunc() {
	return
}
func (StoreGetList) Exec(thedate string, args ...string) string {
	res := getListDb(args, 2).GetList(thedate + args[0])
	return strings.Join(res, DELIMIT)
}

type StoreGetListLen struct{}

var _ IFun = StoreGetListLen{}
var _ ISFun = StoreGetListLen{}
var _ IArgc = StoreGetListLen{}
var _ IDb = StoreGetListLen{}

func (StoreGetListLen) Label() string {
	return "storegetlistlen"
}
func (StoreGetListLen) Argc() int {
	return 1
}
func (StoreGetListLen) DbFunc() {
	return
}
func (StoreGetListLen) Exec(thedate string, args ...string) string {
	return ToString(getListDb(args, 2).GetListLen(thedate + args[0]))
}

type StoreGetListDelimit struct{}

var _ IFun = StoreGetListDelimit{}
var _ ISFun = StoreGetListDelimit{}
var _ IArgc = StoreGetListDelimit{}
var _ IDb = StoreGetListDelimit{}

func (StoreGetListDelimit) Label() string {
	return "storegetlistdelimit"
}
func (StoreGetListDelimit) Argc() int {
	return 2
}
func (StoreGetListDelimit) DbFunc() {
	return
}
func (StoreGetListDelimit) Exec(thedate string, args ...string) string {
	res := getListDb(args, 3).GetList(thedate + args[0])
	return strings.Join(res, args[1])
}

type StoreAppendListExp struct{}

var _ IFun = StoreAppendListExp{}
var _ IDFun = StoreAppendListExp{}
var _ IArgc = StoreAppendListExp{}

func (StoreAppendListExp) Label() string {
	return "storeappendlistexp"
}
func (StoreAppendListExp) Argc() int {
	return 3
}
func (StoreAppendListExp) DbFunc() {
	return
}
func (StoreAppendListExp) Exec(args ...string) string {
	expire := ParseExp(args[0])
	return ToString(getListDb(args, 4).AppendList(args[1], args[2], expire))
}

type StoreGetListExp struct{}

var _ IFun = StoreGetListExp{}
var _ IDFun = StoreGetListExp{}
var _ IArgc = StoreGetListExp{}
var _ IDb = StoreGetListExp{}

func (StoreGetListExp) Label() string {
	return "storegetlistexp"
}
func (StoreGetListExp) Argc() int {
	return 1
}
func (StoreGetListExp) DbFunc() {
	return
}
func (StoreGetListExp) Exec(args ...string) string {
	res := getListDb(args, 2).GetList(args[0])
	return strings.Join(res, DELIMIT)
}

type StoreGetListLenExp struct{}

var _ IFun = StoreGetListLenExp{}
var _ IDFun = StoreGetListLenExp{}
var _ IArgc = StoreGetListLenExp{}
var _ IDb = StoreGetListLenExp{}

func (StoreGetListLenExp) Label() string {
	return "storegetlistlenexp"
}
func (StoreGetListLenExp) Argc() int {
	return 1
}
func (StoreGetListLenExp) DbFunc() {
	return
}
func (StoreGetListLenExp) Exec(args ...string) string {
	return ToString(getListDb(args, 2).GetListLen(args[0]))
}

type StoreGetListDelimitExp struct{}

var _ IFun = StoreGetListDelimitExp{}
var _ IDFun = StoreGetListDelimitExp{}
var _ IArgc = StoreGetListDelimitExp{}
var _ IDb = StoreGetListDelimitExp{}

func (StoreGetListDelimitExp) Label() string {
	return "storegetlistdelimitexp"
}
func (StoreGetListDelimitExp) Argc() int {
	return 2
}
func (StoreGetListDelimitExp) DbFunc() {
	return
}
func (StoreGetListDelimitExp) Exec(args ...string) string {
	res := getListDb(args, 3).GetList(args[0])
	return strings.Join(res, args[1])
}

type StoreAppendListDays struct{}

var _ IFun = StoreAppendListDays{}
var _ ISFun = StoreAppendListDays{}
var _ IArgc = StoreAppendListDays{}
var _ IDb = StoreAppendListDays{}
var storeGetListLenDays = StoreGetListLenDays{}

func (StoreAppendListDays) Label() string {
	return "storeappendlistdays"
}
func (StoreAppendListDays) Argc() int {
	return 3
}
func (StoreAppendListDays) DbFunc() {
	return
}
func (StoreAppendListDays) Exec(thedate string, args ...string) string {
	dayInt := ToString(GetDayInt(thedate))
	maxDay := ToInt(args[0])
	if maxDay <= 0 {
		return False
	}
	getListDb(args, 4).AppendList(dayInt+args[1], args[2], maxDay*DAY)
	if len(args) == 4 {
		return storeGetListLenDays.Exec(thedate, args[0], args[1], args[3])
	} else {
		return storeGetListLenDays.Exec(thedate, args[0], args[1])
	}
}

type StoreGetListDays struct{}

var _ IFun = StoreGetListDays{}
var _ ISFun = StoreGetListDays{}
var _ IArgc = StoreGetListDays{}
var _ IDb = StoreGetListDays{}

func (StoreGetListDays) Label() string {
	return "storegetlistdays"
}
func (StoreGetListDays) Argc() int {
	return 2
}
func (StoreGetListDays) DbFunc() {
	return
}
func (StoreGetListDays) Exec(thedate string, args ...string) string {
	dayInt := GetDayInt(thedate)
	maxDay := ToInt(args[0])
	if maxDay <= 0 {
		return False
	}
	rows := make([]string, 0, 0)
	rowMap := make(map[string]struct{})
	for i := dayInt; i > dayInt-maxDay; i-- {
		res := getListDb(args, 3).GetList(ToString(i) + args[1])
		for _, r := range res {
			if _, ok := rowMap[r]; !ok {
				rows = append(rows, r)
				rowMap[r] = struct{}{}
			}
		}
	}
	if len(rows) == 0 {
		return False
	}
	return strings.Join(rows, DELIMIT)
}

type StoreGetListLenDays struct{}

var _ IFun = StoreGetListLenDays{}
var _ ISFun = StoreGetListLenDays{}
var _ IArgc = StoreGetListLenDays{}
var _ IDb = StoreGetListLenDays{}

func (StoreGetListLenDays) Label() string {
	return "storegetlistlendays"
}
func (StoreGetListLenDays) Argc() int {
	return 2
}
func (StoreGetListLenDays) DbFunc() {
	return
}
func (StoreGetListLenDays) Exec(thedate string, args ...string) string {
	dayInt := GetDayInt(thedate)
	maxDay := ToInt(args[0])
	if maxDay <= 0 {
		return False
	}
	var count int
	rowMap := make(map[string]struct{})
	for i := dayInt; i > dayInt-maxDay; i-- {
		res := getListDb(args, 3).GetList(ToString(i) + args[1])
		for _, r := range res {
			if _, ok := rowMap[r]; !ok {
				rowMap[r] = struct{}{}
				count++
			}
		}
	}
	return ToString(count)
}

type StoreGetListDelimitDays struct{}

var _ IFun = StoreGetListDelimitDays{}
var _ ISFun = StoreGetListDelimitDays{}
var _ IArgc = StoreGetListDelimitDays{}
var _ IDb = StoreGetListDelimitDays{}

func (StoreGetListDelimitDays) Label() string {
	return "storegetlistdelimitdays"
}
func (StoreGetListDelimitDays) Argc() int {
	return 3
}
func (StoreGetListDelimitDays) DbFunc() {
	return
}
func (StoreGetListDelimitDays) Exec(thedate string, args ...string) string {
	dayInt := GetDayInt(thedate)
	maxDay := ToInt(args[0])
	if maxDay <= 0 {
		return False
	}
	rows := make([]string, 0, 0)
	rowMap := make(map[string]struct{})
	for i := dayInt; i > dayInt-maxDay; i-- {
		res := getListDb(args, 4).GetList(ToString(i) + args[1])
		for _, r := range res {
			if _, ok := rowMap[r]; !ok {
				rows = append(rows, r)
				rowMap[r] = struct{}{}
			}
		}
	}
	if len(rows) == 0 {
		return False
	}
	return strings.Join(rows, args[2])
}

type StoreGetReasonList struct{}

var _ IFun = StoreGetReasonList{}
var _ ISFun = StoreGetReasonList{}
var _ IArgc = StoreGetReasonList{}
var _ IDb = StoreGetReasonList{}

func (StoreGetReasonList) Label() string {
	return "storegetreasonlist"
}
func (StoreGetReasonList) Argc() int {
	return 3
}
func (StoreGetReasonList) DbFunc() {
	return
}
func (StoreGetReasonList) Exec(thedate string, args ...string) string {
	res := getListDb(args, 4).GetReasonList(thedate, args[0], args[1], args[2])
	return strings.Join(res, DELIMIT)
}
