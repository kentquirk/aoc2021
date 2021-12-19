package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Probe struct {
	Loc  Point
	Vel  Point
	YMax int
}

type Area struct {
	Xmin int
	Ymin int
	Xmax int
	Ymax int
}

func ParseArea(s string) *Area {
	re := regexp.MustCompile("x=([0-9-]+)..([0-9-]+), y=([0-9-]+)..([0-9-]+)")
	ms := re.FindStringSubmatch(s)
	a := Area{}
	a.Xmin, _ = strconv.Atoi(ms[1])
	a.Xmax, _ = strconv.Atoi(ms[2])
	a.Ymin, _ = strconv.Atoi(ms[3])
	a.Ymax, _ = strconv.Atoi(ms[4])
	return &a
}

func (a *Area) IsBelow(pt Point) bool {
	return pt.Y < a.Ymin
}

func (a *Area) IsBeyond(pt Point) bool {
	return pt.X > a.Xmax
}

func (a *Area) IsInside(pt Point) bool {
	return !(pt.X < a.Xmin || pt.X > a.Xmax || pt.Y < a.Ymin || pt.Y > a.Ymax)
}

// Fire sends one probe and returns the results --
// max height, final x velocity, beyond, and below
// If it's a hit, then both beyond and below will be false
func (a *Area) Fire(dx, dy int) (int, bool, bool) {
	probe := Probe{
		Loc:  Point{},
		Vel:  Point{X: dx, Y: dy},
		YMax: 0,
	}
	for {
		// adjust probe
		probe.Loc.X += probe.Vel.X
		probe.Loc.Y += probe.Vel.Y
		switch {
		case probe.Vel.X > 0:
			probe.Vel.X--
		case probe.Vel.X < 0:
			probe.Vel.X++
		}
		probe.Vel.Y--
		if probe.Loc.Y > probe.YMax {
			probe.YMax = probe.Loc.Y
		}
		// fmt.Printf("current %v %v\n", probe, a)

		if a.IsInside(probe.Loc) {
			// fmt.Println("inside")
			return probe.YMax, false, false
		}

		if a.IsBeyond(probe.Loc) {
			// fmt.Printf("beyond %v", probe)
			return probe.YMax, true, false
		}

		if a.IsBelow(probe.Loc) {
			// fmt.Printf("beyond %v\n", probe)
			return probe.YMax, false, true
		}
	}
}

// We start out by firing almost straight down at the maximum velocity
// that could possibly hit.
// If we drop below without going beyond, we increase x velocity
// If we go beyond without going below, we increase y and decrease x
// If we hit, we record the height and velocities, then increase y and decrease x
// if the first attempt after a hit is beyond, we can't do any better.
// We record all the hit velocities, and afterward, we iterate around those positions
// to try to find more.
func day17a(area *Area) (int, int) {
	hits := make(map[Point]int)
	// seed it with the fastest direct shot that will work
	hits[Point{area.Xmax, area.Ymin}] = 0
	dx := 1
	dy := area.Ymin
	maxheight := 0
	first := true
	for dx != 0 && dy < 1000 {
		// fmt.Println("trying ", dx, dy)
		height, beyond, below := area.Fire(dx, dy)
		switch {
		case !below && !beyond: // a hit
			if height > maxheight {
				maxheight = height
			}
			hits[Point{X: dx, Y: dy}] = maxheight
			// fmt.Printf("hit %d with (%d, %d)", maxheight, dx, dy)
			dy++
			dx--
			first = true
		case below:
			dx++
		case beyond:
			if first {
				return maxheight, len(hits)
			}
			dy++
			dx--
		}
		first = false
	}

	// fmt.Println(len(hits), hits)

	slop := 3
	for {
		newhits := make(map[Point]int)
		for vel := range hits {
			for x := -slop; x <= slop; x++ {
				for y := -slop; y <= slop; y++ {
					trial := Point{vel.X + x, vel.Y + y}
					if _, found := hits[trial]; !found {
						if height, beyond, below := area.Fire(trial.X, trial.Y); !beyond && !below {
							newhits[trial] = height
						}
					}
				}
			}
		}
		if len(newhits) == 0 {
			break
		}
		for p, h := range newhits {
			hits[p] = h
		}
	}

	// fmt.Println(len(hits), hits)

	return maxheight, len(hits)
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
	area := ParseArea(lines[0])
	fmt.Println(day17a(area))
}
