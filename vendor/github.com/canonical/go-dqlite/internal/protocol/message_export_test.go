package protocol

func (m *Message) Body1() ([]byte, int) {
	return m.body1.Bytes, m.body1.Offset
}

func (m *Message) Rewind() {
	m.body1.Offset = 0
	m.body2.Offset = 0
}
