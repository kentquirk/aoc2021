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

func day07a(crabs []int) int {
	sort.Ints(crabs)
	// fmt.Println(crabs)

	// first guess is the pos @ half the length
	ix := len(crabs) / 2
	pos := crabs[ix]

	fmt.Println(getError(crabs, pos))

	return pos
}

func getError(crabs []int, pos int) int {
	// return sum(abs(x - pos) for x in crabs)
	total := 0
	for _, x := range crabs {
		if x > pos {
			total += x - pos
		} else {
			total += pos - x
		}
	}
	return total
}

func day07b(crabs []int) int {
	sort.Ints(crabs)
	// fmt.Println(crabs)

	// first guess is the pos @ half the length
	ix := len(crabs) / 2
	pos := crabs[ix]

	for {
		fuel := getError2(crabs, pos)
		fmt.Printf("trying %d: fuel=%d\n", pos, fuel)

		fuel1 := getError2(crabs, pos+1)
		fuel2 := getError2(crabs, pos-1)

		if fuel1 > fuel && fuel2 > fuel {
			// we've found the optimum
			fmt.Printf("found %d: fuel=%d\n", pos, fuel)
			return fuel
		} else if fuel1 < fuel {
			// increase pos
			pos++
		} else {
			// decrease pos
			pos--
		}
	}
}

func getError2(crabs []int, pos int) int {
	// return sum(error(x) for x in crabs)
	// where error(x) =
	//   let delta = abs(pos - x)
	//   delta * (delta + 1) / 2
	var total int
	var delta int
	for _, x := range crabs {
		if x > pos {
			delta = x - pos
		} else {
			delta = pos - x
		}
		total += delta * (delta + 1) / 2
	}
	return total
}

func main() {
	f, err := os.Open("./lqinput.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Trim(string(b), " \n\t")
	nums := strings.Split(s, ",")
	var crabs []int
	for _, n := range nums {
		v, _ := strconv.Atoi(n)
		crabs = append(crabs, v)
	}
	fmt.Println(day07a(crabs))
	fmt.Println(day07b(crabs))
}
