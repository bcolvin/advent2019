package main

import (
	"fmt"
	"io"
	"log"
)

func main() {
	instrs, err := readProgram("input.txt")
	if err != nil {
		panic(err)
	}
	prog := newProgram(1, instrs)
	go func(prog program) {
		err = prog.execute()
		if err != nil && err != io.EOF {
			log.Fatalf("BOOM %v\n", err)
		}
	}(prog)

	ship := newShip()
	start := ship.getPanel(point{0, 0})
	robot := newRobot(&prog, &ship, start)

	for {
		err := robot.sendStatus()
		if err != nil {
			fmt.Errorf("Error sending status %v\n", err)
			break
		}

		err = robot.getColorAndPaint()
		if err != nil {
			fmt.Errorf("Error painting %v\n", err)
			break
		}

		err = robot.getDirectionAndMove()
		if err != nil {
			fmt.Errorf("Error reading direction %v\n", err)
			break
		}
	}
	fmt.Println(robot.painted)
}

func (r *robot) sendStatus() error {
	fmt.Printf("status: color %d %v\n", r.panel.color, r.panel.point)
	r.p.sendInput(r.panel.color)
	return nil
}

func (r *robot) getColorAndPaint() error {
	val, err := r.p.getOutput()
	if err == nil {
		err = r.paint(val)
	}
	return err
}

func (r *robot) paint(color int) error {
	switch color {
	case BLACK:
		r.panel.color = BLACK
	case WHITE:
		r.panel.color = WHITE
	default:
		return fmt.Errorf("Unknown color %d\n", color)
	}
	r.panel.coats++
	if r.panel.coats == 1 {
		r.painted++
	}
	return nil
}

func (r *robot) getDirectionAndMove() error {
	val, err := r.p.getOutput()
	if err == nil {
		err = r.move(val)
	}
	return err
}

func (r *robot) move(direction int) error {
	err := r.changeDirection(direction)
	if err != nil {
		return err
	}
	var newPos point
	switch r.direction {
	case W:
		newPos = point{r.panel.x, r.panel.y + 1}
	case S:
		newPos = point{r.panel.x, r.panel.y - 1}
	case A:
		newPos = point{r.panel.x - 1, r.panel.y}
	case D:
		newPos = point{r.panel.x + 1, r.panel.y}
	}
	r.panel = r.home.getPanel(newPos)
	return nil
}

func (r *robot) changeDirection(rotate int) error {
	switch rotate {
	case LEFT:
		rotate = 270
	case RIGHT:
		rotate = 90
	default:
		return fmt.Errorf("Unknown direction %d\n", rotate)
	}
	newDir := r.direction + rotate
	if newDir >= 360 {
		newDir -= 360
	}
	r.direction = newDir
	return nil
}

func newRobot(p *program, s *ship, start *panel) robot {
	return robot{
		p,
		s,
		W,
		start,
		0,
	}
}

type robot struct {
	p         *program
	home      *ship
	direction int
	panel     *panel
	painted   int
}

func (s ship) getPanel(pt point) *panel {
	pa := s[pt]
	if pa == nil {
		pa = &panel{pt, BLACK, 0}
		s[pt] = pa
	}
	return pa
}

func newShip() ship {
	return make(map[point]*panel)
}

type ship map[point]*panel

type panel struct {
	point
	color int
	coats int
}

type point struct {
	x, y int
}

const (
	BLACK int = 0
	WHITE int = 1
	LEFT  int = 0
	RIGHT int = 1
	W     int = 0
	S     int = 180
	A     int = 270
	D     int = 90
)
