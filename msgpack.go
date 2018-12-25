package gotcp

// MsgPack : 数据包格式
type MsgPack struct {
	len0 byte
	len1 byte
	len2 byte
	flag byte
	cmd0 byte
	cmd1 byte
	data []byte
}
