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

import (
	"fmt"
	"math"
)

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

// String returns a string representation of p like "(3,4)".
func (p Point) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", p.X, p.Y)
	//return "(" + strconv.FormatFloat(p.X, 'f', 2, 64) + "," + strconv.FormatFloat(p.Y, 'f', 2, 64) + ")"
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point) Mul(k float64) Point {
	return Point{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point) Div(k float64) Point {
	return Point{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p Point) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Angle returns the angle in radians from the origin to p and q
func (p Point) Angle(q Point) float64 {
	angle := math.Acos(p.Dot(q) / (p.Norm() * q.Norm()))
	if p.Cross(q) < 0 {
		angle *= -1
	}
	return angle
}

// Angle3 returns the angle in radians from p to q and r
func (p Point) Angle3(q, r Point) float64 {
	return q.Sub(p).Angle(r.Sub(p))
}

// Append appends a point to seperate x, y float slices
func (p Point) Append(x, y []float64) ([]float64, []float64) {
	return append(x, p.X), append(y, p.Y)
}

// Cross product of p and q
func (p Point) Cross(q Point) float64 {
	return p.X*q.Y - p.Y*q.X
}

// Dist returns distance to another point
func (p Point) Dist(q Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Dist2 returns squared distance to another point
func (p Point) Dist2(q Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return dx*dx + dy*dy
}

// Dot returns the dot product of this vector with another. There are multiple ways
// to describe this value. One is the multiplication of their lengths and cos(theta) where
// theta is the angle between the vectors: v1.v2 = |v1||v2|cos(theta).
//
// The other (and what is actually done) is the sum of the element-wise multiplication of all
// elements. So for instance, two Vec3s would yield v1.x * v2.x + v1.y * v2.y + v1.z * v2.z.
//
// This means that the dot product of a vector and itself is the square of its Len (within
// the bounds of floating points error).
//
// The dot product is roughly a measure of how closely two vectors are to pointing in the same
// direction. If both vectors are normalized, the value will be -1 for opposite pointing,
// one for same pointing, and 0 for perpendicular vectors.
func (p Point) Dot(q Point) float64 {
	return p.X*q.X + p.Y*q.Y
}

// MaxRadius returns the maximum distance to the four corners of a rectangle
func (p Point) MaxRadius(r Rectangle) float64 {
	d := math.Max(
		p.Dist2(r.Min), // upper left
		p.Dist2(r.Max)) // down right
	d = math.Max(d, p.Dist2(Point{float64(r.Min.X), float64(r.Max.Y)})) // down left
	d = math.Max(d, p.Dist2(Point{float64(r.Max.X), float64(r.Min.Y)})) // upper right
	return math.Sqrt(d)
}

// Negate returns a point with the opposite coordinates
func (p Point) Negate() Point {
	return Point{-p.X, -p.Y}
}

// Norm returns the norm/length of this vector
func (p Point) Norm() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Norm2 returns the squared norm/length of this vector
func (p Point) Norm2() float64 {
	return p.X*p.X + p.Y*p.Y
}

// Normalize returns the normalizes a vector to a certain length
func (p Point) Normalize(length ...float64) Point {
	l := 1.0
	if len(length) != 0 {
		l = length[0]
	}
	return p.Div(p.Norm() / l)
}

// Normals are the two vectors perpendicular to line pq
func (p Point) Normals(q Point) (Point, Point) {
	dx := q.X - p.X
	dy := q.Y - p.Y
	return Point{-dy, dx}, Point{dy, -dx}
}

// Polar converts polar to cartesion coordinates
func (p Point) Polar(r, a float64) Point {
	return p.Add(Point{r * math.Cos(a), r * math.Sin(a)})
}

// Prepend prepends a point to seperate x, y float slices
func (p Point) Prepend(x, y []float64) ([]float64, []float64) {
	return append([]float64{p.X}, x...), append([]float64{p.Y}, y...)
}

// Rect creates a rectangle  with the point as center and certain size
func (p Point) Rect(width, height float64) Rectangle {
	half := Point{width / 2, height / 2}
	return Rectangle{p.Sub(half), p.Add(half)}
}

// Tangent returns unit tangent from p to q
func (p Point) Tangent(q Point) Point {
	return q.Sub(p).Normalize()
}

// ZP is the zero Point.
var ZP Point

// Pt is shorthand for Point{X, Y}.
func Pt(X, Y float64) Point {
	return Point{X, Y}
}

// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
type Rectangle struct {
	Min, Max Point
}

// Union returns the smallest rectangle that contains both r and s.
func (r1 Rectangle) Union(r2 Rectangle) Rectangle {
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

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle) String() string {
	return fmt.Sprintf("%s-%s", r.Min, r.Max)
}

// Dx returns r's width.
func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle) Size() Point {
	return Point{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rectangle) Add(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X + p.X, r.Min.Y + p.Y},
		Point{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle) Sub(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X - p.X, r.Min.Y - p.Y},
		Point{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle) Inset(n float64) Rectangle {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Min.X > r.Max.X || r.Min.Y > r.Max.Y {
		return ZR
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Equals reports whether r and s are equal.
func (r Rectangle) Equals(s Rectangle) bool {
	return r.Min.X == s.Min.X && r.Min.Y == s.Min.Y &&
		r.Max.X == s.Max.X && r.Max.Y == s.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle) In(s Rectangle) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rectangle) Canon() Rectangle {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// AddPoint expands the rectangle so it contains the point
func (r Rectangle) AddPoint(p Point) Rectangle {
	return r.Union(Rectangle{p, p})
}

// Center returns the center of the rectangle
func (r Rectangle) Center() Point {
	return r.Min.Add(r.Max).Div(2)
}

// Diagonal returns the diagonal of a rectangle
func (r Rectangle) Diagonal() float64 {
	w := r.Dx()
	h := r.Dy()
	return math.Sqrt(w*w + h*h)
}

// Expand returns an rectangle that has been expanded by dx, dy
func (r Rectangle) Expand(dx, dy float64) Rectangle {
	return r.Center().Rect(r.Dx()+dx, r.Dy()+dy)
}

// MaxRadius returns the maximum distance to the four corners of a rectangle
func (r Rectangle) MaxRadius(p Point) float64 {
	d1 := math.Max(
		p.Dist2(r.Min), // upper left
		p.Dist2(r.Max)) // down right
	d2 := math.Max(
		p.Dist2(Point{r.Min.X, r.Max.Y}), // down left
		p.Dist2(Point{r.Max.X, r.Min.Y}), // upper right
	)
	return math.Sqrt(math.Max(d1, d2))
}

// Offset returns an rectangle that has been expanded by d on all sides
func (r Rectangle) Offset(d float64) Rectangle {
	return r.Expand(2*d, 2*d)
}

// Points return the points of an rectangle counter clock wise
func (r Rectangle) Points() []Point {
	min := r.Min
	max := r.Max
	return []Point{min, Point{min.X, max.Y}, max, Point{max.X, min.Y}}
	// reverse
	// return []Point{Point{max.X, min.Y}, max, Point{min.X, max.Y}, min}
}

// ZR is the zero Rectangle.
var ZR Rectangle

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}.
func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Point{x0, y0}, Point{x1, y1}}
}

// RectFromPoints constructs a rect that contains the given points.
func RectFromPoints(pts ...Point) Rectangle {
	switch len(pts) {
	case 0:
		return Rectangle{}
	case 1:
		return Rectangle{pts[0], pts[0]}
	}

	r := Rect(pts[0].X, pts[0].Y, pts[1].X, pts[1].Y)

	for _, p := range pts[2:] {
		r = r.AddPoint(p)
	}
	return r
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
		bb = bb.Union(c.BoundingBox())
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
