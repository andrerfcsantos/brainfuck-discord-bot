package brainfuck

import "fmt"

// Compile compiles a string representing a brainfuck program into a representation
// that can be executed (see the (*Program).Execute() method)
// The program string may contain characters that are not Brainfuck instructions -
// those characters will be ignored.
// The point of this compilation process is to optimize the instructions that are ran.
// The compilation will fail if the number of jumps forwards '[' does not match the number
// of jump ']' backwards. Otherwise, the compilation process is very lenient, and will not check
// for invalid Brainfuck memory addresses or pointers - such checks are runtime checks and not
// compile-time checks.
func Compile(program string) (*Program, error) {
	var p Program

	progRunes := []rune(program)
	n := len(progRunes)

	var nesting int
	var openBracketsStack []int

	for i := 0; i < n; i++ {
		var ins Instruction

		switch progRunes[i] {
		case '>':
			ins.InstructionType = IncrementDataPointer
			ins.Value = 1
			p.Instructions = append(p.Instructions, ins)
		case '<':
			ins.InstructionType = DecrementDataPointer
			ins.Value = 1
			p.Instructions = append(p.Instructions, ins)
		case '+':
			ins.InstructionType = IncrementData
			ins.Value = 1
			p.Instructions = append(p.Instructions, ins)
		case '-':
			ins.InstructionType = DecrementData
			ins.Value = 1
			p.Instructions = append(p.Instructions, ins)
		case '[':
			ins.InstructionType = JmpForwardIfEqZero
			p.Instructions = append(p.Instructions, ins)
			openBracketsStack = append(openBracketsStack, len(p.Instructions)-1)
			nesting++
		case ']':
			nesting--
			if nesting < 0 {
				return &p, fmt.Errorf("closing ']' at position %v has no matching '['", i)
			}

			// Pop matching open bracket position from stack
			openBracketPos := openBracketsStack[len(openBracketsStack)-1]
			openBracketsStack = openBracketsStack[:len(openBracketsStack)-1]

			// Set jump backwards to the instruction after the matching '['
			ins.InstructionType = JmpBackwardsIfEqNotZero
			ins.Value = openBracketPos + 1
			p.Instructions = append(p.Instructions, ins)

			// Go back to the instruction of '[' and set jump forward to the instruction after
			// the one we are seeing
			p.Instructions[openBracketPos].Value = len(p.Instructions)

		case '.':
			ins.InstructionType = Output
			p.Instructions = append(p.Instructions, ins)
		case ',':
			ins.InstructionType = Input
			p.Instructions = append(p.Instructions, ins)
		default:
		}

	}

	if nesting != 0 {
		return &p, fmt.Errorf("there are %v more '[' than ']': please make sure the number of [ and ] match", nesting)
	}

	return &p, nil
}
