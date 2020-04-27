package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var allowed = []rune{' ', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'ñ', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'á', 'é', 'í', 'ó', 'ú', 'ü'}

type Word struct {
	key string
	val int
	pos int
}

type lessFunc func(p1, p2 *Word) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	changes []Word
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(changes []Word) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

type ByWord []Word

func (a ByWord) Len() int           { return len(a) }
func (a ByWord) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWord) Less(i, j int) bool { return a[i].val > a[j].val }

type ByUnicode []Word

func (a ByUnicode) Len() int           { return len(a) }
func (a ByUnicode) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUnicode) Less(i, j int) bool { return a[i].key < a[j].key }

func main() {
	if len(os.Args) <= 1 {
		panic("no args")
	}

	t := make(map[string]int)

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parsed := ""
		raw := strings.ToLower(scanner.Text())

		for _, c := range raw {
			ok := false
			for _, a := range allowed {
				if c == a {
					ok = true
					break
				}
			}

			if ok {
				parsed += string(c)
			} else {
				parsed += " "
			}
		}

		str := strings.Split(parsed, " ")

		for _, s := range str {
			if utf8.RuneCountInString(s) < 3 {
				continue
			}

			if _, ok := t[s]; ok {
				t[s]++
			} else {
				t[s] = 1
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	byVal := func(c1, c2 *Word) bool {
		return c1.val > c2.val
	}

	equiv := make(map[rune]byte)
	for i := 1; i < len(allowed); i++ {
		equiv[allowed[i]] = byte(i - 1)
	}

	byKey := func(c1, c2 *Word) bool {
		return c1.key < c2.key
	}

	var arr []Word
	pos := 0

	for k := range t {
		arr = append(arr, Word{k, t[k], pos})
		pos++
	}

	OrderedBy(byVal, byKey).Sort(arr)

	var cases int
	fmt.Scanf("%d", &cases)
	for i := 0; i < cases; i++ {
		var input string
		fmt.Scan(&input)

		number, err := strconv.Atoi(input)

		if err == nil {
			fmt.Printf("Case #%d: %s %d\n", i+1, arr[number-1].key, arr[number-1].val)
		} else {
			for j := range arr {
				if arr[j].key == input {
					fmt.Printf("Case #%d: %d #%d\n", i+1, t[input], j+1)
					break
				}
			}
		}
	}
}
