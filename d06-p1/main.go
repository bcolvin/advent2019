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
	root := &planet{id: "SUN"}
	planets := make(map[string]*planet)
	s := bufio.NewScanner(f)
	for s.Scan() {
		sp := strings.Split(s.Text(), ")")
		p := planets[sp[0]]
		if p == nil {
			p = &planet{
				sp[0],
				root,
				nil,
				0,
			}
			root.addPlanet(p)
			planets[p.id] = p
		}
		np := planets[sp[1]]
		if np == nil {
			np = &planet{
				sp[1],
				p,
				nil,
				0,
			}
			planets[np.id] = np
		}
		p.addPlanet(np)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	root = planets["COM"]
	root.parent = nil
	root.countOrbits()
	ret := 0
	for _, v := range planets {
		ret += v.orbits
	}
	fmt.Println(ret)
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

func (p *planet) countOrbits() {
	if p.children != nil {
		p.orbits = len(p.children)
		for _, v := range p.children {
			v.countOrbits()
			p.orbits += v.orbits
		}
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
		return fmt.Sprintf("%s planet orbits %s\nIt has %d direct orbits and %d indirect orbits\n", p.id, p.parent.id, do, p.orbits) // + ret
	}
	return fmt.Sprintf("%s is center\nIt has %d direct orbits and %d indirect orbits\n", p.id, do, p.orbits) // + ret
}

type planet struct {
	id       string
	parent   *planet
	children map[string]*planet
	orbits   int
}
