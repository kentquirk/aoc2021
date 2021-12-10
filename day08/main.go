package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/kentquirk/stringset/v2"
)

type Pattern struct {
	Input     []string
	Output    []string
	Possibles [][]int
	Digits    []int
}

func NewPattern(s string) *Pattern {
	parts := strings.Split(s, "|")
	r := regexp.MustCompile("[a-g]+")
	return &Pattern{
		Input:  r.FindAllString(parts[0], -1),
		Output: r.FindAllString(parts[1], -1),
		Digits: make([]int, 0),
	}
}

func setFromString(s string) *stringset.StringSet {
	ss := stringset.New()
	for i := 0; i < len(s); i++ {
		ss.Add(s[i : i+1])
	}
	return ss
}

func numberSets() []*stringset.StringSet {
	template := setFromString("abcdefg")
	ss := make([]*stringset.StringSet, 0)
	for i := 0; i < 10; i++ {
		ss = append(ss, template.Clone())
	}
	return ss
}

// func segmentCounts() map[int]*stringset.StringSet {
// 	template := setFromString("abcdefg")
// 	ss := make(map[int]*stringset.StringSet)
// 	for i := 0; i < 10; i++ {
// 		ss = append(ss, template.Clone())
// 	}
// 	ss[2] = append(ss, template.Clone())
// 	return ss
// }

func day08a(pats []*Pattern) int {
	count := 0
	for _, p := range pats {
		for _, s := range p.Output {
			if len(s) == 2 || len(s) == 3 || len(s) == 4 || len(s) == 7 {
				count++
			}
		}
	}
	return count
}

// var digitCounts = []int{6, 2, 5, 5, 4, 5, 6, 3, 7, 6}

//   0:      1:      2:      3:      4:
//  0000    ....    0000    0000    ....
// 1    2  .    2  .    2  .    2  1    2
// 1    2  .    2  .    2  .    2  1    2
//  ....    ....    3333    3333    3333
// 4    5  .    5  4    .  .    5  .    5
// 4    5  .    5  4    .  .    5  .    5
//  6666    ....    6666    6666    ....

//  5:      6:      7:      8:      9:
//  0000    0000    0000    0000    0000
// 1    .  1    .  .    2  1    2  1    2
// 1    .  1    .  .    2  1    2  1    2
//  3333    3333    ....    3333    3333
// .    5  4    5  .    5  4    5  .    5
// .    5  4    5  .    5  4    5  .    5
//  6666    6666    ....    6666    6666

// Data structures I need:
// map[int][]int -- Input index -> the possible digits they represent
// map[digit]stringset -- Digits -> the set of characters we know are in them
// map[character] ->

func solve(pat *Pattern) int {
	ns := numberSets()
	var fives []*stringset.StringSet
	// var sixes []*stringset.StringSet
	for _, inp := range pat.Input {
		keeps := setFromString(inp)
		switch len(inp) {
		case 2:
			// digit 1
			ns[1] = ns[1].Intersection(keeps)
		case 3:
			// digit 7
			ns[7] = ns[7].Intersection(keeps)
		case 4:
			// digit 4
			ns[4] = ns[4].Intersection(keeps)
		case 5:
			// either 2, 3, or 5
			fives = append(fives, keeps)
		case 6:
			// either 0, 6, or 9
			// sixes = append(sixes, keeps)
		case 7:
			// 7 is all 8 digits
		default:
			panic("oops")
		}
	}

	var result [7]*stringset.StringSet
	// set difference of '7' and '1' is 0
	result[0] = ns[7].Difference(ns[1])
	// the fives -- '2' and '5' and '3' is 0 3 6, and we can eliminate 0
	// then '4' - '1' - x36 == 1,  '4' and x36 == 3, x36 - 3 == 6
	x36 := fives[0].Intersection(fives[1]).Intersection(fives[2]).Difference(result[0])
	result[1] = ns[4].Difference(ns[1]).Difference(x36)
	result[3] = ns[4].Intersection(x36)
	result[6] = x36.Difference(result[3])
	// we can iterate all the fives and subtract 0136; only '5' will have 1 remaining, and that one is 5
	x0136 := result[0].Union(result[1]).Union(result[3]).Union(result[6])

	for _, t := range fives {
		d := t.Difference(x0136)
		if d.Length() == 1 {
			result[5] = d
			ns[5] = t
			break
		}
	}
	// '1' - 5 is 2
	result[2] = ns[1].Difference(result[5])
	// and '8' - all the rest is 4
	rest := ns[8]
	for _, t := range result {
		if t != nil {
			rest = rest.Difference(t)
		}
	}
	result[4] = rest

	ns[0] = ns[8].Difference(result[3])
	ns[6] = ns[8].Difference(result[2])
	ns[9] = ns[8].Difference(result[4])
	ns[2] = ns[8].Difference(result[1]).Difference(result[5])
	ns[3] = ns[8].Difference(result[1]).Difference(result[4])

	// for i := range result {
	// 	fmt.Printf("%d: %v\n", i, result[i].Strings())
	// }
	// fmt.Println("---")
	// for i, s := range ns {
	// 	fmt.Printf("%d: %s\n", i, s.Strings())
	// }

	output := 0
	for _, o := range pat.Output {
		os := setFromString(o)
		for i, s := range ns {
			if s.Equals(os) {
				output = output*10 + i
				break
			}
		}
	}

	// fmt.Println(output)
	return output
}

func day08b(pats []*Pattern) int {
	total := 0
	for _, p := range pats {
		n := solve(p)
		total += n
	}
	return total
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
	var pats []*Pattern
	for _, l := range lines {
		pats = append(pats, NewPattern(l))
	}
	fmt.Println(day08a(pats))
	fmt.Println(day08b(pats))
}
