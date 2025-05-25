package JsonTool

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestDzJsonMapArray(t *testing.T) {
	f, _ := os.Open("G:\\code\\go\\dztool\\JsonTool\\static\\budget.json")
	all, _ := io.ReadAll(f)
	jsonMaps := NewDzJsonMapArray(all)
	fmt.Println(jsonMaps)
	for _, jsonMap := range jsonMaps {
		fmt.Println(jsonMap.GetString("dimZhName"))
		fmt.Println(jsonMap.GetNumber("budget"))
		fmt.Println(jsonMap.GetString("dimName"))
		fmt.Println(jsonMap.GetNumber("expense"))
	}
}
