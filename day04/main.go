package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Square struct {
	Value  int
	Marked bool
}

type Board struct {
	Squares []Square
	Win     bool
}

func (b Board) String() string {
	s := strings.Builder{}
	s.WriteRune('\n')
	for i := range b.Squares {
		selected := " "
		if b.Squares[i].Marked {
			selected = "*"
		}
		s.WriteString(fmt.Sprintf("%2d%s ", b.Squares[i].Value, selected))
		if i%5 == 4 {
			s.WriteRune('\n')
		}
	}
	return s.String()
}

func (b *Board) MarkMatching(value int) bool {
	if b.Win {
		return false
	}
	for i, s := range b.Squares {
		if s.Value == value {
			b.Squares[i].Marked = true
			return true
		}
	}
	return false
}

var winConditions = [][]int{
	{0, 1, 2, 3, 4},
	{5, 6, 7, 8, 9},
	{10, 11, 12, 13, 14},
	{15, 16, 17, 18, 19},
	{20, 21, 22, 23, 24},
	{0, 5, 10, 15, 20},
	{1, 6, 11, 16, 21},
	{2, 7, 12, 17, 22},
	{3, 8, 13, 18, 23},
	{4, 9, 14, 19, 24},
}

func (b *Board) CheckForWin() bool {
	if b.Win {
		return true
	}
	for _, wc := range winConditions {
		found := true
		for _, ix := range wc {
			if !b.Squares[ix].Marked {
				found = false
				break
			}
		}
		if found {
			b.Win = true
			return true
		}
	}
	return false
}

func (b Board) Score(lastNumber int) int {
	var sum int
	for i := range b.Squares {
		if !b.Squares[i].Marked {
			sum += b.Squares[i].Value
		}
	}
	return sum * lastNumber
}

func parse(input string) []int {
	splitpat := regexp.MustCompile("[, \n\t]+")
	ss := splitpat.Split(input, -1)

	var result []int
	for _, s := range ss {
		if s == "" {
			continue
		}
		x, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		result = append(result, x)
	}
	return result
}

func NewBoard(input string) *Board {
	board := new(Board)
	values := parse(input)
	for _, v := range values {
		board.Squares = append(board.Squares, Square{Value: v})
	}
	if len(board.Squares) != 25 {
		panic("wrong length")
	}
	return board
}

func day04a(draws []int, boards []*Board) int {
	for _, draw := range draws {
		for _, b := range boards {
			b.MarkMatching(draw)
			if b.CheckForWin() {
				return b.Score(draw)
			}
		}
	}
	return -1
}

func day04b(draws []int, boards []*Board) int {
	for _, draw := range draws {
		nonwins := 0
		score := 0
		for _, b := range boards {
			if !b.Win {
				b.MarkMatching(draw)
				if b.CheckForWin() {
					score = b.Score(draw)
				} else {
					nonwins++
				}
			}
		}
		if nonwins == 0 {
			return score
		}
	}
	return -1
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
	blocks := strings.Split(string(b), "\n\n")
	draws := parse(blocks[0])

	var boards []*Board
	for _, block := range blocks[1:] {
		board := NewBoard(block)
		boards = append(boards, board)
	}

	// for _, b := range boards {
	// 	fmt.Println(b)
	// }

	// fmt.Println("A: ", day04a(draws, boards))
	fmt.Println("B: ", day04b(draws, boards))

}
