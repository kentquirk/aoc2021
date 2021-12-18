package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
)

type TypeID uint

const (
	Sum TypeID = iota
	Product
	Minimum
	Maximum
	Literal
	GreaterThan
	LessThan
	EqualTo
)

func (t TypeID) String() string {
	switch t {
	case Sum:
		return "Sum"
	case Product:
		return "Product"
	case Minimum:
		return "Minimum"
	case Maximum:
		return "Maximum"
	case Literal:
		return "Literal"
	case GreaterThan:
		return "GreaterThan"
	case LessThan:
		return "LessThan"
	case EqualTo:
		return "EqualTo"
	}
	return "BUG"
}

type Stream struct {
	Data   big.Int
	Len    uint
	Cursor uint // 0 is at the left
}

func NewStream(s string) *Stream {
	stream := &Stream{}
	_, ok := stream.Data.SetString(s, 16)
	if !ok {
		panic("oops")
	}
	stream.Len = uint(len(s) * 4)

	// fmt.Println(stream.Data.Text(16))
	return stream
}

func (s *Stream) Range(pos uint, n uint) uint {
	var result uint
	for i := uint(0); i < n; i++ {
		ix := s.Len - pos - i - 1
		bit := s.Data.Bit(int(ix))
		// fmt.Println(bit, ix)
		result = (result << 1) | (bit & 1)
	}
	return result
}

// Read returns n bits starting at the cursor, increments the cursor by n
func (s *Stream) Read(n uint) uint {
	u := s.Range(s.Cursor, n)
	s.Cursor += n
	return u
}

// consumes a single complete packet (which might be nested) from the stream
// and returns the number of bits consumed
func (s *Stream) ReadPacket() (*Packet, uint) {
	savedCur := s.Cursor
	p := &Packet{
		Version: s.Read(3),
		Type:    TypeID(s.Read(3)),
	}
	if p.Type == Literal {
		var continueMask uint = 0b10000
		g := continueMask
		for g&continueMask != 0 {
			g = s.Read(5)
			p.LiteralValue = (p.LiteralValue << 4) | (g & 0b1111)
		}
	} else {
		lengthTypeID := s.Read(1)
		if lengthTypeID == 0 {
			length := s.Read(15)
			for length > 0 {
				sub, n := s.ReadPacket()
				p.Subpackets = append(p.Subpackets, sub)
				length -= n
			}
		} else {
			numSubPackets := s.Read(11)
			for i := uint(0); i < numSubPackets; i++ {
				sub, _ := s.ReadPacket()
				p.Subpackets = append(p.Subpackets, sub)
			}
		}
	}
	return p, s.Cursor - savedCur
}

type Packet struct {
	Version      uint
	Type         TypeID
	LiteralValue uint
	Subpackets   []*Packet
}

func (p *Packet) SumVersions() uint {
	sum := p.Version
	for _, sub := range p.Subpackets {
		sum += sub.SumVersions()
	}
	return sum
}

func (p *Packet) Evaluate() uint {
	var value uint
	switch p.Type {
	case Sum:
		for _, sub := range p.Subpackets {
			value += sub.Evaluate()
		}
	case Product:
		value = 1
		for _, sub := range p.Subpackets {
			value *= sub.Evaluate()
		}
	case Minimum:
		value = math.MaxUint
		for _, sub := range p.Subpackets {
			v := sub.Evaluate()
			if v < value {
				value = v
			}
		}
	case Maximum:
		for _, sub := range p.Subpackets {
			v := sub.Evaluate()
			if v > value {
				value = v
			}
		}
	case Literal:
		value = p.LiteralValue
	case GreaterThan:
		v1 := p.Subpackets[0].Evaluate()
		v2 := p.Subpackets[1].Evaluate()
		if v1 > v2 {
			value = 1
		}
	case LessThan:
		v1 := p.Subpackets[0].Evaluate()
		v2 := p.Subpackets[1].Evaluate()
		if v1 < v2 {
			value = 1
		}
	case EqualTo:
		v1 := p.Subpackets[0].Evaluate()
		v2 := p.Subpackets[1].Evaluate()
		if v1 == v2 {
			value = 1
		}
	}
	return value
}

func (p *Packet) Print(n int) {
	fmt.Printf("%sv%d %s (%d)\n", strings.Repeat("  ", n), p.Version, p.Type, p.Evaluate())
	for _, sub := range p.Subpackets {
		sub.Print(n + 1)
	}
}

func day16a(line string) uint {
	s := NewStream(line)
	fmt.Println(" ----\n", line[:4])
	p, _ := s.ReadPacket()
	p.Print(0)
	r := p.SumVersions()
	fmt.Printf("-> %d\n", r)
	return r
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
	for _, line := range lines {
		day16a(line)
	}
}
