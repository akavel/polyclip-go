package polyclip_test

import (
	"fmt"
	"math"
	"sort"
	"testing"
	"time"

	polyclip "github.com/ctessum/polyclip-go"
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
		path := append(c[min:], c[:min]...)

		// give all paths a consistent direction.
		if len(path) > 2 && signedArea(path[0], path[1], path[len(path)-1]) < 0 {
			for l, r := 1, len(path)-1; l < r; l, r = l+1, r-1 {
				path[l], path[r] = path[r], path[l]
			}
		}
		poly[i] = path
	}

	sort.Sort(sorter(poly))
	return poly
}

func signedArea(p0, p1, p2 polyclip.Point) float64 {
	return (p0.X-p2.X)*(p1.Y-p2.Y) - (p1.X-p2.X)*(p0.Y-p2.Y)
}

func dump(poly polyclip.Polygon) string {
	return fmt.Sprintf("%v", normalize(poly))
}

type testCase struct {
	op       polyclip.Op
	subject  polyclip.Polygon
	clipping polyclip.Polygon
	result   polyclip.Polygon
}

type testCases []testCase

func (cases testCases) verify(t *testing.T) {
	t.Helper()
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			expected := c.result
			for i, path := range expected {
				// Remove duplicate endpoints to eliminate zero-length segments that the clipper won't produce.
				for len(path) > 1 && path[0] == path[len(path)-1] {
					path = path[0 : len(path)-1]
				}
				expected[i] = path
			}

			result := dump(c.subject.Construct(c.op, c.clipping))
			if result != dump(expected) {
				t.Errorf("case %d: %v\nsubject:  %v\nclipping: %v\nexpected: %v\ngot:      %v",
					i, c.op, c.subject, c.clipping, c.result, result)
			}
		})
	}
}

func TestBug3(t *testing.T) {
	testCases{
		// original reported github issue #3
		{
			op:      polyclip.UNION,
			subject: polyclip.Polygon{{{1, 1}, {1, 2}, {2, 2}, {2, 1}}}.Simplify(),
			clipping: polyclip.Polygon{
				{{2, 1}, {2, 2}, {3, 2}, {3, 1}},
				{{1, 2}, {1, 3}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 3}, {3, 2}}}.Simplify(),
			result: polyclip.Polygon{{
				{1, 1}, {2, 1}, {3, 1},
				{3, 2}, {3, 3},
				{2, 3}, {1, 3},
				{1, 2}}},
		},
		// simplified variant of issue #3, for easier debugging
		{
			op:      polyclip.UNION,
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{
				{{2, 1}, {2, 2}, {3, 2}},
				{{1, 2}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 2}}}.Simplify(),
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 1}}}.Simplify(),
		},
		{
			op:      polyclip.UNION,
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}}.Simplify(),
			clipping: polyclip.Polygon{
				{{1, 2}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 2}}}.Simplify(),
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 2}, {2, 1}}},
		},
		// another variation, now with single degenerated curve
		{
			op:      polyclip.UNION,
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}}.Simplify(),
			clipping: polyclip.Polygon{
				{{1, 2}, {2, 3}, {2, 2}, {2, 3}, {3, 2}}}.Simplify(),
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 2}, {2, 1}}},
		},
		{
			op:      polyclip.UNION,
			subject: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}}.Simplify(),
			clipping: polyclip.Polygon{
				{{2, 1}, {2, 2}, {2, 3}, {3, 2}},
				{{1, 2}, {2, 3}, {2, 2}}}.Simplify(),
			result: polyclip.Polygon{{{1, 2}, {2, 3}, {3, 2}, {2, 1}}},
		},
		// "union" with effectively empty polygon (wholly self-intersecting)
		{
			op:       polyclip.UNION,
			subject:  polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
			clipping: polyclip.Polygon{{{1, 2}, {2, 2}, {2, 3}, {1, 2}, {2, 2}, {2, 3}}}.Simplify(),
			result:   polyclip.Polygon{{{1, 2}, {2, 2}, {2, 1}}},
		},
	}.verify(t)
}

