package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/logrusorgru/aurora"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const id int = 9

func main() {
	done := make(chan bool)
	botID := make(chan string)

	botIDSent := false

	conn0, err := net.Dial("tcp", "52.49.91.111:2092")
	check(err)

	for i := 0; i < 7; i++ {
		go func(botI int) {
			if botI == 0 {
				conn := conn0
				defer conn.Close()
				connbuf := bufio.NewReader(conn)

				bots := make([]string, 0)

				wrote := false
				began := false

				// teh current round counter
				// var magic int

				promise := 1337

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
						bullshit := fmt.Sprintf("ROOT: %s", line)
						fmt.Printf("%s\n", aurora.Green(bullshit))
						roundData := strings.SplitN(line, ":", 2)

						if len(roundData) < 2 {
							continue
						}

						// round, err := strconv.Atoi(roundData[0][6:])
						// if err != nil {
						// 	continue
						// }

						try := strings.Split(roundData[1], "LEARN")

						if len(try) < 2 || wrote {
							continue
						}

						var begin int
						var end int

						for i, c := range roundData[1] {
							if c == '[' {
								begin = i + 1
							}

							if c == ']' {
								end = i
								break
							}
						}

						ownerSplit := strings.Split(roundData[1], "secret_owner: ")
						owner := ""

						for _, c := range ownerSplit[1] {
							if c == '}' {
								break
							}

							owner = fmt.Sprintf("%s%c", owner, c)
						}

						servers := strings.Split(roundData[1][begin:end], ",")

						if !wrote {
							for _, server := range servers {
								tmp := fmt.Sprintf("{%d,9}", promise)
								prepare := fmt.Sprintf("PREPARE %s -> %s\n", tmp, server)

								fmt.Println(">", prepare[:len(prepare)-1])

								n, err := conn.Write([]byte(prepare))
								check(err)

								if len(prepare) != n {
									fmt.Println("ERROR: could not write message")
									continue
								}
							}

							// promise++
							// val := try[1][:len(try[1])-17]
							val := "{servers: ["

							for i := len(servers) - 1; i >= 0; i-- {
								if i == 0 {
									val = fmt.Sprintf("%s%s", val, servers[i])
								} else {
									val = fmt.Sprintf("%s%s,", val, servers[i])
								}
							}

							if !botIDSent {
								bot := <-botID
								bots = append(bots, bot)

								if len(bots) == 6 {
									botIDSent = true
								}

								val = fmt.Sprintf("%s,%s], secret_owner: %s}", val, bot, owner)
							} else {
								val = fmt.Sprintf("%s], secret_owner: %s}", val, "9")
							}

							for _, server := range servers {
								// accept := fmt.Sprintf("ACCEPT {id: {%d,%d}, value: no_proposal} -> %s\n", 1, 9, server)
								// val := strings.Split(try[1][:len(try[1])-17], " [")
								accept := fmt.Sprintf("ACCEPT {id: {%d,%d}, value: %s} -> %s\n", promise, 9, val, server)

								fmt.Println(">", accept[:len(accept)-1])

								n, err := conn.Write([]byte(accept))
								check(err)

								if len(accept) != n {
									fmt.Println("ERROR: could not write message")
									continue
								}
							}

							// wrote = true
						}
					}
				}

			} else {
				conn, err := net.Dial("tcp", "52.49.91.111:2092")
				check(err)

				defer conn.Close()
				connbuf := bufio.NewReader(conn)

				var currID string

				for {
					line, err := connbuf.ReadString('\n')
					check(err)

					line = line[:len(line)-1]

					fmt.Printf("%s\n", aurora.Cyan(line))

					if len(line) > 0 {
						var begin int

						for i, c := range line {
							if c == ':' {
								begin = i + 2
								break
							}
						}

						currID = line[begin:]
						botID <- currID
					} else {
						continue
					}

					bullshit := fmt.Sprintf("BOT #%s: %s", currID, line)
					fmt.Printf("%s\n", aurora.Cyan(bullshit))
				}
			}
		}(i)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			fmt.Println("exit")
			done <- true
		}
	}()

	<-done
}
