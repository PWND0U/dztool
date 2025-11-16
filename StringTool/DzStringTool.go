package StringTool

import (
	"regexp"
	"strings"

	"github.com/PWND0U/dztool/Algorithm"
	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DzString 定义一个字符串类型，支持级联操作，用于处理字符串相关的操作（如拼接、分割、替换等）
type DzString string

// NewDzString 根据传入的字符串创建一个新的 DzString 对象
// 参数 srcStr: 用于初始化 DzString 的原始字符串内容
// 返回值: 由输入字符串封装的 DzString 对象
func NewDzString(srcStr string) DzString {
	// 直接将输入字符串转换为 DzString 类型
	return DzString(srcStr)
}

// Join 将多个字符串使用指定分隔符连接成一个新的 DzString 对象
// 参数 strList: 待连接的字符串列表（如 []string{"a", "b"}）
// 参数 seqStr: 用于连接的分隔符字符串（如 ","）
// 返回值: 连接后的完整字符串封装为 DzString 对象
func Join(strList []string, seqStr string) DzString {
	// 调用标准库 strings.Join 执行连接操作
	return NewDzString(strings.Join(strList, seqStr))
}

// ToString 将 DzString 对象转换为标准 string 类型
// 返回值: DzString 对应的底层字符串内容
func (ds DzString) ToString() string {
	// 直接将 DzString 类型转换为 string
	return string(ds)
}

// ToBytes 将 DzString 对象转换为字节切片（[]byte）
// 返回值: 字符串对应的字节表示
func (ds DzString) ToBytes() []byte {
	// 直接将 DzString 转换为字节切片
	return []byte(ds)
}

// ReplaceAll 替换当前字符串中所有匹配的旧子串为新子串
// 参数 old: 需要被替换的旧子串（如 "old"）
// 参数 new: 用于替换的新子串（如 "new"）
// 返回值: 替换后的新 DzString 对象
func (ds DzString) ReplaceAll(old, new string) DzString {
	// 调用标准库 strings.ReplaceAll 执行替换
	return NewDzString(strings.ReplaceAll(string(ds), old, new))
}

// RexReplaceAll 使用正则表达式替换当前字符串中所有匹配的子串为新字符串（注：函数名建议修正为 RegexReplaceAll）
// 参数 regexpStr: 正则表达式模式（如 `\d+` 匹配数字）
// 参数 new: 用于替换的新字符串
// 返回值: 替换后的新 DzString 对象（若正则表达式编译失败则返回原对象）
func (ds DzString) RexReplaceAll(regexpStr, new string) DzString {
	// 编译正则表达式，失败时返回原对象
	compile, err := regexp.Compile(regexpStr)
	if err != nil {
		return ds
	}
	// 执行正则替换并封装为 DzString
	return NewDzString(compile.ReplaceAllString(string(ds), new))
}

// ReplaceN 替换当前字符串中最多指定次数的旧子串为新子串
// 参数 old: 需要被替换的旧子串
// 参数 new: 用于替换的新子串
// 参数 count: 最多替换次数（count < 0 表示替换所有匹配项）
// 返回值: 替换后的新 DzString 对象
func (ds DzString) ReplaceN(old, new string, count int) DzString {
	// 调用标准库 strings.Replace 执行有限次数替换
	return NewDzString(strings.Replace(string(ds), old, new, count))
}

// Upper 将当前字符串转换为全大写形式
// 返回值: 转换为大写后的 DzString 对象
func (ds DzString) Upper() DzString {
	// 调用标准库 strings.ToUpper 转换大写
	return NewDzString(strings.ToUpper(string(ds)))
}

// Lower 将当前字符串转换为全小写形式
// 返回值: 转换为小写后的 DzString 对象
func (ds DzString) Lower() DzString {
	// 调用标准库 strings.ToLower 转换小写
	return NewDzString(strings.ToLower(string(ds)))
}

// Find 查找子串在当前字符串中的起始位置索引
// 参数 subStr: 需要查找的子串（如 "sub"）
// 返回值: 子串的起始索引（未找到时返回 -1）
func (ds DzString) Find(subStr string) int {
	// 调用标准库 strings.Index 查找位置
	return strings.Index(string(ds), subStr)
}

// Split 按指定分隔符分割当前字符串为字符串列表（封装为 DzStrings 类型）
// 参数 sepStr: 用于分割的分隔符字符串（如 ","）
// 返回值: 分割后的字符串列表（封装为 DzStrings）
func (ds DzString) Split(sepStr string) DzStrings {
	// 调用标准库 strings.Split 分割并封装为 DzStrings
	return NewDzStrings(strings.Split(string(ds), sepStr))
}

// SplitN 按指定分隔符和最多分割次数分割当前字符串为字符串列表（封装为 DzStrings 类型）
// 参数 sepStr: 用于分割的分隔符字符串
// 参数 count: 最多分割次数（count <= 0 表示全部分割；count > 0 表示分割为最多 count 个片段）
// 返回值: 分割后的字符串列表（封装为 DzStrings）
func (ds DzString) SplitN(sepStr string, count int) DzStrings {
	// 调用标准库 strings.SplitN 分割并封装为 DzStrings
	return NewDzStrings(strings.SplitN(string(ds), sepStr, count))
}

// IsEmpty 检查当前字符串是否为空（长度为 0 或空字符串）
// 返回值: 空字符串返回 true，否则返回 false
func (ds DzString) IsEmpty() bool {
	// 直接判断字符串是否为空或长度为 0
	return ds == "" || len(ds) <= 0
}

// Strip 去除当前字符串首尾（两侧）中包含的指定字符集
// 参数 cutset: 需要去除的字符集合（如 " \t\n" 表示空格、制表符、换行符）
// 返回值: 去除后的新 DzString 对象
func (ds DzString) Strip(cutset string) DzString {
	// 调用标准库 strings.Trim 去除首尾指定字符
	return NewDzString(strings.Trim(string(ds), cutset))
}

// LStrip 去除当前字符串左侧（开头）中包含的指定字符集
// 参数 cutset: 需要去除的字符集合
// 返回值: 去除后的新 DzString 对象
func (ds DzString) LStrip(cutset string) DzString {
	// 调用标准库 strings.TrimLeft 去除左侧指定字符
	return NewDzString(strings.TrimLeft(string(ds), cutset))
}

// RStrip 去除当前字符串右侧（末尾）中包含的指定字符集
// 参数 cutset: 需要去除的字符集合
// 返回值: 去除后的新 DzString 对象
func (ds DzString) RStrip(cutset string) DzString {
	// 调用标准库 strings.TrimRight 去除右侧指定字符
	return NewDzString(strings.TrimRight(string(ds), cutset))
}

// IsContains 检查当前字符串是否包含指定子串
// 参数 subStr: 需要检查的子串
// 返回值: 包含子串返回 true，否则返回 false
func (ds DzString) IsContains(subStr string) bool {
	// 调用标准库 strings.Contains 检查包含关系
	return strings.Contains(string(ds), subStr)
}

// Title 将当前字符串转换为标题格式（每个单词首字母大写，其余字母小写）
// 返回值: 标题格式的新 DzString 对象
func (ds DzString) Title() DzString {
	// 调用 golang.org/x/text/cases 库的 Title 方法转换标题格式
	return NewDzString(cases.Title(language.Und).String(ds.ToString()))
}

// SimilarText 计算当前字符串与目标字符串的相似度（基于 Levenshtein 编辑距离）
// 参数 str: 用于比较的目标字符串
// 返回值: 相似度值（范围 0-1，1 表示完全相同）
func (ds DzString) SimilarText(str string) float64 {
	// 通过 Levenshtein 距离计算相似度（距离越小，相似度越高）
	return 1 - float64(Algorithm.DzLevenshtein(ds.ToString(), str, 1, 1, 1))/float64(len([]rune(ds.ToString())))
}

// ToInt 将当前字符串转换为 int 类型（依赖 cast 库，转换失败返回 0）
// 返回值: 转换后的整数值（转换失败时返回 0）
func (ds DzString) ToInt() int {
	// 调用 cast.ToInt 执行字符串转整数
	return cast.ToInt(ds.ToString())
}

// FStringFormat 格式化当前字符串，替换其中的占位符为指定值
// 参数 values: 包含占位符和对应值的映射（如 {"name": "Alice"}）
// 返回值: 格式化后的新 DzString 对象
func (ds DzString) FStringFormat(values map[string]any) DzString {
	// 调用 cast.ToInt 执行字符串转整数
	result := ds
	for k, v := range values {
		result = result.ReplaceAll("{"+k+"}", cast.ToString(v))
	}
	return result
}
