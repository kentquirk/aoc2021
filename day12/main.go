package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
)

// returns whether node should be considered a 'big' node
// (can be visited multiple times)
func isBig(node string) bool {
	return unicode.IsUpper(rune(node[0]))
}

// add the edge lhs->rhs to the map (but don't add any edges to 'start'
// and don't add any edges from 'end')
func addToMap(cavemap map[string][]string, lhs string, rhs string) {
	if lhs != "end" && rhs != "start" {
		cavemap[lhs] = append(cavemap[lhs], rhs)
	}
}

// parse input lines into an adjacency list representation of the graph.
// since the graph is non-directional, each edge in the input is added twice:
// once in each direction.
func parse(lines []string) map[string][]string {
	cavemap := make(map[string][]string)
	for _, l := range lines {
		splits := strings.Split(strings.Trim(l, "\n \t"), "-")
		// fmt.Println(splits)
		lhs := splits[0]
		rhs := splits[1]

		addToMap(cavemap, lhs, rhs)
		addToMap(cavemap, rhs, lhs)
	}
	return cavemap
}

func isAlreadyVisited(visited []string, node string) bool {
	for _, visitedNode := range visited {
		if visitedNode == node && !isBig(visitedNode) {
			return true
		}
	}
	return false
}

// explore from node, avoiding small rooms we've seen before.
// (explores recursively, depth first)
// returns number of found paths from node to 'end'.
// 'cavemap' is the adjacency-list graph repr
// 'visited' is the list of nodes we've visited in this path so far
func traverse(cavemap map[string][]string, visited []string, canRevisitSmallCave bool, self string) int {
	neighbors := cavemap[self]
	pathCount := 0
	newVisited := append(visited, self)

	for _, n := range neighbors {
		newCanRevisitSmallCave := canRevisitSmallCave
		if n == "end" {
			pathCount++
		} else {
			// new logic:
			// if n is small, and already visited, and can revisit, then we flip can revisit to false
			// and visit it anyway.
			// otherwise we do the below logic

			// test if n already in visited (only applies if n is not already big)
			if !isBig(n) {
				if isAlreadyVisited(visited, n) {
					if !canRevisitSmallCave {
						continue
					}
					// let's take our one-time option to revisit this cave. in doing so,
					// we will use it up!
					newCanRevisitSmallCave = false
				}
			}

			// we have not found a visited node that incurs a skip
			pathCount += traverse(cavemap, newVisited, newCanRevisitSmallCave, n)
		}
	}
	return pathCount
}

func day12a(lines []string) int {
	cavemap := parse(lines)
	fmt.Println(cavemap)

	return traverse(cavemap, []string{}, false, "start")
}

func day12b(lines []string) int {
	cavemap := parse(lines)
	fmt.Println(cavemap)

	return traverse(cavemap, []string{}, true, "start")
}

func main() {
	f, err := os.Open("./inputlq.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	fmt.Println(day12a(lines))
	fmt.Println(day12b(lines))
}
