package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const N int = 300
const X int = 150
const Y int = 150
const Wall rune = '#'
const Princes rune = 'P'

type point [2]int

type stack struct {
	lock sync.Mutex
	s    []point
}

func NewStack() *stack {
	return &stack{sync.Mutex{}, make([]point, 0)}
}

func (s *stack) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.s)
}

func (s *stack) Push(v point) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (point, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return point{0, 0}, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func shortPath(e map[point][]point, v map[point]bool, start, end point) ([]point, error) {
	queue := NewStack()

	discovered := make(map[point]bool)
	parents := make(map[point]point)

	for n := range v {
		discovered[n] = false
	}
	discovered[start] = true

	queue.Push(start)
	found := false

	it := 0

	for {
		current, err := queue.Pop()
		if err != nil {
			break
		}

		if current == end {
			found = true
			break
		}

		for _, neig := range e[current] {
			if discovered[neig] {
				continue
			}

			discovered[neig] = true
			parents[neig] = current
			queue.Push(neig)
		}

		it++
	}

	if !found {
		return nil, errors.New(fmt.Sprintf("could not find a path from %v to %v", start, end))
	}

	var path []point
	current := end

	for {
		path = append([]point{current}, path...)

		if current == start {
			break
		}

		current = parents[current]
	}

	return path, nil
}

func printMaze(maze [5][5]rune) {
	for i := range maze {
		for _, c := range maze[i] {
			fmt.Print(string(c))
		}
		fmt.Printf("\n")
	}
}

func printGrid(grid [][]rune) {
	for i := range grid {
		for _, c := range grid[i] {
			fmt.Print(string(c))
		}
		fmt.Printf("\n")
	}
}

func main() {
	// Possible moves
	moves := []string{
		"2u1l",
		"2u1r",
		"1u2r",
		"1d2r",
		"2d1r",
		"2d1l",
		"1d2l",
		"1u2l",
	}

	movesMap := map[string]point{
		"2u1l": {-2, -1},
		"2u1r": {-2, +1},
		"1u2r": {-1, +2},
		"1d2r": {+1, +2},
		"2d1r": {+2, +1},
		"2d1l": {+2, -1},
		"1d2l": {+1, -2},
		"1u2l": {-1, -2},
	}

	movesBacks := map[string]string{
		"2u1l": "2d1r",
		"2u1r": "2d1l",
		"1u2r": "1d2l",
		"1d2r": "1u2l",
		"2d1r": "2u1l",
		"2d1l": "2u1r",
		"1d2l": "1u2r",
		"1u2l": "1d2r",
	}

	// The grid
	grid := make([][]rune, N)
	for i := 0; i < N; i++ {
		grid[i] = make([]rune, N)

		for j := 0; j < N; j++ {
			grid[i][j] = '#'
		}
	}

	// The graph structures
	graph := make(map[[2]point]string)
	edges := make(map[point][]point, 0)
	vertices := make(map[point]bool, 0)

	// The visited grid
	visited := make([][]bool, N)
	for i := 0; i < N; i++ {
		visited[i] = make([]bool, N)

		for j := 0; j < N; j++ {
			visited[i][j] = false
		}
	}

	// Make the TCP connection
	conn, err := net.Dial("tcp", "52.49.91.111:2003")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// First maze
	buf := make([]byte, 128)
	if _, err := conn.Read(buf); err != nil {
		panic(err)
	}

	data := strings.Split(string(buf), "\n")

	var maze [5][5]rune

	for i := 0; i < 5; i++ {
		j := 0
		for _, c := range data[i] {
			maze[i][j] = c
			j++
		}
	}

	expanded := 0
	princessFound := false
	var princessPos point

	var last point
	hasLast := false

	pos := point{X, Y}
	vertices[pos] = true

	for i := range maze {
		for j := range maze[i] {
			grid[pos[0]-2+i][pos[1]-2+j] = maze[i][j]
		}
	}

	// Init the stack
	stack := NewStack()
	stack.Push(pos)

	its := 0
	elapsed := time.Duration(0)

	for {
		t := time.Now()

		current, err := stack.Pop()
		if err != nil {
			break
		}

		if hasLast {
			path, err := shortPath(edges, vertices, last, current)
			if err != nil {
				panic(err)
			}

			critical := time.Now()

			for i := 1; i < len(path); i++ {
				move := graph[[2]point{path[i-1], path[i]}]

				// Update pos
				d := movesMap[move]
				pos[0], pos[1] = pos[0]+d[0], pos[1]+d[1]

				// Move and update the maze
				if _, err := conn.Write([]byte(move)); err != nil {
					panic(err)
				}

				buf := make([]byte, 128)
				if _, err := conn.Read(buf); err != nil {
					panic(err)
				}

				data := strings.Split(string(buf), "\n")

				for i := 0; i < 5; i++ {
					j := 0
					for _, c := range data[i] {
						maze[i][j] = c
						j++
					}
				}

				for i := range maze {
					for j := range maze[i] {
						grid[pos[0]-2+i][pos[1]-2+j] = maze[i][j]
					}
				}
			}

			if its%100 == 0 {
				fmt.Printf("%d (%v) ", its, time.Since(critical))
			}
		}

		// Get the current K position on the maze
		var kPos point
		for i := range maze {
			for j := range maze[i] {
				if maze[i][j] == 'K' {
					kPos[0], kPos[1] = i, j
				}
			}
		}

		if !visited[current[0]][current[1]] {
			visited[current[0]][current[1]] = true

			for _, move := range moves {
				d := movesMap[move]
				nx, ny := kPos[0]+d[0], kPos[1]+d[1]

				if maze[nx][ny] == Wall {
					continue
				}

				if visited[current[0]+d[0]][current[1]+d[1]] {
					continue
				}

				moved := point{current[0] + d[0], current[1] + d[1]}

				forward := [2]point{
					current,
					moved,
				}
				backward := [2]point{
					moved,
					current,
				}

				graph[forward] = move
				graph[backward] = movesBacks[move]

				edges[current] = append(edges[current], moved)
				edges[moved] = append(edges[moved], current)

				vertices[moved] = true
				stack.Push(moved)

				expanded++

				if maze[nx][ny] == Princes {
					princessFound = true
					princessPos = moved
					break
				}
			}

		}

		last = current
		hasLast = true

		// c := exec.Command("clear")
		// c.Stdout = os.Stdout
		// c.Run()
		// printGrid(grid)

		elapsed += time.Since(t)

		if its%100 == 0 {
			fmt.Println("elapsed", elapsed, "expanded", expanded, "stack", stack.Len(), "last", last, "nodes", len(vertices), "edges", len(edges))
		}

		its++

		if princessFound {
			break
		}
	}

	if !princessFound {
		panic("unlucky....")
	}

	path, err := shortPath(edges, vertices, point{X, Y}, princessPos)
	if err != nil {
		panic(err)
	}

	// Solving the problem

	fmt.Println("found the path to the princess!")
	fmt.Println(path)

	conn.Close()

	conn, err = net.Dial("tcp", "52.49.91.111:2003")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf = make([]byte, 128)
	if _, err := conn.Read(buf); err != nil {
		panic(err)
	}

	data = strings.Split(string(buf), "\n")

	for i := 0; i < 5; i++ {
		j := 0
		for _, c := range data[i] {
			maze[i][j] = c
			j++
		}
	}

	for i := 1; i < len(path); i++ {
		move := graph[[2]point{path[i-1], path[i]}]

		// Update pos
		d := movesMap[move]
		pos[0], pos[1] = pos[0]+d[0], pos[1]+d[1]

		// Move and update the maze
		if _, err := conn.Write([]byte(move)); err != nil {
			panic(err)
		}

		buf := make([]byte, 128)
		if _, err := conn.Read(buf); err != nil {
			panic(err)
		}

		data := strings.Split(string(buf), "\n")
		fmt.Println(data)

		if len(data) < 5 {
			return
		}

		for i := 0; i < 5; i++ {
			j := 0
			for _, c := range data[i] {
				maze[i][j] = c
				j++
			}
		}

		for i := range maze {
			for j := range maze[i] {
				grid[pos[0]-2+i][pos[1]-2+j] = maze[i][j]
			}
		}
	}
}
