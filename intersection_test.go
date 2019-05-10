package polyclip

import "testing"

func TestFindIntersection(t *testing.T) {
	cases := []struct {
		s1, s2           segment
		numIntersections int
		ip1              Point
	}{
		{
			// Almost (but not) parallel lines
			segment{Point{0, 0}, Point{100, 0.0001}},
			segment{Point{1, 0}, Point{100, 0}},
			0, Point{},
		},
		{
			// Almost (but not) parallel lines
			segment{Point{0, 0}, Point{100, 0.0000001}},
			segment{Point{1, 0}, Point{100, 0}},
			0, Point{},
		},
		{
			// Cross
			segment{Point{1, 0}, Point{1, 3}},
			segment{Point{0, 1}, Point{3, 1}},
			1, Point{1, 1},
		},
		{
			// Rays
			segment{Point{0, 1}, Point{1, 3}},
			segment{Point{0, 1}, Point{3, 1}},
			1, Point{0, 1},
		},
		{
			// Colinear rays
			segment{Point{2, 1}, Point{0, 1}},
			segment{Point{2, 1}, Point{1, 1}},
			2, Point{2, 1},
		},
		{
			// Colinear rays
			segment{Point{0, 3}, Point{0, 1}},
			segment{Point{0, 3}, Point{0, 2}},
			2, Point{0, 3},
		},
		{
			// Overlapping segments
			segment{Point{0, 1}, Point{3, 1}},
			segment{Point{1, 1}, Point{2, 1}},
			2, Point{1, 1},
		},
		{
			// Overlapping segments
			segment{Point{0, 1}, Point{0, 4}},
			segment{Point{0, 2}, Point{0, 3}},
			2, Point{0, 2},
		},
		{ // Overlapping segments
			segment{Point{43.2635182233307, 170.15192246987792}, Point{41.57979856674331, 170.60307379214092}},
			segment{Point{43.2635182233307, 170.15192246987792}, Point{42.78116786015871, 170.28116786015872}},
			2, Point{43.2635182233307, 170.15192246987792},
		},
		{ // Overlapping segments
			segment{Point{41.57979856674331, 170.60307379214092}, Point{43.2635182233307, 170.15192246987792}},
			segment{Point{42.78116786015871, 170.28116786015872}, Point{43.2635182233307, 170.15192246987792}},
			2, Point{42.78116786015871, 170.28116786015872},
		},
		{ // Overlapping segments
			segment{Point{43.2635182233307, 170.15192246987792}, Point{41.57979856674331, 170.60307379214092}},
			segment{Point{42.78116786015871, 170.28116786015872}, Point{43.2635182233307, 170.15192246987792}},
			2, Point{43.2635182233307, 170.15192246987792},
		},
		{ // Overlapping segments
			segment{Point{41.57979856674331, 170.60307379214092}, Point{43.2635182233307, 170.15192246987792}},
			segment{Point{43.2635182233307, 170.15192246987792}, Point{42.78116786015871, 170.28116786015872}},
			2, Point{43.2635182233307, 170.15192246987792},
		},
		{ // Identical segments
			segment{Point{66, 160}, Point{67.1242262770966, 147.15003485264717}},
			segment{Point{66, 160}, Point{67.1242262770966, 147.15003485264717}},
			2, Point{66, 160},
		},
	}
	for i, v := range cases {
		num, ip1, _ := findIntersection(v.s1, v.s2, true)
		verify(t, num == v.numIntersections, "Case %d: Expected numIntersections to be %d, but got %d", i, v.numIntersections, num)
		verify(t, ip1.Equals(v.ip1), "Case %d: Expected ip1 to be %v, but got %v", i, v.ip1, ip1)
	}
}
