package funcs

func init() {
	AddData(GetSrc{})
	AddData(SetSrc{})
}

var DataFuncs = []IFun{}

func AddData(fun IFun) {
	DataFuncs = append(DataFuncs, fun)
}

// GetSrc
// 从当前数据中取内容,如果前面的key为空，则依次从后面的key中取。直到取出非空值
// getsrc('cid') == 'midu'
// 这是一个高性能的方法
type GetSrc struct{}

var _ IFun = GetSrc{}
var _ IMFun = GetSrc{}
var _ IHeight = GetSrc{}
var _ IArgcMin = GetSrc{}

func (GetSrc) Height() {}
func (GetSrc) Label() string {
	return "getsrc"
}
func (GetSrc) ArgcMin() int {
	return 1
}
func (g GetSrc) Exec(m map[string]string, args ...string) string {
	for _, k := range args {
		if v, ok := m[k]; ok && v != "" {
			return v
		}
	}
	return ""
}

// GetSrc
// 设置当前数据某个key的值
// setsrc('_spam' , '1') == '1'
// 这是一个高性能的方法
// 限制只能有2个参数
type SetSrc struct{}

var _ IFun = SetSrc{}
var _ IMFun = SetSrc{}
var _ IHeight = SetSrc{}
var _ IArgc = SetSrc{}

func (SetSrc) Height() {}
func (SetSrc) Label() string {
	return "setsrc"
}
func (SetSrc) Argc() int {
	return 2
}
func (g SetSrc) Exec(m map[string]string, args ...string) string {
	if len(args) < 2 {
		return False
	}
	m[args[0]] = args[1]
	return True
}
