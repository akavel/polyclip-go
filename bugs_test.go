package polyclip_test

import (
	"fmt"
	"github.com/akavel/polyclip-go"
	"sort"
	. "testing"
	"time"
)

type sorter polyclip.Polygon

func (s sorter) Len() int      { return len(s) }
func (s sorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sorter) Less(i, j int) bool {
	if len(s[i]) != len(s[j]) {
		return len(s[i]) < len(s[j])
	}
	for k := range s[i] {
		pi, pj := s[i][k], s[j][k]
		if pi.X != pj.X {
			return pi.X < pj.X
		}
		if pi.Y != pj.Y {
			return pi.Y < pj.Y
		}
	}
	return false
}

// basic normalization just for tests; to be improved if needed
func normalize(poly polyclip.Polygon) polyclip.Polygon {
	for i, c := range poly {
		if len(c) == 0 {
			continue
		}

		// find bottom-most of leftmost points, to have fixed anchor
		min := 0
		for j, p := range c {
			if p.X < c[min].X || p.X == c[min].X && p.Y < c[min].Y {
				min = j
			}
		}

		// rotate points to make sure min is first
		poly[i] = append(c[min:], c[:min]...)
	}

	sort.Sort(sorter(poly))
	return poly
}

func dump(poly polyclip.Polygon) string {
	return fmt.Sprintf("%v", normalize(poly))
}

func TestBug3(t *T) {
	cases := []struct{ subject, clipping, result polyclip.Polygon }{
		// original reported github issue #3
		{
			subject: polyclip.Polygon{{{1, 1}, {1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{
				{{2, 1}, {2, 2}, {3, 2}, {3, 1}},
				{{1, 2}, {1, 3}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 3}, {3, 2}}},
			result: polyclip.Polygon{{
				{1, 1}, {2, 1}, {3, 1},
				{3, 2}, {3, 3},
				{2, 3}, {1, 3},
				{1, 2}}},
		},
		// simplified variant of issue #3, for easier debugging
		{
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{
				{{2, 1}, {2, 2}, {3, 2}},
				{{1, 2}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 2}}},
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 1}}},
		},
		{
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{
				{{1, 2}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 2}}},
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 2}, {2, 1}}},
		},
		// another variation, now with single degenerated curve
		{
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{
				{{1, 2}, {2, 3}, {2, 2}, {2, 3}, {3, 2}}},
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 2}, {2, 1}}},
		},
		{
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{
				{{2, 1}, {2, 2}, {2, 3}, {3, 2}},
				{{1, 2}, {2, 3}, {2, 2}}},
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 1}}},
		},
		// "union" with effectively empty polygon (wholly self-intersecting)
		{
			subject:  polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 3}, {1, 2}, {2, 2}, {2, 3}}},
			result:   polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
		},
	}
	for _, c := range cases {
		result := dump(c.subject.Construct(polyclip.UNION, c.clipping))
		if result != dump(c.result) {
			t.Errorf("case UNION:\nsubject:  %v\nclipping: %v\nexpected: %v\ngot:      %v",
				c.subject, c.clipping, c.result, result)
		}
	}
}

func TestBug4(t *T) {
	if Short() {
		return
	}

	cases := []struct{ subject, clipping, result polyclip.Polygon }{
		// original reported github issue #4, resulting in infinte loop
		{
			subject: polyclip.Polygon{{
				{1.427255375e+06, -2.3283064365386963e-10},
				{1.4271285e+06, 134.7111358642578},
				{1.427109e+06, 178.30108642578125}}},
			clipping: polyclip.Polygon{{
				{1.416e+06, -12000},
				{1.428e+06, -12000},
				{1.428e+06, 0},
				{1.416e+06, 0},
				{1.416e+06, -12000}}},
			result: polyclip.Polygon{},
		},
	}
	for _, c := range cases {
		// check that we get a result in finite time

		ch := make(chan polyclip.Polygon)
		go func() {
			ch <- c.subject.Construct(polyclip.UNION, c.clipping)
		}()

		var result polyclip.Polygon
		select {
		case result = <-ch:
		case <-time.After(1 * time.Second):
			// panicking in attempt to get full stacktrace
			panic(fmt.Sprintf("case UNION:\nsubject:  %v\nclipping: %v\ntimed out.", c.subject, c.clipping))
		}
		s := dump(result)
		if s != dump(c.result) {
			t.Errorf("case UNION:\nsubject:  %v\nclipping: %v\nexpected: %v\ngot:      %v",
				c.subject, c.clipping, c.result, s)
		}
	}
}
