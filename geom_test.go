// Copyright (c) 2011 Mateusz Czapliński
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package polyclip

import (
	"fmt"
	"math"
	"sort"
	. "testing"
)

func verify(t *T, cond bool, format string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Errorf(format, args...)
	}
}

func circa(f, g float64) bool {
	//TODO: (f-g)/g < 1e-6  ?
	return math.Abs(f-g) < 1e-6
}

func TestPoint(t *T) {
	verify(t, Point{0, 0}.Equals(Point{0, 0}), "Expected equal points")
	verify(t, Point{1, 2}.Equals(Point{1, 2}), "Expected equal points")
	verify(t, circa(Point{3, 4}.Length(), 5), "Expected length 5")
}

func TestEqualWithin(t *T) {
	p := Point{0, 0}
	verify(t, !p.equalWithin(Point{0, 1e-9}, 1e-10), "Expected not equal")
	verify(t, p.equalWithin(Point{1e-11, 1e-11}, 1e-10), "Expected equal")
}

func TestPointIsBefore(t *T) {
	cases := []struct {
		p1, p2 Point
		before bool
	}{
		{Point{0, 1}, Point{0, 1}, false},
		{Point{0, 1}, Point{1, 1}, true},
		{Point{0, 1}, Point{-1, 1}, false},
		{Point{0, 1}, Point{0, 2}, true},
		{Point{0, 1}, Point{0, -1}, false},
	}
	for i, v := range cases {
		verify(t, v.p1.isBefore(v.p2) == v.before, "Expected %v isBefore(%v)=%v (case %d)", v.p1, v.p2, v.before, i)
	}
}

func rect(x, y, w, h float64) Rectangle {
	return Rectangle{Min: Point{x, y}, Max: Point{x + w, y + h}}
}

func TestRectangleUnion(t *T) {
	cases := []struct{ a, b, result Rectangle }{
		{rect(0, 0, 20, 30), rect(0, 0, 30, 20), rect(0, 0, 30, 30)},
		{rect(10, 10, 10, 10), rect(-10, -10, 10, 10), rect(-10, -10, 30, 30)},
	}
	for i, v := range cases {
		u := v.a.union(v.b)
		r := v.result
		verify(t, u.Min.X == r.Min.X && u.Min.Y == r.Min.Y && u.Max.X == r.Max.X && u.Max.Y == r.Max.Y, "Expected equal rectangles in case %d", i)
	}
}

func TestRectangleIntersects(t *T) {
	r1 := rect(5, 5, 10, 10)
	cases := []struct {
		a, b   Rectangle
		result bool
	}{
		{rect(0, 0, 10, 20), rect(0, 10, 20, 10), true},
		{rect(0, 0, 10, 20), rect(20, 0, 10, 20), false},
		{rect(10, 50, 10, 10), rect(0, 0, 50, 45), false},
		{r1, rect(0, 0, 10, 10), true}, // diagonal intersections
		{r1, rect(10, 0, 10, 10), true},
		{r1, rect(0, 10, 10, 10), true},
		{r1, rect(10, 10, 10, 10), true},
		{r1, rect(-10, -10, 10, 10), false}, // non-intersecting rects on diagonal axes
		{r1, rect(20, -10, 10, 10), false},
		{r1, rect(-10, 20, 10, 10), false},
		{r1, rect(20, 20, 10, 10), false},
	}
	for i, v := range cases {
		verify(t, v.a.Overlaps(v.b) == v.result, "Expected result %v in case %d", v.result, i)
	}
}

func TestContourAdd(t *T) {
	c := Contour{}
	pp := []Point{{1, 2}, {3, 4}, {5, 6}}
	for i := range pp {
		c.Add(pp[i])
	}
	verify(t, len(c) == len(pp), "Expected all points in contour")
	for i := range pp {
		verify(t, c[i].Equals(pp[i]), "Wrong point at position %d", i)
	}
}

func TestContourBoundingBox(t *T) {
	// TODO
}

func TestContourSegment(t *T) {
	c := Contour([]Point{{1, 2}, {3, 4}, {5, 6}})
	segeq := func(s1, s2 segment) bool {
		return s1.start.Equals(s2.start) && s1.end.Equals(s2.end)
	}
	verify(t, segeq(c.segment(0), segment{Point{1, 2}, Point{3, 4}}), "Expected segment 0")
	verify(t, segeq(c.segment(1), segment{Point{3, 4}, Point{5, 6}}), "Expected segment 1")
	verify(t, segeq(c.segment(2), segment{Point{5, 6}, Point{1, 2}}), "Expected segment 2")
}

func TestContourSegmentError1(t *T) {
	c := Contour([]Point{{1, 2}, {3, 4}, {5, 6}})

	defer func() {
		verify(t, recover() != nil, "Expected error")
	}()
	_ = c.segment(3)
}

type pointresult struct {
	p      Point
	result bool
}

func TestContourContains(t *T) {
	var cases1 []pointresult
	c1 := Contour([]Point{{0, 0}, {10, 0}, {0, 10}})
	c2 := Contour([]Point{{0, 0}, {0, 10}, {10, 0}}) // opposite rotation
	cases1 = []pointresult{
		{Point{1, 1}, true},
		{Point{2, .1}, true},
		{Point{10, 10}, false},
		{Point{11, 0}, false},
		{Point{0, 11}, false},
		{Point{-1, -1}, false},
	}
	for i, v := range cases1 {
		verify(t, c1.Contains(v.p) == v.result, "Expected %v for point %d for c1", v.result, i)
		verify(t, c2.Contains(v.p) == v.result, "Expected %v for point %d for c2", v.result, i)
	}
}

func TestContourContains2(t *T) {
	c1 := Contour{{55, 35}, {25, 35}, {25, 119}, {55, 119}}
	c2 := Contour{{145, 35}, {145, 77}, {105, 77}, {105, 119}, {55, 119}, {55, 35}}
	cases := []struct {
		Contour
		Point
		Result bool
	}{
		{c1, Point{54.95, 77}, true},
		{c1, Point{55.05, 77}, false},
		{c2, Point{54.95, 77}, false},
		{c2, Point{55.05, 77}, true},
	}
	for _, c := range cases {
		result := c.Contour.Contains(c.Point)
		if result != c.Result {
			t.Errorf("case %v expected %v, got %v", c, c.Result, result)
		}
	}
}

func ExamplePolygon_Construct() {
	subject := Polygon{{{1, 1}, {1, 2}, {2, 2}, {2, 1}}} // small square
	clipping := Polygon{{{0, 0}, {0, 3}, {3, 0}}}        // overlapping triangle
	result := subject.Construct(INTERSECTION, clipping)

	out := []string{}
	for _, point := range result[0] {
		out = append(out, fmt.Sprintf("%v", point))
	}
	sort.Strings(out)
	fmt.Println(out)
	// Output: [{1 1} {1 2} {2 1}]
}
