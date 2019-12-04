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
	for i, v := range str {
		a, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Invalid value at %d: %s\n", i, v)
		}
		vals[i] = a
	}
	vals[1] = 12
	vals[2] = 2
	for i := 0; i < len(vals) && vals[i] != 99; i += 4 {
		lPos := vals[i+1]
		rPos := vals[i+2]
		pos := vals[i+3]
		if vals[i] == 1 {
			vals[pos] = vals[lPos] + vals[rPos]
		} else if vals[i] == 2 {
			vals[pos] = vals[lPos] * vals[rPos]
		}
	}
	fmt.Printf("%d\n", vals[0])
}
