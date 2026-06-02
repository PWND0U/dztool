package dzUtils

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// ==================== 数值进制转换 ====================

// DecToBin 将十进制整型转换为二进制字符串
func DecToBin(n int64) string {
	return strconv.FormatInt(n, 2)
}

// DecToOct 将十进制整型转换为八进制字符串
func DecToOct(n int64) string {
	return strconv.FormatInt(n, 8)
}

// DecToHex 将十进制整型转换为十六进制字符串（小写）
func DecToHex(n int64) string {
	return strconv.FormatInt(n, 16)
}

// BinToDec 将二进制字符串转换为十进制整型
// 参数 bin: 二进制字符串，支持带"0b"前缀
func BinToDec(bin string) (int64, error) {
	bin = trimPrefix(bin, "0b")
	n, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("dzUtils: 解析二进制字符串失败 %q: %w", bin, err)
	}
	return n, nil
}

// OctToDec 将八进制字符串转换为十进制整型
// 参数 oct: 八进制字符串，支持带"0o"或"0"前缀
func OctToDec(oct string) (int64, error) {
	oct = trimPrefix(oct, "0o", "0")
	n, err := strconv.ParseInt(oct, 8, 64)
	if err != nil {
		return 0, fmt.Errorf("dzUtils: 解析八进制字符串失败 %q: %w", oct, err)
	}
	return n, nil
}

// HexToDec 将十六进制字符串转换为十进制整型
// 参数 hex: 十六进制字符串，支持带"0x"前缀，大小写不敏感
func HexToDec(hex string) (int64, error) {
	hex = trimPrefix(hex, "0x", "0X")
	n, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("dzUtils: 解析十六进制字符串失败 %q: %w", hex, err)
	}
	return n, nil
}

// BinToOct 将二进制字符串转换为八进制字符串
func BinToOct(bin string) (string, error) {
	n, err := BinToDec(bin)
	if err != nil {
		return "", err
	}
	return DecToOct(n), nil
}

// BinToHex 将二进制字符串转换为十六进制字符串（小写）
func BinToHex(bin string) (string, error) {
	n, err := BinToDec(bin)
	if err != nil {
		return "", err
	}
	return DecToHex(n), nil
}

// OctToBin 将八进制字符串转换为二进制字符串
func OctToBin(oct string) (string, error) {
	n, err := OctToDec(oct)
	if err != nil {
		return "", err
	}
	return DecToBin(n), nil
}

// OctToHex 将八进制字符串转换为十六进制字符串（小写）
func OctToHex(oct string) (string, error) {
	n, err := OctToDec(oct)
	if err != nil {
		return "", err
	}
	return DecToHex(n), nil
}

// HexToBin 将十六进制字符串转换为二进制字符串
func HexToBin(hex string) (string, error) {
	n, err := HexToDec(hex)
	if err != nil {
		return "", err
	}
	return DecToBin(n), nil
}

// HexToOct 将十六进制字符串转换为八进制字符串
func HexToOct(hex string) (string, error) {
	n, err := HexToDec(hex)
	if err != nil {
		return "", err
	}
	return DecToOct(n), nil
}

// ==================== 大整数进制转换 ====================

// BigDecToBin 将十进制大整数转换为二进制字符串
func BigDecToBin(n *big.Int) string {
	return n.Text(2)
}

// BigDecToOct 将十进制大整数转换为八进制字符串
func BigDecToOct(n *big.Int) string {
	return n.Text(8)
}

// BigDecToHex 将十进制大整数转换为十六进制字符串（小写）
func BigDecToHex(n *big.Int) string {
	return n.Text(16)
}

// BigBinToDec 将二进制字符串转换为十进制大整数
func BigBinToDec(bin string) (*big.Int, error) {
	bin = trimPrefix(bin, "0b")
	n := new(big.Int)
	_, ok := n.SetString(bin, 2)
	if !ok {
		return nil, fmt.Errorf("dzUtils: 解析二进制大整数失败 %q", bin)
	}
	return n, nil
}

// BigOctToDec 将八进制字符串转换为十进制大整数
func BigOctToDec(oct string) (*big.Int, error) {
	oct = trimPrefix(oct, "0o", "0")
	n := new(big.Int)
	_, ok := n.SetString(oct, 8)
	if !ok {
		return nil, fmt.Errorf("dzUtils: 解析八进制大整数失败 %q", oct)
	}
	return n, nil
}

