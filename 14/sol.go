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

	wrote := false
	began := false

	// var round *int
	connbuf := bufio.NewReader(conn)

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
			// left := roundData[0][6:]

			if len(roundData) < 2 {
				continue
			}

			try := strings.Split(roundData[1], "ACCEPTED")

			if len(try) < 2 || wrote {
				continue
			}

			tmp := "{1,9}"
			prepare := "" // fmt.Sprintf("PREPARE %s -> %d\n", tmp, id)
			accept := fmt.Sprintf("%s{id: %s, value: %s} -> 9\n", prepare, tmp, try[1])

			n, err := conn.Write([]byte(accept))
			check(err)

			if len(accept) != n {
				panic("could not write entire accept message")
			}
			wrote = true
		}

		// buf := make([]byte, 1024)
		// n, err := conn.Read(buf)
		// check(err)

		// if n == 0 {
		// 	break
		// }

		// data := strings.Split(string(buf[:n]), "\n")

		// for _, line := range data {
		// }

		// if !wrote {
		// 	prepare := "PREPARE {1,9} -> 9\n"
		// 	n, err = conn.Write([]byte(prepare))
		// 	check(err)

		// 	if len(prepare) != n {
		// 		panic("could not write entire prepare message")
		// 	}

		// 	accept := "ACCEPT {id: {1,9}, value: {servers: [], secret_owner: 9}} -> 9\n"
		// 	n, err = conn.Write([]byte(accept))
		// 	check(err)

		// 	if len(accept) != n {
		// 		panic("could not write entire accept message")
		// 	}
		// 	wrote = true
		// }
	}
}
