package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	vals := permutations([]int{0, 1, 2, 3, 4})
	max := 0
	for _, v := range vals {
		res := 0
		for _, vv := range v {
			input := bytes.NewBufferString(fmt.Sprintf("%d\n%d\n", vv, res))
			var output bytes.Buffer
			p1 := newProgram(input, &output, instrs)
			err := p1.execute()
			if err != nil && err != io.EOF {
				log.Fatalln("Error running program! ahh %v", v)
			}
			_, err = fmt.Fscanf(&output, "%d", &res)
			if err != nil {
				log.Fatalln("Error reading program output %v", v)
			}
		}
		if res > max {
			max = res
		}
	}
	fmt.Println(max)
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func (p program) execute() error {
	for *p.instPtr < len(p.instructions) {
		op, err := p.createOp()
		if err != nil {
			return err
		}
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
			0,
		}
	}
	return params
}

func (p program) createOp() (op, error) {
	opCode := p.instructions[*p.instPtr]
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
		return op{}, fmt.Errorf("Unsupported OpCode: %d\n", opCode)
	}
	return op{
		opCode,
		action,
		p.getParams(paramCnt, *p.instPtr+1, opCode/100),
	}, nil
}

func (p program) runOp(o op) error {
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
	case EXIT:
		return io.EOF
	}
	return nil
}

func (p program) arithmatic(f func(int, int) int, o op) error {
	if len(o.params) != 3 {
		return fmt.Errorf("not enough params")
	}
	var err error
	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, err = p.readMemory(p1.value)
		if err != nil {
			return fmt.Errorf("Arithmatic: Unable to read param 1 from memory %v\n", err)
		}
	}
	p2 := &o.params[1]
	p2.reg = p2.value
	if p2.isPtr {
		p2.reg, err = p.readMemory(p2.value)
		if err != nil {
			return fmt.Errorf("Arithmatic: Unable to read param 2 from memory %v\n", err)
		}
	}
	res := &o.params[2]
	res.reg = f(p1.reg, p2.reg)
	err = p.writeMemory(res.reg, res.value)
	if err != nil {
		return fmt.Errorf("Arithmatic: Unable to write param to memory %v\n", err)
	}
	*p.instPtr += 4
	return nil
}

func (p program) comparison(f func(int, int) bool, o op) error {
	if len(o.params) != 3 {
		return fmt.Errorf("not enough params")
	}
	var err error
	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, err = p.readMemory(p1.value)
		if err != nil {
			return fmt.Errorf("Comparison: Unable to read param 1 from memory %v\n", err)
		}
	}
	p2 := &o.params[1]
	p2.reg = p2.value
	if p2.isPtr {
		p2.reg, err = p.readMemory(p2.value)
		if err != nil {
			return fmt.Errorf("Comparison: Unable to read param 2 from memory %v\n", err)
		}
	}
	res := &o.params[2]
	res.reg = 0
	if f(p1.reg, p2.reg) {
		res.reg = 1
	}
	err = p.writeMemory(res.reg, res.value)
	if err != nil {
		return fmt.Errorf("Comparison: Unable to write param to memory %v\n", err)
	}
	*p.instPtr += 4
	return nil
}

func (p program) scan(o op) error {
	if len(o.params) != 1 {
		return fmt.Errorf("not enough params")
	}
	var res int
	_, err := fmt.Fscanf(p.inputReader, "%d\n", &res)
	if err != nil {
		return fmt.Errorf("Scan: Unable to read input from stream %v\n", err)
	}
	err = p.writeMemory(res, o.params[0].value)
	if err != nil {
		return fmt.Errorf("Scan: Unable to write param to memory %v\n", err)
	}
	*p.instPtr += 2
	return nil
}

func (p program) print(o op) error {
	if len(o.params) != 1 {
		return fmt.Errorf("not enough params")
	}
	var err error
	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, err = p.readMemory(p1.value)
		if err != nil {
			return fmt.Errorf("Print: Unable to read param 1 from memory %v\n", err)
		}
	}

	_, err = fmt.Fprintln(p.outputWriter, o.params[0].reg)
	if err != nil {
		return fmt.Errorf("Print: Unable to write output to stream %v\n", err)
	}
	*p.instPtr += 2
	return nil
}

func (p program) jump(f func(int) bool, o op) error {
	if len(o.params) != 2 {
		return fmt.Errorf("not enough params")
	}
	var err error
	p1 := &o.params[0]
	p1.reg = p1.value
	if p1.isPtr {
		p1.reg, err = p.readMemory(p1.value)
		if err != nil {
			return fmt.Errorf("Jump: Unable to read param 1 from memory %v\n", err)
		}
	}
	if f(p1.reg) {
		p2 := &o.params[1]
		p2.reg = p2.value
		if p2.isPtr {
			p2.reg, err = p.readMemory(p2.value)
			if err != nil {
				return fmt.Errorf("Jump: Unable to read param 2 from memory %v\n", err)
			}
		}
		*p.instPtr = p2.reg
	} else {
		*p.instPtr += 3
	}
	return nil
}

func (p program) readMemory(memoryAddress int) (int, error) {
	if memoryAddress >= len(p.instructions) {
		return 0, fmt.Errorf("invalid memory address %d", memoryAddress)
	}
	res := p.instructions[memoryAddress]
	return res, nil
}

func (p program) writeMemory(value int, memoryAddress int) error {
	if memoryAddress >= len(p.instructions) {
		return fmt.Errorf("invalid memory address %d", memoryAddress)
	}
	p.instructions[memoryAddress] = value
	return nil
}

func newProgram(reader io.Reader, writer io.Writer, instructions []int) program {
	start := 0
	p := program{
		reader,
		writer,
		&start,
		make([]int, len(instructions)),
	}
	copy(p.instructions, instructions)
	return p
}

type program struct {
	inputReader  io.Reader
	outputWriter io.Writer
	instPtr      *int
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
