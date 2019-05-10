package polyclip

import "testing"

func TestSnap(t *testing.T) {
	p := snap(Point{0, 0}, Point{0, 1e-9}, Point{1e-9, 0}, Point{1e-13, 1e-13})
	verify(t, p.Equals(Point{0, 0}), "Expected no snapping but snapped to %v", p)

	p = snap(Point{0, 0}, Point{0, 1e-9}, Point{1e-9, 0}, Point{1e-15, 1e-15})
	verify(t, p.Equals(Point{1e-15, 1e-15}), "Expected snapping to {1e-15, 1e-15}")
}
