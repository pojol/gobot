package behavior

type APIJump struct {
}

func (aj *APIJump) Prev() int {
	return 0
}

type IndexJump struct {
}

func (ij *IndexJump) Prev() int {
	return 0
}
