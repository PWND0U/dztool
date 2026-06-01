package IOTool

import (
	"os"
	"path/filepath"
	"strings"
)

// BaseName 获取文件名（不含扩展名）
func BaseName(path string) string {
	file := filepath.Base(path)
	ext := filepath.Ext(file)
	return strings.TrimSuffix(file, ext)
}

// Ext 获取文件扩展名（含点号）
func Ext(path string) string {
	return filepath.Ext(path)
}

// FileName 获取文件全名（含扩展名）
func FileName(path string) string {
	return filepath.Base(path)
}

// Parent 获取父目录路径
func Parent(path string) string {
	return filepath.Dir(path)
}

// Join 路径拼接
func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Normalize 路径规范化，消除 . 和 ..，统一分隔符，返回绝对路径
func Normalize(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

// Rel 获取相对路径
func Rel(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

// IsAbs 判断路径是否为绝对路径
func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// Depth 获取路径层级深度，根目录为 0
func Depth(path string) int {
	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return 0
	}
	vol := filepath.VolumeName(cleaned)
	withoutVol := cleaned[len(vol):]
	if withoutVol == "" || withoutVol == string(os.PathSeparator) {
		return 0
	}
	parts := strings.Split(strings.Trim(withoutVol, string(os.PathSeparator)), string(os.PathSeparator))
	count := 0
	for _, p := range parts {
		if p != "" {
			count++
		}
	}
	return count
}

// PathIterator 路径遍历器，返回从根到目标的每一级路径
func PathIterator(path string) []string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil
	}
	cleaned := filepath.Clean(abs)

	vol := filepath.VolumeName(cleaned)
	withoutVol := cleaned[len(vol):]

	result := make([]string, 0)

	root := vol + string(os.PathSeparator)
	result = append(result, root)

	if withoutVol == "" || withoutVol == string(os.PathSeparator) {
		return result
	}

	trimmed := strings.Trim(withoutVol, string(os.PathSeparator))
	parts := strings.Split(trimmed, string(os.PathSeparator))

	current := root
	for _, p := range parts {
		if p == "" {
			continue
		}
		current = current + p
		result = append(result, current)
		current = current + string(os.PathSeparator)
	}

	return result
}

// ReplaceExt 替换文件扩展名
func ReplaceExt(path string, newExt string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext) + newExt
}

// IsSubPath 判断 child 是否在 parent 目录下
func IsSubPath(parent, child string) bool {
	absParent, err := filepath.Abs(parent)
	if err != nil {
		return false
	}
	absChild, err := filepath.Abs(child)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absParent, absChild)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && rel != "."
}

// IsSafePath 路径安全校验，判断路径是否安全（不含路径穿越 ..）
func IsSafePath(path string) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	for _, p := range parts {
		if p == ".." {
			return false
		}
	}
	return true
}

// HomeDir 获取用户主目录
func HomeDir() (string, error) {
	return os.UserHomeDir()
}

// ExecutableDir 获取可执行文件所在目录
func ExecutableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// PathType 路径类型判断，返回 "file"、"dir" 或 "not_exist"
func PathType(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "not_exist"
		}
		return "not_exist"
	}
	if info.IsDir() {
		return "dir"
	}
	return "file"
}
