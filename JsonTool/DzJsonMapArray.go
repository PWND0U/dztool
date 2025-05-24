package JsonTool

import (
	"encoding/json"
	"log"
)

type DzJsonMapArray []DzJsonMap

func NewDzJsonMapArray(jsonData []byte) DzJsonMapArray {
	djm := make(DzJsonMapArray, 0)
	err := json.Unmarshal(jsonData, &djm)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return djm
}
