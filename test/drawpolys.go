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

package test

import (
	"github.com/akavel/polyclip.go"
	"github.com/akavel/polyclip.go/polyutil"
	"image"
	"image/color"
	"image/draw"
)

func max(a float64, b ...float64) float64 {
	for _, f := range b {
		if f > a {
			a = f
		}
	}
	return a
}

func brush(img draw.Image, id int) func(x, y int) {
	colors := [][]uint8{
		{0, 0, 0xff},
		{0xff, 0, 0},
		{0xff, 0xff, 0xff},
	}
	c := colors[id%len(colors)]
	return func(x, y int) {
		img.Set(x, y, &color.NRGBA{R: c[0], G: c[1], B: c[2], A: 0xff})
	}
}

func safebbox(p polyclip.Polygon) polyclip.Rectangle {
	if len(p) == 0 {
		return polyclip.Rectangle{}
	}
	return p.BoundingBox()
}

// Warning: does modify contents of polys
func DrawPolygons(mul float64, polys []polyclip.Polygon) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 0, 0))
	min := polyclip.Point{0, 0}

	for i, polygon := range polys {
		r0 := safebbox(polygon)
		translate := image.Point{}
		if r0.Min.X < min.X {
			translate.X = int(mul*min.X) - int(mul*r0.Min.X)
			min.X = r0.Min.X
		}
		if r0.Min.Y < min.Y {
			translate.Y = int(mul*min.Y) - int(mul*r0.Min.Y)
			min.Y = r0.Min.Y
		}

		for _, c := range polygon {
			for j, p := range c {
				c[j].X = (p.X - min.X) * mul
				c[j].Y = (p.Y - min.Y) * mul
			}
		}
		r := safebbox(polygon)

		img2 := image.NewNRGBA(img.Bounds().Add(translate).Union(image.Rect(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))))
		draw.Draw(img2, img2.Bounds(), image.Black, image.Pt(0, 0), draw.Src)
		draw.Draw(img2, img.Bounds().Add(translate), img, image.Pt(0, 0), draw.Src)
		img = img2

		for _, c := range polygon {
			polyutil.DrawPolyline(c, brush(img, i))
		}
	}

	return img
}
