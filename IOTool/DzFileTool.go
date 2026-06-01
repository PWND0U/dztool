package IOTool

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ReadFile 读取文件全部内容为字节切片
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFileAsString 读取文件全部内容为字符串
func ReadFileAsString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadLines 按行读取文件
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteFile 将字节切片写入文件（覆盖），若文件不存在则创建
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0666)
}

// WriteFileString 将字符串写入文件（覆盖）
func WriteFileString(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0666)
}

// WriteFileSync 写入文件并 Sync 刷盘
func WriteFileSync(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return f.Sync()
}

// AppendFile 追加字节到文件
func AppendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AppendFileString 追加字符串到文件
func AppendFileString(path string, content string) error {
	return AppendFile(path, []byte(content))
}

// CopyFile 复制文件，自动创建目标目录
func CopyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
		return err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// MoveFile 移动文件，自动创建目标目录
func MoveFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
		return err
	}
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

// FileExists 判断文件是否存在
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists 判断目录是否存在
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// MkdirAll 递归创建目录
func MkdirAll(path string) error {
	return os.MkdirAll(path, 0777)
}

// Remove 删除文件或目录（递归）
func Remove(path string) error {
	return os.RemoveAll(path)
}

// FileSize 获取文件大小
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// WalkDir 遍历目录下所有文件
func WalkDir(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}

// ListDir 列出目录下直接子项
func ListDir(path string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	infos := make([]os.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// LastModified 获取文件最后修改时间
func LastModified(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// CreateTempFile 创建临时文件
func CreateTempFile(dir, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

// CreateTempDir 创建临时目录
func CreateTempDir(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}
