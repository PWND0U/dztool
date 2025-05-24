package StringTool

import "strings"

type DzStrings []DzString

func NewDzStrings(srcStr []string) DzStrings {
	// 创建一个空的DzString类型的切片
	var listDzString = make(DzStrings, 0)
	for _, s := range srcStr {
		// 将分割后的字符串转换为DzString类型，并添加到切片中
		listDzString = append(listDzString, NewDzString(s))
	}
	return listDzString
}

func (lds DzStrings) Join(sepStr string) DzString {
	var listString = make([]string, 0)
	for _, ds := range lds {
		listString = append(listString, ds.ToString())
	}
	return DzString(strings.Join(listString, sepStr))
}

func (lds DzStrings) IsContain(containStr string) bool {
	for _, ds := range lds {
		if string(ds) == containStr {
			return true
		}
	}
	return false
}
