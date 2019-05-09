package polyclip

import "testing"

func TestFindIntersection(t *testing.T) {
	cases := []struct {
		s1, s2           segment
		numIntersections int
		ip1, ip2         Point
	}{
		{
			// Almost (but not) parallel lines
			segment{Point{0, 0}, Point{100, 0.0001}},
			segment{Point{1, 0}, Point{100, 0}},
			0, Point{}, Point{},
		},
		{
			// Almost (but not) parallel lines
			segment{Point{0, 0}, Point{100, 0.0000001}},
			segment{Point{1, 0}, Point{100, 0}},
			0, Point{}, Point{},
		},
		{
			// Cross
			segment{Point{1, 0}, Point{1, 3}},
			segment{Point{0, 1}, Point{3, 1}},
			1, Point{1, 1}, Point{},
		},
		{
			// Rays
			segment{Point{0, 1}, Point{1, 3}},
			segment{Point{0, 1}, Point{3, 1}},
			1, Point{0, 1}, Point{},
		},
		{
			// Colinear rays
			segment{Point{2, 1}, Point{0, 1}},
			segment{Point{2, 1}, Point{1, 1}},
			2, Point{2, 1}, Point{}, // Why isn't this 2 intersections at {1,1} {2,1}?
		},
		{
			// Colinear rays
			segment{Point{0, 3}, Point{0, 1}},
			segment{Point{0, 3}, Point{0, 2}},
			2, Point{0, 3}, Point{}, // Why isn't this 2 intersections at {0,2}, {0,3}?
		},
		{
			// Overlapping segments
			segment{Point{0, 1}, Point{3, 1}},
			segment{Point{1, 1}, Point{2, 1}},
			1, Point{3, 1}, Point{}, // Why isn't this 2 intersections at {1,1} {2,1}?
		},
		{
			// Overlapping segments
			segment{Point{0, 1}, Point{0, 4}},
			segment{Point{0, 2}, Point{0, 3}},
			1, Point{0, 4}, Point{}, // Why isn't this 2 intersections at {0,2}, {0,3}?
		},
	}
	for i, v := range cases {
		num, ip1, _ := findIntersection(v.s1, v.s2)
		verify(t, num == v.numIntersections, "Case %d: Expected numIntersections to be %d, but got %d", i, v.numIntersections, num)
		verify(t, ip1.Equals(v.ip1), "Case %d: Expected ip1 to be %v, but got %v", i, v.ip1, ip1)
	}
}