func TestResweepingIntersectingEndpoints(t *testing.T) {
	testCases{
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{160.09449516387843, 200.37992407997774},
				{90.09449516387845, 200.37992407997774},
				{55.094495163878435, 139.75814581506702},
				{90.0944951638784, 79.13636755015634}}},
			clipping: polyclip.Polygon{{
				{82.84661138052363, 131.51881422166852},
				{66.59206311550543, 159.6725176707606},
				{90.09449516387845, 200.37992407997774},
				{160.09449516387843, 200.37992407997774},
			}},
			result: polyclip.Polygon{{
				{82.84661138052363, 131.51881422166852},
				{66.59206311550543, 159.6725176707606},
				{90.09449516387845, 200.37992407997774},
				{160.09449516387843, 200.37992407997774},
			}},
		},
		{
			op: polyclip.DIFFERENCE,
			subject: polyclip.Polygon{{
				{160.09449516387843, 200.37992407997774},
				{90.09449516387845, 200.37992407997774},
				{55.094495163878435, 139.75814581506702},
				{90.0944951638784, 79.13636755015634}}},
			clipping: polyclip.Polygon{{
				{82.84661138052363, 131.51881422166852},
				{66.59206311550543, 159.6725176707606},
				{90.09449516387845, 200.37992407997774},
				{160.09449516387843, 200.37992407997774},
			}},
			result: polyclip.Polygon{{
				{160.09449516387843, 200.37992407997774},
				{82.84661138052363, 131.51881422166852},
				{66.59206311550543, 159.6725176707606},
				{55.094495163878435, 139.75814581506702},
				{90.0944951638784, 79.13636755015634},
			}},
		},
		{
			op: polyclip.UNION,
			subject: polyclip.Polygon{{
				{160.09449516387843, 200.37992407997774},
				{90.09449516387845, 200.37992407997774},
				{55.094495163878435, 139.75814581506702},
				{90.0944951638784, 79.13636755015634}}},
			clipping: polyclip.Polygon{{
				{82.84661138052363, 131.51881422166852},
				{66.59206311550543, 159.6725176707606},
				{90.09449516387845, 200.37992407997774},
				{160.09449516387843, 200.37992407997774},
			}},
			result: polyclip.Polygon{{
				{160.09449516387843, 200.37992407997774},
				{90.09449516387845, 200.37992407997774},
				{66.59206311550543, 159.6725176707606},
				{55.094495163878435, 139.75814581506702},
				{90.0944951638784, 79.13636755015634}}},
		},
		{
			op: polyclip.UNION,
			subject: polyclip.Polygon{{
				{70.78432620601497, -7.668842337087888},
				{42.500054958553065, -19.38457108962598},
				{22.504998288170377, -11.102347436334847},
				{14.215783711091163, -7.668842337087877},
				{2.500054958553072, 20.615428910374025},
				{4.163269713667806, 24.63078452931106},
				{-16.386530327112805, 33.142790410257575},
				{-28.102259079650896, 61.42706165771948},
			}},
			clipping: polyclip.Polygon{{
				{22.504998288170377, -11.102347436334847},
				{14.215783711091163, -7.668842337087877},
				{2.500054958553072, 20.615428910374025},
				{4.163269713667806, 24.63078452931106},
				{-16.386530327112805, 33.142790410257575},
				{-18.453791204657392, 38.133599657789034},
				{-23.270336557414375, 26.505430543378026},
				{16.72966344258562, -13.494569456621978},
				{45.01393469004752, -1.778840704083887},
				{22.504998288170377, -11.102347436334847},
			}},
			result: polyclip.Polygon{{
				{-28.102259079650896, 61.42706165771948},
				{-18.453791204657392, 38.133599657789034},
				{-23.270336557414375, 26.505430543378026},
				{16.72966344258562, -13.494569456621978},
				{22.504998288170377, -11.102347436334847},
				{42.500054958553065, -19.38457108962598},
				{70.78432620601497, -7.668842337087888},
			}},
		},
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{131.59597133486625, 287.9385241571817},
				{100.00000000000004, 273.20508075688775},
				{71.44247806269215, 253.20888862379562},
				{71.44247806269209, -53.20888862379559},
				{99.99999999999991, -73.20508075688767},
				{131.59597133486614, -87.93852415718163},
			}},
			clipping: polyclip.Polygon{{
				{128.55752193730785, -53.208888623795616},
				{100, -73.20508075688772},
				{99.99999999999991, -73.20508075688767},
				{71.44247806269209, -53.20888862379559},
				{71.44247806269215, 253.20888862379562},
				{100.00000000000003, 273.2050807568877},
				{128.55752193730788, 253.20888862379562},
			}},
			result: polyclip.Polygon{{
				{128.55752193730785, -53.208888623795616},
				{100, -73.20508075688772},
				{99.99999999999991, -73.20508075688767},
				{71.44247806269209, -53.20888862379559},
				{71.44247806269215, 253.20888862379562},
				{100.00000000000003, 273.2050807568877},
				{128.55752193730788, 253.20888862379562},
			}},
		},
	}.verify(t)
}

func TestIntersectionFalsePositives(t *testing.T) {
	testCases{
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{100.0000001, 100},
				{98.48077539970159, 117.36481778405785},
				{98.48077539970159, 82.63518221594214},
			}},
			clipping: polyclip.Polygon{{
				{100.0000001, 100},
				{100, 99.99999885699484},
				{99.9999999, 100.00000000000001},
				{100, 100.00000114300516},
				{100.0000001, 100},
			}},
			result: polyclip.Polygon{{
				{100.0000001, 100},
				{100, 99.99999885699484},
				{99.9999999, 100.00000000000001},
				{100, 100.00000114300516},
				{100.0000001, 100},
			}},
		},
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{100.00000001, 100},
				{98.48077531106888, 117.36481776842952},
				{98.48077531106888, 82.63518223157048},
			}},
			clipping: polyclip.Polygon{{
				{100.00000001, 100},
				{100, 99.99999988569955},
				{99.99999999, 100.00000000000001},
				{100, 100.00000011430046},
				{100.00000001, 100},
			}},
			result: polyclip.Polygon{{
				{100.00000001, 100},
				{100, 99.99999988569955},
				{99.99999999, 100.00000000000001},
				{100, 100.00000011430046},
				{100.00000001, 100},
			}},
		},
	}.verify(t)
}

