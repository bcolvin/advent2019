package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	var instRaw []int
	for i, v := range strings.Split(string(b), ",") {
		a, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Invalid value at %d: %s\n", i, v)
		}
		instRaw = append(instRaw, a)
	}
	parseOps(instRaw)
}

func parseOps(vals []int) []op {
	var ops []op
	for instPtr := 0; instPtr < len(vals); {
		op := newOp(instPtr, vals)
		op.execute(&instPtr, vals)
		//fmt.Println(op)
		ops = append(ops, op)
	}
	return ops
}

func newOp(instPtr int, vals []int) op {
	opCode := vals[instPtr]
	action := opCode % 100
	paramCnt := 0

	switch action {
	case ADD: //add
		paramCnt = 3
	case MULT: //mult
		paramCnt = 3
	case SCAN: //scan
		paramCnt = 1
	case PRINT: //print
		paramCnt = 1
	case JUMPT:
		paramCnt = 2
	case JUMPF:
		paramCnt = 2
	case LT:
		paramCnt = 3
	case EQ:
		paramCnt = 3
	case EXIT:
	default:
		return op{err: fmt.Errorf("Unsupported OpCode: %d\n", opCode)}
	}
	return op{
		opCode,
		action,
		getParams(paramCnt, instPtr+1, vals, opCode/100),
		nil,
	}
}

func getParams(paramCnt int, instPtr int, vals []int, lPtrs int) []param {
	params := make([]param, paramCnt)
	for i, ptr := 0, lPtrs; i < paramCnt; i, ptr = i+1, ptr/10 {
		params[i] = param{
			vals[instPtr+i],
			ptr%10 == 0,
			0,
		}
	}
	return params
}

func readMemory(memoryAddress int, vals []int) (int, error) {
	if memoryAddress >= len(vals) {
		return 0, fmt.Errorf("invalid memory address %d", memoryAddress)
	}
	res := vals[memoryAddress]
	return res, nil
}

func writeMemory(value int, memoryAddress int, vals []int) error {
	if memoryAddress >= len(vals) {
		return fmt.Errorf("invalid memory address %d", memoryAddress)
	}
	vals[memoryAddress] = value
	return nil
}

func (o op) String() string {
	var ret string
	switch o.action {
	case ADD:
		ret = "ADD "
	case MULT: //mult
		ret = "MULT "
	case SCAN: //scan
		ret = "SCAN "
	case PRINT: //print
		ret = "PRINT "
	case JUMPT:
		ret = "JUMP IF TRUE "
	case JUMPF:
		ret = "JUMP IF FALSE "
	case LT:
		ret = "LESS THAN "
	case EQ:
		ret = "EQUALS "
	case EXIT:
		ret = "EXIT"
	}
	if len(o.params) > 0 {
		for _, v := range o.params {
			ret += fmt.Sprintf("%s ", v)
		}
	}
	return ret
}

func (o op) execute(instPtr *int, vals []int) {
	switch o.action {
	case ADD:
		o.arithmatic(func(a, b int) int { return a + b }, instPtr, vals)
	case MULT: //mult
		o.arithmatic(func(a, b int) int { return a * b }, instPtr, vals)
	case SCAN: //scan
		o.scan(instPtr, vals)
	case PRINT: //print
		o.print(instPtr, vals)
	case JUMPT:
		o.jump(func(i int) bool { return i > 0 }, instPtr, vals)
	case JUMPF:
		o.jump(func(i int) bool { return i == 0 }, instPtr, vals)
	case LT:
		o.comparison(func(a, b int) bool { return a < b }, instPtr, vals)
	case EQ:
		o.comparison(func(a, b int) bool { return a == b }, instPtr, vals)
	case EXIT:
		os.Exit(0)
	}
}

func (o op) arithmatic(f func(int, int) int, instPtr *int, vals []int) {
	if len(o.params) != 3 {
		o.err = fmt.Errorf("not enough params")
		return
	}

	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, o.err = readMemory(p1.value, vals)
	}
	p2 := &o.params[1]
	p2.reg = p2.value
	if p2.isPtr {
		p2.reg, o.err = readMemory(p2.value, vals)
	}
	res := &o.params[2]
	res.reg = f(p1.reg, p2.reg)
	writeMemory(res.reg, res.value, vals)
	*instPtr += 4
}

func (o op) comparison(f func(int, int) bool, instPtr *int, vals []int) {
	if len(o.params) != 3 {
		o.err = fmt.Errorf("not enough params")
		return
	}

	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, o.err = readMemory(p1.value, vals)
	}
	p2 := &o.params[1]
	p2.reg = p2.value
	if p2.isPtr {
		p2.reg, o.err = readMemory(p2.value, vals)
	}
	res := &o.params[2]
	res.reg = 0
	if f(p1.reg, p2.reg) {
		res.reg = 1
	}
	writeMemory(res.reg, res.value, vals)
	*instPtr += 4
}

func (o op) scan(instPtr *int, vals []int) {
	if len(o.params) != 1 {
		o.err = fmt.Errorf("not enough params")
		return
	}
	fmt.Printf("Enter value: ")
	var res int
	_, o.err = fmt.Scanf("%d\n", &res)
	writeMemory(res, o.params[0].value, vals)
	*instPtr += 2
}

func (o op) print(instPtr *int, vals []int) {
	if len(o.params) != 1 {
		o.err = fmt.Errorf("not enough params")
		return
	}
	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, o.err = readMemory(p1.value, vals)
	}

	fmt.Println(o.params[0].reg)
	*instPtr += 2
}

func (o op) jump(f func(int) bool, instPtr *int, vals []int) {
	if len(o.params) != 2 {
		o.err = fmt.Errorf("not enough params")
		return
	}
	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, o.err = readMemory(p1.value, vals)
	}
	if f(p1.reg) {
		p2 := &o.params[1]
		p2.reg = p2.value
		if p2.isPtr {
			p2.reg, o.err = readMemory(p2.value, vals)
		}
		*instPtr = p2.reg
	} else {
		*instPtr += 3
	}
}

type op struct {
	opCode int
	action int
	params []param
	err    error
}

func (p param) String() string {
	if p.isPtr {
		return fmt.Sprintf("%d(&%d)", p.reg, p.value)
	}
	return fmt.Sprintf("%d", p.reg)
}

type param struct {
	value int
	isPtr bool
	reg   int
}

const (
	MV    int = 0
	ADD   int = 1
	MULT  int = 2
	SCAN  int = 3
	PRINT int = 4
	JUMPT int = 5
	JUMPF int = 6
	LT    int = 7
	EQ    int = 8
	EXIT  int = 99
)
