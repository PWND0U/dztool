package JsonTool

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestDzJsonMap(t *testing.T) {
	f, _ := os.Open("G:\\code\\go\\dztool\\JsonTool\\static\\trend.json")
	all, _ := io.ReadAll(f)
	dzm := NewDzJsonMap(all)
	for _, s := range dzm.GetArray("commodity.data").GetMapArray() {
		fmt.Println(s.GetString("name"))
		fmt.Println(s.GetArray("data").GetStringArray())
	}
	for _, s := range dzm.GetArray("seller.data").GetMapArray() {
		fmt.Println(s.GetString("name"))
		fmt.Println(s.GetArray("data").GetStringArray())
	}
	for _, s := range dzm.GetArray("map.data").GetMapArray() {
		fmt.Println(s.GetString("name"))
		fmt.Println(s.GetArray("data").GetStringArray())
	}
	fmt.Println(dzm.GetInt("code"))
	fmt.Println(dzm.GetInt("status"))
	fmt.Println(dzm.GetInt("ret"))
	fmt.Println(dzm.GetArray("common.month").GetStringArray()[0])
	fmt.Println(dzm.GetArray("map.data").GetMapArray()[0].GetString("name"))
}
