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
	vals := make(map[point]int)
	s := bufio.NewScanner(f)
	start := &point{0, 0}
	dist := 100000
	for line := 1; s.Scan(); line++ {
		line, err := layLine(line, strings.Split(s.Text(), ","))
		if err != nil {
			fmt.Errorf("Error Laying line %d\n", line)
			continue
		}
		if len(line.path) > 0 {
			for _, v := range line.path {
				id := vals[*v]
				if id != 0 && id != line.id {
					d := v.distance(start)
					if d < dist {
						dist = d
					}
				}
				vals[*v] = line.id
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
