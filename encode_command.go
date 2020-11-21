package main

import (
	bf "brainfuck-discord-bot/brainfuck"
	"fmt"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
)

func validateEncodeArgs(args ...string) (bool, error) {
	n := len(args)
	if n == 1 {
		return false, fmt.Errorf("wrong number of arguments to encode: expected 1 `encode <desired_output>`, but got none")
	}
	return true, nil
}

func encodeCommand(args ...string) (*dgo.MessageEmbed, error) {
	var err error
	var ok bool

	if ok, err = validateExecArgs(args...); !ok {
		return &dgo.MessageEmbed{
			Title:       "Invalid number of arguments",
			Description: err.Error(),
			Color:       ErrorColor,
			Type:        dgo.EmbedTypeArticle,
		}, err
	}

	textToEncode := strings.Join(args[1:], " ")
	bfProgram := bf.Encode(textToEncode)

	return &dgo.MessageEmbed{
		Color: SuccessColor,
		Fields: []*dgo.MessageEmbedField{
			{Name: "Target output", Value: textToEncode, Inline: false},
			{Name: "Brainfuck Program", Value: bfProgram, Inline: false},
		},
		Type: dgo.EmbedTypeArticle,
	}, err

}
