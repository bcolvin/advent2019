package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	height, width := 6, 25
	layerSize := height * width
	runes := bytes.Runes(b)
	layer := make([][]int, height)
	for i := 0; i < height; i++ {
		layer[i] = make([]int, width)
		for j := 0; j < width; j++ {
			layer[i][j] = 2
		}
	}
	done := false
	for i := 0; i < len(runes) && !done; i += layerSize {
		cur := runes[i:]
		done = true
		for j := 0; j < height; j++ {
			for k := 0; k < width; k++ {
				ind := (j * width) + k
				in := int(cur[ind] - '0')
				if layer[j][k] > 1 {
					layer[j][k] = in
					if in == 2 {
						done = false
					}
				}
			}
		}
	}
	for _, l := range layer {
		for _, v := range l {
			if v == 0 {
				fmt.Printf(" *")
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Println()
	}
}
