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

func (p Point) Less(other Point) bool {
	if p.X < other.X {
		return true
	}
	if p.X == other.X && p.Y < other.Y {
		return true
	}
	return false
}

type Line struct {
	P1 Point
	P2 Point
}

func NewLine(p1 Point, p2 Point) *Line {
	if p1.Less(p2) {
		return &Line{P1: p1, P2: p2}
	}
	return &Line{P1: p2, P2: p1}
}

func (l Line) IsHorizontal() bool {
	return l.P1.Y == l.P2.Y
}

func (l Line) IsVertical() bool {
	return l.P1.X == l.P2.X
}

func (l Line) IsDiagonal() bool {
	return l.P1.X != l.P2.X && l.P1.Y != l.P2.Y
}

func (l Line) DrawHV(grid map[Point]int) {
	if l.IsHorizontal() {
		for x := l.P1.X; x <= l.P2.X; x++ {
			grid[Point{X: x, Y: l.P1.Y}]++
		}
	} else if l.IsVertical() {
		for y := l.P1.Y; y <= l.P2.Y; y++ {
			grid[Point{X: l.P1.X, Y: y}]++
		}
	}
}

func (l Line) Draw(grid map[Point]int) {
	switch {
	case l.IsHorizontal():
		for x := l.P1.X; x <= l.P2.X; x++ {
			grid[Point{X: x, Y: l.P1.Y}]++
		}
	case l.IsVertical():
		for y := l.P1.Y; y <= l.P2.Y; y++ {
			grid[Point{X: l.P1.X, Y: y}]++
		}
	case l.IsDiagonal():
		y := l.P1.Y
		for x := l.P1.X; x <= l.P2.X; x++ {
			grid[Point{X: x, Y: y}]++
			if y < l.P2.Y {
				y++
			} else {
				y--
			}
		}
	}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func parseLine(s string) *Line {
	pat := regexp.MustCompile("([0-9]+),([0-9]+) -> ([0-9]+),([0-9]+)")
	data := pat.FindStringSubmatch(s)
	return NewLine(
		Point{X: atoi(data[1]), Y: atoi(data[2])},
		Point{X: atoi(data[3]), Y: atoi(data[4])},
	)
}

func day05a(input []string) int {
	grid := make(map[Point]int)

	for _, inp := range input {
		l := parseLine(inp)
		// fmt.Println(l)
		l.DrawHV(grid)
	}
	// fmt.Println(grid)

	count := 0
	for _, v := range grid {
		if v > 1 {
			count++
		}
	}
	return count
}

func day05b(input []string) int {
	grid := make(map[Point]int)

	for _, inp := range input {
		l := parseLine(inp)
		// fmt.Println(l)
		l.Draw(grid)
	}
	// fmt.Println(grid)

	count := 0
	for _, v := range grid {
		if v > 1 {
			count++
		}
	}
	return count
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
	fmt.Println("A: ", day05a(lines))
	fmt.Println("B: ", day05b(lines))
}
