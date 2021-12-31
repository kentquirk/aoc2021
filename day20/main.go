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

type Coord struct {
	R int
	C int
}

type Image map[Coord]int

func (img Image) Extents() (Coord, Coord) {
	var lo, hi Coord
	for c := range img {
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
		if x := img[Coord{R: c.R + row, C: c.C + col}]; x == 1 {
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

// we're going to do a sparse image -- store only the white pixels in a map --
// that way, looking up a nonexistent cell will return a zero (black) pixel
// this also gives us the ability to get the number of lit pixels with a len() call.
func parseImage(s string) Image {
	img := make(Image)
	lines := strings.Split(s, "\n")
	for row, l := range lines {
		for col, c := range l {
			if c == '#' {
				img[Coord{R: row, C: col}] = 1
			}
		}
	}
	return img
}

func (img Image) Enhance(algo *Algorithm) Image {
	result := make(Image)
	lo, hi := img.Extents()
	for r := lo.R - 1; r <= hi.R+1; r++ {
		for c := lo.C - 1; c <= hi.C+1; c++ {
			co := Coord{R: r, C: c}
			bitIndex := img.EnhancePixel(co)
			if algo.GetBit(bitIndex) {
				result[co] = 1
			}
		}
	}
	return result
}

func (img Image) Print() {
	lo, hi := img.Extents()
	for r := lo.R; r <= hi.R; r++ {
		for c := lo.C; c <= hi.C; c++ {
			co := Coord{R: r, C: c}
			if img[co] == 1 {
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
	f, err := os.Open("./inputsample.txt")
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

	img.Print()
	lo, hi := img.Extents()
	fmt.Println(len(img), lo, hi)

	img2 := img.Enhance(algo)
	img2.Print()
	lo, hi = img2.Extents()
	fmt.Println(len(img2), lo, hi)

	img3 := img2.Enhance(algo)
	img3.Print()
	lo, hi = img3.Extents()
	fmt.Println(len(img3), lo, hi)
}
