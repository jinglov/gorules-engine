package channel

type Channel struct {
	Name  string
	IChan chan *Input
	SChan chan *Input
	OChan chan *Output
	Pool  InputPool
}

func NewChannel(name string, ilen, slen, olen int) *Channel {
	return &Channel{
		Name:  name,
		IChan: make(chan *Input, ilen),
		SChan: make(chan *Input, slen),
		OChan: make(chan *Output, olen),
		Pool:  NewInputPool(),
	}
}

func NewPool(name string) *Channel {
	return &Channel{
		Name: name,
		Pool: NewInputPool(),
	}
}
