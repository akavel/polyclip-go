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

func ExamplePoint_Add() {
	p := Point{0.8, 2.5}
	q := Point{0.2, 3}
	fmt.Println(p.Add(q))
	// Output: (1.00,5.50)
}

func ExamplePoint_Sub() {
	p := Point{0.8, 5}
	q := Point{0.2, 0.5}
	fmt.Println(p.Sub(q))
	// Output: (0.60,4.50)
}

func ExamplePoint_Mul() {
	p := Point{0.2, 0.5}
	fmt.Println(p.Mul(2))
	// Output: (0.40,1.00)
}

func ExamplePoint_Div() {
	p := Point{0.2, 0.5}
	fmt.Println(p.Div(2))
	// Output: (0.10,0.25)
}

func ExamplePoint_In() {
	p := Point{3, 4}
	q := Point{0, 0}
	r := Rect(1, 2, 9, 8)
	fmt.Println(p.In(r))
	fmt.Println(q.In(r))
	// Output:
	// true
	// false
}

func ExamplePoint_Equals() {
	p := Point{3, 4}
	q := Point{0, 0}
	r := Pt(3, 4)
	fmt.Println(p.Equals(r))
	fmt.Println(q.Equals(r))
	// Output:
	// true
	// false
}

// Extra functions (not supported by image.Point)

func ExamplePoint_Angle() {
	p := Point{1, 0}
	fmt.Println(p.Angle(Point{1, 1}))  // 45
	fmt.Println(p.Angle(Point{1, -1})) // -45
	// Output:
	// 0.7853981633974484
	// -0.7853981633974484
}

func ExamplePoint_Angle3() {
	p := Point{1, 1}
	q := Point{2, 1}
	fmt.Println(p.Angle3(q, Point{2, 2}))
	fmt.Println(p.Angle3(q, Point{2, 0}))
	// Output:
	// 0.7853981633974484
	// -0.7853981633974484
}

func ExamplePoint_Append() {
	p := Point{1, 2}
	x := []float64{0, 1, 2}
	y := []float64{3, 4, 5}
	x, y = p.Append(x, y)
	fmt.Println(x)
	fmt.Println(y)
	// Output:
	// [0 1 2 1]
	// [3 4 5 2]
}

func ExamplePoint_Cross() {
	p := Point{1, 2}
	q := Point{3, 4}
	fmt.Println(p.Cross(q))
	// Output: -2
}

func ExamplePoint_Dist() {
	p := Point{1, 1}
	q := Point{1, 5}
	fmt.Println(p.Dist(q))
	// Output: 4
}

func ExamplePoint_Dist2() {
	p := Point{1, 1}
	q := Point{1, 5}
	fmt.Println(p.Dist2(q))
	// Output: 16
}

func ExamplePoint_Dot() {
	p := Point{1, 2}
	q := Point{3, 4}
	fmt.Println(p.Dot(q))
	// Output: 11
}

func ExamplePoint_MaxRadius() {
	rect := Rect(1, 2, 6, 9)
	p := Point{1, 2}
	fmt.Println(p.MaxRadius(rect))
	// Output: 8.602325267042627
}

func ExamplePoint_Norm() {
	p := Point{1, 2}
	fmt.Println(p.Norm() == math.Sqrt(5))
	// Output: true
}

func ExamplePoint_Norm2() {
	p := Point{1, 2}
	fmt.Println(p.Norm2())
	// Output: 5
}

func ExamplePoint_Normalize() {
	p := Point{1, 2}
	fmt.Println(p.Normalize())
	fmt.Println(p.Normalize() == p.Normalize(1))
	fmt.Println(p.Normalize(2))
	// Output:
	// (0.45,0.89)
	// true
	// (0.89,1.79)
}

func ExamplePoint_Negate() {
	p := Point{1, 2}
	fmt.Println(p.Negate())
	// Output: (-1.00,-2.00)
}

func ExamplePoint_Normals() {
	p := Point{1, 1}
	q := Point{3, 3}
	fmt.Println(p.Normals(q))
	// Output: (-2.00,2.00) (2.00,-2.00)
}

func ExamplePoint_Polar() {
	p := Point{3, 3}
	fmt.Println(p.Polar(2, math.Pi/2))
	// Output: (3.00,5.00)
}

func ExamplePoint_Prepend() {
	p := Point{1, 2}
	x := []float64{0, 1, 2}
	y := []float64{3, 4, 5}
	x, y = p.Prepend(x, y)
	fmt.Println(x)
	fmt.Println(y)
	// Output:
	// [1 0 1 2]
	// [2 3 4 5]
}

func ExamplePoint_Rect() {
	p := Point{1, 2}
	fmt.Println(p.Rect(8, 6))
	// Output: (-3.00,-1.00)-(5.00,5.00)
}