func TestCorruptionResistanceFromFloatingPointImprecision(t *testing.T) {
	testCases{
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{2.500054958553072, 20.615428910374025},
				{42.500054958553065, -19.38457108962598},
				{82.50005495855308, 20.61542891037402},
				{42.50005495855307, 60.61542891037402},
			}},
			clipping: polyclip.Polygon{{
				{7.604714313123809, 25.720088264944764},
				{11.897740920349097, 21.42706165771947},
				{36.852886624296644, 46.382207361667014},
				{32.55986001707135, 50.6752339688923},
			}},
			result: polyclip.Polygon{{
				{7.604714313123809, 25.720088264944764},
				{32.55986001707135, 50.6752339688923},
				{36.852886624296644, 46.382207361667014},
				{11.897740920349097, 21.42706165771947},
			}},
		},
		{
			op: polyclip.DIFFERENCE,
			subject: polyclip.Polygon{{
				{2.500054958553072, 20.615428910374025},
				{42.500054958553065, -19.38457108962598},
				{82.50005495855308, 20.61542891037402},
				{42.50005495855307, 60.61542891037402},
			}},
			clipping: polyclip.Polygon{{
				{7.604714313123809, 25.720088264944764},
				{32.55986001707135, 50.6752339688923},
				{36.852886624296644, 46.382207361667014},
				{11.897740920349097, 21.42706165771947},
			}},
			result: polyclip.Polygon{{
				{2.500054958553072, 20.615428910374025},
				{42.500054958553065, -19.38457108962598},
				{82.50005495855308, 20.61542891037402},
				{42.50005495855307, 60.61542891037402},
				{32.55986001707135, 50.6752339688923},
				{36.852886624296644, 46.382207361667014},
				{11.897740920349097, 21.42706165771947},
				{7.604714313123809, 25.720088264944764},
			}},
		},
		{
			op: polyclip.UNION,
			subject: polyclip.Polygon{{
				{2.500054958553072, 20.615428910374025},
				{42.500054958553065, -19.38457108962598},
				{82.50005495855308, 20.61542891037402},
				{42.50005495855307, 60.61542891037402},
			}},
			clipping: polyclip.Polygon{{
				{7.604714313123809, 25.720088264944764},
				{32.55986001707135, 50.6752339688923},
				{36.852886624296644, 46.382207361667014},
				{11.897740920349097, 21.42706165771947},
			}},
			result: polyclip.Polygon{{
				{2.500054958553072, 20.615428910374025},
				{42.500054958553065, -19.38457108962598},
				{82.50005495855308, 20.61542891037402},
				{42.50005495855307, 60.61542891037402},
				{32.55986001707135, 50.6752339688923},
				{7.604714313123809, 25.720088264944764},
			}},
		},
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{300, -100},
				{259.8076211353316, 49.99999999999997},
				{-259.80762113533154, 50.000000000000114},
				{-300, -99.99999999999996},
			}},
			clipping: polyclip.Polygon{{
				{273.2050807568877, 7.815970093361102e-14},
				{259.8076211353315, -50.00000000000014},
				{-259.80762113533166, -49.99999999999994},
				{-273.20508075688775, -7.105427357601002e-14},
				{-259.80762113533154, 50.000000000000114},
				{259.8076211353316, 49.99999999999997},
			}},
			result: polyclip.Polygon{{
				{273.2050807568877, 7.815970093361102e-14},
				{259.8076211353315, -50.00000000000014},
				{-259.80762113533166, -49.99999999999994},
				{-273.20508075688775, -7.105427357601002e-14},
				{-259.80762113533154, 50.000000000000114},
				{259.8076211353316, 49.99999999999997},
			}},
		},
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{300, -100},
				{277.163859753386, 14.805029709526934},
				{-277.163859753386, 14.805029709526949},
				{-300, -99.99999999999996},
			}},
			clipping: polyclip.Polygon{{
				{280.1087632620342, 4.618527782440651e-14},
				{277.16385975338596, -14.805029709527119},
				{-277.1638597533861, -14.805029709526906},
				{-280.10876326203424, -1.1368683772161603e-13},
				{-277.163859753386, 14.805029709526949},
				{277.163859753386, 14.805029709526934},
				{280.1087632620342, 4.618527782440651e-14},
			}},
			result: polyclip.Polygon{{
				{280.1087632620342, 4.618527782440651e-14},
				{277.16385975338596, -14.805029709527119},
				{-277.1638597533861, -14.805029709526906},
				{-280.10876326203424, -1.1368683772161603e-13},
				{-277.163859753386, 14.805029709526949},
				{277.163859753386, 14.805029709526934},
				{280.1087632620342, 4.618527782440651e-14},
			}},
		},
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{-196.96155060244163, 65.270364466614},
				{-187.9385241571817, 31.595971334866263},
				{-173.20508075688775, 4.263256414560601e-14},
				{-153.20888862379562, -28.557521937307854},
				{153.2088886237956, -28.55752193730791},
				{173.20508075688767, -8.526512829121202e-14},
				{187.93852415718163, 31.59597133486612},
				{196.96155060244163, 65.27036446661393},
			}},
			clipping: polyclip.Polygon{{
				{196.96155060244163, -65.27036446661393},
				{187.9385241571817, -31.595971334866263},
				{173.20508075688775, -1.4210854715202004e-14},
				{153.20888862379562, 28.557521937307854},
				{-153.2088886237956, 28.557521937307882},
				{-173.20508075688775, -1.4210854715202004e-14},
				{-187.93852415718166, -31.59597133486622},
				{-196.9615506024416, -65.27036446661387},
			}},
			result: polyclip.Polygon{{
				{-173.20508075688775, 4.263256414560601e-14},
				{-153.20888862379562, -28.557521937307854},
				{153.2088886237956, -28.55752193730791},
				{173.20508075688767, -8.526512829121202e-14},
				{173.20508075688775, -1.4210854715202004e-14},
				{153.20888862379562, 28.557521937307854},
				{-153.2088886237956, 28.557521937307882},
			}},
		},
		{
			op: polyclip.INTERSECTION,
			subject: polyclip.Polygon{{
				{128.55752193730788, 253.20888862379562},
				{100.00000000000003, 273.2050807568877},
				{68.40402866513377, 287.9385241571817},
				{68.40402866513364, -87.93852415718172},
				{100, -73.20508075688772},
				{128.55752193730785, -53.208888623795616},
			}},
			clipping: polyclip.Polygon{{
				{131.59597133486625, 287.9385241571817},
				{100.00000000000004, 273.20508075688775},
				{71.44247806269215, 253.20888862379562},
				{71.44247806269209, -53.20888862379559},
				{99.99999999999991, -73.20508075688767},
				{131.59597133486614, -87.93852415718163},
			}},
			result: polyclip.Polygon{{
				{71.44247806269209, -53.20888862379559},
				{99.99999999999991, -73.20508075688767},
				{100, -73.20508075688772},
				{128.55752193730785, -53.208888623795616},
				{128.55752193730788, 253.20888862379562},
				{100.00000000000003, 273.2050807568877},
				{71.44247806269215, 253.20888862379562},
			}},
		},
	}.verify(t)
}

