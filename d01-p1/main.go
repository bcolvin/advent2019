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
	s := bufio.NewScanner(f)
	total := 0
	for s.Scan() {
		var n int
		_, err := fmt.Sscanf(s.Text(), "%d", &n)
		if err != nil {
			log.Fatalf("Error reading into value")
		}
		total += n/3 - 2
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(total)
}
