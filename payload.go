package tinybox

type Payload struct {
	buf []byte
}

func (p Payload) Len() int {
	return len(p.buf)
}

func (p Payload) Bytes() []byte {
	rst := make([]byte, len(p.buf))
	copy(rst, p.buf)
	return rst
}
