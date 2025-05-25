package TimeIntervalTool

import (
	"maps"
	"slices"
	"time"
)

// TimeIntervalTool Timer module

type DzTimeInterval struct {
	timerMap map[string]time.Time
}

func NewDzTimeInterval() *DzTimeInterval {
	return &DzTimeInterval{make(map[string]time.Time)}
}

func (dti *DzTimeInterval) Start(timerName ...string) *DzTimeInterval {
	if len(timerName) <= 0 {
		dti.timerMap["default"] = time.Now()
	} else {
		for _, tn := range timerName {
			dti.timerMap[tn] = time.Now()
		}
	}
	return dti
}

func (dti *DzTimeInterval) ReStart(timerName ...string) *DzTimeInterval {
	dti.Start(timerName...)
	return dti
}

func (dti *DzTimeInterval) IntervalMs(timerName ...string) float64 {
	if len(timerName) <= 0 {
		return float64(time.Since(dti.timerMap["default"]).Milliseconds())
	} else {
		if slices.Contains(slices.Collect(maps.Keys(dti.timerMap)), timerName[0]) {
			return float64(time.Since(dti.timerMap[timerName[0]]).Milliseconds())
		}
	}
	return 0
}

func (dti *DzTimeInterval) IntervalSecond(timerName ...string) float64 {
	return dti.IntervalMs(timerName...) / 1000
}

func (dti *DzTimeInterval) IntervalMinute(timerName ...string) float64 {
	return dti.IntervalSecond(timerName...) / 60
}

func (dti *DzTimeInterval) IntervalHour(timerName ...string) float64 {
	return dti.IntervalMinute(timerName...) / 60
}

func (dti *DzTimeInterval) IntervalDay(timerName ...string) float64 {
	return dti.IntervalHour(timerName...) / 24
}

func (dti *DzTimeInterval) IntervalWeek(timerName ...string) float64 {
	return dti.IntervalHour(timerName...) / 24 / 7
}
