package StringTool

import (
	"fmt"
	"testing"
)

func TestDzStrings_IsContain(t *testing.T) {
	strings := NewDzStrings([]string{
		"hello", "嘿嘿",
	})
	fmt.Println(strings.IsContain("he"))
	fmt.Println(strings.IsContain("嘿嘿"))
}

func TestDzStrings_Join(t *testing.T) {
	strings := NewDzStrings([]string{
		"hello", "嘿嘿",
	})
	fmt.Println(strings.Join(",").ToString())
	fmt.Println(strings.Join("").ToString())
	fmt.Println(strings.Join("/").ToString())
}
