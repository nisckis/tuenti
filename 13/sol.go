package main

import (
	"fmt"
	"math"
	"sort"
)

// sorting code block

type Try struct {
	height, surface, used uint64
}

type lessFunc func(p1, p2 *Try) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	changes []Try
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(changes []Try) {
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

// sorting code block end

func around(n, m uint64) uint64 {
	return 2 * (n + m + 2)
}

func arounds(s uint64) uint64 {
	return 2 * (s + 2)
}

// computes the used packs of a structure of height i
// and x + y = j, where (x, y) is the dimensions of the
// higher most part
func compute(i, j uint64) uint64 {
	// s is the number top level squares
	//
	//     x
	//  ---------
	//  | | | | | y
	//  ---------
	//
	//  s = x + y
	s := uint64(j)

	x := s / 2
	var y uint64
	if s%2 == 0 {
		y = x
	} else {
		y = x + 1
	}

	// total number of packs used
	t := uint64(i) * x * y

	// the current height of the tower
	h := uint64(i)

	// fmt.Println("base: ", h)

	for {
		if h == 2 {
			break
		}

		a1 := arounds(s)
		s += 4

		a2 := arounds(s)
		s += 4

		t += a1 * (h - 2)
		t += a2 * (h - 1)
		// fmt.Printf("around(%d): %d (* %d) %d (* %d)\n", h, a1, h-2, a2, h-1)

		h--
	}

	return t
}

// binary search
// super optimal to find a candidate for used packs
// the complexity is O(log(n))
// so for n = 2^62, its around 20 iterations to find a good candidate
// as compute takes at most 20-30ms for the higher inputs this
// workaround is pretty damn fast
func search(target, s uint64) (uint64, error) {
	l := uint64(3)
	r := uint64(15000000)

	// fmt.Println("-----------------------------------")
	// fmt.Println("searching for target", target, s)
	its := 0

	var m uint64

	for {
		if l > r {
			break
		}

		m = uint64(math.Floor(float64(l+r) / 2))
		x := compute(m, s)

		// fmt.Println("current loop", l, r)
		// fmt.Println("m and its compute(m)", m, x)

		if x < target {
			l = m + 1
		} else if x > target {
			r = m - 1
		} else {
			// fmt.Println(its, "iterations and found")
			return m, nil
		}

		its++
	}

	// fmt.Println(its, "iterations and not found")
	return m, fmt.Errorf("could not find a candidate for %d", target)
}

func main() {
	increasingUsed := func(c1, c2 *Try) bool {
		return c1.used < c2.used
	}

	decreasingHeight := func(c1, c2 *Try) bool {
		return c1.used > c2.used
	}

	var t int
	fmt.Scanf("%d ", &t)

	for tc := 1; tc <= t; tc++ {
		var val uint64
		fmt.Scanf("%d ", &val)

		if val < 43 {
			fmt.Printf("Case #%d: IMPOSSIBLE\n", tc)
			continue
		}

		found := false
		var height uint64
		var surface uint64
		var used uint64

		ms := make([]Try, 8)
		i := uint64(2)

		for {
			if i == 10 {
				break
			}

			c, err := search(val, i)
			ms[i-2] = Try{c, i, compute(c, i)}

			if err == nil {
				height = c
				surface = i
				used = compute(c, i)
				found = true
				break
			}

			i++
		}

		if found {
			fmt.Printf("Case #%d: %d %d\n", tc, height, used)
			continue
		}

		OrderedBy(decreasingHeight, increasingUsed).Sort(ms)

		// fmt.Println(val)
		// fmt.Println(ms)

		// pick a new candidate
		for _, m := range ms {
			height = m.height
			surface = m.surface
			used = m.used

			if m.used < val {
				break
			}
		}

		if used < val {
			// try to find a better candidate up
			for hi := height; ; hi++ {
				best := false

				for si := surface; si < 10; si++ {
					ui := compute(hi, si)
					// fmt.Println(hi, si, ui)

					if ui > val {
						best = true
						break
					}

					height = hi
					surface = si
					used = ui
				}

				if best {
					break
				}

				surface = 2
			}
		} else {
			// try to find a better candidate down
			for hi := height; hi > 2; hi-- {
				best := false

				for si := surface; si > 1; si-- {
					ui := compute(hi, si)
					// fmt.Println(hi, si, ui)
					height = hi
					surface = si
					used = ui

					if ui <= val {
						best = true
						break
					}
				}

				if best {
					break
				}

				surface = 9
			}
		}

		fmt.Printf("Case #%d: %d %d\n", tc, height, used)
	}
}
