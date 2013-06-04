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
