package behavior

type IPOST interface {
	Do([]byte, string) ([]byte, error)
}