func TestSelfIntersectionAvoidance(t *testing.T) {
	testCases{
		{
			op: polyclip.DIFFERENCE,
			subject: polyclip.Polygon{{
				{38.5721239031346, 172.33955556881023},
				{39.99999999999999, 171.3397459621556},
				{41.57979856674331, 170.60307379214092},
				{43.2635182233307, 170.15192246987792},
				{45, 170},
				{46.7364817766693, 170.15192246987792},
				{48.42020143325668, 170.60307379214092},
				{50, 171.3397459621556},
				{51.42787609686539, 172.33955556881023},
			}},
			clipping: polyclip.Polygon{{
				{51.42787609686539, 172.33955556881023},
				{50, 171.3397459621556},
				{48.42020143325668, 170.60307379214092},
				{46.7364817766693, 170.15192246987792},
				{45, 170},
				{43.2635182233307, 170.15192246987792},
				{42.78116786015871, 170.28116786015872},
				{42.65192246987792, 170.7635182233307},
				{42.5, 172},
			}},
			result: polyclip.Polygon{{
				{51.42787609686539, 172.33955556881023},
				{42.5, 172},
				{42.65192246987792, 170.7635182233307},
				{42.78116786015871, 170.28116786015872},
				// Should not contain this point: {43.2635182233307, 170.15192246987792},
				{41.57979856674331, 170.60307379214092},
				{39.99999999999999, 171.3397459621556},
				{38.5721239031346, 172.33955556881023},
			}},
		},
	}.verify(t)
}

