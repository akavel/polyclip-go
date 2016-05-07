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

import (
	"fmt"
	"math"
	"os"
)

//func _DBG(f func()) { f() }
func _DBG(f func()) {}

type polygonType int

const (
	_SUBJECT polygonType = iota
	_CLIPPING
)

type edgeType int

const (
	_EDGE_NORMAL edgeType = iota
	_EDGE_NON_CONTRIBUTING
	_EDGE_SAME_TRANSITION
	_EDGE_DIFFERENT_TRANSITION
)

// This class contains methods for computing clipping operations on polygons.
// It implements the algorithm for polygon intersection given by Francisco Martínez del Río.
// See http://wwwdi.ujaen.es/~fmartin/bool_op.html
type clipper struct {
	subject, clipping Polygon
	eventQueue
}

func (c *clipper) compute(operation Op) Polygon {

	// Test 1 for trivial result case
	if len(c.subject)*len(c.clipping) == 0 {
		switch operation {
		case DIFFERENCE:
			return c.subject.Clone()
		case UNION:
			if len(c.subject) == 0 {
				return c.clipping.Clone()
			}
			return c.subject.Clone()
		}
		return Polygon{}
	}

	// Test 2 for trivial result case
	subjectbb := c.subject.BoundingBox()
	clippingbb := c.clipping.BoundingBox()
	if !subjectbb.Overlaps(clippingbb) {
		switch operation {
		case DIFFERENCE:
			return c.subject.Clone()
		case UNION:
			result := c.subject.Clone()
			for _, cont := range c.clipping {
				result.Add(cont.Clone())
			}
			return result
		}
		return Polygon{}
	}

	numSegments := 0
	// Add each segment to the eventQueue, sorted from left to right.
	numSegments += addPolygonToQueue(&c.eventQueue, c.subject, _SUBJECT)
	numSegments += addPolygonToQueue(&c.eventQueue, c.clipping, _CLIPPING)

	connector := connector{} // to connect the edge solutions

	// This is the sweepline. That is, we go through all the polygon edges
	// by sweeping from left to right.
	S := sweepline{}

	MINMAX_X := math.Min(subjectbb.Max.X, clippingbb.Max.X)

	_DBG(func() {
		e := c.eventQueue.dequeue()
		c.eventQueue.enqueue(e)
		fmt.Print("\nInitial queue:\n")
		for i, e := range c.eventQueue.elements {
			fmt.Println(i, "=", *e)
		}
	})

	i := 0

	//  From the manuscript, the cycle is executed
	// n + 4k times, where n is the number of segments and k is the number of
	// intersections. I believe the maximum k would be about n^2.
	maxPossibleEvents := numSegments + 4*numSegments*numSegments

	for !c.eventQueue.IsEmpty() {

		if i > maxPossibleEvents {
			f, err := os.Create("polyclipError.log")
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(f, "subject: %#v,\nclipping: %#v\n,", c.subject, c.clipping)
			f.Close()
			panic("polyclip.compute: infinite loop. " +
				"Writing geometries to file polyclipError.log. " +
				"Please report this issue at github.com/ctessum/polyclip-go.")
		}
		i++

		var prev, next *endpoint
		e := c.eventQueue.dequeue()
		_DBG(func() { fmt.Printf("\nProcess event: (%d of %d)\n%v\n", i, len(c.eventQueue.elements)+1, *e) })

		// optimization 1
		switch {
		case operation == INTERSECTION && e.p.X > MINMAX_X:
			fallthrough
		case operation == DIFFERENCE && e.p.X > subjectbb.Max.X:
			return connector.toPolygon()
			//case operation == UNION && e.p.X > MINMAX_X:
			//	_DBG(func() { fmt.Print("\nUNION optimization, fast quit\n") })
			//	// add all the non-processed line segments to the result
			//	if !e.left {
			//		connector.add(e.segment())
			//	}
			//
			//	for !c.eventQueue.IsEmpty() {
			//		e = c.eventQueue.dequeue()
			//		if !e.left {
			//			connector.add(e.segment())
			//		}
			//	}
			//	return connector.toPolygon()
		}

		if e.left { // the line segment must be inserted into S
			pos := S.insert(e)
			//e.PosInS = pos

			prev = nil
			if pos > 0 {
				prev = S[pos-1]
			}
			next = nil
			if pos < len(S)-1 {
				next = S[pos+1]
			}

			// Compute the inside and inOut flags
			switch {
			case prev == nil: // there is not a previous line segment in S?
				e.inside, e.inout = false, false
			case prev.edgeType != _EDGE_NORMAL:
				if pos-2 < 0 { // e overlaps with prev
					// Not sure how to handle the case when pos - 2 < 0, but judging
					// from the C++ implementation this looks like how it should be handled.
					e.inside, e.inout = false, false
					if prev.polygonType != e.polygonType { // [MC: where does this come from?]
						e.inside = true
					} else {
						e.inout = true
					}
				} else { // the previous two line segments in S are overlapping line segments
					prevTwo := S[pos-2]
					if prev.polygonType == e.polygonType {
						e.inout = !prev.inout
						e.inside = !prevTwo.inout
					} else {
						e.inout = !prevTwo.inout
						e.inside = !prev.inout
					}
				}
			case e.polygonType == prev.polygonType: // previous line segment in S belongs to the same polygon that "e" belongs to
				e.inside = prev.inside
				e.inout = !prev.inout
			default: // previous line segment in S belongs to a different polygon that "e" belongs to
				e.inside = !prev.inout
				e.inout = prev.inside
			}

			_DBG(func() {
				fmt.Println("Status line after left insertion: ")
				for _, e := range S {
					fmt.Println(*e)
				}
			})

			// Process a possible intersection between "e" and its next neighbor in S
			if next != nil {
				c.possibleIntersection(e, next)
			}
			// Process a possible intersection between "e" and its previous neighbor in S
			if prev != nil {
				c.possibleIntersection(prev, e)
				//c.possibleIntersection(&e, prev)
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

			// Check if the line segment belongs to the Boolean operation
			switch e.edgeType {
			case _EDGE_NORMAL:
				switch operation {
				case INTERSECTION:
					if e.other.inside {
						connector.add(e.segment())
					}
				case UNION:
					if !e.other.inside {
						connector.add(e.segment())
					}
				case DIFFERENCE:
					if (e.polygonType == _SUBJECT && !e.other.inside) ||
						(e.polygonType == _CLIPPING && e.other.inside) {
						connector.add(e.segment())
					}
				case XOR:
					connector.add(e.segment())
				}
			case _EDGE_SAME_TRANSITION:
				if operation == INTERSECTION || operation == UNION {
					connector.add(e.segment())
				}
			case _EDGE_DIFFERENT_TRANSITION:
				if operation == DIFFERENCE {
					connector.add(e.segment())
				}
			}

			// delete line segment associated to e from S and check for intersection between the neighbors of "e" in S
			if otherPos != -1 {
				S.remove(S[otherPos])
			}

			if next != nil && prev != nil {
				c.possibleIntersection(next, prev)
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
	return connector.toPolygon()
}

var nanPoint Point

func init() {
	nanPoint = Point{X: math.NaN(), Y: math.NaN()}
}

func findIntersection(seg0, seg1 segment) (int, Point, Point) {
	pi0 := nanPoint
	pi1 := nanPoint
	p0 := seg0.start
	d0 := Point{seg0.end.X - p0.X, seg0.end.Y - p0.Y}
	p1 := seg1.start
	d1 := Point{seg1.end.X - p1.X, seg1.end.Y - p1.Y}
	sqrEpsilon := 0. // was 1e-3 earlier
	E := Point{p1.X - p0.X, p1.Y - p0.Y}
	kross := d0.X*d1.Y - d0.Y*d1.X
	sqrKross := kross * kross
	sqrLen0 := d0.Length()
	sqrLen1 := d1.Length()

	if sqrKross > sqrEpsilon*sqrLen0*sqrLen1 {
		// lines of the segments are not parallel
		s := (E.X*d1.Y - E.Y*d1.X) / kross
		if s < 0 || s > 1 {
			return 0, Point{}, Point{}
		}
		t := (E.X*d0.Y - E.Y*d0.X) / kross
		if t < 0 || t > 1 {
			return 0, nanPoint, nanPoint
		}
		// intersection of lines is a point an each segment [MC: ?]
		pi0.X = p0.X + s*d0.X
		pi0.Y = p0.Y + s*d0.Y

		// [MC: commented fragment removed]

		return 1, pi0, nanPoint
	}

	// lines of the segments are parallel
	sqrLenE := E.Length()
	kross = E.X*d0.Y - E.Y*d0.X
	sqrKross = kross * kross
	if sqrKross > sqrEpsilon*sqrLen0*sqrLenE {
		// lines of the segment are different
		return 0, nanPoint, nanPoint
	}

	// Lines of the segment are the same. Need to test for overlap of segments.
	// s0 = Dot (D0, E) * sqrLen0
	s0 := (d0.X*E.X + d0.Y*E.Y) / sqrLen0
	// s1 = s0 + Dot (D0, D1) * sqrLen0
	s1 := s0 + (d0.X*d1.X+d0.Y*d1.Y)/sqrLen0
	smin := math.Min(s0, s1)
	smax := math.Max(s0, s1)
	w := make([]float64, 0, 2)
	imax := findIntersection2(0.0, 1.0, smin, smax, &w)

	if imax > 0 {
		pi0.X = p0.X + w[0]*d0.X
		pi0.Y = p0.Y + w[0]*d0.Y

		// [MC: commented fragment removed]

		if imax > 1 {
			pi1.X = p0.X + w[1]*d0.X
			pi1.Y = p0.Y + w[1]*d0.Y
		}
	}

	return imax, pi0, pi1
}

func findIntersection2(u0, u1, v0, v1 float64, w *[]float64) int {
	if u1 < v0 || u0 > v1 {
		return 0
	}
	if u1 == v0 {
		*w = append(*w, u1)
		return 1
	}

	// u1 > v0

	if u0 == v1 {
		*w = append(*w, u0)
		return 1
	}

	// u0 < v1

	if u0 < v0 {
		*w = append(*w, v0)
	} else {
		*w = append(*w, u0)
	}
	if u1 > v1 {
		*w = append(*w, v1)
	} else {
		*w = append(*w, u1)
	}
	return 2
}

func (c *clipper) possibleIntersection(e1, e2 *endpoint) {
	// [MC]: commented fragment removed

	numIntersections, ip1, _ := findIntersection(e1.segment(), e2.segment())

	if numIntersections == 0 {
		return
	}

	if numIntersections == 1 && (e1.p.Equals(e2.p) || e1.other.p.Equals(e2.other.p)) {
		return // the line segments intersect at an endpoint of both line segments
	}

	//if numIntersections == 2 && e1.p.Equals(e2.p) {
	if numIntersections == 2 && e1.polygonType == e2.polygonType {
		return // the line segments overlap, but they belong to the same polygon
	}

	if numIntersections == 1 {
		if !e1.p.Equals(ip1) && !e1.other.p.Equals(ip1) {
			// if ip1 is not an endpoint of the line segment associated to e1 then divide "e1"
			c.divideSegment(e1, ip1)
		}
		if !e2.p.Equals(ip1) && !e2.other.p.Equals(ip1) {
			// if ip1 is not an endpoint of the line segment associated to e2 then divide "e2"
			c.divideSegment(e2, ip1)
		}
		return
	}

	// The line segments overlap and belong to different polygons
	sortedEvents := make([]*endpoint, 0, 4)
	switch {
	case e1.p.Equals(e2.p) || e1.p.Equals(e2.other.p):
		sortedEvents = append(sortedEvents, nil) // WTF [MC: WTF]
	case endpointLess(e1, e2):
		sortedEvents = append(sortedEvents, e2, e1)
	default:
		sortedEvents = append(sortedEvents, e1, e2)
	}

	switch {
	case e1.other.p.Equals(e2.other.p) || e1.other.p.Equals(e2.p):
		sortedEvents = append(sortedEvents, nil)
	case endpointLess(e1.other, e2.other):
		sortedEvents = append(sortedEvents, e2.other, e1.other)
	default:
		sortedEvents = append(sortedEvents, e1.other, e2.other)
	}

	if len(sortedEvents) == 2 { // are both line segments equal?
		e1.edgeType, e1.other.edgeType = _EDGE_NON_CONTRIBUTING, _EDGE_NON_CONTRIBUTING
		if e1.inout == e2.inout {
			e2.edgeType, e2.other.edgeType = _EDGE_SAME_TRANSITION, _EDGE_SAME_TRANSITION
		} else {
			e2.edgeType, e2.other.edgeType = _EDGE_DIFFERENT_TRANSITION, _EDGE_DIFFERENT_TRANSITION
		}
		return
	}

	if len(sortedEvents) == 3 { // the line segments share an endpoint
		sortedEvents[1].edgeType, sortedEvents[1].other.edgeType = _EDGE_NON_CONTRIBUTING, _EDGE_NON_CONTRIBUTING
		var idx int
		// is the right endpoint the shared point?
		if sortedEvents[0] != nil {
			idx = 0
		} else { // the shared point is the left endpoint
			idx = 2
		}
		if e1.inout == e2.inout {
			sortedEvents[idx].other.edgeType = _EDGE_SAME_TRANSITION
		} else {
			sortedEvents[idx].other.edgeType = _EDGE_DIFFERENT_TRANSITION
		}
		if sortedEvents[0] != nil {
			c.divideSegment(sortedEvents[0], sortedEvents[1].p)
		} else {
			c.divideSegment(sortedEvents[2].other, sortedEvents[1].p)
		}
		return
	}

	if sortedEvents[0] != sortedEvents[3].other {
		// no line segment includes totally the OtherEnd one
		sortedEvents[1].edgeType = _EDGE_NON_CONTRIBUTING
		if e1.inout == e2.inout {
			sortedEvents[2].edgeType = _EDGE_SAME_TRANSITION
		} else {
			sortedEvents[2].edgeType = _EDGE_DIFFERENT_TRANSITION
		}
		c.divideSegment(sortedEvents[0], sortedEvents[1].p)
		c.divideSegment(sortedEvents[1], sortedEvents[2].p)
		return
	}

	// one line segment includes the other one
	sortedEvents[1].edgeType, sortedEvents[1].other.edgeType = _EDGE_NON_CONTRIBUTING, _EDGE_NON_CONTRIBUTING
	c.divideSegment(sortedEvents[0], sortedEvents[1].p)
	if e1.inout == e2.inout {
		sortedEvents[3].other.edgeType = _EDGE_SAME_TRANSITION
	} else {
		sortedEvents[3].other.edgeType = _EDGE_DIFFERENT_TRANSITION
	}
	c.divideSegment(sortedEvents[3].other, sortedEvents[2].p)
}

func (c *clipper) divideSegment(e *endpoint, p Point) {
	// "Right event" of the "left line segment" resulting from dividing e (the line segment associated to e)
	r := &endpoint{p: p, left: false, polygonType: e.polygonType, other: e, edgeType: e.edgeType}
	// "Left event" of the "right line segment" resulting from dividing e (the line segment associated to e)
	l := &endpoint{p: p, left: true, polygonType: e.polygonType, other: e.other, edgeType: e.other.edgeType}

	if endpointLess(l, e.other) { // avoid a rounding error. The left event would be processed after the right event
		// println("Oops")
		e.other.left = true
		e.left = false
	}

	e.other.other = l
	e.other = r

	c.eventQueue.enqueue(l)
	c.eventQueue.enqueue(r)
}

type empty struct{}

// a polygonGraph holds the points of a polygon in a graph struct.
// The index of the first map is the starting point of each segment
// in the polygon and the index of the second map is the ending point
// of each segment.
type polygonGraph map[Point]map[Point]empty

// addToGraph adds the segments of the polygon to the graph in a
// way that ensures the same segment is not included twice in the
// polygon.
func addToGraph(g *polygonGraph, seg segment) {
	if seg.start.Equals(seg.end) {
		// The starting and ending points are the same, so this is
		// not in fact a segment.
		return
	}

	if _, ok := (*g)[seg.end][seg.start]; ok {
		// This polygonGraph already has a segment end -> start, adding
		// start -> end would make the polygon degenerate, so we delete both.
		delete((*g)[seg.end], seg.start)
		return
	}

	if _, ok := (*g)[seg.start]; !ok {
		(*g)[seg.start] = make(map[Point]empty)
	}

	// Add the segment.
	(*g)[seg.start][seg.end] = empty{}
}

// addPolygonToQueue	adds p to the event queue, returning the number of
// segments that were added.
func addPolygonToQueue(q *eventQueue, p Polygon, polyType polygonType) int {
	g := make(polygonGraph)
	for _, cont := range p {
		for i := range cont {
			addToGraph(&g, cont.segment(i))
		}
	}
	numSegments := 0
	for start, gg := range g {
		for end := range gg {
			addProcessedSegment(q, segment{start: start, end: end}, polyType)
			numSegments++
		}
	}
	return numSegments
}

func addProcessedSegment(q *eventQueue, segment segment, polyType polygonType) {

	e1 := &endpoint{p: segment.start, left: true, polygonType: polyType}
	e2 := &endpoint{p: segment.end, left: true, polygonType: polyType, other: e1}
	e1.other = e2

	switch {
	case e1.p.X < e2.p.X:
		e2.left = false
	case e1.p.X > e2.p.X:
		e1.left = false
	case e1.p.Y < e2.p.Y:
		// the line segment is vertical. The bottom endpoint is the left endpoint
		e2.left = false
	default:
		e1.left = false
	}

	// Pushing it so the queue is sorted from left to right, with object on the left having the highest priority
	q.enqueue(e1)
	q.enqueue(e2)
}
