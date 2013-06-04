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
	"bytes"
	"github.com/akavel/polyclip-go"
	. "testing"
)

func verify(t *T, cond bool, format string, args ...interface{}) {
	if !cond {
		t.Errorf(format, args...)
	}
}

func TestPolygonDecodeEncode(t *T) {
	txt1 := "1\n3 1\n\t0 0\n\t1 1\n\t0.5 0.5\n"

	p, err := DecodePolygon(bytes.NewBufferString(txt1))
	verify(t, err == nil, "Expected no error decoding, got: %v", err)
	verify(t, len(*p) == 1, "Expected 1 contour")
	c := (*p)[0]
	verify(t, len(c) == 3, "Expected 3 points")
	verify(t, c[0].Equals(polyclip.Point{0, 0}), "Expected p0")
	verify(t, c[1].Equals(polyclip.Point{1, 1}), "Expected p1")
	verify(t, c[2].Equals(polyclip.Point{.5, .5}), "Expected p2")

	buf := &bytes.Buffer{}
	err = EncodePolygon(buf, *p)
	verify(t, err == nil, "Expected no error encoding, got: %v", err)
	verify(t, buf.String() == txt1, "Expected: %v, got: %v", txt1, buf)
}
