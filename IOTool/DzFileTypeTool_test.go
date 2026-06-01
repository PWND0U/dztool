package IOTool

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFileTypeByBytes_JPEG(t *testing.T) {
	data := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	result := DetectFileTypeByBytes(data)
	if result != "JPEG" {
		t.Errorf("期望 JPEG，实际 %s", result)
	}
}

func TestDetectFileTypeByBytes_PNG(t *testing.T) {
	data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00}
	result := DetectFileTypeByBytes(data)
	if result != "PNG" {
		t.Errorf("期望 PNG，实际 %s", result)
	}
}

func TestDetectFileTypeByBytes_PDF(t *testing.T) {
	data := []byte{0x25, 0x50, 0x44, 0x46, 0x2D}
	result := DetectFileTypeByBytes(data)
	if result != "PDF" {
		t.Errorf("期望 PDF，实际 %s", result)
	}
}

func TestDetectFileTypeByBytes_ZIP(t *testing.T) {
	data := []byte{0x50, 0x4B, 0x03, 0x04, 0x00, 0x00}
	result := DetectFileTypeByBytes(data)
	if result != "ZIP" {
		t.Errorf("期望 ZIP，实际 %s", result)
	}
}

func TestDetectFileTypeByBytes_Unknown(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x00}
	result := DetectFileTypeByBytes(data)
	if result != "" {
		t.Errorf("期望空字符串，实际 %s", result)
	}
}

func TestRegisterFileType(t *testing.T) {
	RegisterFileType("CustomType", [][]byte{{0xAA, 0xBB, 0xCC}}, []string{".custom"}, 0)
	data := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	result := DetectFileTypeByBytes(data)
	if result != "CustomType" {
		t.Errorf("期望 CustomType，实际 %s", result)
	}
}

func TestGetRegisteredTypes(t *testing.T) {
	types := GetRegisteredTypes()
	if len(types) == 0 {
		t.Error("已注册类型不应为空")
	}
	found := false
	for _, ft := range types {
		if ft.Name == "JPEG" {
			found = true
			break
		}
	}
	if !found {
		t.Error("应包含 JPEG 类型")
	}
}

func TestDetectFileType(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.jpg")
	data := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	err := os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	result, err := DetectFileType(tmpFile)
	if err != nil {
		t.Fatalf("检测文件类型失败: %v", err)
	}
	if result != "JPEG" {
		t.Errorf("期望 JPEG，实际 %s", result)
	}
}

func TestIsFileType(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.png")
	data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	err := os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	isPNG, err := IsFileType(tmpFile, "PNG")
	if err != nil {
		t.Fatalf("判断文件类型失败: %v", err)
	}
	if !isPNG {
		t.Error("期望文件为 PNG 类型")
	}
	isJPEG, err := IsFileType(tmpFile, "JPEG")
	if err != nil {
		t.Fatalf("判断文件类型失败: %v", err)
	}
	if isJPEG {
		t.Error("期望文件不是 JPEG 类型")
	}
}