func TestNonReductiveSegmentDivisions(t *testing.T) {
	if testing.Short() {
		return
	}

	cases := []struct{ subject, clipping polyclip.Polygon }{
		{
			// original reported github issue #4, resulting in infinite loop
			subject: polyclip.Polygon{{
				{X: 1.427255375e+06, Y: -2.3283064365386963e-10},
				{X: 1.4271285e+06, Y: 134.7111358642578},
				{X: 1.427109e+06, Y: 178.30108642578125}}},
			clipping: polyclip.Polygon{{
				{X: 1.416e+06, Y: -12000},
				{X: 1.428e+06, Y: -12000},
				{X: 1.428e+06, Y: 0},
				{X: 1.416e+06, Y: 0},
				{X: 1.416e+06, Y: -12000}}},
		},
		// Test cases from https://github.com/ctessum/polyclip-go/blob/master/bugs_test.go
		{
			subject: polyclip.Polygon{{
				{X: 1.7714672107465276e+06, Y: -102506.68254093888},
				{X: 1.7713768917571804e+06, Y: -102000.75485953009},
				{X: 1.7717109214841307e+06, Y: -101912.19625031832}}},
			clipping: polyclip.Polygon{{
				{X: 1.7714593229229522e+06, Y: -102470.35230830211},
				{X: 1.7714672107465276e+06, Y: -102506.68254093867},
				{X: 1.771439738086082e+06, Y: -102512.92027456204}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: -1.8280000000000012e+06, Y: -492999.99999999953},
				{X: -1.8289999999999995e+06, Y: -494000.0000000006},
				{X: -1.828e+06, Y: -493999.9999999991},
				{X: -1.8280000000000012e+06, Y: -492999.99999999953}}},
			clipping: polyclip.Polygon{{
				{X: -1.8280000000000005e+06, Y: -495999.99999999977},
				{X: -1.8280000000000007e+06, Y: -492000.0000000014},
				{X: -1.8240000000000007e+06, Y: -492000.0000000014},
				{X: -1.8280000000000005e+06, Y: -495999.99999999977}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: -2.0199999999999988e+06, Y: -394999.99999999825},
				{X: -2.0199999999999988e+06, Y: -392000.0000000009},
				{X: -2.0240000000000012e+06, Y: -395999.9999999993},
				{X: -2.0199999999999988e+06, Y: -394999.99999999825}}},
			clipping: polyclip.Polygon{{
				{X: -2.0199999999999988e+06, Y: -394999.99999999825},
				{X: -2.020000000000001e+06, Y: -394000.0000000001},
				{X: -2.0190000000000005e+06, Y: -394999.9999999997},
				{X: -2.0199999999999988e+06, Y: -394999.99999999825}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: -47999.99999999992, Y: -23999.999999998756},
				{X: 0, Y: -24000.00000000017},
				{X: 0, Y: 24000.00000000017},
				{X: -48000.00000000014, Y: 24000.00000000017},
				{X: -47999.99999999992, Y: -23999.999999998756}}},
			clipping: polyclip.Polygon{{
				{X: -48000, Y: -24000},
				{X: 0, Y: -24000},
				{X: 0, Y: 24000},
				{X: -48000, Y: 24000},
				{X: -48000, Y: -24000}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: -2.137000000000001e+06, Y: -122000.00000000093},
				{X: -2.1360000000000005e+06, Y: -121999.99999999907},
				{X: -2.1360000000000014e+06, Y: -121000.00000000186}}},
			clipping: polyclip.Polygon{{
				{X: -2.1120000000000005e+06, Y: -120000},
				{X: -2.136000000000001e+06, Y: -120000.00000000093},
				{X: -2.1360000000000005e+06, Y: -144000}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: 1.556e+06, Y: -1.139999999999999e+06},
				{X: 1.5600000000000002e+06, Y: -1.140000000000001e+06},
				{X: 1.56e+06, Y: -1.136000000000001e+06}}},
			clipping: polyclip.Polygon{{
				{X: 1.56e+06, Y: -1.127999999999999e+06},
				{X: 1.5600000000000002e+06, Y: -1.151999999999999e+06}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: 1.0958876176594219e+06, Y: -567467.5197556159},
				{X: 1.0956330600760083e+06, Y: -567223.72588934},
				{X: 1.0958876176594219e+06, Y: -567467.5197556159}}},
			clipping: polyclip.Polygon{{
				{X: 1.0953516248896217e+06, Y: -564135.1861293605},
				{X: 1.0959085007300845e+06, Y: -568241.1879245406},
				{X: 1.0955136237022132e+06, Y: -581389.3748769956}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: 608000, Y: -113151.36476426799},
				{X: 608000, Y: -114660.04962779157},
				{X: 612000, Y: -115414.39205955336},
				{X: 1.616e+06, Y: -300000},
				{X: 1.608e+06, Y: -303245.6575682382},
				{X: 0, Y: 0}}},
			clipping: polyclip.Polygon{{
				{X: 1.612e+06, Y: -296000}}},
		},
		{
			subject: polyclip.Polygon{{
				{X: 1.1458356382266793e+06, Y: -251939.4635597784},
				{X: 1.1460824662209095e+06, Y: -251687.86194535438},
				{X: 1.1458356382266793e+06, Y: -251939.4635597784}}},
			clipping: polyclip.Polygon{{
				{X: 1.1486683769211173e+06, Y: -251759.06331944838},
				{X: 1.1468807511323579e+06, Y: -251379.90576799586},
				{X: 1.1457914974731328e+06, Y: -251816.31287551578}}},
		},
		{
			// From https://github.com/ctessum/polyclip-go/commit/6614925d6d7087b7afcd4c55571554f67efd2ec3
			subject: polyclip.Polygon{{
				{X: 426694.6365274183, Y: -668547.1611580737},
				{X: 426714.57523030025, Y: -668548.9238652373},
				{X: 426745.39648089616, Y: -668550.4651249861}}},
			clipping: polyclip.Polygon{{
				{X: 426714.5752302991, Y: -668548.9238652373},
				{X: 426744.63718662335, Y: -668550.0591896093},
				{X: 426745.3964821229, Y: -668550.4652243527}}},
		},
		{
			// Produces invalid divisions that would otherwise continually generate new segments.
			subject: polyclip.Polygon{{
				{X: 99.67054939325573, Y: 23.50752393246498},
				{X: 99.88993946188153, Y: 20.999883973365655},
				{X: 100.01468418889, Y: 20.53433031419374}}},
			clipping: polyclip.Polygon{{
				{X: 100.15374164547939, Y: 20.015360821030836},
				{X: 95.64222842284941, Y: 36.85255738690467},
				{X: 100.15374164547939, Y: -14.714274712355238}}},
		},
	}

	for _, c := range cases {
		const rotations = 360
		// Test multiple rotations of each case to catch any orientation assumptions.
		for i := 0; i < rotations; i++ {
			angle := 2 * math.Pi * float64(i) / float64(rotations)
			subject := rotate(c.subject, angle)
			clipping := rotate(c.clipping, angle)

			for _, op := range []polyclip.Op{polyclip.UNION, polyclip.INTERSECTION, polyclip.DIFFERENCE} {
				ch := make(chan polyclip.Polygon)
				go func() {
					ch <- subject.Construct(op, clipping)
				}()

				select {
				case <-ch:
					// check that we get a result in finite time
				case <-time.After(1 * time.Second):
					// panicking in attempt to get full stacktrace
					panic(fmt.Sprintf("case %v:\nsubject:  %v\nclipping: %v\ntimed out.", op, subject, clipping))
				}
			}
		}
	}
}

