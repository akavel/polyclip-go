package polyclip_test

import (
	"fmt"
	"github.com/akavel/polyclip-go"
	"sort"
	. "testing"
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
	subject := polyclip.Polygon{{{1, 1}, {1, 2}, {2, 2}, {2, 1}}}
	clipping := polyclip.Polygon{
		{{2, 1}, {2, 2}, {3, 2}, {3, 1}},
		{{1, 2}, {1, 3}, {2, 3}, {2, 2}},
		{{2, 2}, {2, 3}, {3, 3}, {3, 2}}}
	result := dump(subject.Construct(polyclip.UNION, clipping))

	exp := dump(polyclip.Polygon{{
		{1, 1}, {2, 1}, {3, 1},
		{3, 2}, {3, 3},
		{2, 3}, {1, 3},
		{1, 2}}})
	if result != exp {
		t.Errorf("expected %s, got %s", exp, result)
	}
}
