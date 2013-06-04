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

func connopen(openchains ...[]Point) connector {
	c := connector{openPolys: []chain{}}
	for _, pts := range openchains {
		c.openPolys = append(c.openPolys, chain{points: pts})
	}
	return c
}

func TestConnectorAddClosing1(t *T) {
	cases := []struct {
		c      connector
		add    segment
		length int
	}{
		{
			c: connopen(
				[]Point{{0.527105, 0.24687}, {0.2705720799269327, 0.2795780221218095}, {0.262624807729291, 0.30113844655235167}, {0.43093, 0.407828}, {0.48944187037949144, 0.6116041332606713}, {0.502984, 0.612599}},
				[]Point{{0.5813234786695596, 0.6602679842620749}, {0.569772, 0.46489}}),
			add:    segment{Point{0.5813234786695596, 0.6602679842620749}, Point{0.502984, 0.612599}},
			length: 8,
		},
		{ // simplified version of the above case
			c: connopen(
				[]Point{{0, 1}, {0, 2}, {0, 3}},
				[]Point{{1, 1}, {1, 2}}),
			add:    segment{Point{1, 1}, Point{0, 3}},
			length: 5,
		},
	}

	for i, x := range cases {
		x.c.add(x.add)
		verify(t, len(x.c.openPolys[0].points) == x.length, "Case %d, expected len(openPolys[0])==%d, got: %v", i, x.length, x.c)
	}

}
