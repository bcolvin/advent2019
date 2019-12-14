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
	var planets []*planet
	s := bufio.NewScanner(f)
	for s.Scan() {
		var p point
		_, err := fmt.Sscanf(s.Text(), "<x=%d, y=%d, z=%d>", &p.x, &p.y, &p.z)
		if err != nil {
			fmt.Errorf("Error reading planet %v\n", err)
		}
		planets = append(planets, &planet{p, point{}, point{}})
	}
	timePass(1000, planets)

	var e int
	for _, planet := range planets {
		p := abs(planet.pos.x) + abs(planet.pos.y) + abs(planet.pos.z)
		k := abs(planet.vel.x) + abs(planet.vel.y) + abs(planet.vel.z)
		e += p * k
	}
	fmt.Println(e)
}

func timePass(steps int, planets []*planet) {
	for step := 0; step < steps; step++ {
		for i := 0; i < len(planets); i++ {
			p1 := planets[i]
			for j := i; j < len(planets); j++ {
				p2 := planets[j]
				if p1.pos.x < p2.pos.x {
					p1.grv.x++
					p2.grv.x--
				} else if p1.pos.x > p2.pos.x {
					p1.grv.x--
					p2.grv.x++
				}
				if p1.pos.y < p2.pos.y {
					p1.grv.y++
					p2.grv.y--
				} else if p1.pos.y > p2.pos.y {
					p1.grv.y--
					p2.grv.y++
				}
				if p1.pos.z < p2.pos.z {
					p1.grv.z++
					p2.grv.z--
				} else if p1.pos.z > p2.pos.z {
					p1.grv.z--
					p2.grv.z++
				}
			}
			p1.vel.x += p1.grv.x
			p1.vel.y += p1.grv.y
			p1.vel.z += p1.grv.z
			p1.grv.x = 0
			p1.grv.y = 0
			p1.grv.z = 0
			p1.pos.x += p1.vel.x
			p1.pos.y += p1.vel.y
			p1.pos.z += p1.vel.z
		}
	}
}

func (p planet) String() string {
	return fmt.Sprintf("pos=<x=%d, y=%d, z= %d>, vel=<x=%d, y=%d, z=%d>", p.pos.x, p.pos.y, p.pos.z, p.vel.x, p.vel.y, p.vel.z)
}

type planet struct {
	pos point
	vel point
	grv point
}

type point struct {
	x, y, z int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
