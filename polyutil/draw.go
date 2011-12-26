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

// Package polyutil contains very simple utility functions for drawing and (de)serialization of polygons.
package polyutil

import "bitbucket.org/akavel/polyclip.go"

// Putpixel describes a function expected to draw a point on a bitmap at (x, y) coordinates.
type Putpixel func(x, y int)

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// Bresenham's algorithm, http://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
// TODO: handle int overflow etc.
func drawline(x0, y0, x1, y1 int, brush Putpixel) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		brush(x0, y0)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// DrawPolyline is a simple function for drawing a series of rasterized lines,
// connecting the points specified in pts (treated as a closed polygon).
// This function uses a very basic implementation of the Bresenham's algorithm
// (http://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm),
// thus with no antialiasing. Moreover, the coordinates of the nodes are rounded
// down towards the nearest integer.
// The computed points are passed to brush function for final rendering.
func DrawPolyline(pts []polyclip.Point, brush Putpixel) {
	last := len(pts) - 1
	for i := 0; i < len(pts); i++ {
		drawline(int(pts[last].X), int(pts[last].Y), int(pts[i].X), int(pts[i].Y), brush)
		last = i
	}
}
