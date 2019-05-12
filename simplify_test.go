// Copyright (c) 2019 Chris Tessum
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

package polyclip_test

import (
	"testing"

	polyclip "github.com/ctessum/polyclip-go"
)

type testCaseSimplify struct {
	name   string
	poly   polyclip.Polygon
	result polyclip.Polygon
}

type testCasesSimplify []testCaseSimplify

func (cases testCasesSimplify) verify(t *testing.T) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := dump(c.poly.Simplify())
			if result != dump(c.result) {
				t.Errorf("%s:\npolygon:  %v\nexpected: %v\ngot:      %v",
					c.name, c.poly, c.result, result)
			}
		})
	}
}

func TestSimplify(t *testing.T) {
	testCasesSimplify{
		{
			name: "Self-intersecting polygon",
			poly: polyclip.Polygon{{{0, 0}, {1, 1}, {1, 0}, {0, 1}}},
			result: polyclip.Polygon{
				{{0, 0}, {0.5, 0.5}, {0, 1}},
				{{0.5, 0.5}, {1, 1}, {1, 0}},
			},
		},
		{
			name: "Polygon with repeated edge",
			poly: polyclip.Polygon{{{0, 0}, {1, 0}, {1, 1}, {2, 1}, {2, 0}, {1, 0},
				{1, 1}, {0, 1}}},
			result: polyclip.Polygon{{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {1, 1}, {0, 1}}},
		},
		{
			name: "Polygon with partially repeated edge",
			poly: polyclip.Polygon{{{0, 0}, {1, 0}, {1, 0.75}, {2, 0.75}, {2, 0.25}, {1, 0.25},
				{1, 1}, {0, 1}}},
			result: polyclip.Polygon{{{0, 0}, {1, 0}, {1, 0.25}, {2, 0.25}, {2, 0.75},
				{1, 0.75}, {1, 1}, {0, 1}}},
		},
		{
			name: "Polygon with repeated edge in opposite direction",
			poly: polyclip.Polygon{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
				{{1, 0}, {2, 0}, {2, 1}, {1, 1}},
			},
			result: polyclip.Polygon{
				{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {1, 1}, {0, 1}},
			},
		},
		{
			name: "Polygon with partially repeated edge in opposite direction",
			poly: polyclip.Polygon{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
				{{1, 0.25}, {2, 0.25}, {2, 0.75}, {1, 0.75}},
			},
			result: polyclip.Polygon{
				{{0, 0}, {1, 0}, {1, 0.25}, {2, 0.25},
					{2, 0.75}, {1, 0.75}, {1, 1}, {0, 1}},
			},
		},
		{
			name:   "Completely degenerate",
			poly:   polyclip.Polygon{{{1, 2}, {2, 2}, {2, 3}, {1, 2}, {2, 2}, {2, 3}}},
			result: polyclip.Polygon{},
		},
	}.verify(t)
}
