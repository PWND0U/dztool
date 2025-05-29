package StringTool

import (
	"github.com/PWND0U/dztool/Algorithm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

// DzString 支持级联操作的String utils
type DzString string

// NewDzString 创建一个新的DzString对象，传入一个字符串参数
func NewDzString(srcStr string) DzString {
	// 返回一个新的DzString对象，传入的参数为srcStr
	return DzString(srcStr)
}

func Join(strList []string, seqStr string) DzString {
	// 返回分割后的字符串列表
	return NewDzString(strings.Join(strList, seqStr))

}

// ToString 将DzString类型转换为string类型
func (ds DzString) ToString() string {
	// 将DzString类型转换为string类型
	return string(ds)
}

// ReplaceAll 用于将DzString类型的ds中的所有old字符串替换为new字符串，并返回一个新的DzString类型
func (ds DzString) ReplaceAll(old, new string) DzString {
	// 将DzString类型的ds转换为string类型
	return NewDzString(strings.ReplaceAll(string(ds), old, new))
}

// ReplaceN 用于替换字符串中的指定字符，并返回替换后的字符串
func (ds DzString) ReplaceN(old, new string, count int) DzString {
	// 将DzString类型转换为string类型
	return NewDzString(strings.Replace(string(ds), old, new, count))
}

// Upper 用于将字符串转换为大写，并返回转换后的字符串
func (ds DzString) Upper() DzString {
	// 将DzString类型转换为string类型
	return NewDzString(strings.ToUpper(string(ds)))
}

// Lower 用于将字符串转换为小写，并返回转换后的字符串
func (ds DzString) Lower() DzString {
	// 将DzString类型转换为string类型
	return NewDzString(strings.ToLower(string(ds)))
}

// Find 用于查找子字符串在字符串中的位置，并返回位置索引
func (ds DzString) Find(subStr string) int {
	// 将DzString类型转换为string类型
	return strings.Index(string(ds), subStr)
}

// Spilt 用于将字符串按照指定的分隔符进行分割，并返回分割后的字符串列表
func (ds DzString) Spilt(sepStr string) DzStrings {
	// 返回分割后的字符串列表
	return NewDzStrings(strings.Split(string(ds), sepStr))

}

func (ds DzString) SpiltN(sepStr string, count int) DzStrings {
	// 返回分割后的字符串列表
	return NewDzStrings(strings.SplitN(string(ds), sepStr, count))
}

func (ds DzString) IsEmpty() bool {
	// 返回分割后的字符串列表
	return ds == "" || len(ds) <= 0
}

func (ds DzString) Strip(cutset string) DzString {
	// 将DzString类型转换为string类型
	return NewDzString(strings.Trim(string(ds), cutset))
}

func (ds DzString) LStrip(cutset string) DzString {
	// 将DzString类型转换为string类型
	return NewDzString(strings.TrimLeft(string(ds), cutset))
}

func (ds DzString) RStrip(cutset string) DzString {
	// 将DzString类型转换为string类型
	return NewDzString(strings.TrimRight(string(ds), cutset))
}

func (ds DzString) IsContains(subStr string) bool {
	// 将DzString类型转换为string类型
	return strings.Contains(string(ds), subStr)
}

func (ds DzString) Title() DzString {
	// 将DzString类型转换为string类型
	return NewDzString(cases.Title(language.Und).String(ds.ToString()))
}

func (ds DzString) SimilarText(str string) float64 {
	// 计算字符串相似度
	return 1 - float64(Algorithm.DzLevenshtein(ds.ToString(), str, 1, 1, 1))/float64(len([]rune(ds.ToString())))
}
