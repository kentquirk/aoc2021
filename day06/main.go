package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Trim(string(b), " \n\t")
	nums := strings.Split(s, ",")
	var states []int
	for _, n := range nums {
		states = append(states, int(n[0]-'0'))
	}
	ocean := make(map[int]int)
	for _, s := range states {
		ocean[s]++
	}

	numDays := 256
	for i := 1; i <= numDays; i++ {
		for days := 8; days >= 0; days-- {
			count := ocean[days]
			if count == 0 {
				continue
			}
			if days == 0 {
				ocean[7] += count
				ocean[9] = count
			}
		}
		for days := 0; days <= 9; days++ {
			ocean[days] = ocean[days+1]
		}
		// fmt.Println(ocean)
		// total := 0
		// for days := 0; days < 9; days++ {
		// 	total += ocean[days]
		// }
		// fmt.Printf("Total after %d days: %d\n", i, total)
	}

	total := 0
	for days := 0; days < 9; days++ {
		total += ocean[days]
	}
	fmt.Printf("Total after %d days: %d\n", numDays, total)
}
