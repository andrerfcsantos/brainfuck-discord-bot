package main

import (
	bf "brainfuck-discord-bot/brainfuck"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	dgo "github.com/bwmarrin/discordgo"
)

func validateExecArgs(args ...string) (bool, error) {
	n := len(args)
	if n == 1 || n > 3 {
		return false, fmt.Errorf("wrong number of arguments to exec: expected 1 `exec <program>` or 2 `exec [input] <program>`, but got %v", n-1)
	}
	return true, nil
}

func execCommand(args ...string) (*dgo.MessageEmbed, error) {
	if ok, err := validateExecArgs(args...); !ok {
		return &dgo.MessageEmbed{
			Title:       "Invalid number of arguments",
			Description: err.Error(),
			Color:       ErrorColor,
			Type:        dgo.EmbedTypeArticle,
		}, err
	}

	var p *bf.Program
	var err error

	nArgs := len(args)

	start := time.Now()
	if nArgs == 3 {
		p, err = bf.Compile(args[2])
	} else {
		p, err = bf.Compile(args[1])
	}
	elapsedCompilation := time.Now().Sub(start)

	if err != nil {
		return &dgo.MessageEmbed{
			Title:       "Compilation Error",
			Description: err.Error(),
			Color:       ErrorColor,
			Type:        dgo.EmbedTypeArticle,
		}, fmt.Errorf("compilation error: %v", err)
	}

	var out *bf.ExecutionResult
	var elapsedExecute time.Duration
	if nArgs == 3 {
		var inputs []int
		strInputs := strings.Split(args[1], ",")
		for _, inp := range strInputs {
			intInput, err := strconv.Atoi(inp)
			if err != nil {
				annotatedError := fmt.Errorf("could not parse input %v as int: %v", inp, err)
				return &dgo.MessageEmbed{
					Title:       "Input parsing error",
					Description: annotatedError.Error(),
					Color:       ErrorColor,
					Type:        dgo.EmbedTypeArticle,
				}, annotatedError
			}
			inputs = append(inputs, intInput)
		}

		start := time.Now()
		out, err = p.Execute(inputs...)
		elapsedExecute = time.Now().Sub(start)
	} else {
		start := time.Now()
		out, err = p.Execute()
		elapsedExecute = time.Now().Sub(start)
	}

	if err != nil {
		return &dgo.MessageEmbed{
			Title:       "Execution error",
			Description: err.Error(),
			Color:       ErrorColor,
			Type:        dgo.EmbedTypeArticle,
		}, fmt.Errorf("execution error: %v", err)
	}

	finalOutput := out.Output
	description := "Program ran successfully."

	// TODO: see better ways to check last rune
	if len(out.Output) > 0 && unicode.IsSpace(rune(out.Output[len(out.Output)-1])) {
		finalOutput += "<EOF>"
		description = "Program ran successfully. Since the output ends in whitespace, an explicit <EOF> was introduced for you"
	}

	if len(out.Output) == 0 {
		finalOutput = "No output"
		description = "Program ran successfully, but produced no output"
	}

	return &dgo.MessageEmbed{
		Title:       "Execution successful",
		Description: description,
		Color:       SuccessColor,
		Fields: []*dgo.MessageEmbedField{
			{Name: "Output", Value: finalOutput, Inline: false},
			{Name: "Compilation in", Value: elapsedCompilation.String(), Inline: true},
			{Name: "Execution in", Value: elapsedExecute.String(), Inline: true},
			{Name: "Total", Value: (elapsedCompilation + elapsedExecute).String(), Inline: true},
			{Name: "Cells used", Value: strconv.Itoa(out.MemoryCellsUsed), Inline: true},
			{Name: "Instructions", Value: strconv.Itoa(out.InstructionsExecuted), Inline: true},
		},
		Type: dgo.EmbedTypeArticle,
	}, nil
}
