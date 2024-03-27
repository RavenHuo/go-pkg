package redis

type Member struct {
	T      string
	M      []byte
	Offset int64 // 当前轮次的index
	ts     int64 // 当前轮次的时间戳
}

func (p *Member) Topic() string {
	return p.T
}

func (p *Member) Message() []byte {
	return p.M
}

func (p *Member) Ack() error {
	return nil
}

func (p *Member) Sequence() int64 {
	return p.ts + p.Offset
}
