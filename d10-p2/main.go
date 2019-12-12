package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

func getAsteroids(filename string) []*point {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	row := 0
	var pts []*point
	for s.Scan() {
		for k, ru := range []rune(s.Text()) {
			if '#' == ru {
				pts = append(pts, &point{k, row})
			}
		}
		row++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return pts
}

func main() {
	roids := getAsteroids("input.txt")
	station := getStation(roids)
	fmt.Println(station)
	var slops []*asteroid
	for id, as := range roids {
		if station != as {
			sl := &asteroid{id, *as, station.slope(*as), station.distance(*as)}
			slops = append(slops, sl)
		}
	}

	sort.Slice(slops, func(i, j int) bool {
		a := slops[i]
		b := slops[j]
		if a.v.region != b.v.region {
			return a.v.region < b.v.region
		}
		if a.v.slope == b.v.slope {
			return a.d < b.d
		}
		switch a.v.region {
		case 1:
			return a.v.slope > b.v.slope
		case 2:
			return a.v.slope < b.v.slope
		case 3:
			return a.v.slope > b.v.slope
		case 4:
			return a.v.slope < b.v.slope
		}
		return true
	})

	fire := make(map[int]*asteroid)
	dedupe := make(map[int]struct{})
	for i := 1; i <= len(slops); {
		curSlop := math.NaN()
		for _, val := range slops {
			_, here := dedupe[val.id]
			if val.v.slope != curSlop && !here {
				dedupe[val.id] = struct{}{}
				fire[i] = val
				curSlop = val.v.slope
				i++
			}
		}
	}
	v := fire[200]
	fmt.Printf("%d : %v %v\n", 200, v.pt, v.v)
}

type asteroid struct {
	id int
	pt point
	v  vector
	d  int
}

func getStation(roids []*point) *point {
	var max *point
	maxView := 0
	for i, cur := range roids {
		neighbors := 0
		for j, prospect := range roids {
			if i != j {
				isViewable := true
				for k, check := range roids {
					if i != k && j != k && cur.inLine(check, prospect) {
						isViewable = false
						break
					}
				}
				if isViewable {
					neighbors++
				}
			}
		}
		if neighbors > maxView {
			maxView = neighbors
			max = cur
		}
	}
	return max
}

func (v vector) String() string {
	return fmt.Sprintf("%f(%d/%d) in %d", v.slope, v.rise, v.run, v.region)
}

func getRegion(run, rise int) int {
	if run >= 0 && rise < 0 {
		return 1
	} else if run >= 0 && rise >= 0 {
		return 2
	} else if run < 0 && rise >= 0 {
		return 3
	} else {
		return 4
	}
}

func newVector(run, rise int) vector {
	r := getRegion(run, rise)
	if run == 0 && rise < 0 {
		return vector{0, -1, math.Inf(1), r}
	} else if run == 0 {
		return vector{0, 1, math.Inf(1), r}
	} else if rise == 0 && run < 0 {
		return vector{-1, 0, 0, r}
	} else if rise == 0 {
		return vector{1, 0, 0, r}
	}
	return vector{run, rise, math.Abs(float64(rise) / float64(run)), r}
}

type vector struct {
	run, rise int
	slope     float64
	region    int
}

func (a *point) inLine(b, c *point) bool {
	return ((b.x-a.x)*(c.y-a.y) == (c.x-a.x)*(b.y-a.y)) && a.inBetween(b, c)
}

func (a *point) inBetween(b, c *point) bool {
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

func (p point) slope(q point) vector {
	xD := q.x - p.x
	yD := q.y - p.y
	return newVector(xD, yD)
}

func (p point) distance(q point) int {
	return abs(p.x-q.x) + abs(p.y-q.y)
}

func (a point) String() string {
	return fmt.Sprintf("{%d,%d}", a.x, a.y)
}

type point struct {
	x, y int
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
