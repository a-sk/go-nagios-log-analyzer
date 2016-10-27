package interval

import (
	//"fmt"
	"math"
)

type Interval struct {
	Start int64
	End   int64
}

func (i Interval) Len() int64 {
	return i.End - i.Start
}
func (a Interval) Include(b Interval) bool {
	return b.End <= a.End && b.Start >= b.Start
}

func (i Interval) Union(other Interval) Interval {
	if !i.Overlaps(other) {
		return Interval{}
	}
	newStart := int64(math.Min(float64(i.Start), float64(other.Start)))
	newEnd := int64(math.Max(float64(i.End), float64(other.End)))
	return Interval{newStart, newEnd}
}
func (i Interval) Overlaps(other Interval) bool {
	if i.Start <= other.Start && i.End >= other.Start {
		return true
	}
	if i.Start >= other.Start && i.Start <= other.End {
		return true
	}
	return false
}

type IntervalSet struct {
	Intervals []Interval
}

func (is *IntervalSet) Add(other Interval) {
	if len(is.Intervals) == 0 {
		is.Intervals = append(is.Intervals, other)
	} else {
		for idx, interval := range is.Intervals {
			if interval.Overlaps(other) {
				is.Intervals[idx] = interval.Union(other)
				return
			}
		}
		is.Intervals = append(is.Intervals, other)
	}
}

func (is IntervalSet) Len() int64 {
	if len(is.Intervals) == 0 {
		return -1
	}
	intervalLen := int64(0)
	for _, interval := range is.Intervals {
		intervalLen += interval.Len()
	}
	return intervalLen
}
