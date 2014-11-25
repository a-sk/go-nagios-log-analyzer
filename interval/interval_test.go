package interval

import (
	//"fmt"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestInterval(t *testing.T) {
	a := Interval{1, 5}
	expect(t, a.Start, 1)
	expect(t, a.End, 5)

}
func TestIntervalLen(t *testing.T) {
	a := Interval{1, 5}
	expect(t, a.Len(), 4)
}

func TestIntervalInclude(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 4}
	expect(t, a.Include(b), true)
}

func TestIntervalIncludeFalse(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 4}
	expect(t, b.Include(a), false)
}

func TestIntervalOverlaps(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 6}
	expect(t, b.Overlaps(a), true)
}
func TestIntervalOverlapsReverse(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 6}
	expect(t, a.Overlaps(b), true)
}
func TestIntervalOverlapsFalse(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{8, 10}
	expect(t, a.Overlaps(b), false)
}
func TestIntervalOverlapsFalseReverse(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{8, 10}
	expect(t, b.Overlaps(a), false)
}
func TestIntervalSetAdd(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 6}
	c := IntervalSet{}
	c.Add(a)
	c.Add(b)
	expect(t, c.Len(), 5)
}

func TestIntervalSetAddPointersMap(t *testing.T) {
	dict := make(map[string]*IntervalSet)
	a := Interval{1, 5}
	b := Interval{2, 6}
	dict["a"] = &IntervalSet{}
	c := dict["a"]
	c.Add(a)
	c.Add(b)

	expect(t, c.Len(), 5)
}
func TestIntervalSetAddMore(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 7}
	c := Interval{2, 3}
	d := IntervalSet{}
	d.Add(a)
	d.Add(b)
	d.Add(c)
	expect(t, d.Len(), 6)
}
func TestIntervalSetAddEvenMore(t *testing.T) {
	a := IntervalSet{}
	a.Add(Interval{1405225899, 1405243000})
	a.Add(Interval{1405240900, 1405243000})
	a.Add(Interval{1405277920, 1405278522})
	a.Add(Interval{1405279960, 1405289558})
	expect(t, a.Len(), 27301)
}

func TestIntervalUnion(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 6}
	expect(t, a.Union(b), Interval{1, 6})
}

func TestIntervalUnionReverse(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{2, 6}
	expect(t, b.Union(a), Interval{1, 6})
}
func TestIntervalUnionNotOverlaps(t *testing.T) {
	a := Interval{1, 5}
	b := Interval{8, 10}
	expect(t, a.Union(b), Interval{})
}

//func TestIntervalSetLen(t *testing.T) {
//a := Interval{1, 5}
//b := Interval{2, 6}
//c := IntervalSet{a, b}
//expect(t, c.Len(), 5)
//}

//func TestIntervalSetLenEmpty(t *testing.T) {
//c := IntervalSet{}
//expect(t, c.Len(), -1)
//}
