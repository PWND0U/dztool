package IOTool

import (
	"bytes"
	"strings"
	"testing"
)

func TestFastBufferNew(t *testing.T) {
	buf := NewFastBuffer()
	if buf == nil {
		t.Fatal("NewFastBuffer 返回 nil")
	}
	if buf.blockSize != 4096 {
		t.Fatalf("默认块大小应为 4096，实际为 %d", buf.blockSize)
	}
	if len(buf.blocks) != 1 {
		t.Fatalf("应初始分配 1 个块，实际为 %d", len(buf.blocks))
	}
}

func TestFastBufferNewWithSize(t *testing.T) {
	buf := NewFastBufferWithSize(1024)
	if buf.blockSize != 1024 {
		t.Fatalf("块大小应为 1024，实际为 %d", buf.blockSize)
	}

	buf2 := NewFastBufferWithSize(0)
	if buf2.blockSize != 4096 {
		t.Fatalf("无效块大小应回退为 4096，实际为 %d", buf2.blockSize)
	}

	buf3 := NewFastBufferWithSize(-1)
	if buf3.blockSize != 4096 {
		t.Fatalf("负数块大小应回退为 4096，实际为 %d", buf3.blockSize)
	}
}

func TestFastBufferWriteSmall(t *testing.T) {
	buf := NewFastBufferWithSize(64)
	data := []byte("hello")
	n := buf.Write(data)
	if n != 5 {
		t.Fatalf("写入字节数应为 5，实际为 %d", n)
	}
	if buf.totalLen != 5 {
		t.Fatalf("totalLen 应为 5，实际为 %d", buf.totalLen)
	}
	if buf.Len() != 5 {
		t.Fatalf("Len 应为 5，实际为 %d", buf.Len())
	}
}

func TestFastBufferWriteCrossBlock(t *testing.T) {
	buf := NewFastBufferWithSize(8)
	data := []byte("hello world, this is a cross-block test")
	n := buf.Write(data)
	if n != len(data) {
		t.Fatalf("写入字节数应为 %d，实际为 %d", len(data), n)
	}
	if buf.totalLen != len(data) {
		t.Fatalf("totalLen 应为 %d，实际为 %d", len(data), buf.totalLen)
	}
	if buf.Len() != len(data) {
		t.Fatalf("Len 应为 %d，实际为 %d", len(data), buf.Len())
	}
	if len(buf.blocks) < 5 {
		t.Fatalf("跨块写入应分配多个块，实际为 %d", len(buf.blocks))
	}
}

func TestFastBufferRead(t *testing.T) {
	buf := NewFastBufferWithSize(8)
	buf.Write([]byte("hello world"))

	data := buf.Read(5)
	if string(data) != "hello" {
		t.Fatalf("读取数据应为 'hello'，实际为 '%s'", string(data))
	}
	if buf.readPos != 5 {
		t.Fatalf("readPos 应为 5，实际为 %d", buf.readPos)
	}
	if buf.Len() != 6 {
		t.Fatalf("Len 应为 6，实际为 %d", buf.Len())
	}

	data2 := buf.Read(100)
	if string(data2) != " world" {
		t.Fatalf("读取数据应为 ' world'，实际为 '%s'", string(data2))
	}
	if buf.Len() != 0 {
		t.Fatalf("读取完毕后 Len 应为 0，实际为 %d", buf.Len())
	}
}

func TestFastBufferReadAll(t *testing.T) {
	buf := NewFastBufferWithSize(8)
	buf.Write([]byte("hello world"))

	data := buf.ReadAll()
	if string(data) != "hello world" {
		t.Fatalf("ReadAll 数据应为 'hello world'，实际为 '%s'", string(data))
	}
	if buf.Len() != 0 {
		t.Fatalf("ReadAll 后 Len 应为 0，实际为 %d", buf.Len())
	}
	if buf.totalLen != 0 {
		t.Fatalf("ReadAll 后 totalLen 应为 0，实际为 %d", buf.totalLen)
	}
	if buf.readPos != 0 {
		t.Fatalf("ReadAll 后 readPos 应为 0，实际为 %d", buf.readPos)
	}
}

