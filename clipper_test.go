package polyclip

import "testing"

func TestSnap(t *testing.T) {
	cases := []struct {
		pt, l1, r1, l2, r2, snap Point
	}{
		{
			pt:   Point{0, 0},
			l1:   Point{0, 1e-9},
			r1:   Point{1e-9, 0},
			l2:   Point{1e-13, 1e-13},
			r2:   Point{1e-13, 0},
			snap: Point{0, 0},
		},
		{
			pt:   Point{0, 0},
			l1:   Point{0, 1e-9},
			r1:   Point{1e-9, 0},
			l2:   Point{1e-15, 1e-15},
			r2:   Point{1e-13, 0},
			snap: Point{1e-15, 1e-15},
		},
		{
			pt:   Point{0, 0},
			l1:   Point{1e-15, 1e-15},
			r1:   Point{1e-9, 0},
			l2:   Point{1e-15, 2e-15},
			r2:   Point{1, 0},
			snap: Point{1e-15, 2e-15}, // Should choose l2 instead of l1.
		},
		{
			pt:   Point{0, 0},
			l1:   Point{1e-15, 2e-15},
			r1:   Point{1e-9, 0},
			l2:   Point{1e-15, 1e-15},
			r2:   Point{1, 0},
			snap: Point{1e-15, 2e-15}, // Should choose l1 instead of l2.
		},
		{
			pt:   Point{0, 0},
			l1:   Point{1e-15, 1e-15},
			r1:   Point{1e-9, 0},
			l2:   Point{1e-15, 2e-15},
			r2:   Point{1, 0},
			snap: Point{1e-15, 2e-15}, // Should choose l2 instead of l1.
		},
		{
			pt:   Point{0, 0},
			l1:   Point{-1, 0},
			r1:   Point{1e-15, 3e-15},
			l2:   Point{-1, -1},
			r2:   Point{1e-15, 1e-15},
			snap: Point{1e-15, 1e-15}, // Should choose r1 instead of r2.
		},
		{
			pt:   Point{0, 0},
			l1:   Point{-1, 0},
			r1:   Point{1e-15, 1e-15},
			l2:   Point{-1, -1},
			r2:   Point{1e-15, 3e-15},
			snap: Point{1e-15, 1e-15}, // Should choose r2 instead of r1.
		},
	}
	for i, v := range cases {
		e1 := &endpoint{p: v.l1, left: true, other: &endpoint{p: v.r1, left: false}}
		e2 := &endpoint{p: v.l2, left: true, other: &endpoint{p: v.r2, left: false}}
		p := snap(v.pt, e1, e2)
		verify(t, p.Equals(v.snap), "Case %d: Expected snap to return %v but returned %v", i, v.snap, p)
	}
}
