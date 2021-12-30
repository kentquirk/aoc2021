package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point [3]int

func (p Point) Add(rhs Point) Point {
	return Point{rhs[0] + p[0], rhs[1] + p[1], rhs[2] + p[2]}
}

func (p Point) VectorTo(rhs Point) Point {
	return Point{rhs[0] - p[0], rhs[1] - p[1], rhs[2] - p[2]}
}

func (p Point) Less(rhs Point) bool {
	return p[0] < rhs[0] ||
		p[0] == rhs[0] && p[1] < rhs[1] ||
		p[0] == rhs[0] && p[1] == rhs[1] && p[2] < rhs[2]
}

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p Point) Manhattan(rhs Point) int {
	return iabs(rhs[0]-p[0]) + iabs(rhs[1]-p[1]) + iabs(rhs[2]-p[2])
}

func (p Point) Dist2(rhs Point) int {
	dx := rhs[0] - p[0]
	dy := rhs[1] - p[1]
	dz := rhs[2] - p[2]
	return dx*dx + dy*dy + dz*dz
}

var Selections = [][3]int{
	{0, 1, 2},
	{0, 2, 1},
	{1, 2, 0},
	{1, 0, 2},
	{2, 0, 1},
	{2, 1, 0},
}

func flip(x int, b bool) int {
	if b {
		return -x
	}
	return x
}

func (p Point) Reorient(selection int, signs int) Point {
	return Point{
		flip(p[Selections[selection][0]], signs&1 != 0),
		flip(p[Selections[selection][1]], signs&2 != 0),
		flip(p[Selections[selection][2]], signs&4 != 0),
	}
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p[0], p[1], p[2])
}

type PointSlice []Point

func (p PointSlice) Len() int           { return len(p) }
func (p PointSlice) Less(i, j int) bool { return p[i].Less(p[j]) }
func (p PointSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p PointSlice) Contains(pt Point) bool {
	for i := range p {
		if p[i] == pt {
			return true
		}
	}
	return false
}

type Scanner struct {
	Points    PointSlice
	Selection int // 0 to 5
	Signs     int // 0 to 7
	Location  Point
}

func NewScanner() *Scanner {
	return &Scanner{
		Points:    make(PointSlice, 0),
		Selection: 0,
		Signs:     0,
		Location:  Point{0, 0, 0},
	}
}

func (s *Scanner) Reorient(sel, signs int) {
	for i := range s.Points {
		s.Points[i] = s.Points[i].Reorient(sel, signs)
	}
}

func MakePoint(s string) Point {
	var pt Point
	vs := strings.Split(s, ",")
	for i := range vs {
		pt[i], _ = strconv.Atoi(vs[i])
	}
	return pt
}

// We have two point clouds (a & b) for a pair of scanners.
// For each possible rotation:
//   for each point in A
//		for each point in B
//			calculate an offset to move B to A
// 			for all points in B
//				if you move them by offset see if they're in A; if so, increment counter
//			if counter is 12 we found a match

// optimization -- return bool, exit as soon as count hits 12/
func countOverlap(p1, p2 PointSlice, offset Point) int {
	count := 0
	for i := range p2 {
		if p1.Contains(p2[i].Add(offset)) {
			count++
		}
	}
	return count
}

func comparePointSlices(p1, p2 PointSlice) (Point, bool) {
	for i := range p1 {
		for j := range p2 {
			offset := p2[j].VectorTo(p1[i])
			if countOverlap(p1, p2, offset) >= 12 {
				return offset, true
			}
		}
	}
	return Point{}, false
}

// iterate through every possible orientation to try to find a match
// if we find one, reorient and position the target scanner and return the offset and true
func (s *Scanner) SearchForMatch(base *Scanner) (Point, bool) {
	for sel := 0; sel < 6; sel++ {
		for sign := 0; sign < 8; sign++ {
			pts := make(PointSlice, len(s.Points))
			for i := range s.Points {
				pts[i] = s.Points[i].Reorient(sel, sign)
			}
			if offset, found := comparePointSlices(base.Points, pts); found {
				s.Reorient(sel, sign)
				return offset, true
			}
		}
	}

	return Point{}, false
}

func day19a(lines []string) (int, int) {
	unmatched := make(map[int]*Scanner)
	matched := make(map[int]*Scanner)
	var scanner *Scanner
	var key = -1
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		if strings.HasPrefix(l, "---") {
			key++
			scanner = NewScanner()
			unmatched[key] = scanner
			continue
		}
		scanner.Points = append(scanner.Points, MakePoint(l))
	}
	matched[0] = unmatched[0]
	delete(unmatched, 0)
	sort.Sort(PointSlice(matched[0].Points))

	// while there are still items in the unmatched list
	for len(unmatched) > 0 {
	outer:
		for baseKey, baseScanner := range matched {
			for unmatchedKey, unmatchedScanner := range unmatched {
				if offset, found := unmatchedScanner.SearchForMatch(baseScanner); found {
					fmt.Printf("Found match: base(%d) unmatched(%d)\n", baseKey, unmatchedKey)
					unmatchedScanner.Location = offset.Add(baseScanner.Location)
					matched[unmatchedKey] = unmatched[unmatchedKey]
					delete(unmatched, unmatchedKey)
					break outer
				}
			}
		}
	}
	beacons := make(map[Point]struct{})
	for k, s := range matched {
		fmt.Printf("%d: %v\n", k, s.Location)
		for _, pt := range s.Points {
			beacons[pt.Add(s.Location)] = struct{}{}
		}
	}
	// we know that the keys are successive integers, so we're going to iterate the matched list
	maxDist := 0
	for i := 0; i < len(matched)-1; i++ {
		for j := i + 1; j < len(matched); j++ {
			d := matched[i].Location.Manhattan(matched[j].Location)
			if d > maxDist {
				maxDist = d
			}
		}
	}
	return len(beacons), maxDist
}

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	fmt.Println(day19a(lines))
}
