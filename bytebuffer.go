package gotcp

const (
	presize  = 0
	initsize = 16
)

// ByteBuffer : 缓冲区类
type ByteBuffer struct {
	_buffer      []byte
	_prependSize int
	_readerIndex int
	_writerIndex int
}

// NewByteBuffer : 缓冲区类的构造函数
func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{
		_buffer:      make([]byte, presize+initsize),
		_prependSize: presize,
		_readerIndex: presize,
		_writerIndex: presize,
	}
}

// Append : 向缓冲区追加数据
func (bytebuffer *ByteBuffer) Append(buff []byte) {
	size := len(buff)
	if size == 0 {
		return
	}
	bytebuffer.WrGrow(size)
	copy(bytebuffer._buffer[bytebuffer._writerIndex:], buff)
	bytebuffer.WrFlip(size)
}

// WrBuf : 获取缓冲区可写入区域
func (bytebuffer *ByteBuffer) WrBuf() []byte {
	if bytebuffer._writerIndex >= len(bytebuffer._buffer) {
		return nil
	}
	return bytebuffer._buffer[bytebuffer._writerIndex:]
}

// WrSize : 获取缓冲区可写入区域大小
func (bytebuffer *ByteBuffer) WrSize() int {
	return len(bytebuffer._buffer) - bytebuffer._writerIndex
}

// WrFlip : 完成写入 size 字节数据
func (bytebuffer *ByteBuffer) WrFlip(size int) {
	bytebuffer._writerIndex += size
}

// WrGrow : 如果缓冲区可写入区域大小小于 size ，扩大缓冲区
func (bytebuffer *ByteBuffer) WrGrow(size int) {
	if size > bytebuffer.WrSize() {
		bytebuffer.wrreserve(size)
	}
}

// RdBuf :  获取缓冲区可读取区域
func (bytebuffer *ByteBuffer) RdBuf() []byte {
	if bytebuffer._readerIndex >= len(bytebuffer._buffer) {
		return nil
	}
	return bytebuffer._buffer[bytebuffer._readerIndex:]
}

// RdReady : 是否有数据可读
func (bytebuffer *ByteBuffer) RdReady() bool {
	return bytebuffer._writerIndex > bytebuffer._readerIndex
}

// RdSize : 获取缓冲区可写入区域大小
func (bytebuffer *ByteBuffer) RdSize() int {
	return bytebuffer._writerIndex - bytebuffer._readerIndex
}

// RdFlip : 完成读取 size 字节数据
func (bytebuffer *ByteBuffer) RdFlip(size int) {
	if size < bytebuffer.RdSize() {
		bytebuffer._readerIndex += size
	} else {
		bytebuffer.Reset()
	}
}

// Reset : 重置缓冲区
func (bytebuffer *ByteBuffer) Reset() {
	bytebuffer._readerIndex = bytebuffer._prependSize
	bytebuffer._writerIndex = bytebuffer._prependSize
}

func (bytebuffer *ByteBuffer) wrreserve(size int) {
	if bytebuffer.WrSize()+bytebuffer._readerIndex < size+bytebuffer._prependSize {
		newsize := bytebuffer.RdSize() + bytebuffer.WrSize()
		for newsize < bytebuffer._writerIndex+size {
			newsize <<= 1
		}
		tmpbuff := make([]byte, newsize+bytebuffer._prependSize)
		copy(tmpbuff, bytebuffer._buffer)
		bytebuffer._buffer = tmpbuff
	} else {
		readable := bytebuffer.RdSize()
		copy(bytebuffer._buffer[bytebuffer._prependSize:], bytebuffer._buffer[bytebuffer._readerIndex:bytebuffer._writerIndex])
		bytebuffer._readerIndex = bytebuffer._prependSize
		bytebuffer._writerIndex = bytebuffer._readerIndex + readable
	}
}
