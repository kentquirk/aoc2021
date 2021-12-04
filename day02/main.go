package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Submarine struct {
	aim      int
	position int
	depth    int
}

func (s *Submarine) Report() int {
	return s.position * s.depth
}

func (s *Submarine) day02a(lines []string) {
	for _, line := range lines {
		splits := strings.Split(line, " ")
		n, err := strconv.Atoi(splits[1])
		if err != nil {
			panic(err)
		}
		switch splits[0] {
		case "forward":
			s.position += n
		case "down":
			s.depth += n
		case "up":
			s.depth -= n
		default:
			panic(line)
		}
	}
}

func (s *Submarine) day02b(lines []string) {
	for _, line := range lines {
		splits := strings.Split(line, " ")
		x, err := strconv.Atoi(splits[1])
		if err != nil {
			panic(err)
		}
		switch splits[0] {
		case "forward":
			s.position += x
			s.depth += s.aim * x
		case "down":
			s.aim += x
		case "up":
			s.aim -= x
		default:
			panic(line)
		}
	}
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
	subA := new(Submarine)
	subA.day02a(lines)
	fmt.Printf("Sub A: %d\n", subA.Report())
	subB := new(Submarine)
	subB.day02b(lines)
	fmt.Printf("Sub B: %d\n", subB.Report())
}
