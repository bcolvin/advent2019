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
	var planets system
	s := bufio.NewScanner(f)
	for s.Scan() {
		var p point
		_, err := fmt.Sscanf(s.Text(), "<x=%d, y=%d, z=%d>", &p.x, &p.y, &p.z)
		if err != nil {
			fmt.Errorf("Error reading planet %v\n", err)
		}
		planets = append(planets, &planet{p, point{}, point{}})
	}
	e := gogogo(planets)
	fmt.Println(e)
}

func gogogo(planets system) int {
	init := make(system, len(planets))
	copy(init, planets)
	for step := 0; ; step++ {
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
		if step%100000 == 0 {
			fmt.Println(step)
		}

		if init.equals(planets) {
			return step
		}
	}
}

func (s system) equals(p system) bool {
	for i := 0; i < len(p); i++ {
		if !s[i].equals(p[i]) {
			return false
		}
	}
	return true
}

func (p system) String() string {
	var str string
	for _, p := range p {
		str += fmt.Sprintf("pos=<x=%d, y=%d, z= %d>, vel=<x=%d, y=%d, z=%d>", p.pos.x, p.pos.y, p.pos.z, p.vel.x, p.vel.y, p.vel.z)
	}
	return str
}

type system []*planet

func (p *planet) equals(p1 *planet) bool {
	return p.pos.x == p1.pos.x && p.pos.y == p1.pos.y && p.pos.z == p1.pos.z &&
		p.vel.x == p1.vel.x && p.vel.y == p1.vel.y && p.vel.z == p1.vel.z
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
