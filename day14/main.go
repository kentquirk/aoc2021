package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
)

type Pair struct {
	Left  byte
	Right byte
}

func NewPair(a, b byte) Pair {
	return Pair{Left: a, Right: b}
}

type Insertion struct {
	Deltas map[Pair]int
}

func NewInsertion(src Pair, insert byte) Insertion {
	d := make(map[Pair]int)
	d[NewPair(src.Left, insert)]++
	d[NewPair(insert, src.Right)]++
	return Insertion{Deltas: d}
}

type Polymer struct {
	Pairs map[Pair]int
	Ends  Pair
}

func NewPolymer(s string) *Polymer {
	pairs := make(map[Pair]int)
	for i := 0; i < len(s)-1; i++ {
		pairs[NewPair(s[i], s[i+1])]++
	}
	return &Polymer{
		Pairs: pairs,
		Ends:  NewPair(s[0], s[len(s)-1]),
	}
}

func (p *Polymer) Print() {
	for k, v := range p.Pairs {
		fmt.Printf("%c%c: %d ", k.Left, k.Right, v)
	}
	fmt.Println()
}

func (p *Polymer) ApplyInsertions(insertions map[Pair]Insertion) {
	newPairs := make(map[Pair]int)
	for pair, count := range p.Pairs {
		for dpair, delta := range insertions[pair].Deltas {
			newPairs[dpair] += delta * count
		}
	}
	p.Pairs = newPairs
}

func (p *Polymer) Score() int {
	// this counts twice as high because everything is double-counted,
	// so divide by 2 at the end
	scores := make(map[byte]int)
	for pair, count := range p.Pairs {
		scores[pair.Left] += count
		scores[pair.Right] += count
	}
	scores[p.Ends.Left]++
	scores[p.Ends.Right]++
	max := 0
	min := math.MaxInt
	for _, v := range scores {
		if max < v {
			max = v
		}
		if min > v {
			min = v
		}
	}
	return (max - min) / 2
}

func parseInsertions(lines []string) map[Pair]Insertion {
	insertions := make(map[Pair]Insertion)

	for _, l := range lines {
		parts := strings.Split(l, "->")
		s := strings.Trim(parts[0], " ")
		pair := NewPair(s[0], s[1])
		i := NewInsertion(pair, strings.Trim(parts[1], " ")[0])
		insertions[pair] = i
	}
	return insertions
}

func day14a(iterations int, polymer *Polymer, insertions map[Pair]Insertion) int {
	for step := 0; step < iterations; step++ {
		polymer.ApplyInsertions(insertions)
		// fmt.Printf("%s\n", polymer)
	}
	return polymer.Score()
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
	polymer := NewPolymer(lines[0])
	insertions := parseInsertions(lines[2:])
	fmt.Println(day14a(10, polymer, insertions))
	polymer2 := NewPolymer(lines[0])
	fmt.Println(day14a(40, polymer2, insertions))
}
