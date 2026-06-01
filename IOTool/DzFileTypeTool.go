package IOTool

import (
	"bytes"
	"io"
	"os"
	"sync"
)

// FileType 文件类型定义
type FileType struct {
	Name         string   // 类型名称
	MagicNumbers [][]byte // 魔数字节列表
	Extensions   []string // 常见扩展名
	Offset       int      // 魔数偏移量
	matchFunc    func([]byte) bool
}

var (
	fileTypeRegistry []FileType
	registryMu       sync.RWMutex
)

// DetectFileType 通过文件路径判断文件类型
func DetectFileType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return DetectFileTypeFromReader(f)
}

// DetectFileTypeByBytes 通过字节切片判断文件类型
func DetectFileTypeByBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	registryMu.RLock()
	defer registryMu.RUnlock()
	for _, ft := range fileTypeRegistry {
		if ft.matchFunc != nil {
			if ft.matchFunc(data) {
				return ft.Name
			}
			continue
		}
		if matchMagic(data, ft) {
			return ft.Name
		}
	}
	return ""
}

// DetectFileTypeFromReader 通过 Reader 判断文件类型
func DetectFileTypeFromReader(reader io.Reader) (string, error) {
	buf := make([]byte, 8192)
	n, err := io.ReadFull(reader, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}
	if n == 0 {
		return "", nil
	}
	return DetectFileTypeByBytes(buf[:n]), nil
}

// IsFileType 判断文件是否为指定类型
func IsFileType(path string, typeName string) (bool, error) {
	detected, err := DetectFileType(path)
	if err != nil {
		return false, err
	}
	return detected == typeName, nil
}

// RegisterFileType 注册自定义文件类型
func RegisterFileType(name string, magicNumbers [][]byte, extensions []string, offset int) {
	registryMu.Lock()
	defer registryMu.Unlock()
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         name,
		MagicNumbers: magicNumbers,
		Extensions:   extensions,
		Offset:       offset,
	})
}

// GetRegisteredTypes 获取所有已注册的文件类型
func GetRegisteredTypes() []FileType {
	registryMu.RLock()
	defer registryMu.RUnlock()
	result := make([]FileType, len(fileTypeRegistry))
	copy(result, fileTypeRegistry)
	return result
}

func matchMagic(data []byte, ft FileType) bool {
	if len(data) < ft.Offset {
		return false
	}
	data = data[ft.Offset:]
	for _, magic := range ft.MagicNumbers {
		if len(data) < len(magic) {
			continue
		}
		if bytes.HasPrefix(data, magic) {
			return true
		}
	}
	return false
}

