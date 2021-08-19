package marshal

type Options struct {
}

type Option func(*Options)

type IMarshal interface {
	// 对传入的数据进行编码
	Marshal(interface{}) ([]byte, error)
}
