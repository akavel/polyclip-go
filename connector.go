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

// Holds intermediate results (pointChains) of the clipping operation and forms them into
// the final polygon.
type connector struct {
	openPolys   []chain
	closedPolys []chain
}

func (c *connector) add(s segment) {
	// j iterates through the openPolygon chains.
	for j := range c.openPolys {
		chain := &c.openPolys[j]
		if !chain.linkSegment(s) {
			continue
		}

		if chain.closed {
			if len(chain.points) == 2 {
				// We tried linking the same segment (but flipped end and start) to
				// a chain. (i.e. chain was <p0, p1>, we tried linking Segment(p1, p0)
				// so the chain was closed illegally.
				chain.closed = false
				return
			}
			// move the chain from openPolys to closedPolys
			c.closedPolys = append(c.closedPolys, c.openPolys[j])
			c.openPolys = append(c.openPolys[:j], c.openPolys[j+1:]...)
			return
		}

		// !chain.closed
		k := len(c.openPolys)
		for i := j + 1; i < k; i++ {
			// Try to connect this open link to the rest of the chains.
			// We won't be able to connect this to any of the chains preceding this one
			// because we know that linkSegment failed on those.
			if chain.linkChain(&c.openPolys[i]) {
				// delete
				c.openPolys = append(c.openPolys[:i], c.openPolys[i+1:]...)
				return
			}
		}
		return
	}

	// The segment cannot be connected with any open polygon
	c.openPolys = append(c.openPolys, *newChain(s))
}

func (c *connector) toPolygon() Polygon {
	poly := Polygon{}
	for _, chain := range c.closedPolys {
		con := Contour{}
		for _, p := range chain.points {
			con.Add(p)
		}
		poly.Add(con)
	}
	return poly
}
