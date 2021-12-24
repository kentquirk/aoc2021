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

type Value struct {
	Num  int
	Pair *Pair
}

func PairValue(l int, r int) *Value {
	return &Value{Pair: &Pair{Left: &Value{Num: l}, Right: &Value{Num: r}}}
}

func (v *Value) String() string {
	if v.IsRegular() {
		return strconv.Itoa(v.Num)
	}
	return v.Pair.String()
}

func (v *Value) PPrint(depth int) string {
	if v.IsRegular() {
		if v.Num >= 10 {
			return fmt.Sprintf("\x1b[31m%d\x1b[0m", v.Num)
		}
		return strconv.Itoa(v.Num)
	}
	return v.Pair.PPrint(depth)
}

func (v *Value) IsRegular() bool {
	return v.Pair == nil
}

func (v *Value) IsPair() bool {
	return v.Pair != nil
}

type Pair struct {
	Left  *Value
	Right *Value
}

func (p *Pair) Set(v *Value, left bool) {
	if left {
		p.Left = v
	} else {
		p.Right = v
	}
}

func (p *Pair) IsRegularPair() bool {
	return p.Left.IsRegular() && p.Right.IsRegular()
}

func (p *Pair) String() string {
	return fmt.Sprintf("[%s,%s]", p.Left.String(), p.Right.String())
}

func (p *Pair) PPrint(depth int) string {
	colors := map[int]int{
		1: 34,
		2: 36,
		3: 35,
		4: 32,
		5: 33,
		6: 31,
	}

	depth++
	return fmt.Sprintf("\x1b[%dm[\x1b[0m%s,%s\x1b[%dm]\x1b[0m",
		colors[depth], p.Left.PPrint(depth), p.Right.PPrint(depth), colors[depth])
}

func (p *Pair) Magnitude() int {
	var lt, rt int
	if p.Left.IsPair() {
		lt = p.Left.Pair.Magnitude()
	} else {
		lt = p.Left.Num
	}
	if p.Right.IsPair() {
		rt = p.Right.Pair.Magnitude()
	} else {
		rt = p.Right.Num
	}
	return 3*lt + 2*rt
}

func (p *Pair) Add(rhs *Pair) *Pair {
	sum := &Pair{Left: &Value{Pair: p}, Right: &Value{Pair: rhs}}
	return sum.Reduce()
}

func (p *Pair) Reduce() *Pair {
	// fmt.Println("reducing      ", p.PPrint(0))
	var didSomething bool
	for {
		p, didSomething = p.explode()
		// if didSomething {
		// 	fmt.Println("  exploded to ", p.PPrint(0))
		// }
		if !didSomething {
			p, didSomething = p.split()
			// if didSomething {
			// 	fmt.Println("     split to ", p.PPrint(0))
			// }
		}
		if !didSomething {
			return p
		}
	}
}

func addTo(s string, left int, right int, add int) string {
	v, _ := strconv.Atoi(s[left:right])
	vr := strconv.Itoa(v + add)
	return s[:left] + vr + s[right:]
}

func addToPrev(s string, add int) string {
	valuePat := regexp.MustCompile(`[0-9]+`)
	nums := valuePat.FindAllStringIndex(s, -1)
	if len(nums) > 0 {
		n := len(nums) - 1
		return addTo(s, nums[n][0], nums[n][1], add)
	}
	return s
}

func addToNext(s string, add int) string {
	valuePat := regexp.MustCompile(`[0-9]+`)
	nums := valuePat.FindAllStringIndex(s, -1)
	if len(nums) > 0 {
		return addTo(s, nums[0][0], nums[0][1], add)
	}
	return s
}

// text-based explode
func (p *Pair) explode() (*Pair, bool) {
	// find an explodeable pair
	regPairPat := regexp.MustCompile(`\[([0-9]+),([0-9]+)\]`)
	s := p.String()
	pi := regPairPat.FindAllStringSubmatchIndex(s, -1)
	for _, ixs := range pi {
		// count the depth up to but not including the pair we found
		leftPart := s[:ixs[0]]
		ltn := strings.Count(leftPart, "[")
		rtn := strings.Count(leftPart, "]")
		if ltn-rtn < 4 {
			// it's not deep enough, so explore farther in the string
			continue
		}
		// find the values in the pair
		ltv, _ := strconv.Atoi(s[ixs[2]:ixs[3]])
		rtv, _ := strconv.Atoi(s[ixs[4]:ixs[5]])
		leftPart = addToPrev(leftPart, ltv)
		rightPart := addToNext(s[ixs[1]:], rtv)
		p2 := Parse(leftPart + "0" + rightPart)
		return p2, true
	}
	return p, false
}

func (p *Pair) split() (*Pair, bool) {
	valuePat := regexp.MustCompile(`[0-9]{2}`)
	s := p.String()
	nums := valuePat.FindAllStringIndex(s, -1)
	if len(nums) > 0 {
		left := nums[0][0]
		right := nums[0][1]
		v, _ := strconv.Atoi(s[left:right])
		newS := fmt.Sprintf("%s[%d,%d]%s", s[:left], v/2, int(float64(v)/2+.6), s[right:])
		return Parse(newS), true
	}
	return p, false
}

// Parses a string to read a complete pair from it, expecting
// that the first character is '['; terminates after reading
// the corresponding closing ']', returning the number of characters
// processed
func Parse(s string) *Pair {
	p, _ := parse(s)
	return p
}

func parse(s string) (*Pair, int) {
	pair := &Pair{}
	if s[0] != '[' {
		panic("oops")
	}
	isLeft := true
	for ix := 1; ix < len(s); ix++ {
		ch := s[ix]
		switch {
		case ch == '[':
			p, n := parse(s[ix:])
			pair.Set(&Value{Pair: p}, isLeft)
			ix += n
		case strings.IndexByte("0123456789", ch) != -1:
			n := int(ch - '0')
			// because explode might create 2-digit numbers, we need to handle this case
			if strings.IndexByte("0123456789", s[ix+1]) != -1 {
				n = n*10 + int(s[ix+1]-'0')
				ix++
			}
			pair.Set(&Value{Num: n}, isLeft)
		case ch == ',':
			isLeft = false
		case ch == ']':
			return pair, ix
		default:
			panic("oopsie")
		}
	}
	panic("whoops")
}

func SumList(lines []string) *Pair {
	lhs := Parse(lines[0])
	for _, l := range lines[1:] {
		// fmt.Println(l)
		rhs := Parse(l)
		lhs = lhs.Add(rhs)
		// fmt.Printf("sum: %s\n", lhs)
	}
	return lhs
}

func day18a(lines []string) int {
	p := SumList(lines)
	return p.Magnitude()
}

func day18b(lines []string) int {
	var pairs []*Pair
	for _, l := range lines {
		pairs = append(pairs, Parse(l))
	}
	largest := 0
	for i := 0; i < len(pairs)-1; i++ {
		for j := i; j < len(pairs)-1; j++ {
			t := pairs[i].Add(pairs[j])
			mag := t.Magnitude()
			if mag > largest {
				largest = mag
			}
			t = pairs[j].Add(pairs[i])
			mag = t.Magnitude()
			if mag > largest {
				largest = mag
			}
		}
	}
	return largest
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
	fmt.Println(day18a(lines))
	fmt.Println(day18b(lines))
}
