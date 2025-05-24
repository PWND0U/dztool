package main

import (
	"fmt"
	"github.com/PWND0U/dztool/BytesTool"
)

func main() {
	fmt.Println("create pro")
	fmt.Println(BytesTool.NewDzString("你好世界").Find("哈哈"))
}
