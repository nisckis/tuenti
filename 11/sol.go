package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func dp(v []int, amount int) int {
	solution := make([][]int, len(v)+1)

	for i := range solution {
		solution[i] = make([]int, amount+1)
		solution[i][0] = 1
	}

	for i := range solution {
		solution[0][i] = 0
	}

	for i := 1; i <= len(v); i++ {
		for j := 1; j <= amount; j++ {
			if v[i-1] <= j {
				solution[i][j] = solution[i-1][j] + solution[i][j-v[i-1]]
			} else {
				solution[i][j] = solution[i-1][j]
			}
		}
	}

	return solution[len(v)][amount]
}

func main() {
	var t int
	fmt.Scanf("%d", &t)

	tCase := 1
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	re := regexp.MustCompile(`\r?\n`)

	for {
		if tCase > t {
			break
		}

		line, _ := reader.ReadString('\n')
		input := strings.Split(line, " ")
		cleaned := re.ReplaceAllString(input[0], "")
		n, err := strconv.Atoi(cleaned)
		if err != nil {
			panic(err)
		}

		allowed := make([]bool, n-1)

		for i := range allowed {
			allowed[i] = true
		}

		for i := 1; i < len(input); i++ {
			cleaned := re.ReplaceAllString(input[i], "")
			tmp, err := strconv.Atoi(cleaned)
			if err != nil {
				panic(err)
			}

			if tmp < n {
				allowed[tmp-1] = false
			}
		}

		var using []int

		for i := range allowed {
			if allowed[i] {
				using = append(using, i+1)
			}
		}

		fmt.Printf("Case #%d: %d\n", tCase, dp(using, n))
		tCase++
	}
}
