// Copyright (c) 2011 Mateusz Czapliński (Go port)
// Copyright (c) 2011 Mahir Iqbal (as3 version)
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

// based on http://code.google.com/p/as3polyclip/ (MIT licensed)
// and code by Martínez et al: http://wwwdi.ujaen.es/~fmartin/bool_op.html (public domain)

package polyutil

import (
	"bitbucket.org/akavel/polyclip.go"
	"fmt"
	"io"
)

// EncodeContour serializes all points of a specified contour using a simple textual format.
func EncodeContour(w io.Writer, c polyclip.Contour) error {
	_, err := fmt.Fprint(w, len(c), " 1\n")
	if err != nil {
		return err
	}
	for _, p := range c {
		_, err = fmt.Fprint(w, "\t", p.X, " ", p.Y, "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// EncodePolygon serializes all contours of a polygon using a simple textual format.
func EncodePolygon(w io.Writer, p polyclip.Polygon) error {
	_, err := fmt.Fprint(w, len(p), "\n")
	if err != nil {
		return err
	}
	for _, c := range p {
		err = EncodeContour(w, c)
		if err != nil {
			return err
		}
	}
	return nil
}

// DecodePolygon loads a polygon saved using EncodePolygon function.
func DecodePolygon(in io.Reader) (*polyclip.Polygon, error) {
	var ncontours int
	_, err := fmt.Fscan(in, &ncontours)
	if err != nil {
		return nil, err
	}
	polygon := polyclip.Polygon{}
	for i := 0; i < ncontours; i++ {
		var npoints, level int
		_, err = fmt.Fscan(in, &npoints, &level)
		if err != nil {
			return nil, err
		}
		c := polyclip.Contour{}
		for j := 0; j < npoints; j++ {
			p := polyclip.Point{}
			_, err = fmt.Fscan(in, &p.X, &p.Y)
			if err != nil {
				return nil, err
			}
			if j > 0 && p.Equals(c[len(c)-1]) {
				continue
			}
			if j == npoints-1 && p.Equals(c[0]) {
				continue
			}
			c.Add(p)
		}
		if len(c) < 3 {
			continue
		}
		polygon.Add(c)
	}
	return &polygon, nil
}
