package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func parse(data []string) ([]int64, int) {
	values := make([]int64, 0)
	nbits := len(data[0])
	for _, d := range data {
		v, err := strconv.ParseInt(d, 2, nbits+1)
		if err != nil {
			panic(err)
		}
		values = append(values, v)
	}
	return values, nbits
}

func day03a(data []string) (int64, int64) {
	values, nbits := parse(data)

	var gamma, epsilon int64
	for b := 0; b < nbits; b++ {
		var mask int64 = 1 << b
		numZeroes, numOnes := count(values, mask)
		if numOnes > numZeroes {
			gamma |= mask
		} else {
			epsilon |= mask
		}
	}

	return gamma, epsilon
}

func count(values []int64, mask int64) (int, int) {
	numOnes := 0
	for _, v := range values {
		if v&mask != 0 {
			numOnes++
		}
	}
	numZeros := len(values) - numOnes
	return numZeros, numOnes
}

func filter(values []int64, predicate func(int64) bool) []int64 {
	result := make([]int64, 0)
	for _, v := range values {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func day03b(data []string, mostCommon bool) int64 {
	values, nbits := parse(data)

	for b := 0; b < nbits; b++ {
		var mask int64 = 1 << (nbits - b - 1)
		numZeroes, numOnes := count(values, mask)
		switch {
		case numOnes > numZeroes && mostCommon,
			numOnes == numZeroes && mostCommon,
			numOnes < numZeroes && !mostCommon:
			values = filter(values, func(v int64) bool {
				return (v & mask) == mask
			})
		case numOnes < numZeroes && mostCommon,
			numOnes == numZeroes && !mostCommon,
			numOnes > numZeroes && !mostCommon:
			values = filter(values, func(v int64) bool {
				return (v & mask) == 0
			})
		}
		if len(values) == 1 {
			return values[0]
		}
		if len(values) == 0 {
			panic("none left")
		}
	}
	fmt.Println("oops", values)
	return 0
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
	gamma, epsilon := day03a(lines)
	fmt.Printf("gamma: %d, epsilon: %d, result: %d\n",
		gamma, epsilon, gamma*epsilon)

	oxygen := day03b(lines, true)
	co2 := day03b(lines, false)
	fmt.Printf("oxygen: %d, co2: %d, result: %d\n",
		oxygen, co2, oxygen*co2)
}
