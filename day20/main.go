package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"
)

type Algorithm big.Int

func parseAlgorithm(s string) *Algorithm {
	r := regexp.MustCompile("[^.#]+")
	s = r.ReplaceAllString(s, "")
	s = strings.Replace(s, ".", "0", -1)
	s = strings.Replace(s, "#", "1", -1)
	if len(s) != 512 {
		fmt.Println(len(s))
		panic("oops")
	}
	i := new(big.Int)
	i.SetString(s, 2)
	return (*Algorithm)(i)
}

func (a *Algorithm) GetBit(index int) bool {
	return ((*big.Int)(a)).Bit(511-index) == 1
}

func (a *Algorithm) NeedsInvert() bool {
	return a.GetBit(0) && !a.GetBit(511)
}

type Coord struct {
	R int
	C int
}

type Image struct {
	m        map[Coord]bool
	inverted bool
}

func MakeImage(inverted bool) Image {
	return Image{
		m:        make(map[Coord]bool),
		inverted: inverted,
	}
}

func (img Image) IsInverted() bool {
	return img.inverted
}

func (img Image) Get(c Coord) bool {
	p := img.m[c]
	return p != img.inverted
}

func (img Image) Set(c Coord, v bool) {
	if img.inverted != v {
		img.m[c] = true
	}
}

func (img Image) Extents() (Coord, Coord) {
	var lo, hi Coord
	for c := range img.m {
		if c.C < lo.C {
			lo.C = c.C
		}
		if c.R < lo.R {
			lo.R = c.R
		}
		if c.C > hi.C {
			hi.C = c.C
		}
		if c.R > hi.R {
			hi.R = c.R
		}
	}
	return lo, hi
}

func (img Image) EnhancePixel(c Coord) int {
	mask := 256
	row := -1
	col := -1
	bits := 0
	for mask != 0 {
		if img.Get(Coord{R: c.R + row, C: c.C + col}) {
			bits |= mask
		}
		col++
		if col > 1 {
			col = -1
			row++
		}
		mask >>= 1
	}
	return bits
}

func (img Image) Count() int {
	return len(img.m)
}

// we're going to do a sparse image -- store only the white pixels in a map --
// that way, looking up a nonexistent cell will return a zero (black) pixel
// this also gives us the ability to get the number of lit pixels with a len() call.
func parseImage(s string) Image {
	img := MakeImage(false)
	lines := strings.Split(s, "\n")
	for row, l := range lines {
		for col, c := range l {
			if c == '#' {
				img.Set(Coord{R: row, C: col}, c == '#')
			}
		}
	}
	return img
}

func (img Image) Enhance(algo *Algorithm) Image {
	shouldInvert := algo.NeedsInvert() && !img.IsInverted()
	result := MakeImage(shouldInvert)
	lo, hi := img.Extents()
	for r := lo.R - 1; r <= hi.R+1; r++ {
		for c := lo.C - 1; c <= hi.C+1; c++ {
			co := Coord{R: r, C: c}
			bitIndex := img.EnhancePixel(co)
			result.Set(co, algo.GetBit(bitIndex))
		}
	}
	return result
}

func (img Image) Print() {
	lo, hi := img.Extents()
	for r := lo.R; r <= hi.R; r++ {
		for c := lo.C; c <= hi.C; c++ {
			co := Coord{R: r, C: c}
			if img.Get(co) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
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
	parts := strings.Split(string(b), "\n\n")
	algo := parseAlgorithm(parts[0])
	// algo := parseAlgorithm(strings.Repeat(".#", 256))

	// fmt.Printf("%x\n", (*big.Int)(algo))
	// fmt.Println(algo.GetBit(0), algo.GetBit(1), algo.GetBit(510), algo.GetBit(511))

	img := parseImage(parts[1])
	nTimes := 50

	// img.Print()
	for i := 0; i < nTimes; i++ {
		img = img.Enhance(algo)
	}
	lo, hi := img.Extents()
	fmt.Printf("After %d enhancements, there were %d pixels lit with extents of (%v, %v)\n", nTimes, img.Count(), lo, hi)
}
