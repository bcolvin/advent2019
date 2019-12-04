package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	vals := make(map[point]*intersection)
	s := bufio.NewScanner(f)
	dist := 100000
	for line := 1; s.Scan(); line++ {
		line, err := layLine(line, strings.Split(s.Text(), ","))
		if err != nil {
			fmt.Errorf("Error Laying line %d\n", line)
			continue
		}
		if len(line.path) > 0 {
			for i, v := range line.path {
				inter := vals[*v]
				if inter == nil {
					inter = &intersection{
						intersect:  nil,
						wireLength: make(map[int]int),
					}
					vals[*v] = inter
				}

				inter.wireLength[line.id] = i + 1

				if len(inter.wireLength) > 1 {
					inter.intersect = v
					l := inter.sum()
					//fmt.Printf("Cross path %s, %d\n", *v, l)
					if l < dist {
						dist = l
					}
				}
			}
		}
	}
	fmt.Println(dist)
}

func layLine(line int, directions []string) (wire, error) {
	l := wire{
		id:  line,
		pos: point{0, 0},
	}
	for _, v := range directions {
		var dir rune
		var steps int
		var pts []*point
		_, err := fmt.Sscanf(v, "%c%d", &dir, &steps)
		if err != nil {
			return l, err
		}
		l.pos, pts = move(l.pos, dir, steps)
		l.path = append(l.path, pts...)
	}
	return l, nil
}

func move(start point, direction rune, steps int) (point, []*point) {
	var pts []*point
	cur := &start
	for i := 0; i < steps; i++ {
		switch direction {
		case 'U':
			cur = &point{
				cur.x,
				cur.y + 1,
			}
		case 'D':
			cur = &point{
				cur.x,
				cur.y - 1,
			}
		case 'L':
			cur = &point{
				cur.x - 1,
				cur.y,
			}
		case 'R':
			cur = &point{
				cur.x + 1,
				cur.y,
			}
		}
		pts = append(pts, cur)
	}
	return *cur, pts
}

func (w wire) String() string {
	v := fmt.Sprintf("Wire %d ends at %s with points: \n", w.id, w.pos)
	for _, pt := range w.path {
		v = v + fmt.Sprintf("%s\n", pt)
	}
	return v
}

type intersection struct {
	intersect  *point
	wireLength map[int]int
}

func (i intersection) sum() int {
	l := 0
	if i.wireLength != nil {
		for _, v := range i.wireLength {
			l += v
		}
	}
	return l
}

type wire struct {
	id   int
	pos  point
	path []*point
}

func (p *point) distance(q *point) int {
	return int(math.Abs(float64(p.x-q.x)) + math.Abs(float64(p.y-q.y)))
}
func (p *point) distanceFromZero() int {
	return int(math.Abs(float64(p.x)) + math.Abs(float64(p.y)))
}

func (p point) String() string {
	return fmt.Sprintf("{%d,%d}", p.x, p.y)
}

type point struct {
	x, y int
}
