package IOTool

import (
	"bufio"
	"bytes"
	"io"
)

// ReadAll 从 Reader 读取全部内容为字节切片
func ReadAll(reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}

// ReadAllString 从 Reader 读取全部内容为字符串
func ReadAllString(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteAll 将字节写入 Writer
func WriteAll(writer io.Writer, data []byte) (int, error) {
	return writer.Write(data)
}

// WriteString 将字符串写入 Writer
func WriteString(writer io.Writer, s string) (int, error) {
	return io.WriteString(writer, s)
}

// Copy 流到流的拷贝
func Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// CopyBuffer 带缓冲区的流拷贝
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	return io.CopyBuffer(dst, src, buf)
}

// MultiReader 合并多个 Reader
func MultiReader(readers ...io.Reader) io.Reader {
	return io.MultiReader(readers...)
}

// MultiWriter 合并多个 Writer
func MultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}

// LimitReader 限制读取字节数
func LimitReader(reader io.Reader, n int64) io.Reader {
	return io.LimitReader(reader, n)
}

// ReadAt 带偏移读取
func ReadAt(reader io.ReaderAt, off int64, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := reader.ReadAt(buf, off)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadToBytes 将 Reader 内容写入字节缓冲
func ReadToBytes(reader io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ReadLinesFromReader 从 Reader 按行读取
func ReadLinesFromReader(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Pipe 管道操作
func Pipe(src io.Reader, dst io.Writer) (int64, error) {
	return io.Copy(dst, src)
}
