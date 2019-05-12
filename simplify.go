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

package polyclip

import "fmt"

// Simplify removes self-intersections and degenerate (repeated)
// edges from polygons.
func (p Polygon) Simplify() Polygon {
	c := new(clipper)
	var edges int
	for _, cont := range p {
		for i := range cont {
			addProcessedSegment(&c.eventQueue, cont.segment(i), _SUBJECT)
			edges++
		}
	}

	connector := connector{operation: UNION} // to connect the edge solutions

	// This is the sweepline. That is, we go through all the polygon edges
	// by sweeping from left to right.
	S := sweepline{}

	endpoints := make([]*endpoint, 0, edges)

	for !c.eventQueue.IsEmpty() {
		var prev, next *endpoint
		e := c.eventQueue.dequeue()
		_DBG(func() { fmt.Printf("\nProcess event: (of %d)\n%v\n", len(c.eventQueue.elements)+1, *e) })

		if e.left { // the line segment must be inserted into S
			pos := S.insert(e)

			prev = nil
			if pos > 0 {
				prev = S[pos-1]
			}
			next = nil
			if pos < len(S)-1 {
				next = S[pos+1]
			}

			_DBG(func() {
				fmt.Println("Status line after insertion: ")
				for _, e := range S {
					fmt.Println(*e)
				}
			})

			// Process a possible intersection between "e" and its next neighbor in S
			if next != nil {
				c.processIntersectionSimplify(e, next)
			}
			// Process a possible intersection between "e" and its previous neighbor in S
			if prev != nil {
				divided := c.processIntersectionSimplify(prev, e)
				// If [prev] was divided, the context (sweep line S) for [e] may have changed,
				// altering what e.inout and e.inside should be. [e] must thus be reenqueued to
				// recompute e.inout and e.inside.
				//
				// (This should not be done if [e] was also divided; in that case
				//  the divided segments are already enqueued).
				if len(divided) == 1 && divided[0] == prev {
					S.remove(e)
					c.eventQueue.enqueue(e)
				}
			}
		} else { // the line segment must be removed from S
			otherPos := -1
			for i := range S {
				if S[i].equals(e.other) {
					otherPos = i
					break
				}
			}
			// otherPos := S.IndexOf(e.other)
			// [or:] otherPos := e.other.PosInS

			if otherPos != -1 {
				prev = nil
				if otherPos > 0 {
					prev = S[otherPos-1]
				}
				next = nil
				if otherPos < len(S)-1 {
					next = S[otherPos+1]
				}
			}

			endpoints = append(endpoints, e)

			// delete line segment associated to e from S and check for intersection between the neighbors of "e" in S
			if otherPos != -1 {
				S.remove(S[otherPos])
			}

			if next != nil && prev != nil {
				c.processIntersectionSimplify(next, prev)
			}

			_DBG(func() { fmt.Print("Connector:\n", connector, "\n") })
		}
		_DBG(func() {
			fmt.Println("Status line after processing intersections: ")
			for _, e := range S {
				fmt.Println(*e)
			}
		})
	}

	for i, e := range endpoints {
		if i == 0 || i == len(endpoints)-1 {
			connector.add(e.segment())
		} else if !(e.p.Equals(endpoints[i+1].p) && e.other.p.Equals(endpoints[i+1].other.p)) &&
			!(e.p.Equals(endpoints[i-1].p) && e.other.p.Equals(endpoints[i-1].other.p)) {
			connector.add(e.segment())
		}
	}
	return connector.toPolygon()
}

func (c *clipper) processIntersectionSimplify(e1, e2 *endpoint) []*endpoint {
	numIntersections, ip1, ip2 := findIntersection(e1.segment(), e2.segment(), true)

	if numIntersections == 0 {
		return nil
	}

	// Adjust for floating point imprecision when intersections are created at endpoints, which
	// otherwise has the tendency to corrupt the original polygons with new, almost-parallel segments.
	ip1 = snap(ip1, e1.p, e2.p, e1.other.p, e2.other.p)

	if numIntersections == 1 {
		ep := make([]*endpoint, 0, 2)
		if !ip1.Equals(e1.p) && !ip1.Equals(e1.other.p) {
			// e2 divides e1.
			ep = append(ep, c.divideSegment(e1, ip1))
		}
		if !ip1.Equals(e2.p) && !ip1.Equals(e2.other.p) {
			// e1 divides e2/
			ep = append(ep, c.divideSegment(e2, ip1))
		}
		return ep
	}

	// The line segements overlap.
	ip2 = snap(ip2, e1.p, e2.p, e1.other.p, e2.other.p)
	ep := make([]*endpoint, 0, 2)
	if !ip1.Equals(e1.p) && !ip2.Equals(e1.other.p) {
		ep = append(ep, c.divideSegment(e1, ip1))
	}
	if !ip1.Equals(e2.p) && !ip2.Equals(e2.other.p) {
		ep = append(ep, c.divideSegment(e2, ip1))
	}
	return ep
}