func rotate(p polyclip.Polygon, radians float64) polyclip.Polygon {
	result := p.Clone()
	for i, contour := range p {
		result[i] = make(polyclip.Contour, len(contour))
		for j, point := range contour {
			result[i][j] = polyclip.Point{
				X: point.X*math.Cos(radians) - point.Y*math.Sin(radians),
				Y: point.Y*math.Cos(radians) + point.X*math.Sin(radians),
			}
		}
	}
	return result
}

func TestBug5(t *testing.T) {
	rect := polyclip.Polygon{{{24, 7}, {36, 7}, {36, 23}, {24, 23}}}
	circle := polyclip.Polygon{{{24, 7}, {24.83622770614123, 7.043824837053814}, {25.66329352654208, 7.174819194129555}, {26.472135954999587, 7.391547869638773}, {27.253893144606412, 7.691636338859195}, {28.00000000000001, 8.071796769724493}, {28.702282018339798, 8.527864045000424}, {29.35304485087088, 9.054841396180851}, {29.94515860381917, 9.646955149129141}, {30.472135954999597, 10.297717981660224}, {30.92820323027553, 11.00000000000001}, {31.308363661140827, 11.746106855393611}, {31.60845213036125, 12.527864045000435}, {31.825180805870467, 13.33670647345794}, {31.95617516294621, 14.16377229385879}, {32.00000000000002, 15.00000000000002}, {31.95617516294621, 15.83622770614125}, {31.825180805870467, 16.6632935265421}, {31.60845213036125, 17.472135954999604}, {31.308363661140827, 18.25389314460643}, {30.92820323027553, 19.00000000000003}, {30.472135954999597, 19.702282018339815}, {29.94515860381917, 20.353044850870898}, {29.35304485087088, 20.945158603819188}, {28.702282018339798, 21.472135954999615}, {28.00000000000001, 21.928203230275546}, {27.253893144606412, 22.308363661140845}, {26.472135954999587, 22.608452130361268}, {25.66329352654208, 22.825180805870485}, {24.83622770614123, 22.956175162946227}, {24, 23.00000000000004}, {23.16377229385877, 22.956175162946227}, {22.33670647345792, 22.825180805870485}, {21.527864045000413, 22.608452130361268}, {20.746106855393588, 22.308363661140845}, {19.99999999999999, 21.928203230275546}, {19.297717981660202, 21.472135954999615}, {18.64695514912912, 20.945158603819188}, {18.05484139618083, 20.353044850870898}, {17.527864045000403, 19.702282018339815}, {17.07179676972447, 19.00000000000003}, {16.691636338859173, 18.25389314460643}, {16.39154786963875, 17.472135954999604}, {16.174819194129533, 16.6632935265421}, {16.04382483705379, 15.83622770614125}, {15.999999999999977, 15.00000000000002}, {16.04382483705379, 14.16377229385879}, {16.174819194129533, 13.33670647345794}, {16.39154786963875, 12.527864045000435}, {16.691636338859173, 11.746106855393611}, {17.07179676972447, 11.00000000000001}, {17.527864045000403, 10.297717981660224}, {18.05484139618083, 9.646955149129141}, {18.64695514912912, 9.054841396180851}, {19.297717981660202, 8.527864045000424}, {19.99999999999999, 8.071796769724493}, {20.746106855393588, 7.691636338859194}, {21.527864045000413, 7.391547869638772}, {22.33670647345792, 7.1748191941295545}, {23.16377229385877, 7.043824837053813}}}

	testCases{
		{
			polyclip.UNION, rect, circle,
			polyclip.Polygon{{{36, 23}, {36, 7}, {24, 7}, {23.16377229385877, 7.043824837053813}, {22.33670647345792, 7.1748191941295545}, {21.527864045000413, 7.391547869638772}, {20.746106855393588, 7.691636338859194}, {19.99999999999999, 8.071796769724493}, {19.297717981660202, 8.527864045000424}, {18.64695514912912, 9.054841396180851}, {18.05484139618083, 9.646955149129141}, {17.527864045000403, 10.297717981660224}, {17.07179676972447, 11.00000000000001}, {16.691636338859173, 11.746106855393611}, {16.39154786963875, 12.527864045000435}, {16.174819194129533, 13.33670647345794}, {16.04382483705379, 14.16377229385879}, {15.999999999999977, 15.00000000000002}, {16.04382483705379, 15.83622770614125}, {16.174819194129533, 16.6632935265421}, {16.39154786963875, 17.472135954999604}, {16.691636338859173, 18.25389314460643}, {17.07179676972447, 19.00000000000003}, {17.527864045000403, 19.702282018339815}, {18.05484139618083, 20.353044850870898}, {18.64695514912912, 20.945158603819188}, {19.297717981660202, 21.472135954999615}, {19.99999999999999, 21.928203230275546}, {20.746106855393588, 22.308363661140845}, {21.527864045000413, 22.608452130361268}, {22.33670647345792, 22.825180805870485}, {23.16377229385877, 22.956175162946227}, {24, 23.00000000000004}}},
		},
		{
			polyclip.INTERSECTION, rect, circle,
			polyclip.Polygon{{{31.95617516294621, 15.83622770614125}, {31.825180805870467, 16.6632935265421}, {31.60845213036125, 17.472135954999604}, {31.308363661140827, 18.25389314460643}, {30.92820323027553, 19.00000000000003}, {30.472135954999597, 19.702282018339815}, {29.94515860381917, 20.353044850870898}, {29.35304485087088, 20.945158603819188}, {28.702282018339798, 21.472135954999615}, {28.00000000000001, 21.928203230275546}, {27.253893144606412, 22.308363661140845}, {26.472135954999587, 22.608452130361268}, {25.66329352654208, 22.825180805870485}, {24.83622770614123, 22.956175162946227}, {24, 23.00000000000004}, {24, 23}, {24, 7}, {24.83622770614123, 7.043824837053814}, {25.66329352654208, 7.174819194129555}, {26.472135954999587, 7.391547869638773}, {27.253893144606412, 7.691636338859195}, {28.00000000000001, 8.071796769724493}, {28.702282018339798, 8.527864045000424}, {29.35304485087088, 9.054841396180851}, {29.94515860381917, 9.646955149129141}, {30.472135954999597, 10.297717981660224}, {30.92820323027553, 11.00000000000001}, {31.308363661140827, 11.746106855393611}, {31.60845213036125, 12.527864045000435}, {31.825180805870467, 13.33670647345794}, {31.95617516294621, 14.16377229385879}, {32.00000000000002, 15.00000000000002}}},
		},
		{
			polyclip.DIFFERENCE, rect, circle,
			polyclip.Polygon{{{24, 23.00000000000004}, {24.83622770614123, 22.956175162946227}, {25.66329352654208, 22.825180805870485}, {26.472135954999587, 22.608452130361268}, {27.253893144606412, 22.308363661140845}, {28.00000000000001, 21.928203230275546}, {28.702282018339798, 21.472135954999615}, {29.35304485087088, 20.945158603819188}, {29.94515860381917, 20.353044850870898}, {30.472135954999597, 19.702282018339815}, {30.92820323027553, 19.00000000000003}, {31.308363661140827, 18.25389314460643}, {31.60845213036125, 17.472135954999604}, {31.825180805870467, 16.6632935265421}, {31.95617516294621, 15.83622770614125}, {32.00000000000002, 15.00000000000002}, {31.95617516294621, 14.16377229385879}, {31.825180805870467, 13.33670647345794}, {31.60845213036125, 12.527864045000435}, {31.308363661140827, 11.746106855393611}, {30.92820323027553, 11.00000000000001}, {30.472135954999597, 10.297717981660224}, {29.94515860381917, 9.646955149129141}, {29.35304485087088, 9.054841396180851}, {28.702282018339798, 8.527864045000424}, {28.00000000000001, 8.071796769724493}, {27.253893144606412, 7.691636338859195}, {26.472135954999587, 7.391547869638773}, {25.66329352654208, 7.174819194129555}, {24.83622770614123, 7.043824837053814}, {24, 7}, {36, 7}, {36, 23}}},
		},
		{
			polyclip.XOR, rect, circle,
			polyclip.Polygon{
				{{24, 23}, {24, 7}, {23.16377229385877, 7.043824837053813}, {22.33670647345792, 7.1748191941295545}, {21.527864045000413, 7.391547869638772}, {20.746106855393588, 7.691636338859194}, {19.99999999999999, 8.071796769724493}, {19.297717981660202, 8.527864045000424}, {18.64695514912912, 9.054841396180851}, {18.05484139618083, 9.646955149129141}, {17.527864045000403, 10.297717981660224}, {17.07179676972447, 11.00000000000001}, {16.691636338859173, 11.746106855393611}, {16.39154786963875, 12.527864045000435}, {16.174819194129533, 13.33670647345794}, {16.04382483705379, 14.16377229385879}, {15.999999999999977, 15.00000000000002}, {16.04382483705379, 15.83622770614125}, {16.174819194129533, 16.6632935265421}, {16.39154786963875, 17.472135954999604}, {16.691636338859173, 18.25389314460643}, {17.07179676972447, 19.00000000000003}, {17.527864045000403, 19.702282018339815}, {18.05484139618083, 20.353044850870898}, {18.64695514912912, 20.945158603819188}, {19.297717981660202, 21.472135954999615}, {19.99999999999999, 21.928203230275546}, {20.746106855393588, 22.308363661140845}, {21.527864045000413, 22.608452130361268}, {22.33670647345792, 22.825180805870485}, {23.16377229385877, 22.956175162946227}, {24, 23.00000000000004}},
				{{24, 23.00000000000004}, {24.83622770614123, 22.956175162946227}, {25.66329352654208, 22.825180805870485}, {26.472135954999587, 22.608452130361268}, {27.253893144606412, 22.308363661140845}, {28.00000000000001, 21.928203230275546}, {28.702282018339798, 21.472135954999615}, {29.35304485087088, 20.945158603819188}, {29.94515860381917, 20.353044850870898}, {30.472135954999597, 19.702282018339815}, {30.92820323027553, 19.00000000000003}, {31.308363661140827, 18.25389314460643}, {31.60845213036125, 17.472135954999604}, {31.825180805870467, 16.6632935265421}, {31.95617516294621, 15.83622770614125}, {32.00000000000002, 15.00000000000002}, {31.95617516294621, 14.16377229385879}, {31.825180805870467, 13.33670647345794}, {31.60845213036125, 12.527864045000435}, {31.308363661140827, 11.746106855393611}, {30.92820323027553, 11.00000000000001}, {30.472135954999597, 10.297717981660224}, {29.94515860381917, 9.646955149129141}, {29.35304485087088, 9.054841396180851}, {28.702282018339798, 8.527864045000424}, {28.00000000000001, 8.071796769724493}, {27.253893144606412, 7.691636338859195}, {26.472135954999587, 7.391547869638773}, {25.66329352654208, 7.174819194129555}, {24.83622770614123, 7.043824837053814}, {24, 7}, {36, 7}, {36, 23}},
			},
		},
	}.verify(t)
}

