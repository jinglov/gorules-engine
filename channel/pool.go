package channel

import "sync"

type InputPool struct {
	sync.Pool
}

func NewInputPool() InputPool {
	return InputPool{
		sync.Pool{New: func() interface{} {
			ipt := new(Input)
			ipt.Data = make(map[string]string)
			return ipt
		}}}
}

func (mp *InputPool) Get() *Input {
	return mp.Pool.Get().(*Input).reset()
}

func (mp *InputPool) Put(m *Input) {
	mp.Pool.Put(m)
}
