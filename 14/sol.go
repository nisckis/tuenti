package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const id int = 9

func main() {
	conn, err := net.Dial("tcp", "52.49.91.111:2092")
	check(err)

	defer conn.Close()
	connbuf := bufio.NewReader(conn)

	wrote := false
	began := false

	// teh current round counter
	// var round int

	for {
		line, err := connbuf.ReadString('\n')
		check(err)

		line = line[:len(line)-1]

		// wait this real communication
		if !began {
			if len(line) < 5 {
				continue
			}

			if strings.Compare(line[:5], "ROUND") == 0 {
				began = true
			}
		}

		if began && len(line) > 0 {
			fmt.Println(line)
			roundData := strings.SplitN(line, ":", 2)

			if len(roundData) < 2 {
				continue
			}

			// round, err := strconv.Atoi(roundData[0][6:])
			try := strings.Split(roundData[1], "ACCEPTED")

			if len(try) < 2 || wrote {
				continue
			}

			// this is the id(?)
			tmp := "{1,9}"

			// prepare command ?
			prepare := ""
			// prepare := fmt.Sprintf("PREPARE %s -> %d\n", tmp, id)

			// accept command ?
			accept := fmt.Sprintf("%s{id: %s, value: %s} -> 9\n", prepare, tmp, try[1])

			n, err := conn.Write([]byte(accept))
			check(err)

			if len(accept) != n {
				fmt.Println("ERROR: could not write message")
				continue
			}

			wrote = true
		}
	}
}
