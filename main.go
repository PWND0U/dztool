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
	fmt.Println(dzm.GetNumber("code"))
	fmt.Println(dzm.GetNumber("status"))
	fmt.Println(dzm.GetNumber("ret"))
	fmt.Println(dzm.GetArray("common.month").GetStringArray()[0])
	fmt.Println(dzm.GetArray("map.data").GetMapArray()[0].GetString("name"))
	f, _ = os.Open("G:\\code\\go\\dztool\\JsonTool\\static\\budget.json")
	all, _ = io.ReadAll(f)
	jsonMaps := JsonTool.NewDzJsonMapArray(all)
	fmt.Println(jsonMaps)
	for _, jsonMap := range jsonMaps {
		fmt.Println(jsonMap.GetString("dimZhName"))
		fmt.Println(jsonMap.GetNumber("budget"))
		fmt.Println(jsonMap.GetString("dimName"))
		fmt.Println(jsonMap.GetNumber("expense"))
	}
}