func TestSelfIntersect(t *testing.T) {
	rect1 := polyclip.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}}, {{1, 0}, {2, 0}, {2, 1}, {1, 1}}}
	rect2 := polyclip.Polygon{{{0, 0.25}, {3, 0.25}, {3, 0.75}, {0, 0.75}}}

	expected := []struct {
		name   string
		op     polyclip.Op
		result polyclip.Polygon
	}{
		{
			"union",
			polyclip.UNION,
			polyclip.Polygon{{{0, 0}, {1, 0}, {2, 0}, {2, 0.25}, {3, 0.25}, {3, 0.75}, {2, 0.75}, {2, 1}, {1, 1}, {0, 1}, {0, 0.75}, {0, 0.25}}},
		},
		{
			"intersection",
			polyclip.INTERSECTION,
			polyclip.Polygon{{{0, 0.25}, {2, 0.25}, {2, 0.75}, {0, 0.75}}},
		},
		{
			"difference",
			polyclip.DIFFERENCE,
			polyclip.Polygon{{{0, 0}, {1, 0}, {2, 0}, {2, 0.25}, {0, 0.25}}, {{0, 0.75}, {2, 0.75}, {2, 1}, {1, 1}, {0, 1}}},
		},
		{
			"xor",
			polyclip.XOR,
			// TODO: This one is a little weird.  It probably shouldn't be self-intersecting.
			polyclip.Polygon{{{0, 0}, {1, 0}, {2, 0}, {2, 0.25}, {0, 0.25}}, {{0, 0.75}, {2, 0.75}, {2, 0.25}, {3, 0.25}, {3, 0.75}, {2, 0.75}, {2, 1}, {1, 1}, {0, 1}}},
		},
	}

	for _, e := range expected {
		t.Run(e.name, func(t *testing.T) {
			result := rect1.Simplify().Construct(e.op, rect2.Simplify())
			if dump(result) != dump(e.result) {
				t.Errorf("case %d expected:\n%v\ngot:\n%v", e.op, dump(e.result), dump(result))
			}
		})
	}
}

