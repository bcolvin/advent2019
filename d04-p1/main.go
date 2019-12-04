package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	str := strings.Split(string(b), "-")
	if len(str) != 2 {
		fmt.Errorf("Incorrect input format: range required")
	}
	start, err := strconv.Atoi(str[0])
	if err != nil {
		fmt.Errorf("First item must be a valid number, got %s", str[0])
	}
	end, err := strconv.Atoi(str[1])
	if err != nil {
		fmt.Errorf("Last item must be a valid number, got %s", str[1])
	}
	i := 0
	for ; start <= end; start++ {
		cur := strconv.Itoa(start)
		if passRules(cur) {
			i++
		}
	}
	fmt.Println(i)
}

func passRules(str string) bool {
	length := len(str)
	if length != 6 {
		return false
	}
	match := false
	b := []byte(str)
	for i := 1; i < length; i++ {
		if b[i-1] > b[i] {
			return false
		}
		if b[i-1] == b[i] {
			match = true
		}
	}
	return match
}
