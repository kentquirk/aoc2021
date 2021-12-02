package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func day01a(data []int) int {
	count := 0
	for i := 0; i < len(data)-1; i++ {
		diff := data[i+1] - data[i]
		if diff > 0 {
			count++
		}
	}
	return count
}

func day01b(data []int) int {
	count := 0
	lastsum := data[0] + data[1] + data[2]
	for i := 1; i < len(data)-2; i++ {
		sum := data[i] + data[i+1] + data[i+2]
		if sum > lastsum {
			count++
		}
		lastsum = sum
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
	ints := make([]int, len(lines))
	for i := range lines {
		ints[i], _ = strconv.Atoi(lines[i])
	}
	fmt.Println(day01a(ints))
	fmt.Println(day01b(ints))
}
