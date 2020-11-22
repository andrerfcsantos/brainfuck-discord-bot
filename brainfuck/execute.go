package brainfuck

import (
	"fmt"
	"strings"
)

// Max Memory cells a Brainfuck program is allowed to use
const MaxMemory = 30_000

// Max number of instructions that will be executed before giving up.
// This is a protection against infinite loops or programs that will
// run for very long otherwise.
const MaxExecInstructions = 10_000_000

// InstructionType represents a Brainfuck instruction
type InstructionType uint8

const (
	// No operation
	Nop InstructionType = iota
	// Increment data pointer operation ('>')
	IncrementDataPointer
	// Decrement data pointer operation or ('<')
	DecrementDataPointer
	// Increment data at current Memory pointer ('+')
	IncrementData
	// Decrement data at current Memory pointer ('-')
	DecrementData
	// Output data at current Memory pointer ('.')
	Output
	// Fetch input and store it a current Memory pointer (',')
	Input
	// Jump forward if value at current Memory address is 0 ('[')
	JmpForwardIfEqZero
	// Jump backwards if value at current Memory address is 0 (']')
	JmpBackwardsIfEqNotZero
	// End the program
	End
)

// Instruction represents an Brainfuck program instruction
type Instruction struct {
	InstructionType
	Value int
}

// Program is a compiled Brainfuck program that can be ran
type Program struct {
	Source       string
	Instructions []Instruction
	Memory       map[int]int8
}

// GetMemValue gets the current value in memory at the given address.
// If the address was never visited before, it will be initialized with 0 (zero).
func (p *Program) GetMemValue(addr int) int8 {
	var v int8
	var ok bool

	if v, ok = p.Memory[addr]; !ok {
		p.Memory[addr] = 0
		return 0
	}

	return v
}

// GetMemValue increments a value in memory at the given address.
// If the address was never visited before, it will be initialized with 0 (zero)
// and then the increment will be applied.
func (p *Program) IncMemValue(addr int, inc int8) {
	oldVal := p.GetMemValue(addr)
	p.Memory[addr] = oldVal + inc
}

// DecMemValue decrements a value in memory at the given address.
// If the address was never visited before, it will be initialized with 0 (zero)
// and then the decrement will be applied.
func (p *Program) DecMemValue(addr int, dec int8) {
	oldVal := p.GetMemValue(addr)
	p.Memory[addr] = oldVal - dec
}

// SetMemValue sets the memory value at the given address to the value specified.
func (p *Program) SetMemValue(addr int, val int8) {
	p.Memory[addr] = val
}

// Execute executes the Brainfuck program returning *ExecutionResult that contains the output
// amongst other stats about the execution.
// Inputs can optionally be given to Execute and will be used to feed the program when an input
// instruction happens (','). The number of input can be fewer than the number of input instructions,
// in which case case the inputs will be fed in a cyclic manner.
// Giving no inputs to a program that has
func (p *Program) Execute(inputs ...int) (*ExecutionResult, error) {
	p.Memory = make(map[int]int8)

	var out strings.Builder

	programSize := len(p.Instructions)
	nInputs, currInput := len(inputs), 0
	insExec := 0
	currentMemSize := 0
	pc := 0
	ap := 0

	for pc < programSize &&
		insExec <= MaxExecInstructions &&
		currentMemSize <= MaxMemory {
		i := p.Instructions[pc]
		switch i.InstructionType {
		case Nop:
		case IncrementDataPointer:
			ap += i.Value
		case DecrementDataPointer:
			ap -= i.Value
		case IncrementData:
			p.IncMemValue(ap, int8(i.Value))
		case DecrementData:
			p.DecMemValue(ap, int8(i.Value))
		case Output:
			out.WriteRune(rune(p.GetMemValue(ap)))
		case Input:
			if nInputs == 0 {
				return nil, fmt.Errorf("there is an input instruction at position %v, but no inputs were given: please provide at least 1 input to this program", pc)
			}
			p.SetMemValue(ap, int8(inputs[currInput%nInputs]))
			currInput++
		case JmpForwardIfEqZero:
			if p.GetMemValue(ap) == 0 {
				pc = i.Value
				continue
			}
		case JmpBackwardsIfEqNotZero:
			if p.GetMemValue(ap) != 0 {
				pc = i.Value
				continue
			}
		case End:
			return &ExecutionResult{
				Output:               out.String(),
				InstructionsExecuted: insExec + 1,
				MemoryCellsUsed:      len(p.Memory),
			}, nil
		default:
		}
		pc++
		insExec++
		currentMemSize = len(p.Memory)
	}

	if insExec > MaxExecInstructions {
		return nil, fmt.Errorf("the program reached the maximum number of instructions allowed (%v) and so it was stopped", MaxExecInstructions)
	}

	if currentMemSize > MaxMemory {
		return nil, fmt.Errorf("the program reached the maximum number of Memory cells allowed (%v)", MaxMemory)
	}

	return &ExecutionResult{
		Output:               out.String(),
		InstructionsExecuted: insExec,
		MemoryCellsUsed:      len(p.Memory),
	}, nil
}

// ExecutionResult contains information about a successful execution of a Brainfuck program
type ExecutionResult struct {
	Output               string `json:"output"`
	InstructionsExecuted int    `json:"instructions_executed"`
	MemoryCellsUsed      int    `json:"memory_cells_used"`
}
