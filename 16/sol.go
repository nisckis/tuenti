package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Group struct {
	emp    int
	floors []int
}

type Floor struct {
	freeWC int
	groups []int
}

func main() {
	if len(os.Args) <= 1 {
		panic("no args")
	}

	input, err := os.Open(os.Args[1])
	check(err)
	defer input.Close()

	scanner := bufio.NewScanner(input)
	check(scanner.Err())

	scanner.Scan()
	first := scanner.Text()
	cases, err := strconv.Atoi(first)
	check(err)

	for testCase := 1; testCase <= cases; testCase++ {
		var numGroups int
		var numFloors int

		scanner.Scan()
		line := strings.Split(scanner.Text(), " ")

		parsedFLoors, err := strconv.Atoi(line[0])
		check(err)
		numFloors = parsedFLoors

		floorsTmp := make([][2]int, numFloors)
		floorsByGroup := make([][]bool, numFloors)
		floorsByGroupInt := make([][]int, numFloors)

		parsedGroups, err := strconv.Atoi(line[1])
		check(err)
		numGroups = parsedGroups

		floors := make([]Floor, numFloors)

		for i := 0; i < numFloors; i++ {
			floors[i] = Floor{0, make([]int, 0)}

			floorsTmp[i] = [2]int{0, 0}
			floorsByGroup[i] = make([]bool, numGroups)
			floorsByGroupInt[i] = make([]int, numGroups)

			for j := 0; j < numGroups; j++ {
				floorsByGroup[i][j] = false
				floorsByGroupInt[i][j] = 0
			}
		}

		groups := make([]Group, numGroups)

		var smallest int

		for i := 0; i < numGroups; i++ {
			scanner.Scan()
			groupInfoLine := strings.Split(scanner.Text(), " ")

			numEmp, err := strconv.Atoi(groupInfoLine[0])
			check(err)

			numFloors, err := strconv.Atoi(groupInfoLine[1])
			check(err)

			groupInfo := make([]int, numFloors)

			scanner.Scan()
			groupFloors := strings.Split(scanner.Text(), " ")

			for j, groupFloor := range groupFloors {
				tmp, err := strconv.Atoi(groupFloor)
				check(err)
				floorsTmp[tmp][0] += numEmp
				floorsTmp[tmp][1] += numEmp / numFloors
				floorsByGroup[tmp][i] = true
				floors[tmp].groups = append(floors[tmp].groups, i)
				floorsByGroupInt[tmp][i] += numEmp
				groupInfo[j] = tmp
			}

			groups[i] = Group{numEmp, groupInfo}

			if i == 0 {
				smallest = numEmp
			} else {
				if numEmp < smallest {
					smallest = numEmp
				}
			}
		}

		largest := 1

		persons := 0
		for i := range groups {
			persons += groups[i].emp
		}

		locGroups := make([]Group, len(groups))
		locFloors := make([]Floor, len(floors))

		for {

			for i := range locFloors {
				locFloors[i].freeWC = 1
				locFloors[i].groups = make([]int, len(floors[i].groups))
				copy(locFloors[i].groups, floors[i].groups)
			}

			for i := range locGroups {
				locGroups[i].emp = groups[i].emp
				locGroups[i].floors = make([]int, len(groups[i].floors))
				copy(locGroups[i].floors, groups[i].floors)
			}

			iteration := 0
			lastZeta := [3]int{0, 0, 0}

			for {
				for i := range locFloors {
					locFloors[i].freeWC = 1
				}

				for i := range locGroups {
					locGroups[i].emp = 0
				}

				// find current tree
				treeGroups, treeFloors := make([]bool, len(locGroups)), make([]bool, len(locFloors))
				newGroups, newFloors := make(map[int]bool), make(map[int]bool)

				for i := range treeFloors {
					treeFloors[i] = false
				}

				for i := range treeGroups {
					treeGroups[i] = true
				}

				for i := range locGroups {
					for j := range locGroups[i].floors {
						if locFloors[locGroups[i].floors[j]].freeWC == 1 {
							// assignement done
							locGroups[i].emp = 1
							locFloors[locGroups[i].floors[j]].freeWC = 0
							break
						}
					}

					// add the rows with no assignement
					if locGroups[i].emp == 0 {
						treeGroups[i] = false
						newGroups[i] = true
					}
				}

				isColumn := true

				for {
					changed := false

					if isColumn {
						for i := range locFloors {
							for j := range locFloors[i].groups {
								if _, ok := newGroups[locFloors[i].groups[j]]; ok {
									treeFloors[i] = true
									newFloors[i] = true
									changed = true
								}
							}
						}

						newGroups = make(map[int]bool)
					} else {
						for i := range locGroups {
							for j := range locGroups[i].floors {
								if _, ok := newFloors[locGroups[i].floors[j]]; ok {
									treeGroups[i] = false
									newGroups[i] = true
									changed = true
								}
							}
						}

						newFloors = make(map[int]bool)
					}

					if !changed {
						break
					}
				}

				zeta := 0
				eta := 0

				for i := range locGroups {
					locGroups[i].emp = groups[i].emp

					if treeGroups[i] {
						zeta += locGroups[i].emp
					}

					eta += locGroups[i].emp
				}

				for i := range locFloors {
					locFloors[i].freeWC = largest

					if treeFloors[i] {
						zeta += largest
					}
				}

				// check if optimal matrix situation
				// if is optimal

				// if this is true we are in the optimal case
				if zeta >= eta {
					break
				}

				if iteration > 3 && zeta == lastZeta[(iteration-2)%3] {
					break
				}

				lastZeta[iteration%3] = zeta

				// add (append) to everything that's outside
				// the tree and the double crossings that are
				// inside delete them

				for i := range locGroups {
					for j := range locFloors {
						if treeGroups[i] && treeFloors[j] {
							// delete elements
							for k, floor := range locFloors[j].groups {
								if floor == i {
									locFloors[j].groups = locFloors[j].groups[:k+copy(locFloors[j].groups[k:], locFloors[j].groups[k+1:])]
									break
								}
							}

							for k, grp := range locGroups[i].floors {
								if grp == i {
									locGroups[i].floors = locGroups[i].floors[:k+copy(locGroups[i].floors[k:], locGroups[i].floors[k+1:])]
									break
								}
							}
						} else {
							// add element
							locFloors[j].groups = append(locFloors[j].groups, i)
							locGroups[i].floors = append(locGroups[i].floors, j)
						}
					}
				}

				iteration++
			}

			success := false

			for {
				minI, minVal, isGroup := 0, 800001, true

				for i := range locGroups {
					if locGroups[i].emp > 0 && len(locGroups[i].floors) <= minVal {
						minVal = len(locGroups[i].floors)
						minI = i
					}
				}

				for i := range locFloors {
					if locFloors[i].freeWC > 0 && len(locFloors[i].groups) <= minVal {
						minVal = len(locFloors[i].groups)
						minI = i
						isGroup = false
					}
				}

				if isGroup {
					row := locGroups[minI]

					for i := len(row.floors) - 1; i >= 0; i-- {
						var locMin int

						floor := locFloors[row.floors[i]]
						magic := floor.freeWC

						if row.emp < magic {
							locMin = row.emp
						} else {
							locMin = magic
						}

						row.emp -= locMin
						locGroups[minI].emp -= locMin
						locFloors[row.floors[i]].freeWC -= locMin

						for j, grp := range floor.groups {
							if grp == minI {
								locFloors[row.floors[i]].groups = floor.groups[:j+copy(floor.groups[j:], floor.groups[j+1:])]
								break
							}
						}
					}

					if locGroups[minI].emp != 0 {
						success = false
						break
					}
				} else {
					column := locFloors[minI]

					for i := range column.groups {
						var locMin int

						grp := locGroups[column.groups[i]]
						magic := grp.emp

						if column.freeWC < magic {
							locMin = column.freeWC
						} else {
							locMin = magic
						}

						column.freeWC -= locMin
						locFloors[minI].freeWC -= locMin
						locGroups[column.groups[i]].emp -= locMin
					}

					locFloors[minI].freeWC = 0
				}

				peopleSum := 0
				for i := range locGroups {
					peopleSum += locGroups[i].emp
				}

				if peopleSum == 0 {
					success = true
					break
				}
			}

			if success {
				break
			}

			largest++
		}

		fmt.Printf("Case #%d: %d\n", testCase, largest)
	}
}
