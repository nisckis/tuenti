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
	pos int
	b   byte
}

const MOD int64 = 4294967295

func main() {
	initCrc32Table()

	crc := ^uint32(1341234123)

	fmt.Printf("before: %08x %08b\n", ^crc, byte(^crc))

	for i := int64(0); i < MOD; i++ {
		crc = Globaltable[byte(crc)^byte(0)] ^ (crc >> 8)
	}

	fmt.Printf("after: %08x vs %08x %08b\n", ^crc, 1341234123, byte(^crc))

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
		// tn := time.Now()

		line := strings.Split(scanner.Text(), " ")
		fName := line[0]
		nMods, _ := strconv.Atoi(line[1])

		mods := make([]mod, nMods)

		for i := 0; i < nMods; i++ {
			scanner.Scan()
			modLine := strings.Split(scanner.Text(), " ")
			modByte, _ := strconv.Atoi(modLine[1])
			pos, _ := strconv.Atoi(modLine[0])
			mods[i] = mod{pos, byte(modByte)}
		}

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
				magic := hdr.Size % MOD
				crc := ^uint32(0)

				for i := int64(0); i < magic; i++ {
					crc = Globaltable[byte(crc)^byte(0)] ^ (crc >> 8)
				}

				crc = ^crc

				fmt.Printf("%s %d: %08x\n", fName, 0, crc)

				for i, mod := range mods {
					crc = Globaltable[byte(crc)^mod.b] ^ (crc >> 8)
					crc = ^crc

					fmt.Printf("%s %d: %08x %08b\n", fName, i+1, crc, byte(crc))
				}

				break
			}
		}

		currCase++
	}

	fmt.Printf("%d done -> total elapsed %v\n", currCase, time.Since(t0))
}
