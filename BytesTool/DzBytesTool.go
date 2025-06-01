package BytesTool

import (
	"bytes"
	"github.com/PWND0U/dztool/StringTool"
)

// DzBytes 定义一个字节切片类型，支持级联操作，用于处理字节相关的操作（如拼接、分割、替换等）
type DzBytes []byte

// NewDzBytes 根据传入的字符串创建一个新的 DzBytes 对象
// 参数 srcStr: 用于初始化 DzBytes 的字符串内容
// 返回值: 由字符串转换而来的 DzBytes 对象
func NewDzBytes(srcStr string) DzBytes {
	// 将传入的字符串转换为字节切片并返回
	return DzBytes(srcStr)
}

// NewDzBytesByBytes 根据传入的字节切片创建一个新的 DzBytes 对象
// 参数 bs: 用于初始化 DzBytes 的原始字节切片
// 返回值: 直接封装后的 DzBytes 对象
func NewDzBytesByBytes(bs []byte) DzBytes {
	// 直接封装传入的字节切片为 DzBytes 类型
	return bs
}

// Join 将多个字节切片使用指定分隔符连接成一个新的 DzBytes 对象
// 参数 bsList: 待连接的字节切片列表（如 [][]byte{[]byte("a"), []byte("b")}）
// 参数 bs: 用于连接的分隔符字节切片（如 []byte(",")）
// 返回值: 连接后的完整字节切片封装为 DzBytes 对象
func Join(bsList [][]byte, bs []byte) DzBytes {
	// 调用标准库 bytes.Join 执行连接操作
	joinBytes := bytes.Join(bsList, bs)
	return joinBytes
}

// ToString 将 DzBytes 对象转换为字符串类型
// 返回值: 字节切片对应的字符串表示
func (db DzBytes) ToString() string {
	// 将字节切片转换为字符串并返回
	return string(db)
}

// ToDzString 将 DzBytes 对象转换为 StringTool 包中的 DzString 类型（字符串工具类型）
// 返回值: 由当前字节切片转换而来的 DzString 对象
func (db DzBytes) ToDzString() StringTool.DzString {
	// 直接将字节切片转换为 DzString 类型
	return StringTool.DzString(db)
}

// ReplaceAll 替换当前字节切片中所有匹配的旧字节片段为新字节片段
// 参数 old: 需要被替换的旧字节片段（如 []byte("old")）
// 参数 new: 用于替换的新字节片段（如 []byte("new")）
// 返回值: 替换后的新 DzBytes 对象
func (db DzBytes) ReplaceAll(old, new []byte) DzBytes {
	return bytes.ReplaceAll(db, old, new)
}

// ReplaceN 替换当前字节切片中最多指定次数的旧字节片段为新字节片段
// 参数 old: 需要被替换的旧字节片段
// 参数 new: 用于替换的新字节片段
// 参数 count: 最多替换次数（count < 0 表示替换所有匹配项）
// 返回值: 替换后的新 DzBytes 对象
func (db DzBytes) ReplaceN(old, new []byte, count int) DzBytes {
	return bytes.Replace(db, old, new, count)
}

// Find 查找子字节片段在当前字节切片中的起始位置索引
// 参数 subBytes: 需要查找的子字节片段（如 []byte("sub")）
// 返回值: 子片段的起始索引（未找到返回 -1）
func (db DzBytes) Find(subBytes []byte) int {
	return bytes.Index(db, subBytes)
}

// Split 按指定分隔符字节片段分割当前字节切片（注意函数名拼写应为 Split）
// 参数 subBytes: 用于分割的分隔符字节片段（如 []byte(",")）
// 返回值: 分割后的 DzBytes 切片列表（空片段会被保留）
func (db DzBytes) Split(subBytes []byte) []DzBytes {
	re := make([]DzBytes, 0)
	// IsContains 检查当前字节切片是否包含指定的子字节片段
	// 参数 subBytes: 需要检查的子字节片段
	// 返回值: 包含返回 true，否则返回 false
	// SplitN 按指定分隔符和最多分割次数分割当前字节切片（注意函数名拼写应为 SplitN）
	// 参数 subBytes: 用于分割的分隔符字节片段
	// 参数 count: 最多分割次数（count <= 0 表示全部分割；count > 0 表示分割为最多 count 个片段）
	// 返回值: 分割后的 DzBytes 切片列表
	for _, s := range bytes.Split(db, subBytes) {
		re = append(re, s)
	}
	return re

}

func (db DzBytes) SplitN(subBytes []byte, count int) []DzBytes {
	// 返回分割后的字符串列表
	re := make([]DzBytes, 0)
	for _, s := range bytes.SplitN(db, subBytes, count) {
		re = append(re, s)
	}
	return re
}

func (db DzBytes) IsContains(subBytes []byte) bool {
	// 将DzBytes类型转换为string类型
	return bytes.Contains(db, subBytes)
}
