package StringTool

import (
	"fmt"
	"testing"
)

func TestDzString_SimilarText(t *testing.T) {
	fmt.Println(NewDzString("你好中国").SimilarText("你好世界"))
	fmt.Println(len("你好世界"))
	fmt.Println(len([]rune("你好世界")))
	fmt.Println(len("你好中国"))
	fmt.Println(NewDzString("ni hao ha").Title())
	fmt.Println(NewDzString("你好中国").Title())
}
