package standard_parser

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	res := IndexAllSubstrings("ababac", "aba")
	for _, pos := range res {
		fmt.Printf("%d: %v\n", pos[0], pos[1])
	}
}
func TestGenerator(t *testing.T) {
	nodes := Parse("aaa{{int(a)}}aa", NewTagDefine("fuzztag", "=>", "<="))
	generator := NewGenerator(nodes, map[string]func(string) []string{
		"int": func(s string) []string {
			return []string{s}
		},
	})
	for {
		if v, ok := generator.Generate(); ok {
			println(v)
		} else {
			break
		}
	}
}
