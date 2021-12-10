package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func day09a(heightmap []string) int {
	risk := 0
	var lowpoints [][2]int
	for i, row := range heightmap {
		for j := 0; j < len(row); j++ {
			var candidates []byte
			if i > 0 {
				candidates = append(candidates, (heightmap[i-1][j]))
			}
			if i < len(heightmap)-1 {
				candidates = append(candidates, (heightmap[i+1][j]))
			}
			if j > 0 {
				candidates = append(candidates, row[j-1])
			}
			if j < len(row)-1 {
				candidates = append(candidates, row[j+1])
			}
			lowpoint := true
			for _, c := range candidates {
				if row[j] >= c {
					lowpoint = false
					break
				}
			}
			if lowpoint {
				lowpoints = append(lowpoints, [2]int{i, j})
				risk += int(row[j]-byte('0')) + 1
			}
		}
	}
	// fmt.Println(lowpoints)
	return risk
}

func floodfill(floor [][]byte, r int, c int, index byte) int {
	if r < 0 || r >= len(floor) || c < 0 || c >= len(floor[r]) {
		return 0
	}
	v := floor[r][c]
	if v >= 9 {
		return 0
	}
	floor[r][c] = index
	size := 1
	size += floodfill(floor, r, c-1, index)
	size += floodfill(floor, r, c+1, index)
	size += floodfill(floor, r-1, c, index)
	size += floodfill(floor, r+1, c, index)

	return size
}

// the second half is just a floodfill problem, so we'll do an inefficient recursive floodfill
func day09b(heightmap []string) int {
	// convert to a byte array
	var floor [][]byte
	for _, row := range heightmap {
		b := make([]byte, len(row))
		for i := range row {
			b[i] = row[i] - byte('0')
		}
		floor = append(floor, b)
	}

	var basins []int
	var index byte = 10
	for r := range floor {
		for c := range floor[r] {
			size := floodfill(floor, r, c, index)
			if size != 0 {
				basins = append(basins, size)
				index++
			}
		}
	}
	sort.Ints(basins)
	return basins[len(basins)-1] * basins[len(basins)-2] * basins[len(basins)-3]
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
	fmt.Println(day09a(lines))
	fmt.Println(day09b(lines))
}
