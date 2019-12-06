package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	planets := make(map[string]*planet)
	s := bufio.NewScanner(f)
	for s.Scan() {
		sp := strings.Split(s.Text(), ")")
		p := planets[sp[0]]
		if p == nil {
			p = &planet{
				sp[0],
				nil,
				nil,
				0,
				false,
			}
			planets[p.id] = p
		}
		np := planets[sp[1]]
		if np == nil {
			np = &planet{
				sp[1],
				p,
				nil,
				0,
				false,
			}
			planets[np.id] = np
		}
		p.addPlanet(np)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	me := planets["YOU"]
	if me == nil {
		log.Fatal("I'm undefined")
		return
	}
	san := planets["SAN"]
	if san == nil {
		log.Fatal("SAN undefined")
		return
	}
	me.findPath(san)
	fmt.Println(san.dist - 2)
}

func (p *planet) addPlanet(orb *planet) error {
	if p.children == nil {
		p.children = make(map[string]*planet)
	}
	np := p.children[orb.id]
	if np != nil {
		return fmt.Errorf("%s already orbits %s\n", orb.id, p.id)
	}
	p.children[orb.id] = orb

	if orb.parent != nil && orb.parent != p && orb.parent.children != nil {
		orb.parent.children[orb.id] = nil
	}
	orb.parent = p
	return nil
}

func (p *planet) findPath(p1 *planet) {
	if p == p1 || p.visited {
		return
	}
	p.visited = true
	if p.children != nil {
		for _, v := range p.children {
			if v.dist == 0 {
				v.dist = 1 + p.dist
				v.findPath(p1)
			}
		}
	}
	if p.parent != nil && !p.parent.visited {
		p.parent.dist = 1 + p.dist
		p.parent.findPath(p1)
	}
}

func (p *planet) String() string {
	var ret string
	do := 0
	if p.children != nil {
		do = len(p.children)
		for _, v := range p.children {
			ret += " " + v.String()
		}
	}
	if p.parent != nil {
		return fmt.Sprintf("%s planet orbits %s\nIt is %d distance from me\n", p.id, p.parent.id, do, p.dist) // + ret
	}
	return fmt.Sprintf("%s is center\nIt is %d distance from me\n", p.id, do, p.dist) // + ret
}

type planet struct {
	id       string
	parent   *planet
	children map[string]*planet
	dist     int
	visited  bool
}