func ExamplePoint_Tangent() {
	p := Point{1, 2}
	q := Point{3, 2}
	fmt.Println(p.Tangent(q))
	// Output: (1.00,0.00)
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
		u := v.a.Union(v.b)
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

func ExampleRectangle_String() {
	fmt.Println(Rect(1, 2, 3, 4))
	// Output: (1.00,2.00)-(3.00,4.00)
}

func ExampleRectFromPoints() {
	fmt.Println(RectFromPoints())
	fmt.Println(RectFromPoints(Pt(2, 4)))
	fmt.Println(RectFromPoints(Pt(0, 0), Pt(2, 4), Pt(-1, -2)))
	// Output:
	// (0.00,0.00)-(0.00,0.00)
	// (2.00,4.00)-(2.00,4.00)
	// (-1.00,-2.00)-(2.00,4.00)
}

func ExampleRectangle_Dx() {
	rect := Rect(1, 2, 3, 4)
	fmt.Println(rect.Dx())
	// Output: 2
}

func ExampleRectangle_Dy() {
	rect := Rect(1, 2, 3, 5)
	fmt.Println(rect.Dy())
	// Output: 3
}

func ExampleRectangle_Size() {
	rect := Rect(1, 2, 3, 5)
	fmt.Println(rect.Size())
	// Output: (2.00,3.00)
}

func ExampleRectangle_Add() {
	rect := Rect(1, 2, 3, 5)
	p := Point{6, 7}
	fmt.Println(rect.Add(p))
	// Output: (7.00,9.00)-(9.00,12.00)
}

func ExampleRectangle_Sub() {
	rect := Rect(1, 2, 3, 5)
	p := Point{6, 7}
	fmt.Println(rect.Sub(p))
	// Output: (-5.00,-5.00)-(-3.00,-2.00)
}

func ExampleRectangle_Inset() {
	rect := Rect(1, 2, 3, 5)
	fmt.Println(rect.Inset(0.25))
	fmt.Println(rect.Inset(2))
	// Output:
	// (1.25,2.25)-(2.75,4.75)
	// (2.00,3.50)-(2.00,3.50)
}

func ExampleRectangle_Intersect() {
	r := Rect(1, 2, 3, 5)
	s := Rect(0, 0, 2, 3)
	t := Rect(6, 7, 8, 9)
	fmt.Println(r.Intersect(s))
	fmt.Println(r.Intersect(t) == ZR)
	// Output:
	// (1.00,2.00)-(2.00,3.00)
	// true
}

func ExampleRectangle_Union() {
	r := Rect(1, 2, 3, 5)
	s := Rect(0, 0, 3, 3)
	t := Rect(6, 7, 8, 9)
	fmt.Println(r.Union(s))
	fmt.Println(r.Union(t))
	// Output:
	// (0.00,0.00)-(3.00,5.00)
	// (1.00,2.00)-(8.00,9.00)
}

func ExampleRectangle_Empty() {
	fmt.Println(Rectangle{Min: Point{7, 8}, Max: Point{1, 2}}.Empty())
	fmt.Println(Rect(7, 8, 1, 2).Empty())
	// Output:
	// true
	// false
}

func ExampleRectangle_Eq() {
	r := Rectangle{Min: Point{1, 2}, Max: Point{7, 8}}
	s := Rect(7, 8, 1, 2)
	fmt.Println(r.Equals(s))
	// Output:
	// true
}

func ExampleRectangle_Overlaps() {
	r := Rect(1, 2, 3, 5)
	s := Rect(0, 0, 2, 3)
	t := Rect(6, 7, 8, 9)
	fmt.Println(r.Overlaps(s))
	fmt.Println(r.Overlaps(t))
	// Output:
	// true
	// false
}

func ExampleRectangle_In() {
	r := Rect(1, 2, 4, 5)
	s := Rect(2, 3, 3, 4)
	t := Rectangle{Min: Point{7, 8}, Max: Point{1, 2}}
	fmt.Println(s.In(r))
	fmt.Println(t.In(r))
	// Output:
	// true
	// true
}

func ExampleRectangle_Canon() {
	r := Rectangle{Min: Point{7, 8}, Max: Point{1, 2}}
	fmt.Println(r.Canon())
	// Output: (1.00,2.00)-(7.00,8.00)
}

func ExampleRectangle_Diagonal() {
	r := Rect(0, 0, 2, 2)
	fmt.Println(r.Diagonal())
	// Output: 2.8284271247461903
}

func ExampleRectangle_MaxRadius() {
	rect := Rect(1, 2, 6, 9)
	p := Point{1, 2}
	fmt.Println(rect.MaxRadius(p))
	// Output: 8.602325267042627
}

func ExampleRectangle_Offset() {
	rect := Rect(1, 2, 6, 9)
	fmt.Println(rect.Offset(4))
	// Output: (-3.00,-2.00)-(10.00,13.00)
}

func ExampleRectangle_Points() {
	rect := Rect(1, 2, 6, 9)
	fmt.Println(rect.Points())
	// Output: [(1.00,2.00) (1.00,9.00) (6.00,9.00) (6.00,2.00)]
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
	// [(1.00,1.00) (1.00,2.00) (2.00,1.00)]
}
