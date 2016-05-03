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

package polyclip

import "sort"

// This is the data structure that simulates the sweepline as it parses through
// eventQueue, which holds the events sorted from left to right (x-coordinate).
// TODO: optimizations? use sort.Search()?
type sweepline []*endpoint

func (s *sweepline) remove(key *endpoint) {
	for i, el := range *s {
		if el.equals(key) {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}

func (s *sweepline) insert(item *endpoint) int {
	length := len(*s)
	if length == 0 {
		*s = append(*s, item)
		return 0
	}

	// Search for the correct location to insert item.
	i := sort.Search(len(*s), func(i int) bool {
		return segmentCompare(item, (*s)[i])
	})

	// Insert item in the correct location.
	*s = append(*s, nil)
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = item

	return i
}

// segmentCompare returns whether e1 is considered less than e2.
func segmentCompare(e1, e2 *endpoint) bool {
	switch {
	case e1 == e2:
		return false
	case signedArea(e1.p, e1.other.p, e2.p) != 0:
		fallthrough
	case signedArea(e1.p, e1.other.p, e2.other.p) != 0:
		// Segments are not collinear
		// If they share their left endpoint use the right endpoint to sort
		if e1.p.Equals(e2.p) {
			return e1.below(e2.other.p)
		}
		// Different points
		if endpointLess(e1, e2) { // has the line segment associated to e1 been inserted into S after the line segment associated to e2 ?
			return e2.above(e1.p)
		}
		// The line segment associated to e2 has been inserted into S after the line segment associated to e1
		return e1.below(e2.p)
	// Segments are collinear. Just a consistent criterion is used
	case e1.p.Equals(e2.p):
		//return e1 < e2
		return false
	}
	return endpointLess(e1, e2)
}
