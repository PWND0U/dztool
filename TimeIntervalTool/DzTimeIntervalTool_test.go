package TimeIntervalTool

import (
	"fmt"
	"testing"
	"time"
)

func TestDzTimeInterval(t *testing.T) {
	timer := NewDzTimeInterval().Start("a")
	time.Sleep(time.Second)
	fmt.Println(timer.IntervalSecond("a"))
	timer.Start("b")
	time.Sleep(time.Second)
	fmt.Println(timer.IntervalHour("a"))
	fmt.Println(timer.IntervalSecond("b"))
	timer.ReStart("a", "b")
	timer.Start()
	time.Sleep(time.Second)
	fmt.Println(timer.IntervalSecond("a"))
	fmt.Println(timer.IntervalHour("b"))
	fmt.Println(timer.IntervalSecond())
}
