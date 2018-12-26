package gotcp

// MSGPACKTYPE : 数据包类型
type MSGPACKTYPE int

const (
	// HeaderAndData : Header、 Data 字段均有效
	HeaderAndData MSGPACKTYPE = iota
	// DataOnly : 只有 Data 字段有效
	DataOnly
	// HeaderNoCmdAndData : Header、 Data 字段均有效；并排除 Header 中的 cmd0 cmd1
	HeaderNoCmdAndData
)

// MsgPack : 数据包格式
type MsgPack struct {
	Header [6]byte     // 元素依次为：len0 len1 len2 flag cmd0 cmd1
	Data   []byte      // data
	Flag   MSGPACKTYPE // 0 Header + Data；1 Data Only；2 Header(but no cmd0 cmd1) + Data
}
