package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func parse(s string) (int, int) {
	errscores := map[rune]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}
	fixscores := map[rune]int{
		')': 1,
		']': 2,
		'}': 3,
		'>': 4,
	}

	var expected []rune
	for _, r := range s {
		switch r {
		case '(':
			expected = append(expected, ')')
		case '<':
			expected = append(expected, '>')
		case '[':
			expected = append(expected, ']')
		case '{':
			expected = append(expected, '}')
		case '}', ')', ']', '>':
			want := expected[len(expected)-1]
			expected = expected[:len(expected)-1]
			if want != r {
				return errscores[r], 0
			}
		default:
			panic("oops")
		}
	}

	fixscore := 0
	for i := len(expected) - 1; i >= 0; i-- {
		fixscore = 5*fixscore + fixscores[expected[i]]
	}

	return 0, fixscore
}

func day10a(lines []string) int {
	errscore := 0
	for _, l := range lines {
		e, _ := parse(l)
		errscore += e
	}
	return errscore
}

func day10b(lines []string) int {
	var fixscores []int
	for _, l := range lines {
		_, fixscore := parse(l)
		if fixscore != 0 {
			fixscores = append(fixscores, fixscore)
		}
	}

	sort.Ints(fixscores)
	return fixscores[len(fixscores)/2]
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
	fmt.Println(day10a(lines))
	fmt.Println(day10b(lines))
}
