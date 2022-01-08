package funcs

var StoreGetFuncs = []IFun{}
var StoreSetFuncs = []IFun{}

func AddGetStore(fun IFun) {
	StoreGetFuncs = append(StoreGetFuncs, fun)
}

func AddSetStore(fun IFun) {
	StoreSetFuncs = append(StoreSetFuncs, fun)
}
