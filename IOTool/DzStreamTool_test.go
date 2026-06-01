package IOTool

import (
	"bytes"
	"strings"
	"testing"
)

func TestDzStreamReadAll(t *testing.T) {
	data := []byte("hello world")
	reader := bytes.NewReader(data)
	result, err := ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll 失败: %v", err)
	}
	if string(result) != "hello world" {
		t.Fatalf("ReadAll 结果不匹配，期望 hello world，实际 %s", string(result))
	}
}

func TestDzStreamReadAllString(t *testing.T) {
	reader := strings.NewReader("你好世界")
	result, err := ReadAllString(reader)
	if err != nil {
		t.Fatalf("ReadAllString 失败: %v", err)
	}
	if result != "你好世界" {
		t.Fatalf("ReadAllString 结果不匹配，期望 你好世界，实际 %s", result)
	}
}

func TestDzStreamWriteAll(t *testing.T) {
	var buf bytes.Buffer
	data := []byte("write test")
	n, err := WriteAll(&buf, data)
	if err != nil {
		t.Fatalf("WriteAll 失败: %v", err)
	}
	if n != len(data) {
		t.Fatalf("WriteAll 写入字节数不匹配，期望 %d，实际 %d", len(data), n)
	}
	if buf.String() != "write test" {
		t.Fatalf("WriteAll 内容不匹配，期望 write test，实际 %s", buf.String())
	}
}

func TestDzStreamWriteString(t *testing.T) {
	var buf bytes.Buffer
	n, err := WriteString(&buf, "字符串写入")
	if err != nil {
		t.Fatalf("WriteString 失败: %v", err)
	}
	if buf.String() != "字符串写入" {
		t.Fatalf("WriteString 内容不匹配，期望 字符串写入，实际 %s", buf.String())
	}
	_ = n
}

func TestDzStreamReadLines(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/lines.txt"
	content := "第一行\n第二行\n第三行"
	if err := WriteFileString(path, content); err != nil {
		t.Fatalf("WriteFileString 失败: %v", err)
	}
	lines, err := ReadLines(path)
	if err != nil {
		t.Fatalf("ReadLines 失败: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("ReadLines 行数不匹配，期望 3，实际 %d", len(lines))
	}
	if lines[0] != "第一行" || lines[1] != "第二行" || lines[2] != "第三行" {
		t.Fatalf("ReadLines 内容不匹配，实际 %v", lines)
	}
}

func TestDzStreamCopy(t *testing.T) {
	src := strings.NewReader("copy content")
	var dst bytes.Buffer
	n, err := Copy(&dst, src)
	if err != nil {
		t.Fatalf("Copy 失败: %v", err)
	}
	if n != int64(len("copy content")) {
		t.Fatalf("Copy 字节数不匹配，期望 %d，实际 %d", len("copy content"), n)
	}
	if dst.String() != "copy content" {
		t.Fatalf("Copy 内容不匹配，期望 copy content，实际 %s", dst.String())
	}
}

func TestDzStreamCopyBuffer(t *testing.T) {
	src := strings.NewReader("buffer copy")
	var dst bytes.Buffer
	buf := make([]byte, 4)
	n, err := CopyBuffer(&dst, src, buf)
	if err != nil {
		t.Fatalf("CopyBuffer 失败: %v", err)
	}
	if dst.String() != "buffer copy" {
		t.Fatalf("CopyBuffer 内容不匹配，期望 buffer copy，实际 %s", dst.String())
	}
	_ = n
}

func TestDzStreamMultiReader(t *testing.T) {
	r1 := strings.NewReader("hello ")
	r2 := strings.NewReader("world")
	mr := MultiReader(r1, r2)
	result, err := ReadAll(mr)
	if err != nil {
		t.Fatalf("MultiReader 读取失败: %v", err)
	}
	if string(result) != "hello world" {
		t.Fatalf("MultiReader 结果不匹配，期望 hello world，实际 %s", string(result))
	}
}

func TestDzStreamMultiWriter(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	mw := MultiWriter(&buf1, &buf2)
	_, err := WriteString(mw, "multi")
	if err != nil {
		t.Fatalf("MultiWriter 写入失败: %v", err)
	}
	if buf1.String() != "multi" || buf2.String() != "multi" {
		t.Fatalf("MultiWriter 结果不匹配，buf1=%s, buf2=%s", buf1.String(), buf2.String())
	}
}

func TestDzStreamLimitReader(t *testing.T) {
	reader := strings.NewReader("hello world")
	lr := LimitReader(reader, 5)
	result, err := ReadAll(lr)
	if err != nil {
		t.Fatalf("LimitReader 读取失败: %v", err)
	}
	if string(result) != "hello" {
		t.Fatalf("LimitReader 结果不匹配，期望 hello，实际 %s", string(result))
	}
}

func TestDzStreamReadAt(t *testing.T) {
	data := []byte("hello world")
	reader := bytes.NewReader(data)
	result, err := ReadAt(reader, 6, 5)
	if err != nil {
		t.Fatalf("ReadAt 失败: %v", err)
	}
	if string(result) != "world" {
		t.Fatalf("ReadAt 结果不匹配，期望 world，实际 %s", string(result))
	}
}

func TestDzStreamReadLinesFromReader(t *testing.T) {
	reader := strings.NewReader("第一行\n第二行\n第三行")
	lines, err := ReadLinesFromReader(reader)
	if err != nil {
		t.Fatalf("ReadLinesFromReader 失败: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("ReadLinesFromReader 行数不匹配，期望 3，实际 %d", len(lines))
	}
	if lines[0] != "第一行" || lines[1] != "第二行" || lines[2] != "第三行" {
		t.Fatalf("ReadLinesFromReader 内容不匹配，实际 %v", lines)
	}
}

func TestDzStreamPipe(t *testing.T) {
	src := strings.NewReader("pipe data")
	var dst bytes.Buffer
	n, err := Pipe(src, &dst)
	if err != nil {
		t.Fatalf("Pipe 失败: %v", err)
	}
	if dst.String() != "pipe data" {
		t.Fatalf("Pipe 内容不匹配，期望 pipe data，实际 %s", dst.String())
	}
	_ = n
}
