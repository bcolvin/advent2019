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
	row := 0
	id := 1
	var roids []*asteroid
	for s.Scan() {
		for k, ru := range []rune(s.Text()) {
			if '#' == ru {
				roids = append(roids, &asteroid{id, k, row, 0})
				id++
			}
		}
		row++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	var max *asteroid
	for i := 0; i < len(roids); i++ {
		for j := 0; j < len(roids); j++ {
			if i != j {
				isViewable := true
				for k := 0; k < len(roids); k++ {
					if i != k && j != k && inLine(roids[i], roids[k], roids[j]) {
						isViewable = false
						break
					}
				}
				if isViewable {
					roids[i].neighbors++
				}
			}
		}
		if max == nil || roids[i].neighbors > max.neighbors {
			max = roids[i]
		}
	}
	fmt.Println(max.neighbors)
	fmt.Println(max)
}

func inLine(a, b, c *asteroid) bool {
	return ((b.x-a.x)*(c.y-a.y) == (c.x-a.x)*(b.y-a.y)) && inBetween(a, b, c)
}

func inBetween(a, b, c *asteroid) bool {
	var x, y bool
	if a.x <= b.x {
		x = a.x <= c.x && c.x <= b.x
	} else {
		x = b.x <= c.x && c.x <= a.x
	}
	if a.y <= b.y {
		y = a.y <= c.y && c.y <= b.y
	} else {
		y = b.y <= c.y && c.y <= a.y
	}
	return x && y
}

func (a asteroid) String() string {
	return fmt.Sprintf("%d:{%d,%d}", a.id, a.x, a.y)
}

type asteroid struct {
	id        int
	x, y      int
	neighbors int
}
