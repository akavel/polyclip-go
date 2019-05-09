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
	. "testing"
)

func TestAbove(t *T) {
	cases := []struct {
		left, right Point
		result      bool
		x           Point
	}{
		{Point{0, 1}, Point{2, 1}, true, Point{1, 0}},
		{Point{0, 1}, Point{2, 1}, false, Point{1, 3}},
	}
	for i, v := range cases {
		e := &endpoint{p: v.left, left: true, other: &endpoint{p: v.right, left: false}}
		verify(t, e.above(v.x) == v.result, "Expected %v above %v (case %d/a)", e, v.x, i)
		e = &endpoint{p: v.right, left: false, other: &endpoint{p: v.left, left: true}}
		verify(t, e.above(v.x) == v.result, "Expected %v above %v (case %d/b)", e, v.x, i)
	}
}

const (
	left    = true
	right   = false
	valid   = true
	invalid = false
)

func TestEndpointIsValidDirection(t *T) {
	cases := []struct {
		left, right Point
		dir         bool
		isValid     bool
	}{
		{Point{0, 1}, Point{0, 1}, left, invalid},  // Zero-length
		{Point{0, 1}, Point{0, 1}, right, invalid}, // Zero-length
		{Point{0, 1}, Point{1, 1}, left, valid},    // Horizontally valid
		{Point{0, 1}, Point{-1, 1}, right, valid},  // Horizontally valid
		{Point{0, 1}, Point{-1, 1}, left, invalid}, // Horizontally invalid
		{Point{0, 1}, Point{1, 1}, right, invalid}, // Horizontally invalid
		{Point{0, 1}, Point{0, 2}, left, valid},    // Vertically valid
		{Point{0, 1}, Point{0, -1}, right, valid},  // Vertically valid
		{Point{0, 1}, Point{0, -1}, left, invalid}, // Vertically invalid
		{Point{0, 1}, Point{0, 2}, right, invalid}, // Vertically invalid
	}
	for i, v := range cases {
		e := &endpoint{p: v.left, left: v.dir == left, other: &endpoint{p: v.right, left: v.dir != left}}
		verify(t, e.isValidDirection() == v.isValid, "Expected %v isValidDirection()=%v (case %d)", e, v.isValid, i)
	}
}

func TestInvalidSingleIntersection(t *T) {
	cases := []struct {
		l1, r1, l2, r2, intersection Point
		isValid                      bool
	}{
		{
			Point{0, 1.00000000000000}, Point{1, 1},
			Point{0, 1.00000000000001}, Point{3, 2},
			Point{0, 1.00000000000002},
			invalid,
		},
		{
			Point{0, 1.00000000000001}, Point{1, 1},
			Point{0, 1.00000000000002}, Point{3, 2},
			Point{0, 1.00000000000000},
			invalid,
		},
		{
			Point{0, 1.00000000000000}, Point{1, 1},
			Point{0, 1.00000000000001}, Point{3, 2},
			Point{0, 1.00000000000002},
			invalid,
		},
		{
			Point{0, 1.00000000000000}, Point{1, 1},
			Point{0, 1.00000000000002}, Point{3, 2},
			Point{0, 1.00000000000001},
			valid,
		},
		{
			Point{1.00000000000000, 0}, Point{1, 1},
			Point{1.00000000000001, 0}, Point{3, 2},
			Point{1.00000000000002, 0},
			invalid,
		},
		{
			Point{1.00000000000001, 0}, Point{1, 1},
			Point{1.00000000000002, 0}, Point{3, 2},
			Point{1.00000000000000, 0},
			invalid,
		},
		{
			Point{1.00000000000000, 0}, Point{1, 1},
			Point{1.00000000000001, 0}, Point{3, 2},
			Point{1.00000000000002, 0},
			invalid,
		},
		{
			Point{1.00000000000000, 0}, Point{1, 1},
			Point{1.00000000000002, 0}, Point{3, 2},
			Point{1.00000000000001, 0},
			valid,
		},
	}
	for i, v := range cases {
		e1 := &endpoint{p: v.l1, left: true, other: &endpoint{p: v.r1, left: false}}
		e2 := &endpoint{p: v.l2, left: true, other: &endpoint{p: v.r2, left: false}}
		verify(t, isValidSingleIntersection(e1, e2, v.intersection) == v.isValid,
			"Case %d: Expected intersection at %v of (%v, %v) to be isValidSingleIntersection()=%v", i, v.intersection, e1, e2, v.isValid)
		e3 := &endpoint{p: v.r1, left: false, other: &endpoint{p: v.l1, left: true}}
		e4 := &endpoint{p: v.r2, left: false, other: &endpoint{p: v.l2, left: true}}
		verify(t, isValidSingleIntersection(e3, e4, v.intersection) == v.isValid,
			"Case %d: Expected intersection at %v of (%v, %v) to be isValidSingleIntersection()=%v", i, v.intersection, e3, e4, v.isValid)
		e5 := &endpoint{p: v.l1, left: false, other: &endpoint{p: v.r1, left: true}}
		e6 := &endpoint{p: v.l2, left: false, other: &endpoint{p: v.r2, left: true}}
		verify(t, isValidSingleIntersection(e5, e6, v.intersection) == v.isValid,
			"Case %d: Expected intersection at %v of (%v, %v) to be isValidSingleIntersection()=%v", i, v.intersection, e5, e6, v.isValid)
		e7 := &endpoint{p: v.r1, left: true, other: &endpoint{p: v.l1, left: false}}
		e8 := &endpoint{p: v.r2, left: true, other: &endpoint{p: v.l2, left: false}}
		verify(t, isValidSingleIntersection(e7, e8, v.intersection) == v.isValid,
			"Case %d: Expected intersection at %v of (%v, %v) to be isValidSingleIntersection()=%v", i, v.intersection, e7, e8, v.isValid)
	}

}
