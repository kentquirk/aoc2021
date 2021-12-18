package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"

	"github.com/beefsack/go-astar"
)

type Point struct {
	X int
	Y int
}

type Position struct {
	Loc       Point
	Risk      int
	Neighbors []*Position
}

func NewPosition(pt Point, risk byte) *Position {
	return &Position{
		Loc:  pt,
		Risk: int(risk),
	}
}

// we need *position to implement the Pather interface

func (p *Position) PathNeighbors() []astar.Pather {
	pathers := make([]astar.Pather, 0)
	for _, n := range p.Neighbors {
		pathers = append(pathers, n)
	}
	return pathers
}

func (p *Position) PathNeighborCost(to astar.Pather) float64 {
	return float64(p.Risk)
}

func (p *Position) PathEstimatedCost(to astar.Pather) float64 {
	return math.Abs(float64(p.Loc.X-to.(*Position).Loc.X)) +
		math.Abs(float64(p.Loc.Y-to.(*Position).Loc.Y))
}

type Cave struct {
	Positions [][]*Position
	Score     int
}

func (c *Cave) Width() int {
	return len(c.Positions[0])
}

func (c *Cave) Height() int {
	return len(c.Positions)
}

func (c *Cave) AddNeighbors() {
	// add the appropriate number of neighbors
	for row, psns := range c.Positions {
		for col, pos := range psns {
			if row > 0 {
				pos.Neighbors = append(pos.Neighbors, c.Positions[row-1][col])
			}
			if row < len(c.Positions)-1 {
				pos.Neighbors = append(pos.Neighbors, c.Positions[row+1][col])
			}
			if col > 0 {
				pos.Neighbors = append(pos.Neighbors, c.Positions[row][col-1])
			}
			if col < len(c.Positions[row])-1 {
				pos.Neighbors = append(pos.Neighbors, c.Positions[row][col+1])
			}
		}
	}
}

func NewCave(input []string) *Cave {
	var positions [][]*Position
	for y, row := range input {
		oct := make([]*Position, len(row))
		for x := range row {
			oct[x] = NewPosition(Point{X: x, Y: y}, row[x]-byte('0'))
		}
		positions = append(positions, oct)
	}

	cave := &Cave{Positions: positions}
	cave.AddNeighbors()
	return cave
}

func NewCave5(input []string) *Cave {
	positions := make([][]*Position, 0)
	for i := 0; i < 5; i++ {
		for y, row := range input {
			oct := make([]*Position, 0)
			for j := 0; j < 5; j++ {
				for x := range row {
					b := ((row[x] - byte('1') + byte(i+j)) % 9) + 1
					oct = append(oct, NewPosition(Point{X: x + j*len(row), Y: y + i*len(input)}, b))
				}
			}
			positions = append(positions, oct)
		}
	}

	cave := &Cave{Positions: positions}
	// cave.PrintWithPath(make([]astar.Pather, 0))
	cave.AddNeighbors()
	return cave
}

func (c *Cave) PrintWithPath(rawpath []astar.Pather) {
	path := make(map[Point]struct{})
	for _, r := range rawpath {
		path[r.(*Position).Loc] = struct{}{}
	}
	for _, psns := range c.Positions {
		for _, pos := range psns {
			if _, found := path[pos.Loc]; found {
				fmt.Printf("\x1b[0;34m%d\x1b[0m", pos.Risk)
			} else {
				fmt.Printf("%d", pos.Risk)
			}
		}
		fmt.Println()
	}
}

func day15a(cave *Cave) int {
	start := cave.Positions[0][0]
	exit := cave.Positions[cave.Height()-1][cave.Width()-1]

	path, distance, found := astar.Path(exit, start)
	if !found {
		return -1
	}
	for _, p := range path {
		fmt.Println(p.(*Position).Loc, p.(*Position).Risk)
	}
	// cave.PrintWithPath(path)
	return int(distance)
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
	cave := NewCave(lines)
	fmt.Println(day15a(cave))

	cave2 := NewCave5(lines)
	fmt.Println(day15a(cave2))
}
