package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	str := strings.Split(string(b), ",")
	vals := make([]int, len(str))
	for i := 0; i < len(vals); i++ {
		for j := 0; j < len(vals); j++ {
			reset(str, vals)
			vals[1] = i
			vals[2] = j
			res := execute(vals)
			if res == 19690720 {
				fmt.Printf("100 * noun(%d) + verb(%d) = %d\n", i, j, 100*i+j)
				break
			}
		}
	}
}

func printVals(ints []int) {
	for _, v := range ints {
		fmt.Printf("%d,", v)
	}
	fmt.Println()
}

func reset(str []string, vals []int) {
	for i, v := range str {
		a, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Invalid value at %d: %s\n", i, v)
		}
		vals[i] = a
	}
}

func op(instPtr int, vals []int) (int, bool, error) {
	switch vals[instPtr] {
	case 1: //add
		noun := vals[instPtr+1]
		verb := vals[instPtr+2]
		res := vals[instPtr+3]
		vals[res] = vals[noun] + vals[verb]
		return 4, false, nil
	case 2: //mult
		noun := vals[instPtr+1]
		verb := vals[instPtr+2]
		res := vals[instPtr+3]
		vals[res] = vals[noun] * vals[verb]
		return 4, false, nil
	case 99:
		return 1, true, nil
	}
	return 0, false, fmt.Errorf("Unsupported opcode: %d\n", vals[instPtr])
}

func execute(vals []int) int {
	var err error
	for insPtr, halt, inc := 0, false, 0; !halt && insPtr < len(vals); insPtr += inc {
		inc, halt, err = op(insPtr, vals)
		if err != nil {
			panic(err)
		}
	}
	return vals[0]
}