// Bug test from b4c12673bc80394c472b18f168a042f904ec948.
func TestDifference_bug(t *testing.T) {
	p1 := polyclip.Polygon{
		polyclip.Contour{
			polyclip.Point{X: 99, Y: 164}, polyclip.Point{X: 114, Y: 108},
			polyclip.Point{X: 121, Y: 164},
		},
	}

	p2 := polyclip.Polygon{polyclip.Contour{
		polyclip.Point{X: 114, Y: 0}, polyclip.Point{X: 161, Y: 0},
		polyclip.Point{X: 114, Y: 168},
	}}
	want := polyclip.Polygon{
		polyclip.Contour{{114, 168}, {114, 164}, {115.11904761904762, 164}},
		polyclip.Contour{{114, 0}, {161, 0}, {119.18382352941177, 149.47058823529412}, {114, 108}},
	}
	result := p2.Construct(polyclip.DIFFERENCE, p1)
	if dump(want) != dump(result) {
		t.Errorf("expected:\n%v\ngot:\n%v", dump(want), dump(result))
	}
}

// Bug test from 6614925d6d7087b7afcd4c55571554f67efd2ec3
func TestInfiniteLoopBug(t *testing.T) {
	subject := polyclip.Polygon{polyclip.Contour{
		polyclip.Point{X: 426694.6365274183, Y: -668547.1611580737},
		polyclip.Point{X: 426714.57523030025, Y: -668548.9238652373},
		polyclip.Point{X: 426745.39648089616, Y: -668550.4651249861},
	}}
	clipping := polyclip.Polygon{polyclip.Contour{
		polyclip.Point{X: 426714.5752302991, Y: -668548.9238652373},
		polyclip.Point{X: 426744.63718662335, Y: -668550.0591896093},
		polyclip.Point{X: 426745.3964821229, Y: -668550.4652243527},
	}}

	want := polyclip.Polygon{polyclip.Contour{
		polyclip.Point{X: 426731.5895193888, Y: -668549.5664294426},
		polyclip.Point{X: 426714.57523030025, Y: -668548.9238652373},
		polyclip.Point{X: 426694.6365274183, Y: -668547.1611580737},
	},
		polyclip.Contour{
			polyclip.Point{X: 426745.39648089616, Y: -668550.4651249861},
			polyclip.Point{X: 426745.3962772624, Y: -668550.4651148032},
			polyclip.Point{X: 426745.39627072256, Y: -668550.4651113059},
		},
	}

	result := subject.Construct(polyclip.DIFFERENCE, clipping)
	if dump(want) != dump(result) {
		t.Errorf("expected:\n%v\ngot:\n%v", dump(want), dump(result))
	}
}
