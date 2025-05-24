package main

import (
	"fmt"
	"github.com/PWND0U/dztool/JsonTool"
	"io"
	"os"
)

func main() {
	//fmt.Println("create pro")
	//fmt.Println(StringTool.NewDzString("1").ToString())
	f, _ := os.Open("G:\\code\\go\\dztool\\JsonTool\\static\\trend.json")
	all, _ := io.ReadAll(f)
	dzm := JsonTool.NewDzJsonMap(all)
	for _, s := range dzm.GetMap("commodity").GetArray("data").GetMapArray() {
		fmt.Println(s.GetString("name"))
		fmt.Println(s.GetArray("data").GetStringArray())
	}
	for _, s := range dzm.GetMap("seller").GetArray("data").GetMapArray() {
		fmt.Println(s.GetString("name"))
		fmt.Println(s.GetArray("data").GetStringArray())
	}
	for _, s := range dzm.GetMap("map").GetArray("data").GetMapArray() {
		fmt.Println(s.GetString("name"))
		fmt.Println(s.GetArray("data").GetStringArray())
	}
}