func TestFastBufferLen(t *testing.T) {
	buf := NewFastBufferWithSize(16)
	if buf.Len() != 0 {
		t.Fatalf("初始 Len 应为 0，实际为 %d", buf.Len())
	}
	buf.Write([]byte("abc"))
	if buf.Len() != 3 {
		t.Fatalf("写入后 Len 应为 3，实际为 %d", buf.Len())
	}
	buf.Read(1)
	if buf.Len() != 2 {
		t.Fatalf("读取 1 字节后 Len 应为 2，实际为 %d", buf.Len())
	}
}

func TestFastBufferReset(t *testing.T) {
	buf := NewFastBufferWithSize(8)
	buf.Write([]byte("hello world"))
	buf.Read(3)

	buf.Reset()
	if buf.Len() != 0 {
		t.Fatalf("Reset 后 Len 应为 0，实际为 %d", buf.Len())
	}
	if buf.totalLen != 0 {
		t.Fatalf("Reset 后 totalLen 应为 0，实际为 %d", buf.totalLen)
	}
	if buf.readPos != 0 {
		t.Fatalf("Reset 后 readPos 应为 0，实际为 %d", buf.readPos)
	}
	if len(buf.blocks) != 1 {
		t.Fatalf("Reset 后应保留 1 个块，实际为 %d", len(buf.blocks))
	}

	buf.Write([]byte("new data"))
	data := buf.ReadAll()
	if string(data) != "new data" {
		t.Fatalf("Reset 后重新写入读取应为 'new data'，实际为 '%s'", string(data))
	}
}

func TestFastBufferReadFrom(t *testing.T) {
	buf := NewFastBufferWithSize(16)
	r := strings.NewReader("hello from reader")

	n, err := buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("ReadFrom 不应返回错误: %v", err)
	}
	if n != int64(len("hello from reader")) {
		t.Fatalf("ReadFrom 读取字节数应为 %d，实际为 %d", len("hello from reader"), n)
	}
	data := buf.ReadAll()
	if string(data) != "hello from reader" {
		t.Fatalf("ReadFrom 数据应为 'hello from reader'，实际为 '%s'", string(data))
	}
}

func TestFastBufferWriteTo(t *testing.T) {
	buf := NewFastBufferWithSize(8)
	buf.Write([]byte("hello to writer"))

	var w bytes.Buffer
	n, err := buf.WriteTo(&w)
	if err != nil {
		t.Fatalf("WriteTo 不应返回错误: %v", err)
	}
	if n != int64(len("hello to writer")) {
		t.Fatalf("WriteTo 写入字节数应为 %d，实际为 %d", len("hello to writer"), n)
	}
	if w.String() != "hello to writer" {
		t.Fatalf("WriteTo 数据应为 'hello to writer'，实际为 '%s'", w.String())
	}
	if buf.Len() != 0 {
		t.Fatalf("WriteTo 后 Len 应为 0，实际为 %d", buf.Len())
	}
}

func TestFastBufferWriteReadCycle(t *testing.T) {
	buf := NewFastBufferWithSize(8)

	buf.Write([]byte("part1"))
	data := buf.Read(5)
	if string(data) != "part1" {
		t.Fatalf("第一次读取应为 'part1'，实际为 '%s'", string(data))
	}

	buf.Write([]byte("part2"))
	data2 := buf.ReadAll()
	if string(data2) != "part2" {
		t.Fatalf("第二次读取应为 'part2'，实际为 '%s'", string(data2))
	}

	buf.Write([]byte("part3"))
	buf.Write([]byte("part4"))
	data3 := buf.ReadAll()
	if string(data3) != "part3part4" {
		t.Fatalf("第三次读取应为 'part3part4'，实际为 '%s'", string(data3))
	}
}
