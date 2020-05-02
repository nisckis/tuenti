package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Globaltable [256]uint32

func initCrc32Table() {
	for i := range Globaltable {
		word := uint32(i)
		for j := 0; j < 8; j++ {
			if word&1 == 1 {
				word = (word >> 1) ^ 0xedb88320
			} else {
				word >>= 1
			}
		}
		Globaltable[i] = word
	}
}

func newCrc32Table() [256]uint32 {
	var table [256]uint32

	for i := range table {
		word := uint32(i)
		for j := 0; j < 8; j++ {
			if word&1 == 1 {
				word = (word >> 1) ^ 0xedb88320
			} else {
				word >>= 1
			}
		}
		table[i] = word
	}

	return table
}

type mod struct {
	pos int64
	b   byte
}

const MOD int64 = 4294967295

var CrcToIndex []uint32
var IndexToCrc []uint32

func initCrcMagicShit() {
	// WARNING!
	// 32 GB of RAM :/
	CrcToIndex = make([]uint32, MOD+1)
	IndexToCrc = make([]uint32, MOD+1)

	if int64(len(CrcToIndex)) != MOD+1 {
		panic("CrcToIndex is not good")
	}

	if int64(len(IndexToCrc)) != MOD+1 {
		panic("IndexToCrc is not good")
	}

	fmt.Println("allocation made :3")

	crc := ^uint32(0)
	for i := int64(0); i <= MOD; i++ {
		IndexToCrc[i] = crc
		CrcToIndex[crc] = uint32(i)
		crc = Globaltable[byte(crc)^byte(0)] ^ (crc >> 8)
	}
}

func doNCrcs(crc uint32, n int64) uint32 {
	index := CrcToIndex[crc]
	return IndexToCrc[(int64(index)+n)%MOD]
}

func main() {
	initCrc32Table()

	tDrama := time.Now()
	initCrcMagicShit()
	fmt.Println("elapsed for the drama init", time.Since(tDrama))

	var wg sync.WaitGroup
	solutions := make([]string, 500)

	if len(os.Args) <= 1 {
		panic("no args")
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	defer input.Close()

	scanner := bufio.NewScanner(input)
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	currCase := 0
	t0 := time.Now()

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		fName := line[0]
		nMods, _ := strconv.Atoi(line[1])

		mods := make([]mod, nMods)

		for i := 0; i < nMods; i++ {
			scanner.Scan()
			modLine := strings.Split(scanner.Text(), " ")
			modByte, _ := strconv.Atoi(modLine[1])
			pos, _ := strconv.ParseInt(modLine[0], 10, 64)
			mods[i] = mod{pos, byte(modByte)}
		}

		wg.Add(1)

		go func(currCase int, mods []mod) {
			table := newCrc32Table()

			defer wg.Done()

			// fmt.Printf("Case #%d executing!\n", currCase)
			file, err := os.Open("animals.tar.gz")
			if err != nil {
				panic(err)
			}

			archive, err := gzip.NewReader(file)
			if err != nil {
				panic(err)
			}

			tr := tar.NewReader(archive)

			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(err)
				}

				target := hdr.Name[10:]

				if strings.Compare(target, fName) == 0 {
					solution := ""

					sizeOrg := hdr.Size
					ordered := make([][]mod, len(mods)+1)

					for x := 0; x <= len(mods); x++ {
						size := sizeOrg
						ordered[x] = make([]mod, x)
						inserted := false

						for j := 0; j < x-1; j++ {
							if ordered[x-1][j].pos < mods[x-1].pos {
								ordered[x][j] = ordered[x-1][j]
							} else {
								if !inserted {
									ordered[x][j] = mods[x-1]
									inserted = true
								}
								ordered[x][j+1] = ordered[x-1][j]
								ordered[x][j+1].pos++
							}
						}

						if !inserted && x > 0 {
							ordered[x][x-1] = mods[x-1]
						}

						crc := ^uint32(0)
						position := int64(0)
						bullsMade := 0

						bulls := ordered[x]

						for {
							if bullsMade >= x {
								break
							}

							diff := int64((bulls[bullsMade].pos - position) % MOD)
							crc = doNCrcs(crc, diff)

							// for i := int64(0); i < diff; i++ {
							// 	crc = table[byte(crc)^byte(0)] ^ (crc >> 8)
							// }

							crc = table[byte(crc)^bulls[bullsMade].b] ^ (crc >> 8)

							position = bulls[bullsMade].pos + 1
							size++
							bullsMade++
						}

						diff := int64((size - position) % MOD)
						crc = doNCrcs(crc, diff)

						// for i := int64(0); i < diff; i++ {
						// 	crc = table[byte(crc)^byte(0)] ^ (crc >> 8)
						// }

						// fmt.Printf("Case #%d - %d done\n", currCase, x)
						solution = fmt.Sprintf("%s%s %d: %08x\n", solution, fName, x, ^crc)
					}

					solutions[currCase] = solution
					break
				}
			}

			// fmt.Printf("Case #%d done!\n", currCase)
		}(currCase, mods)

		currCase++
	}

	wg.Wait()
	fmt.Printf("%d done -> total elapsed %v\n", currCase, time.Since(t0))

	for i := 0; i < currCase; i++ {
		fmt.Printf("%s", solutions[i])
	}
}
