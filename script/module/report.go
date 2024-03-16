package script

// Report 网络层在处理消息时的信息
//
//	因为流式消息无法保证 req/res 配对，因此不再记录消费时间
type Report struct {
	MsgID string
	Err   string
}
