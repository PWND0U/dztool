package IOTool

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteFileAndReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	data := []byte("hello world")

	if err := WriteFile(path, data); err != nil {
		t.Fatalf("WriteFile 失败: %v", err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile 失败: %v", err)
	}
	if string(got) != string(data) {
		t.Fatalf("期望 %q, 实际 %q", string(data), string(got))
	}
}

func TestWriteFileStringAndReadFileAsString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	content := "你好世界"

	if err := WriteFileString(path, content); err != nil {
		t.Fatalf("WriteFileString 失败: %v", err)
	}

	got, err := ReadFileAsString(path)
	if err != nil {
		t.Fatalf("ReadFileAsString 失败: %v", err)
	}
	if got != content {
		t.Fatalf("期望 %q, 实际 %q", content, got)
	}
}

func TestWriteFileSync(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sync.txt")
	data := []byte("sync data")

	if err := WriteFileSync(path, data); err != nil {
		t.Fatalf("WriteFileSync 失败: %v", err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile 失败: %v", err)
	}
	if string(got) != string(data) {
		t.Fatalf("期望 %q, 实际 %q", string(data), string(got))
	}
}

func TestAppendFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "append.txt")

	if err := WriteFile(path, []byte("hello")); err != nil {
		t.Fatalf("WriteFile 失败: %v", err)
	}
	if err := AppendFile(path, []byte(" world")); err != nil {
		t.Fatalf("AppendFile 失败: %v", err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile 失败: %v", err)
	}
	if string(got) != "hello world" {
		t.Fatalf("期望 %q, 实际 %q", "hello world", string(got))
	}
}

func TestAppendFileString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "append.txt")

	if err := WriteFileString(path, "foo"); err != nil {
		t.Fatalf("WriteFileString 失败: %v", err)
	}
	if err := AppendFileString(path, "bar"); err != nil {
		t.Fatalf("AppendFileString 失败: %v", err)
	}

	got, err := ReadFileAsString(path)
	if err != nil {
		t.Fatalf("ReadFileAsString 失败: %v", err)
	}
	if got != "foobar" {
		t.Fatalf("期望 %q, 实际 %q", "foobar", got)
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "sub", "dst.txt")

	if err := WriteFile(src, []byte("copy me")); err != nil {
		t.Fatalf("WriteFile 失败: %v", err)
	}
	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile 失败: %v", err)
	}

	got, err := ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile 失败: %v", err)
	}
	if string(got) != "copy me" {
		t.Fatalf("期望 %q, 实际 %q", "copy me", string(got))
	}
}

func TestMoveFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "sub", "dst.txt")

	if err := WriteFile(src, []byte("move me")); err != nil {
		t.Fatalf("WriteFile 失败: %v", err)
	}
	if err := MoveFile(src, dst); err != nil {
		t.Fatalf("MoveFile 失败: %v", err)
	}

	if FileExists(src) {
		t.Fatal("源文件应该已被移动")
	}
	got, err := ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile 失败: %v", err)
	}
	if string(got) != "move me" {
		t.Fatalf("期望 %q, 实际 %q", "move me", string(got))
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "exist.txt")

	if FileExists(path) {
		t.Fatal("文件不应存在")
	}
	if err := WriteFile(path, []byte("x")); err != nil {
		t.Fatalf("WriteFile 失败: %v", err)
	}
	if !FileExists(path) {
		t.Fatal("文件应该存在")
	}
	if FileExists(dir) {
		t.Fatal("目录不应被判断为文件")
	}
}

func TestDirExists(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")

	if DirExists(sub) {
		t.Fatal("目录不应存在")
	}
	if err := MkdirAll(sub); err != nil {
		t.Fatalf("MkdirAll 失败: %v", err)
	}
	if !DirExists(sub) {
		t.Fatal("目录应该存在")
	}
	if DirExists(filepath.Join(dir, "notexist")) {
		t.Fatal("不存在的路径不应被判断为目录")
	}
}

func TestMkdirAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c")

	if err := MkdirAll(path); err != nil {
		t.Fatalf("MkdirAll 失败: %v", err)
	}
	if !DirExists(path) {
		t.Fatal("目录应该存在")
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "to_remove")
	if err := MkdirAll(filepath.Join(path, "sub")); err != nil {
		t.Fatalf("MkdirAll 失败: %v", err)
	}
	if err := WriteFile(filepath.Join(path, "sub", "file.txt"), []byte("x")); err != nil {
		t.Fatalf("WriteFile 失败: %v", err)
	}
	if err := Remove(path); err != nil {
		t.Fatalf("Remove 失败: %v", err)
	}
	if DirExists(path) {
		t.Fatal("目录应该已被删除")
	}
}

func TestFileSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "size.txt")
	data := []byte("12345")
	if err := WriteFile(path, data); err != nil {
		t.Fatalf("WriteFile 失败: %v", err)
	}
	size, err := FileSize(path)
	if err != nil {
		t.Fatalf("FileSize 失败: %v", err)
	}
	if size != int64(len(data)) {
		t.Fatalf("期望 %d, 实际 %d", len(data), size)
	}
}

func TestWalkDir(t *testing.T) {
	dir := t.TempDir()
	_ = MkdirAll(filepath.Join(dir, "a"))
	_ = WriteFile(filepath.Join(dir, "a", "f1.txt"), []byte("1"))
	_ = WriteFile(filepath.Join(dir, "f2.txt"), []byte("2"))

	var files []string
	err := WalkDir(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir 失败: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("期望 2 个文件, 实际 %d", len(files))
	}
}

func TestListDir(t *testing.T) {
	dir := t.TempDir()
	_ = WriteFile(filepath.Join(dir, "a.txt"), []byte("a"))
	_ = WriteFile(filepath.Join(dir, "b.txt"), []byte("b"))

	infos, err := ListDir(dir)
	if err != nil {
		t.Fatalf("ListDir 失败: %v", err)
	}
	if len(infos) != 2 {
		t.Fatalf("期望 2 个子项, 实际 %d", len(infos))
	}
}

func TestLastModified(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mod.txt")
	_ = WriteFile(path, []byte("x"))

	mod, err := LastModified(path)
	if err != nil {
		t.Fatalf("LastModified 失败: %v", err)
	}
	if mod.IsZero() {
		t.Fatal("修改时间不应为零值")
	}
	if mod.After(time.Now()) {
		t.Fatal("修改时间不应在未来")
	}
}

func TestCreateTempFile(t *testing.T) {
	dir := t.TempDir()
	f, err := CreateTempFile(dir, "tmp_*.txt")
	if err != nil {
		t.Fatalf("CreateTempFile 失败: %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if !FileExists(f.Name()) {
		t.Fatal("临时文件应该存在")
	}
}

func TestCreateTempDir(t *testing.T) {
	dir := t.TempDir()
	tmp, err := CreateTempDir(dir, "tmpdir_*")
	if err != nil {
		t.Fatalf("CreateTempDir 失败: %v", err)
	}
	defer os.RemoveAll(tmp)

	if !DirExists(tmp) {
		t.Fatal("临时目录应该存在")
	}
}
