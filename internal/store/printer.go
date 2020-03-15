package store

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func (s *Site) PrintToFile(file *os.File) {
	dw := bufio.NewWriter(file)
	defer dw.Flush()
	s.printHelper("/", 0, map[string]bool{}, map[string]bool{}, func(level int, url string) {
		_, err := dw.WriteString(strings.Repeat(" ", level) + url + "\n")
		if err != nil {
			log.Err(err)
		}
	})
}

// PrintToStdOut for debugging purposes
func (s *Site) PrintToStdOut() {
	s.printHelper("/", 0, map[string]bool{}, map[string]bool{}, func(level int, url string) {
		fmt.Println(strings.Repeat(" ", level) + url)
	})
}

func (s *Site) printHelper(node string, level int, visited map[string]bool, skip map[string]bool, printFn func(int, string)) {
	if visited[node] {
		return
	}
	visited[node] = true
	printFn(level, node)
	for c := range s.store[node] {
		if skip[c] == false {
			s.printHelper(c, level+1, visited, merge(skip, s.store[node]), printFn)
		}

	}
}

func merge(m1, m2 map[string]bool) map[string]bool {
	res := map[string]bool{}
	for k, v := range m1 {
		res[k] = v
	}
	for k, v := range m2 {
		res[k] = v
	}
	return res
}
