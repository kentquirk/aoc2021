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

type Point struct {
	X int
	Y int
}

type Fold struct {
	Direction  string
	Coordinate int
}

type Paper struct {
	P map[Point]struct{}
	F []Fold
}

func (p *Paper) Fold(f Fold) {
	foldedPoints := make(map[Point]struct{})

	for pt := range p.P {
		switch f.Direction {
		case "x":
			if pt.X > f.Coordinate {
				pt.X = f.Coordinate - (pt.X - f.Coordinate)
			}
		case "y":
			if pt.Y > f.Coordinate {
				pt.Y = f.Coordinate - (pt.Y - f.Coordinate)
			}
		default:
			panic("oops")
		}
		foldedPoints[pt] = struct{}{}
	}
	p.P = foldedPoints
}

type PointSlice []Point

func (s PointSlice) Len() int      { return len(s) }
func (s PointSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s PointSlice) Less(i, j int) bool {
	if s[i].Y < s[j].Y {
		return true
	} else if s[i].Y > s[j].Y {
		return false
	} else {
		return s[i].X < s[j].X
	}
}

func (p *Paper) Print() {
	// var output []string
	var points []Point

	for k := range p.P {
		points = append(points, k)
	}

	sort.Sort(PointSlice(points))
	line := 0
	col := 0
	for _, pt := range points {
		for pt.Y > line {
			fmt.Println()
			col = 0
			line++
		}
		for pt.X >= col {
			fmt.Print(" ")
			col++
		}
		fmt.Print("#")
		col++
	}
	fmt.Println()
}

func parse(points string, folds string) Paper {
	paper := Paper{P: make(map[Point]struct{})}
	for _, p := range strings.Split(points, "\n") {
		coords := strings.Split(p, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		paper.P[Point{X: x, Y: y}] = struct{}{}
	}
	for _, f := range strings.Split(folds, "\n") {
		parts := strings.Split(f, "=")
		coord, _ := strconv.Atoi(parts[1])
		paper.F = append(paper.F, Fold{
			Direction:  parts[0][len(parts[0])-1:],
			Coordinate: coord,
		})
	}
	return paper
}

func day13a(paper Paper) int {
	// fmt.Println(paper)
	paper.Fold(paper.F[0])
	// fmt.Println(paper)
	return len(paper.P)
}

func day13b(paper Paper) {
	for _, f := range paper.F {
		paper.Fold(f)
	}
	paper.Print()
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
	parts := strings.Split(string(b), "\n\n")
	paper := parse(parts[0], parts[1])

	fmt.Println(day13a(paper))
	day13b(paper)
}