// BigHexToDec 将十六进制字符串转换为十进制大整数
func BigHexToDec(hex string) (*big.Int, error) {
	hex = trimPrefix(hex, "0x", "0X")
	n := new(big.Int)
	_, ok := n.SetString(hex, 16)
	if !ok {
		return nil, fmt.Errorf("dzUtils: 解析十六进制大整数失败 %q", hex)
	}
	return n, nil
}

// ==================== 字符串与进制转换 ====================

// StrToBin 将字符串转换为二进制表示字符串（每个字节用空格分隔）
func StrToBin(s string) string {
	var builder strings.Builder
	for i := 0; i < len(s); i++ {
		if i > 0 {
			builder.WriteByte(' ')
		}
		fmt.Fprintf(&builder, "%08b", s[i])
	}
	return builder.String()
}

// StrToOct 将字符串转换为八进制表示字符串（每个字节用空格分隔）
func StrToOct(s string) string {
	var builder strings.Builder
	for i := 0; i < len(s); i++ {
		if i > 0 {
			builder.WriteByte(' ')
		}
		fmt.Fprintf(&builder, "%03o", s[i])
	}
	return builder.String()
}

// StrToHex 将字符串转换为十六进制表示字符串（小写，连续输出）
func StrToHex(s string) string {
	return fmt.Sprintf("%x", s)
}

// BinToStr 将二进制表示字符串还原为原始字符串
// 参数 bin: 二进制字符串，支持空格分隔或连续
func BinToStr(bin string) (string, error) {
	bin = strings.ReplaceAll(bin, " ", "")
	if len(bin)%8 != 0 {
		return "", fmt.Errorf("dzUtils: 二进制字符串长度 %d 不是8的倍数", len(bin))
	}
	var builder strings.Builder
	for i := 0; i < len(bin); i += 8 {
		b, err := strconv.ParseUint(bin[i:i+8], 2, 8)
		if err != nil {
			return "", fmt.Errorf("dzUtils: 解析二进制字节失败 %q: %w", bin[i:i+8], err)
		}
		builder.WriteByte(byte(b))
	}
	return builder.String(), nil
}

// OctToStr 将八进制表示字符串还原为原始字符串
// 参数 oct: 八进制字符串，支持空格分隔或连续
func OctToStr(oct string) (string, error) {
	oct = strings.ReplaceAll(oct, " ", "")
	if len(oct)%3 != 0 {
		return "", fmt.Errorf("dzUtils: 八进制字符串长度 %d 不是3的倍数", len(oct))
	}
	var builder strings.Builder
	for i := 0; i < len(oct); i += 3 {
		b, err := strconv.ParseUint(oct[i:i+3], 8, 8)
		if err != nil {
			return "", fmt.Errorf("dzUtils: 解析八进制字节失败 %q: %w", oct[i:i+3], err)
		}
		builder.WriteByte(byte(b))
	}
	return builder.String(), nil
}

// HexToStr 将十六进制表示字符串还原为原始字符串
// 参数 hex: 十六进制字符串，连续输出，大小写不敏感
func HexToStr(hex string) (string, error) {
	hex = trimPrefix(hex, "0x", "0X")
	if len(hex)%2 != 0 {
		return "", fmt.Errorf("dzUtils: 十六进制字符串长度 %d 不是2的倍数", len(hex))
	}
	var builder strings.Builder
	for i := 0; i < len(hex); i += 2 {
		b, err := strconv.ParseUint(hex[i:i+2], 16, 8)
		if err != nil {
			return "", fmt.Errorf("dzUtils: 解析十六进制字节失败 %q: %w", hex[i:i+2], err)
		}
		builder.WriteByte(byte(b))
	}
	return builder.String(), nil
}

// ==================== 字节数组与进制转换 ====================

// BytesToBin 将字节数组转换为二进制表示字符串（每个字节用空格分隔）
func BytesToBin(data []byte) string {
	var builder strings.Builder
	for i, b := range data {
		if i > 0 {
			builder.WriteByte(' ')
		}
		fmt.Fprintf(&builder, "%08b", b)
	}
	return builder.String()
}

// BytesToOct 将字节数组转换为八进制表示字符串（每个字节用空格分隔）
func BytesToOct(data []byte) string {
	var builder strings.Builder
	for i, b := range data {
		if i > 0 {
			builder.WriteByte(' ')
		}
		fmt.Fprintf(&builder, "%03o", b)
	}
	return builder.String()
}

// BytesToHex 将字节数组转换为十六进制表示字符串（小写，连续输出）
func BytesToHex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

