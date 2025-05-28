package Algorithm

import (
	"fmt"
	"testing"
)

func TestLevenshtein(t *testing.T) {
	fmt.Println(Levenshtein("你好世界", "你好中国", 1, 1, 1))
}
