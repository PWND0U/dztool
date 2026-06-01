package JsonTool

import (
	"github.com/bytedance/sonic"
)

type DzJsonMapArray []DzJsonMap

func NewDzJsonMapArray(jsonData []byte) DzJsonMapArray {
	djm := make(DzJsonMapArray, 0)
	err := sonic.Unmarshal(jsonData, &djm)
	if err != nil {
		//log.Fatal(err)
		return djm
	}
	return djm
}