func init() {
	// 图片类型
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "JPEG",
		MagicNumbers: [][]byte{{0xFF, 0xD8, 0xFF}},
		Extensions:   []string{".jpg", ".jpeg"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "PNG",
		MagicNumbers: [][]byte{{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
		Extensions:   []string{".png"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "GIF",
		MagicNumbers: [][]byte{[]byte("GIF87a"), []byte("GIF89a")},
		Extensions:   []string{".gif"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "BMP",
		MagicNumbers: [][]byte{{0x42, 0x4D}},
		Extensions:   []string{".bmp"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "WebP",
		Extensions: []string{".webp"},
		matchFunc: func(data []byte) bool {
			if len(data) < 12 {
				return false
			}
			return bytes.HasPrefix(data, []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP"))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "ICO",
		MagicNumbers: [][]byte{{0x00, 0x00, 0x01, 0x00}},
		Extensions:   []string{".ico"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "TIFF",
		MagicNumbers: [][]byte{{0x49, 0x49, 0x2A, 0x00}, {0x4D, 0x4D, 0x00, 0x2A}},
		Extensions:   []string{".tiff", ".tif"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "SVG",
		Extensions: []string{".svg"},
		matchFunc: func(data []byte) bool {
			return bytes.Contains(data[:min(len(data), 512)], []byte("<svg"))
		},
	})

	// 文档类型
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "PDF",
		MagicNumbers: [][]byte{{0x25, 0x50, 0x44, 0x46}},
		Extensions:   []string{".pdf"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "DOCX",
		Extensions: []string{".docx"},
		matchFunc: func(data []byte) bool {
			if !bytes.HasPrefix(data, []byte{0x50, 0x4B, 0x03, 0x04}) {
				return false
			}
			return bytes.Contains(data[:min(len(data), 8192)], []byte("word/"))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "XLSX",
		Extensions: []string{".xlsx"},
		matchFunc: func(data []byte) bool {
			if !bytes.HasPrefix(data, []byte{0x50, 0x4B, 0x03, 0x04}) {
				return false
			}
			return bytes.Contains(data[:min(len(data), 8192)], []byte("xl/"))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "PPTX",
		Extensions: []string{".pptx"},
		matchFunc: func(data []byte) bool {
			if !bytes.HasPrefix(data, []byte{0x50, 0x4B, 0x03, 0x04}) {
				return false
			}
			return bytes.Contains(data[:min(len(data), 8192)], []byte("ppt/"))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "DOC",
		MagicNumbers: [][]byte{{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
		Extensions:   []string{".doc"},
	})

	// 压缩类型
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "ZIP",
		MagicNumbers: [][]byte{{0x50, 0x4B, 0x03, 0x04}},
		Extensions:   []string{".zip"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "RAR",
		MagicNumbers: [][]byte{{0x52, 0x61, 0x72, 0x21}},
		Extensions:   []string{".rar"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "7Z",
		MagicNumbers: [][]byte{{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}},
		Extensions:   []string{".7z"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "GZIP",
		MagicNumbers: [][]byte{{0x1F, 0x8B}},
		Extensions:   []string{".gz"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "BZ2",
		MagicNumbers: [][]byte{{0x42, 0x5A}},
		Extensions:   []string{".bz2"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "TAR",
		MagicNumbers: [][]byte{[]byte("ustar")},
		Extensions:   []string{".tar"},
		Offset:       257,
	})

	// 音频类型
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "MP3",
		MagicNumbers: [][]byte{{0xFF, 0xFB}, {0xFF, 0xF3}, {0xFF, 0xF2}, {0x49, 0x44, 0x33}},
		Extensions:   []string{".mp3"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "WAV",
		Extensions: []string{".wav"},
		matchFunc: func(data []byte) bool {
			if len(data) < 12 {
				return false
			}
			return bytes.HasPrefix(data, []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WAVE"))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "FLAC",
		MagicNumbers: [][]byte{{0x66, 0x4C, 0x61, 0x43}},
		Extensions:   []string{".flac"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "OGG",
		MagicNumbers: [][]byte{{0x4F, 0x67, 0x67, 0x53}},
		Extensions:   []string{".ogg"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "AAC",
		MagicNumbers: [][]byte{{0xFF, 0xF1}, {0xFF, 0xF9}},
		Extensions:   []string{".aac"},
	})

	// 视频类型
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "MP4",
		Extensions: []string{".mp4"},
		matchFunc: func(data []byte) bool {
			if len(data) < 8 {
				return false
			}
			return bytes.Equal(data[4:8], []byte("ftyp"))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "AVI",
		Extensions: []string{".avi"},
		matchFunc: func(data []byte) bool {
			if len(data) < 12 {
				return false
			}
			return bytes.HasPrefix(data, []byte("RIFF")) && bytes.Equal(data[8:12], []byte("AVI "))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "MKV",
		MagicNumbers: [][]byte{{0x1A, 0x45, 0xDF, 0xA3}},
		Extensions:   []string{".mkv"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "MOV",
		MagicNumbers: [][]byte{[]byte("moov")},
		Extensions:   []string{".mov"},
		Offset:       4,
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "WMV",
		MagicNumbers: [][]byte{{0x30, 0x26, 0xB2, 0x75, 0x8E, 0x66, 0xCF, 0x11}},
		Extensions:   []string{".wmv"},
	})

	// 可执行类型
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "EXE",
		MagicNumbers: [][]byte{{0x4D, 0x5A}},
		Extensions:   []string{".exe"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "ELF",
		MagicNumbers: [][]byte{{0x7F, 0x45, 0x4C, 0x46}},
		Extensions:   []string{""},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name: "Mach-O",
		MagicNumbers: [][]byte{
			{0xCE, 0xFA, 0xED, 0xFE},
			{0xCF, 0xFA, 0xED, 0xFE},
			{0xFE, 0xED, 0xFA, 0xCE},
			{0xFE, 0xED, 0xFA, 0xCF},
		},
		Extensions: []string{""},
	})

	// 文本类型
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:         "XML",
		MagicNumbers: [][]byte{[]byte("<?xml")},
		Extensions:   []string{".xml"},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "HTML",
		Extensions: []string{".html", ".htm"},
		matchFunc: func(data []byte) bool {
			lower := bytes.ToLower(data[:min(len(data), 512)])
			return bytes.Contains(lower, []byte("<!doctype html")) || bytes.Contains(lower, []byte("<html"))
		},
	})
	fileTypeRegistry = append(fileTypeRegistry, FileType{
		Name:       "JSON",
		Extensions: []string{".json"},
		matchFunc: func(data []byte) bool {
			for _, b := range data {
				switch b {
				case ' ', '\t', '\n', '\r':
					continue
				case '{', '[':
					return true
				default:
					return false
				}
			}
			return false
		},
	})
}
