package behavior

type Blackboard struct {
	nods []INod
}

func (b *Blackboard) GetOpenNods() []INod {
	return b.nods
}

func (b *Blackboard) Append(nods []INod) {
	b.nods = append(b.nods, nods...)
}

func (b *Blackboard) Reset() {
	b.nods = b.nods[:0]
}
