package DateTool

import (
	"fmt"
	"testing"
	"time"
)

// ==================== 基本计时测试 ====================

func TestNewDzStopwatch(t *testing.T) {
	sw := NewDzStopwatch()
	fmt.Println("IsRunning:", sw.IsRunning())
	time.Sleep(50 * time.Millisecond)
	fmt.Println("ElapsedMs (>0):", sw.ElapsedMs())
	fmt.Println("ElapsedSeconds:", sw.ElapsedSeconds())
	fmt.Println("String:", sw.String())
}

func TestStopwatch_Stop(t *testing.T) {
	sw := NewDzStopwatch()
	time.Sleep(50 * time.Millisecond)
	sw.Stop()
	fmt.Println("After Stop IsRunning:", sw.IsRunning())
	ms1 := sw.ElapsedMs()
	fmt.Println("ElapsedMs after stop:", ms1)

	// 等一段时间，已停止的计时器不应继续计时
	time.Sleep(50 * time.Millisecond)
	ms2 := sw.ElapsedMs()
	fmt.Println("ElapsedMs after wait (should be same):", ms2)
	fmt.Println("Same:", ms1 == ms2)
}

func TestStopwatch_Start(t *testing.T) {
	sw := NewDzStopwatch()
	time.Sleep(30 * time.Millisecond)
	sw.Stop()
	ms1 := sw.ElapsedMs()
	fmt.Println("ElapsedMs after first stop:", ms1)

	// 继续计时
	sw.Start()
	time.Sleep(30 * time.Millisecond)
	ms2 := sw.ElapsedMs()
	fmt.Println("ElapsedMs after resume:", ms2)
	fmt.Println("Increased:", ms2 > ms1)
}

func TestStopwatch_Reset(t *testing.T) {
	sw := NewDzStopwatch()
	time.Sleep(50 * time.Millisecond)
	sw.Reset()
	fmt.Println("After Reset IsRunning:", sw.IsRunning())
	fmt.Println("After Reset ElapsedMs:", sw.ElapsedMs())
}

func TestStopwatch_Restart(t *testing.T) {
	sw := NewDzStopwatch()
	time.Sleep(50 * time.Millisecond)
	sw.Restart()
	fmt.Println("After Restart IsRunning:", sw.IsRunning())
	time.Sleep(50 * time.Millisecond)
	fmt.Println("After Restart ElapsedMs:", sw.ElapsedMs())
}

// ==================== 各种单位测试 ====================

func TestStopwatch_Units(t *testing.T) {
	sw := NewDzStopwatch()
	time.Sleep(100 * time.Millisecond)
	sw.Stop()

	fmt.Println("Elapsed:", sw.Elapsed())
	fmt.Println("ElapsedMs:", sw.ElapsedMs())
	fmt.Println("ElapsedSeconds:", sw.ElapsedSeconds())
	fmt.Println("ElapsedMinutes:", sw.ElapsedMinutes())
	fmt.Println("ElapsedHours:", sw.ElapsedHours())
	fmt.Println("ElapsedDays:", sw.ElapsedDays())
	fmt.Println("ElapsedWeeks:", sw.ElapsedWeeks())
}

// ==================== 暂停/继续累积测试 ====================

func TestStopwatch_Cumulative(t *testing.T) {
	sw := NewDzStopwatch()

	// 第一段
	time.Sleep(50 * time.Millisecond)
	sw.Stop()
	first := sw.ElapsedMs()
	fmt.Println("First segment ms:", first)

	// 第二段
	sw.Start()
	time.Sleep(50 * time.Millisecond)
	sw.Stop()
	total := sw.ElapsedMs()
	fmt.Println("Total ms after two segments:", total)
	fmt.Println("Total > first:", total > first)
}

// ==================== 格式化测试 ====================

func TestStopwatch_Format(t *testing.T) {
	sw := NewDzStopwatch()
	time.Sleep(100 * time.Millisecond)
	sw.Stop()

	fmt.Println("Format default:", sw.Format(""))
	fmt.Println("Format ms:", sw.Format("ms"))
	fmt.Println("Format sec:", sw.Format("sec"))
	fmt.Println("Format full:", sw.Format("full"))
}

// ==================== 链式调用测试 ====================

func TestStopwatch_Chained(t *testing.T) {
	sw := NewDzStopwatch().Stop().Reset().Restart()
	time.Sleep(50 * time.Millisecond)
	sw.Stop()
	fmt.Println("Chained ElapsedMs:", sw.ElapsedMs())
	fmt.Println("Chained IsRunning:", sw.IsRunning())
}
