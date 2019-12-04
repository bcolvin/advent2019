package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	vals := make(map[int]int)
	s := bufio.NewScanner(f)
	total := 0
	for s.Scan() {
		var n int
		_, err := fmt.Sscanf(s.Text(), "%d", &n)
		if err != nil {
			log.Fatalf("Error reading into value")
		}
		total += calculateWeight(n, vals)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(total)
}

func calculateWeight(val int, weights map[int]int) int {
	v := weights[val]
	if v == -1 {
		return 0
	} else if v == 0 {
		v = val/3 - 2
		nv := 0
		if v > 0 {
			nv = v + calculateWeight(v, weights)
		}
		weights[val] = nv
		return nv
	} else {
		return v
	}
}
