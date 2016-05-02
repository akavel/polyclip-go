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

type eventQueue struct {
	elements []*endpoint
	sorted   bool
}

func (q *eventQueue) enqueue(e *endpoint) {
	if !q.sorted {
		q.elements = append(q.elements, e)
		return
	}

	// If already sorted, search for the correct location to insert e.
	i := sort.Search(len(q.elements), func(i int) bool {
		return endpointLess(e, q.elements[i])
	})

	// Insert e in the correct location.
	q.elements = append(q.elements, nil)
	copy(q.elements[i+1:], q.elements[i:])
	q.elements[i] = e
}

// The ordering is reversed because push and pop are faster.
// [MC: fragment from .c source below:]
// Return true means that [...] e1 is processed by the algorithm after e2
func endpointLess(e1, e2 *endpoint) bool {
	// Different x coordinate
	if e1.p.X != e2.p.X {
		return e1.p.X > e2.p.X
	}

	// Same x coordinate. The event with lower y coordinate is processed first
	if e1.p.Y != e2.p.Y {
		return e1.p.Y > e2.p.Y
	}

	// Same point, but one is a left endpoint and the other a right endpoint. The right endpoint is processed first
	if e1.left != e2.left {
		return e1.left
	}

	// Same point, both events are left endpoints or both are right endpoints. The event associate to the bottom segment is processed first
	return e1.above(e2.other.p)
}

func (q *eventQueue) dequeue() *endpoint {
	if !q.sorted {
		sort.Sort(queueComparer(q.elements))
		q.sorted = true
	}

	// pop
	x := q.elements[len(q.elements)-1]
	q.elements = q.elements[:len(q.elements)-1]
	return x
}

type queueComparer []*endpoint

func (q queueComparer) Len() int { return len(q) }
func (q queueComparer) Less(i, j int) bool {
	return endpointLess(q[i], q[j])
}
func (q queueComparer) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *eventQueue) IsEmpty() bool {
	return len(q.elements) == 0
}
