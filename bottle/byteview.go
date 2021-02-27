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

func (bv ByteView) Clone() []byte {
	return cloneBytes(bv.b)
}

func cloneBytes(b []byte) []byte {
	s := make([]byte, len(b))
	copy(s, b)
	return s
}


