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

// Package polyclip provides implementation of algorithm for Boolean operations on 2D polygons.
// For further details, consult the description of Polygon.Construct method.
package polyclip

import "math"

type Point struct {
	X, Y float64
}

// Equals returns true if both p1 and p2 describe exactly the same point.
func (p1 Point) Equals(p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

// Length returns distance from p to point (0, 0).
func (p Point) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

type Rectangle struct {
	Min, Max Point
}

func (r1 Rectangle) union(r2 Rectangle) Rectangle {
	return Rectangle{
		Min: Point{
			X: math.Min(r1.Min.X, r2.Min.X),
			Y: math.Min(r1.Min.Y, r2.Min.Y),
		},
		Max: Point{
			X: math.Max(r1.Max.X, r2.Max.X),
			Y: math.Max(r1.Max.Y, r2.Max.Y),
		}}
}

// Overlaps returns whether r1 and r2 have a non-empty intersection.
func (r1 Rectangle) Overlaps(r2 Rectangle) bool {
	return r1.Min.X <= r2.Max.X && r1.Max.X >= r2.Min.X &&
		r1.Min.Y <= r2.Max.Y && r1.Max.Y >= r2.Min.Y
}

// Used to represent an edge of a polygon.
type segment struct {
	start, end Point
}

// Contour represents a sequence of vertices connected by line segments, forming a closed shape.
type Contour []Point

// Add is a convenience method for appending a point to a contour.
func (c *Contour) Add(p Point) {
	*c = append(*c, p)
}

// BoundingBox finds minimum and maximum coordinates of points in a contour.
func (c Contour) BoundingBox() Rectangle {
	bb := Rectangle{}
	bb.Min.X = math.Inf(1)
	bb.Min.Y = math.Inf(1)
	bb.Max.X = math.Inf(-1)
	bb.Max.Y = math.Inf(-1)

	for _, p := range c {
		if p.X > bb.Max.X {
			bb.Max.X = p.X
		}
		if p.X < bb.Min.X {
			bb.Min.X = p.X
		}
		if p.Y > bb.Max.Y {
			bb.Max.Y = p.Y
		}
		if p.Y < bb.Min.Y {
			bb.Min.Y = p.Y
		}
	}
	return bb
}

func (c Contour) segment(index int) segment {
	if index == len(c)-1 {
		return segment{c[len(c)-1], c[0]}
	}
	return segment{c[index], c[index+1]}
	// if out-of-bounds, we expect panic detected by runtime
}

// Checks if a point is inside a contour using the "point in polygon" raycast method.
// This works for all polygons, whether they are clockwise or counter clockwise,
// convex or concave.
// See: http://en.wikipedia.org/wiki/Point_in_polygon#Ray_casting_algorithm
// Returns true if p is inside the polygon defined by contour.
func (c Contour) Contains(p Point) bool {
	// Cast ray from p.x towards the right
	intersections := 0
	for i := range c {
		curr := c[i]
		ii := i + 1
		if ii == len(c) {
			ii = 0
		}
		next := c[ii]

		// Is the point out of the edge's bounding box?
		// bottom vertex is inclusive (belongs to edge), top vertex is
		// exclusive (not part of edge) -- i.e. p lies "slightly above
		// the ray"
		bottom, top := curr, next
		if bottom.Y > top.Y {
			bottom, top = top, bottom
		}
		if p.Y < bottom.Y || p.Y >= top.Y {
			continue
		}
		// Edge is from curr to next.

		if p.X >= math.Max(curr.X, next.X) ||
			next.Y == curr.Y {
			continue
		}

		// Find where the line intersects...
		xint := (p.Y-curr.Y)*(next.X-curr.X)/(next.Y-curr.Y) + curr.X
		if curr.X != next.X && p.X > xint {
			continue
		}

		intersections++
	}

	return intersections%2 != 0
}

// Clone returns a copy of a contour.
func (c Contour) Clone() Contour {
	return append([]Point{}, c...)
}

// Polygon is carved out of a 2D plane by a set of (possibly disjoint) contours.
// It can thus contain holes, and can be self-intersecting.
type Polygon []Contour

// NumVertices returns total number of all vertices of all contours of a polygon.
func (p Polygon) NumVertices() int {
	num := 0
	for _, c := range p {
		num += len(c)
	}
	return num
}

// BoundingBox finds minimum and maximum coordinates of points in a polygon.
func (p Polygon) BoundingBox() Rectangle {
	bb := p[0].BoundingBox()
	for _, c := range p[1:] {
		bb = bb.union(c.BoundingBox())
	}

	return bb
}

// Add is a convenience method for appending a contour to a polygon.
func (p *Polygon) Add(c Contour) {
	*p = append(*p, c)
}

// Clone returns a duplicate of a polygon.
func (p Polygon) Clone() Polygon {
	r := Polygon(make([]Contour, len(p)))
	for i := range p {
		r[i] = p[i].Clone()
	}
	return r
}

// Op describes an operation which can be performed on two polygons.
type Op int

const (
	UNION Op = iota
	INTERSECTION
	DIFFERENCE
	XOR
)

// Construct computes a 2D polygon, which is a result of performing
// specified Boolean operation on the provided pair of polygons (p <Op> clipping).
// It uses algorithm described by F. Martínez, A. J. Rueda, F. R. Feito
// in "A new algorithm for computing Boolean operations on polygons"
// - see: http://wwwdi.ujaen.es/~fmartin/bool_op.html
// The paper describes the algorithm as performing in time O((n+k) log n),
// where n is number of all edges of all polygons in operation, and
// k is number of intersections of all polygon edges.
func (p Polygon) Construct(operation Op, clipping Polygon) Polygon {
	c := clipper{
		subject:  p,
		clipping: clipping,
	}
	return c.compute(operation)
}
