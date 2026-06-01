package DateTool

import (
	"fmt"
	"time"
)

// DzStopwatch 简单的计时器，用于计算某段代码的执行时间。
// 支持暂停/继续，提供毫秒、秒、分、时、天、周等各种单位的已用时长计算。
// 使用指针类型，因为内部持有可变状态，不可复制。
type DzStopwatch struct {
	startTime time.Time
	running   bool
	elapsed   time.Duration // 暂停前累积的时间
}

// ==================== 构造函数 ====================

// NewDzStopwatch 创建并自动启动一个计时器
func NewDzStopwatch() *DzStopwatch {
	return &DzStopwatch{
		startTime: time.Now(),
		running:   true,
		elapsed:   0,
	}
}

// ==================== 链式控制方法 ====================

// Start 启动或继续计时器，返回自身支持链式调用
func (s *DzStopwatch) Start() *DzStopwatch {
	if !s.running {
		s.startTime = time.Now()
		s.running = true
	}
	return s
}

// Stop 暂停计时器，返回自身支持链式调用
func (s *DzStopwatch) Stop() *DzStopwatch {
	if s.running {
		s.elapsed += time.Since(s.startTime)
		s.running = false
	}
	return s
}

// Reset 重置计时器（清零并停止），返回自身支持链式调用
func (s *DzStopwatch) Reset() *DzStopwatch {
	s.elapsed = 0
	s.running = false
	return s
}

// Restart 重新开始计时器（清零并启动），返回自身支持链式调用
func (s *DzStopwatch) Restart() *DzStopwatch {
	s.elapsed = 0
	s.startTime = time.Now()
	s.running = true
	return s
}

// ==================== 读取方法 ====================

// Elapsed 返回已用时间
func (s *DzStopwatch) Elapsed() time.Duration {
	if s.running {
		return s.elapsed + time.Since(s.startTime)
	}
	return s.elapsed
}

// ElapsedMs 返回已用毫秒数
func (s *DzStopwatch) ElapsedMs() int64 {
	return s.Elapsed().Milliseconds()
}

// ElapsedSeconds 返回已用秒数
func (s *DzStopwatch) ElapsedSeconds() float64 {
	return s.Elapsed().Seconds()
}

// ElapsedMinutes 返回已用分钟数
func (s *DzStopwatch) ElapsedMinutes() float64 {
	return s.Elapsed().Minutes()
}

// ElapsedHours 返回已用小时数
func (s *DzStopwatch) ElapsedHours() float64 {
	return s.Elapsed().Hours()
}

// ElapsedDays 返回已用天数
func (s *DzStopwatch) ElapsedDays() float64 {
	return s.Elapsed().Hours() / 24
}

// ElapsedWeeks 返回已用周数
func (s *DzStopwatch) ElapsedWeeks() float64 {
	return s.Elapsed().Hours() / (24 * 7)
}

// IsRunning 返回计时器是否正在运行
func (s *DzStopwatch) IsRunning() bool {
	return s.running
}

// ==================== 格式化 ====================

// String 实现 fmt.Stringer 接口
func (s *DzStopwatch) String() string {
	return s.Elapsed().String()
}

// Format 格式化输出已用时间，使用 Go 标准 Duration 格式
func (s *DzStopwatch) Format(layout string) string {
	d := s.Elapsed()
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	sec := int(d.Seconds()) % 60
	ms := d.Milliseconds() % 1000

	switch layout {
	case "ms":
		return fmt.Sprintf("%dms", d.Milliseconds())
	case "sec":
		return fmt.Sprintf("%.3fs", d.Seconds())
	case "full":
		return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, sec, ms)
	default:
		return d.String()
	}
}
