package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Octopus struct {
	Energy    int
	Flashed   bool
	Neighbors []*Octopus
}

func NewOctopus(state byte) *Octopus {
	return &Octopus{
		Energy: int(state),
	}
}

func (o *Octopus) Step() {
	o.Flashed = false
	o.Energy++
}

func (o *Octopus) Nudge() {
	if !o.Flashed {
		o.Energy++
	}
}

func (o *Octopus) MaybeFlash() int {
	if o.Energy > 9 {
		o.Flashed = true
		o.Energy = 0
		for _, oct := range o.Neighbors {
			oct.Nudge()
		}
		return 1
	}
	return 0
}

type OctopusGarden struct {
	Octopuses [][]*Octopus
	Score     int
}

func NewOctopusGarden(octomap []string) *OctopusGarden {
	var octopuses [][]*Octopus
	for _, row := range octomap {
		oct := make([]*Octopus, len(row))
		for i := range row {
			oct[i] = NewOctopus(row[i] - byte('0'))
		}
		octopuses = append(octopuses, oct)
	}
	// add the appropriate number of neighbors
	for row, octs := range octopuses {
		for col, oct := range octs {
			for nr := -1; nr <= 1; nr++ {
				if row+nr < 0 || row+nr > len(octopuses)-1 {
					continue
				}
				for nc := -1; nc <= 1; nc++ {
					if col+nc < 0 || col+nc > len(octopuses[row])-1 {
						continue
					}
					if nc == 0 && nr == 0 {
						continue
					}
					oct.Neighbors = append(oct.Neighbors, octopuses[row+nr][col+nc])
				}
			}
		}
	}
	return &OctopusGarden{Octopuses: octopuses}
}

func (g *OctopusGarden) Step() {
	for _, octs := range g.Octopuses {
		for _, oct := range octs {
			oct.Step()
		}
	}
}

func (g *OctopusGarden) MaybeFlash() int {
	score := 0
	for _, octs := range g.Octopuses {
		for _, oct := range octs {
			score += oct.MaybeFlash()
		}
	}
	return score
}

func (g *OctopusGarden) Print() {
	for _, octs := range g.Octopuses {
		for _, oct := range octs {
			fmt.Printf("%d", oct.Energy)
		}
		fmt.Println()
	}
}

func day11a(octomap []string) int {
	const nSteps = 100
	var score int
	octopuses := NewOctopusGarden(octomap)

	// octopuses.Print()

	for i := 1; i <= nSteps; i++ {
		octopuses.Step()
		n := 1
		for n != 0 {
			n = octopuses.MaybeFlash()
			score += n
			// fmt.Println(n)
		}
		// fmt.Printf("step %d, score %d\n", i, score)
		// octopuses.Print()
		// fmt.Println()
	}
	return score
}

func day11b(octomap []string) int {
	octopuses := NewOctopusGarden(octomap)

	// octopuses.Print()

	for i := 1; ; i++ {
		octopuses.Step()
		score := 0
		n := 1
		for n != 0 {
			n = octopuses.MaybeFlash()
			score += n
			// fmt.Println(n)
			if score == 100 {
				return i
			}
		}
		// if i > 190 {
		// 	fmt.Printf("step %d, n %d\n", i, n)
		// 	octopuses.Print()
		// 	fmt.Println()
		// }
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
	fmt.Println(day11a(lines))
	fmt.Println(day11b(lines))
}
