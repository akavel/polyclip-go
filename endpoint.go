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
)

// A container for endpoint data. A endpoint represents a location of interest (vertex between two polygon edges)
// as the sweep line passes through the polygons.
type endpoint struct {
	p           Point
	left        bool      // Is the point the left endpoint of the segment (p, other->p)?
	polygonType           // polygonType to which this event belongs to
	other       *endpoint // Event associated to the other endpoint of the segment

	// Does the segment (p, other->p) represent an inside-outside transition
	// in the polygon for a vertical ray from (p.x, -infinite) that crosses the segment?
	inout bool
	edgeType
	inside bool // Only used in "left" events. Is the segment (p, other->p) inside the other polygon?
}

func (e endpoint) String() string {
	sleft := map[bool]string{true: "left", false: "right"}
	return fmt.Sprint("{", e.p, " ", sleft[e.left], " type:", e.polygonType,
		" other:", e.other.p, " inout:", e.inout, " inside:", e.inside, " edgeType:", e.edgeType, "}")
}

func (e1 *endpoint) equals(e2 *endpoint) bool {
	return e1.p.Equals(e2.p) &&
		e1.left == e2.left &&
		e1.polygonType == e2.polygonType &&
		e1.other == e2.other &&
		e1.inout == e2.inout &&
		e1.edgeType == e2.edgeType &&
		e1.inside == e2.inside
}

func (se *endpoint) segment() segment {
	return segment{se.p, se.other.p}
}

func (e1 *endpoint) segmentsEqual(e2 *endpoint) bool {
	return e1.segment() == e2.segment()
}

func signedArea(p0, p1, p2 Point) float64 {
	return (p0.X-p2.X)*(p1.Y-p2.Y) -
		(p1.X-p2.X)*(p0.Y-p2.Y)
}

// Checks if this sweep event is below point p.
func (se *endpoint) below(x Point) bool {
	if se.left {
		return signedArea(se.p, se.other.p, x) > 0
	}
	return signedArea(se.other.p, se.p, x) > 0
}

func (se *endpoint) above(x Point) bool {
	return !se.below(x)
}

// leftRight() returns the left and right endpoints, in that order.
func (se *endpoint) leftRight() (Point, Point) {
	if se.left {
		return se.p, se.other.p
	}
	return se.other.p, se.p
}

// isValid() is true if the segment has the correct direction.
// Note that segments of zero length have no direction and are thus not considered valid.
func (se *endpoint) isValidDirection() bool {
	lp, rp := se.leftRight()
	return lp.isBefore(rp)
}

// Floating point imprecision in findIntersection() can create "non-reductive"
// divisions that result in infinite recursion.
//
// One class of non-reductive divisions can be detected at segment division time;
// divisions that result in segments going in the wrong direction
// (including zero-length segments, which have no direction) are non-reductive.
// This is detected by checking if endpoint.isValidDirection() in divideSegment().
//
// The other class of non-reductive divisions does create "valid" segments; one
// being infinitesimally small, and the other which recurses into a similar
// non-reductive division. This happens when the left or right endpoints of the
// two segments are very close but not equal, the problematic division for which manifests
// as an intersection point falling outside of the endpoints on a horizontal or vertical line.
//
// Note that theoretically, both classes of non-reductive division could be detected by
// comparing the length of the original segment against the length of the resulting segments,
// the latter of which should always be less than the former. Unfortunately, computation of
// segment length is subject to floating point imprecision and can introduce false positives.
// The rationale behind the current approach (isValidDirection and isValidSingleIntersection)
// is to rely only on boolean comparisons.
func isValidSingleIntersection(e1, e2 *endpoint, ip Point) bool {
	switch {
	case e1.p.X == ip.X && e2.p.X == ip.X: // e1.p, ip, e2.p on a vertical line
		return (ip.Y-e1.p.Y > 0) != (ip.Y-e2.p.Y > 0) // ip is above (or below) both e1.p and e2.p
	case e1.p.Y == ip.Y && e2.p.Y == ip.Y: // e1.p, ip, e2.p on a horizontal line
		return (ip.X-e1.p.X > 0) != (ip.X-e2.p.X > 0) // ip is to the left (or right) of both e1.p and e2.p
	case e1.other.p.X == ip.X && e2.other.p.X == ip.X: // e1.other.p, ip, e2.other.p on a vertical line
		return (ip.Y-e1.other.p.Y > 0) != (ip.Y-e2.other.p.Y > 0) // ip is above (or below) both e1.other.p and e2.other.p
	case e1.other.p.Y == ip.Y && e2.other.p.Y == ip.Y: // e1.other.p, ip, e2.other.p on a horizontal line
		return (ip.X-e1.other.p.X > 0) != (ip.X-e2.other.p.X > 0) // ip is to the left (or right) of both e1.other.p and e2.other.p
	}
	return true
}
