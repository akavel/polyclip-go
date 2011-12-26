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

package polyutil

import (
	. "testing"
)

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func TestLine(t *T) {
	cases := []struct{ x0, y0, x1, y1, sx, sy int }{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 1, 0},
		{0, 0, -1, 1, -1, 1},
		{-2, -2, 2, 2, 1, 1},
		{10, 2, 20, 4, 1, 1},
		{2, 10, 10, 5, 1, -1},
		{1000, -50, -100, -3, -1, 1},
		{10, 20, 0, -30, -1, -1},
	}
	for i, c := range cases {
		p := 0
		lastx, lasty := c.x0, c.y0

		drawline(c.x0, c.y0, c.x1, c.y1, func(x, y int) {
			verify(t, x == lastx || x == (lastx+c.sx), "Line %d, point %d, got x=%d, expected %d or %d", i, p, x, lastx, lastx+c.sx)
			verify(t, y == lasty || y == (lasty+c.sy), "Line %d, point %d, got y=%d, expected %d or %d", i, p, y, lasty, lasty+c.sy)

			lastx, lasty = x, y
			p++
		})

		verify(t, lastx == c.x1, "After line %d, got x==%d, expected %d", i, lastx, c.x1)
		verify(t, lasty == c.y1, "After line %d, got y==%d, expected %d", i, lasty, c.y1)
		pn := max(abs(c.x1-c.x0), abs(c.y1-c.y0)) + 1
		verify(t, p == pn, "After line %d, got %d points, expected %d", i, p, pn)
	}

}
