package behavior

type POST interface {
	Exec() error
}

type Jump struct {
}

type Delay struct {
}
