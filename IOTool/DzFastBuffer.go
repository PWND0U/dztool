package IOTool

import (
	"io"
)

type DzFastBuffer struct {
	blocks    [][]byte // 缓冲块集合
	curBlock  int      // 当前写入块索引
	curOffset int      // 当前块内写入偏移
	blockSize int      // 每个块的大小
	totalLen  int      // 总数据长度
	readPos   int      // 读取位置（全局偏移）
}

// NewFastBuffer 创建默认块大小(4096)的缓冲区
func NewFastBuffer() *DzFastBuffer {
	return NewFastBufferWithSize(4096)
}

// NewFastBufferWithSize 创建指定块大小的缓冲区
func NewFastBufferWithSize(blockSize int) *DzFastBuffer {
	if blockSize <= 0 {
		blockSize = 4096
	}
	buf := &DzFastBuffer{
		blocks:    make([][]byte, 1),
		blockSize: blockSize,
	}
	buf.blocks[0] = make([]byte, blockSize)
	return buf
}

// Write 写入数据到缓冲区，返回写入字节数
func (b *DzFastBuffer) Write(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	written := 0
	remaining := data

	for len(remaining) > 0 {
		avail := b.blockSize - b.curOffset
		if avail == 0 {
			b.appendBlock()
			avail = b.blockSize
		}
		n := avail
		if len(remaining) < n {
			n = len(remaining)
		}
		copy(b.blocks[b.curBlock][b.curOffset:], remaining[:n])
		b.curOffset += n
		b.totalLen += n
		written += n
		remaining = remaining[n:]
	}

	return written
}

// Read 从缓冲区读取最多 n 字节（从 readPos 位置），返回字节切片
func (b *DzFastBuffer) Read(n int) []byte {
	if n <= 0 || b.readPos >= b.totalLen {
		return nil
	}
	remaining := b.totalLen - b.readPos
	if n > remaining {
		n = remaining
	}
	result := make([]byte, 0, n)

	globalPos := b.readPos
	blockIdx := globalPos / b.blockSize
	offset := globalPos % b.blockSize
	left := n

	for left > 0 && blockIdx < len(b.blocks) {
		blockEnd := b.totalLen - blockIdx*b.blockSize
		if blockEnd > b.blockSize {
			blockEnd = b.blockSize
		}
		avail := blockEnd - offset
		if avail <= 0 {
			break
		}
		readLen := avail
		if readLen > left {
			readLen = left
		}
		result = append(result, b.blocks[blockIdx][offset:offset+readLen]...)
		left -= readLen
		blockIdx++
		offset = 0
	}

	b.readPos += n
	return result
}

// ReadAll 读取缓冲区全部数据（从 readPos 到末尾），缓冲区被清空（readPos 重置）
func (b *DzFastBuffer) ReadAll() []byte {
	data := b.Read(b.totalLen - b.readPos)
	b.blocks = b.blocks[:1]
	b.blocks[0] = make([]byte, b.blockSize)
	b.curBlock = 0
	b.curOffset = 0
	b.totalLen = 0
	b.readPos = 0
	return data
}

// Len 返回缓冲区当前未读数据长度
func (b *DzFastBuffer) Len() int {
	return b.totalLen - b.readPos
}

// Reset 清空缓冲区数据，保留已分配的内存
func (b *DzFastBuffer) Reset() {
	for i := range b.blocks {
		b.blocks[i] = b.blocks[i][:b.blockSize]
	}
	b.blocks = b.blocks[:1]
	b.curBlock = 0
	b.curOffset = 0
	b.totalLen = 0
	b.readPos = 0
}

// ReadFrom 从 Reader 读取数据到缓冲区
func (b *DzFastBuffer) ReadFrom(r io.Reader) (int64, error) {
	var total int64
	buf := make([]byte, b.blockSize)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			b.Write(buf[:n])
			total += int64(n)
		}
		if err != nil {
			if err == io.EOF {
				return total, nil
			}
			return total, err
		}
	}
}

// WriteTo 将缓冲区数据写入 Writer
func (b *DzFastBuffer) WriteTo(w io.Writer) (int64, error) {
	var total int64
	globalPos := b.readPos
	blockIdx := globalPos / b.blockSize
	offset := globalPos % b.blockSize
	remaining := b.totalLen - b.readPos

	for remaining > 0 && blockIdx < len(b.blocks) {
		blockEnd := b.totalLen - blockIdx*b.blockSize
		if blockEnd > b.blockSize {
			blockEnd = b.blockSize
		}
		avail := blockEnd - offset
		if avail <= 0 {
			break
		}
		writeLen := avail
		if writeLen > remaining {
			writeLen = remaining
		}
		n, err := w.Write(b.blocks[blockIdx][offset : offset+writeLen])
		total += int64(n)
		remaining -= writeLen
		if err != nil {
			b.readPos = b.totalLen - remaining
			return total, err
		}
		blockIdx++
		offset = 0
	}

	b.readPos = b.totalLen
	return total, nil
}

func (b *DzFastBuffer) appendBlock() {
	newBlock := make([]byte, b.blockSize)
	b.blocks = append(b.blocks, newBlock)
	b.curBlock++
	b.curOffset = 0
}
