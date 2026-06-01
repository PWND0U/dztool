package IOTool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBaseName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{filepath.Join("foo", "bar", "test.txt"), "test"},
		{filepath.Join("foo", "bar") + string(os.PathSeparator), "bar"},
		{"test.txt", "test"},
		{filepath.Join("foo", "bar", "noext"), "noext"},
	}
	for _, tt := range tests {
		got := BaseName(tt.input)
		if got != tt.want {
			t.Errorf("BaseName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExt(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{filepath.Join("foo", "bar", "test.txt"), ".txt"},
		{filepath.Join("foo", "bar", "test.tar.gz"), ".gz"},
		{"noext", ""},
	}
	for _, tt := range tests {
		got := Ext(tt.input)
		if got != tt.want {
			t.Errorf("Ext(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFileName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{filepath.Join("foo", "bar", "test.txt"), "test.txt"},
		{filepath.Join("foo", "bar") + string(os.PathSeparator), "bar"},
		{"test.txt", "test.txt"},
	}
	for _, tt := range tests {
		got := FileName(tt.input)
		if got != tt.want {
			t.Errorf("FileName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{filepath.Join("foo", "bar", "test.txt"), filepath.Join("foo", "bar")},
		{filepath.Join("foo", "bar"), "foo"},
		{"test.txt", "."},
	}
	for _, tt := range tests {
		got := Parent(tt.input)
		if got != tt.want {
			t.Errorf("Parent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestJoin(t *testing.T) {
	got := Join("foo", "bar", "test.txt")
	want := filepath.Join("foo", "bar", "test.txt")
	if got != want {
		t.Errorf("Join() = %q, want %q", got, want)
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input string
		check func(string) bool
	}{
		{filepath.Join("foo", "..", "bar"), func(got string) bool {
			return filepath.IsAbs(got) && !strings.Contains(got, "..")
		}},
		{filepath.Join(".", "foo", "bar"), func(got string) bool {
			return filepath.IsAbs(got) && !strings.Contains(got, ".")
		}},
	}
	for _, tt := range tests {
		got, err := Normalize(tt.input)
		if err != nil {
			t.Errorf("Normalize(%q) error: %v", tt.input, err)
			continue
		}
		if !tt.check(got) {
			t.Errorf("Normalize(%q) = %q, check failed", tt.input, got)
		}
	}
}

func TestRel(t *testing.T) {
	base := filepath.Join("foo", "bar")
	target := filepath.Join("foo", "bar", "baz", "test.txt")
	rel, err := Rel(base, target)
	if err != nil {
		t.Fatalf("Rel() error: %v", err)
	}
	want := filepath.Join("baz", "test.txt")
	if rel != want {
		t.Errorf("Rel(%q, %q) = %q, want %q", base, target, rel, want)
	}
}

func TestIsAbs(t *testing.T) {
	if IsAbs("foo/bar") {
		t.Error("IsAbs(\"foo/bar\") = true, want false")
	}
	if !IsAbs("C:" + string(os.PathSeparator) + "foo" + string(os.PathSeparator) + "bar") {
		t.Error("IsAbs(\"C:\\foo\\bar\") = false, want true")
	}
}

func TestDepth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{filepath.Join("foo", "bar"), 2},
		{filepath.Join("foo", "bar", "baz"), 3},
		{".", 0},
	}
	for _, tt := range tests {
		got := Depth(tt.input)
		if got != tt.want {
			t.Errorf("Depth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestPathIterator(t *testing.T) {
	path := filepath.Join("foo", "bar", "test")
	result := PathIterator(path)
	if len(result) < 2 {
		t.Fatalf("PathIterator(%q) 返回结果过短: %v", path, result)
	}
	root := result[0]
	if !filepath.IsAbs(root) {
		t.Errorf("PathIterator 首个元素应为根目录， got %q", root)
	}
	last := result[len(result)-1]
	abs, _ := filepath.Abs(path)
	if filepath.Clean(last) != filepath.Clean(abs) {
		t.Errorf("PathIterator 最后元素 = %q, want %q", last, abs)
	}
}

func TestReplaceExt(t *testing.T) {
	tests := []struct {
		path   string
		newExt string
		want   string
	}{
		{filepath.Join("foo", "test.txt"), ".md", filepath.Join("foo", "test.md")},
		{"noext", ".go", "noext.go"},
		{filepath.Join("foo", "test.tar.gz"), ".zip", filepath.Join("foo", "test.tar.zip")},
	}
	for _, tt := range tests {
		got := ReplaceExt(tt.path, tt.newExt)
		if got != tt.want {
			t.Errorf("ReplaceExt(%q, %q) = %q, want %q", tt.path, tt.newExt, got, tt.want)
		}
	}
}

func TestIsSubPath(t *testing.T) {
	parent := filepath.Join("foo", "bar")
	child := filepath.Join("foo", "bar", "baz", "test.txt")
	if !IsSubPath(parent, child) {
		t.Errorf("IsSubPath(%q, %q) = false, want true", parent, child)
	}

	sibling := filepath.Join("foo", "other")
	if IsSubPath(parent, sibling) {
		t.Errorf("IsSubPath(%q, %q) = true, want false", parent, sibling)
	}
}

func TestIsSafePath(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{filepath.Join("foo", "bar", "test.txt"), true},
		{"foo" + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "bar", false},
		{"safe" + string(os.PathSeparator) + "path", true},
	}
	for _, tt := range tests {
		got := IsSafePath(tt.input)
		if got != tt.want {
			t.Errorf("IsSafePath(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestHomeDir(t *testing.T) {
	home, err := HomeDir()
	if err != nil {
		t.Fatalf("HomeDir() error: %v", err)
	}
	if home == "" {
		t.Error("HomeDir() 返回空字符串")
	}
	if !filepath.IsAbs(home) {
		t.Errorf("HomeDir() = %q, 应为绝对路径", home)
	}
}

func TestExecutableDir(t *testing.T) {
	dir, err := ExecutableDir()
	if err != nil {
		t.Fatalf("ExecutableDir() error: %v", err)
	}
	if dir == "" {
		t.Error("ExecutableDir() 返回空字符串")
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("ExecutableDir() = %q, 应为绝对路径", dir)
	}
}

func TestPathType(t *testing.T) {
	tmpDir := t.TempDir()
	if PathType(tmpDir) != "dir" {
		t.Errorf("PathType(%q) 应为 \"dir\"", tmpDir)
	}

	tmpFile := filepath.Join(tmpDir, "testfile.txt")
	err := os.WriteFile(tmpFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	if PathType(tmpFile) != "file" {
		t.Errorf("PathType(%q) 应为 \"file\"", tmpFile)
	}

	notExist := filepath.Join(tmpDir, "not_exist.txt")
	if PathType(notExist) != "not_exist" {
		t.Errorf("PathType(%q) 应为 \"not_exist\"", notExist)
	}
}
