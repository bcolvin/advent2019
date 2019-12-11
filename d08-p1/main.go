package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	height, width := 6, 25
	layerSize := height * width
	runes := bytes.Runes(b)
	var layers []layer
	var wg sync.WaitGroup
	for i := 0; i < len(runes); i += layerSize {
		wg.Add(1)
		id := (i / layerSize) + 1
		l := layer{id: id, nums: make([]int, 10)}
		go l.readLayer(&wg, width, height, runes[i:])
		layers = append(layers, l)
	}
	wg.Wait()
	top, max := 0, 100000
	for _, l := range layers {
		if l.nums[0] < max {
			max = l.nums[0]
			top = l.nums[1] * l.nums[2]
		}
	}
	fmt.Println(top)
}

func (l layer) readLayer(wg *sync.WaitGroup, width, height int, runes []rune) {
	l.level = make([][]int, height)
	for i := 0; i < height; i++ {
		l.level[i] = make([]int, width)
		for j := 0; j < width; j++ {
			in := int(runes[(i*width)+j] - '0')
			//fmt.Printf("id:%d - (%d,%d)%d = %d\n",l.id,i,j,ind, in)
			l.level[i][j] = in
			l.nums[in]++
		}
	}
	wg.Done()
}

type layer struct {
	id    int
	level [][]int
	nums  []int
}