// BinToBytes 将二进制表示字符串还原为字节数组
// 参数 bin: 二进制字符串，支持空格分隔或连续
func BinToBytes(bin string) ([]byte, error) {
	bin = strings.ReplaceAll(bin, " ", "")
	if len(bin)%8 != 0 {
		return nil, fmt.Errorf("dzUtils: 二进制字符串长度 %d 不是8的倍数", len(bin))
	}
	result := make([]byte, 0, len(bin)/8)
	for i := 0; i < len(bin); i += 8 {
		b, err := strconv.ParseUint(bin[i:i+8], 2, 8)
		if err != nil {
			return nil, fmt.Errorf("dzUtils: 解析二进制字节失败 %q: %w", bin[i:i+8], err)
		}
		result = append(result, byte(b))
	}
	return result, nil
}

// OctToBytes 将八进制表示字符串还原为字节数组
// 参数 oct: 八进制字符串，支持空格分隔或连续
func OctToBytes(oct string) ([]byte, error) {
	oct = strings.ReplaceAll(oct, " ", "")
	if len(oct)%3 != 0 {
		return nil, fmt.Errorf("dzUtils: 八进制字符串长度 %d 不是3的倍数", len(oct))
	}
	result := make([]byte, 0, len(oct)/3)
	for i := 0; i < len(oct); i += 3 {
		b, err := strconv.ParseUint(oct[i:i+3], 8, 8)
		if err != nil {
			return nil, fmt.Errorf("dzUtils: 解析八进制字节失败 %q: %w", oct[i:i+3], err)
		}
		result = append(result, byte(b))
	}
	return result, nil
}

// HexToBytes 将十六进制表示字符串还原为字节数组
// 参数 hex: 十六进制字符串，连续输出，大小写不敏感
func HexToBytes(hex string) ([]byte, error) {
	hex = trimPrefix(hex, "0x", "0X")
	if len(hex)%2 != 0 {
		return nil, fmt.Errorf("dzUtils: 十六进制字符串长度 %d 不是2的倍数", len(hex))
	}
	result := make([]byte, 0, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		b, err := strconv.ParseUint(hex[i:i+2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("dzUtils: 解析十六进制字节失败 %q: %w", hex[i:i+2], err)
		}
		result = append(result, byte(b))
	}
	return result, nil
}

// ==================== 通用进制转换 ====================

// ConvertBase 将字符串从源进制转换为目标进制（支持2-36进制）
// 参数 src: 源字符串
// 参数 fromBase: 源进制（2-36）
// 参数 toBase: 目标进制（2-36）
func ConvertBase(src string, fromBase, toBase int) (string, error) {
	if fromBase < 2 || fromBase > 36 {
		return "", fmt.Errorf("dzUtils: 源进制 %d 不在有效范围 [2, 36] 内", fromBase)
	}
	if toBase < 2 || toBase > 36 {
		return "", fmt.Errorf("dzUtils: 目标进制 %d 不在有效范围 [2, 36] 内", toBase)
	}
	n, err := strconv.ParseInt(src, fromBase, 64)
	if err != nil {
		return "", fmt.Errorf("dzUtils: 进制转换失败，无法以 %d 进制解析 %q: %w", fromBase, src, err)
	}
	return strconv.FormatInt(n, toBase), nil
}

// ConvertBaseBig 将大整数字符串从源进制转换为目标进制（支持2-62进制）
// 参数 src: 源字符串
// 参数 fromBase: 源进制（2-62）
// 参数 toBase: 目标进制（2-62）
func ConvertBaseBig(src string, fromBase, toBase int) (string, error) {
	if fromBase < 2 || fromBase > 62 {
		return "", fmt.Errorf("dzUtils: 源进制 %d 不在有效范围 [2, 62] 内", fromBase)
	}
	if toBase < 2 || toBase > 62 {
		return "", fmt.Errorf("dzUtils: 目标进制 %d 不在有效范围 [2, 62] 内", toBase)
	}
	n := new(big.Int)
	_, ok := n.SetString(src, fromBase)
	if !ok {
		return "", fmt.Errorf("dzUtils: 大整数进制转换失败，无法以 %d 进制解析 %q", fromBase, src)
	}
	return n.Text(toBase), nil
}

// ==================== 内部辅助函数 ====================

// trimPrefix 去除字符串的任意一个前缀（匹配第一个命中的前缀）
func trimPrefix(s string, prefixes ...string) string {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) && len(s) > len(p) {
			return strings.TrimPrefix(s, p)
		}
	}
	return s
}
