package main

import (
	"regexp"
	"testing"
)

func TestExplode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		wantBoom bool
	}{
		{"a", "[[[[[9,8],1],2],3],4]", "[[[[0,9],2],3],4]", true},
		{"b", "[7,[6,[5,[4,[3,2]]]]]", "[7,[6,[5,[7,0]]]]", true},
		{"c", "[[6,[5,[4,[3,2]]]],1]", "[[6,[5,[7,0]]],3]", true},
		{"d", "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]", true},
		{"e", "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[7,0]]]]", true},
		{"f", "[3,[2,[8,0]]]", "[3,[2,[8,0]]]", false},
		{"g", "[7,[6,[5,[4,[8,8]]]]]", "[7,[6,[5,[12,0]]]]", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parse(tt.input)
			got, exploded := p.explode()
			s := got.String()
			if exploded != tt.wantBoom {
				t.Errorf("Got exploded=%v, want %v", exploded, tt.wantBoom)
			}
			if exploded && s != tt.want {
				t.Errorf("Pair.Reduce() = %v, want %v", s, tt.want)
			}
		})
	}
}

func Test_addToPrev(t *testing.T) {
	type args struct {
		s   string
		add int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"a", args{"[[6,[5,[4,", 3}, "[[6,[5,[7,"},
		{"b", args{"[[[[,", 3}, "[[[[,"},
		{"c", args{"[[6,5],[4,", 3}, "[[6,5],[7,"},
		{"d", args{"[[6,5],", 3}, "[[6,8],"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addToPrev(tt.args.s, tt.args.add); got != tt.want {
				t.Errorf("addToLast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPair_split(t *testing.T) {
	tests := []struct {
		name  string
		p     string
		want  string
		want1 bool
	}{
		{"a", "[[[[0,7],4],[15,[0,13]]],[1,1]]", "[[[[0,7],4],[[7,8],[0,13]]],[1,1]]", true},
		{"b", "[[[[0,7],4],[[7,8],[0,13]]],[1,1]]", "[[[[0,7],4],[[7,8],[0,[6,7]]]],[1,1]]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parse(tt.p)
			got, got1 := p.split()
			if got.String() != tt.want {
				t.Errorf("Pair.split() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Pair.split() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPair_Reduce(t *testing.T) {
	tests := []struct {
		name string
		p    string
		want string
	}{
		{"a", "[[[[0,7],4],[15,[0,13]]],[1,1]]", "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]"},
		{"b", "[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]", "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parse(tt.p)
			got := p.Reduce()
			if got.String() != tt.want {
				t.Errorf("Pair.Reduce() \n got = '%v', \nwant = '%v'", got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		lhs  string
		rhs  string
		want string
	}{
		{"a", "[[[[4,3],4],4],[7,[[8,4],9]]]", "[1,1]", "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]"},
		{"b", "[[[[7,0],[7,7]],[[7,7],[7,8]]],[[[7,7],[8,8]],[[7,7],[8,7]]]]",
			"[7,[5,[[3,8],[1,4]]]]",
			"[[[[7,7],[7,8]],[[9,5],[8,7]]],[[[6,8],[0,8]],[[9,9],[9,0]]]]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lhs := Parse(tt.lhs)
			rhs := Parse(tt.rhs)
			got := lhs.Add(rhs)

			if got.String() != tt.want {
				t.Errorf("Pair.Reduce() \n got = '%v', \nwant = '%v'", got, tt.want)
			}
		})
	}
}

func splitList(s string) []string {
	re := regexp.MustCompile("[\n\t\r ]+")
	return re.Split(s, -1)
}

func TestSumList(t *testing.T) {
	tests := []struct {
		name string
		list string
		want string
	}{
		{"a", `[1,1]
		[2,2]
		[3,3]
		[4,4]`, "[[[[1,1],[2,2]],[3,3]],[4,4]]"},
		{"b", `[1,1]
		[2,2]
		[3,3]
		[4,4]
		[5,5]`, "[[[[3,0],[5,3]],[4,4]],[5,5]]"},
		{"c", `[1,1]
		[2,2]
		[3,3]
		[4,4]
		[5,5]
		[6,6]`, "[[[[5,0],[7,4]],[5,5]],[6,6]]"},
		{"d", `[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]
		[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]
		[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]
		[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]
		[7,[5,[[3,8],[1,4]]]]
		[[2,[2,2]],[8,[8,1]]]
		[2,9]
		[1,[[[9,3],9],[[9,0],[0,7]]]]
		[[[5,[7,4]],7],1]
		[[[[4,2],2],6],[8,7]]`, "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := splitList(tt.list)
			got := SumList(list)

			if got.String() != tt.want {
				t.Errorf("Pair.Reduce() \n got = '%v', \nwant = '%v'", got, tt.want)
			}
		})
	}
}
