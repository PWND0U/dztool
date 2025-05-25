package JsonTool

import (
	"encoding/json"
)

type DzJsonMapArray []DzJsonMap

func NewDzJsonMapArray(jsonData []byte) DzJsonMapArray {
	djm := make(DzJsonMapArray, 0)
	err := json.Unmarshal(jsonData, &djm)
	if err != nil {
		//log.Fatal(err)
		return djm
	}
	return djm
}
