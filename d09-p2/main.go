package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func readProgram(filename string) []int {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var instrs []int
	for i, v := range strings.Split(string(b), ",") {
		a, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Invalid value at %d: %s\n", i, v)
		}
		instrs = append(instrs, a)
	}
	return instrs
}

func main() {
	instrs := readProgram("input.txt")
	prog := newProgram(1, instrs)
	err := prog.execute()
	if err != nil && err != io.EOF {
		log.Fatalf("BOOM %v\n", err)
	}
}

func (p *program) execute() error {
	for {
		op, err := p.createOp()
		if err != nil {
			return err
		}
		//fmt.Println(op)
		if err := p.runOp(op); err != nil {
			return err
		}
	}
	return nil
}

func (p program) getParams(paramCnt int, offset int, lPtrs int) []param {
	params := make([]param, paramCnt)
	for i, ptr := 0, lPtrs; i < paramCnt; i, ptr = i+1, ptr/10 {
		params[i] = param{
			p.instructions[offset+i],
			ptr%10 == 0,
			ptr%10 == 2,
			0,
		}
	}
	return params
}

func (p program) createOp() (op, error) {
	opCode := p.instructions[p.instPtr]
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
	case BASE:
		paramCnt = 1
	case EXIT:
	default:
		return op{}, fmt.Errorf("Unsupported OpCode: %d\n", opCode)
	}
	return op{
		opCode,
		action,
		p.getParams(paramCnt, p.instPtr+1, opCode/100),
	}, nil
}

func (p *program) runOp(o op) error {
	switch o.action {
	case ADD:
		return p.arithmatic(func(a, b int) int { return a + b }, o)
	case MULT: //mult
		return p.arithmatic(func(a, b int) int { return a * b }, o)
	case SCAN: //scan
		return p.scan(o)
	case PRINT: //print
		return p.print(o)
	case JUMPT:
		return p.jump(func(i int) bool { return i > 0 }, o)
	case JUMPF:
		return p.jump(func(i int) bool { return i == 0 }, o)
	case LT:
		return p.comparison(func(a, b int) bool { return a < b }, o)
	case EQ:
		return p.comparison(func(a, b int) bool { return a == b }, o)
	case BASE:
		return p.moveBase(o)
	case EXIT:
		return io.EOF
	}
	return nil
}

func (p *program) arithmatic(f func(int, int) int, o op) error {
	if len(o.params) != 3 {
		return fmt.Errorf("not enough params")
	}
	var err error
	p1, err := p.getValue(&o.params[0])
	if err != nil {
		return err
	}
	p2, err := p.getValue(&o.params[1])
	if err != nil {
		return err
	}
	res := f(p1, p2)
	p.setValue(&o.params[2], res)
	if err != nil {
		return err
	}
	p.instPtr += 4
	return nil
}

func (p *program) comparison(f func(int, int) bool, o op) error {
	if len(o.params) != 3 {
		return fmt.Errorf("not enough params")
	}
	var err error
	p1, err := p.getValue(&o.params[0])
	if err != nil {
		return err
	}
	p2, err := p.getValue(&o.params[1])
	if err != nil {
		return err
	}
	res := 0
	if f(p1, p2) {
		res = 1
	}
	err = p.setValue(&o.params[2], res)
	if err != nil {
		return err
	}
	p.instPtr += 4
	return nil
}

func (p *program) scan(o op) error {
	if len(o.params) != 1 {
		return fmt.Errorf("not enough params")
	}
	var res int
	_, err := fmt.Fscanf(p.input, "%d\n", &res)
	if err != nil {
		return err
	}
	err = p.setValue(&o.params[0], res)
	if err != nil {
		return err
	}
	p.instPtr += 2
	return nil
}

func (p *program) print(o op) error {
	if len(o.params) != 1 {
		return fmt.Errorf("not enough params")
	}
	var err error
	val, err := p.getValue(&o.params[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(p.output, "%d\n", val)
	p.instPtr += 2
	return nil
}

func (p *program) jump(f func(int) bool, o op) error {
	if len(o.params) != 2 {
		return fmt.Errorf("not enough params")
	}
	var err error
	val, err := p.getValue(&o.params[0])
	if err != nil {
		return err
	}
	if f(val) {
		val, err = p.getValue(&o.params[1])
		if err != nil {
			return err
		}
		p.instPtr = val
	} else {
		p.instPtr += 3
	}
	return nil
}

func (p *program) moveBase(o op) error {
	if len(o.params) != 1 {
		return fmt.Errorf("not enough params")
	}
	val, err := p.getValue(&o.params[0])
	if err != nil {
		return err
	}
	p.relativeBase += val
	p.instPtr += 2
	return nil
}

func (p program) getOutput() (int, error) {
	var ret int
	_, err := fmt.Fscanf(p.output, "%d", &ret)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Output: %d\n", ret)
	return ret, nil
}

func (p program) sendInput(in int) {
	fmt.Printf("Input: %d\n", in)
	fmt.Fprintf(p.input, "%d", in)
}

func (p program) getValue(p1 *param) (int, error) {
	p1.reg = p1.value
	var err error
	if p1.isRel {
		p1.value = p1.value + p.relativeBase
		p1.isPtr = true
	}
	if p1.isPtr {
		p1.reg, err = p.readMemory(p1.value)
		if err != nil {
			return 0, fmt.Errorf("Unable to read param %v from memory %v\n", p1, err)
		}
	}
	return p1.reg, nil
}

func (p program) readMemory(memoryAddress int) (int, error) {
	if memoryAddress < len(p.instructions) && memoryAddress >= 0 {
		return p.instructions[memoryAddress], nil
	}
	return 0, fmt.Errorf("invalid memory address %d", memoryAddress)
}

func (p program) setValue(p1 *param, val int) error {
	p1.reg = val
	if p1.isRel {
		p1.value = p1.value + p.relativeBase
		p1.isPtr = true
	}
	err := p.writeMemory(p1.reg, p1.value)
	if err != nil {
		return fmt.Errorf("Unable to write param %v to memory %v\n", p, err)
	}
	return nil
}

func (p program) writeMemory(value int, memoryAddress int) error {
	if memoryAddress < len(p.instructions) && memoryAddress >= 0 {
		p.instructions[memoryAddress] = value
		return nil
	}
	return fmt.Errorf("invalid memory address %d", memoryAddress)
}

func newProgram(id int, instructions []int) program {
	inp := os.Stdin  //bytes.NewBuffer(make([]byte,100))
	out := os.Stdout //bytes.NewBuffer(make([]byte,100))
	p := program{
		id,
		inp,
		out,
		0,
		0,
		make([]int, len(instructions)*10),
	}
	copy(p.instructions, instructions)
	return p
}

type program struct {
	id           int
	input        io.ReadWriter
	output       io.ReadWriter
	instPtr      int
	relativeBase int
	instructions []int
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
	case BASE:
		ret = "BASE "
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

type op struct {
	opCode int
	action int
	params []param
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
	isRel bool
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
	BASE  int = 9
	EXIT  int = 99
)
