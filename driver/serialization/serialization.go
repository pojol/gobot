package serialization

type Options struct {
}

type Option func(*Options)

type ISerialization interface {
	// 编码
	Marshal(interface{}) ([]byte, error)

	// 解码
	Unmarshal([]byte, interface{}) error
}
