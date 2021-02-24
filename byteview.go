package bottle

type ByteView struct {
	b []byte
}

func (bv ByteView) Len() int {
	return len(bv.b)
}

func (bv ByteView) toString() string {
	return string(bv.b)
}

func (bv ByteView) cloneBytes() []byte {
	s := make([]byte, len(bv.b))
	copy(s, bv.b)
	return s
}


