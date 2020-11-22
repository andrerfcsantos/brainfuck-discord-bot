package main

import (
	bf "brainfuck-discord-bot/brainfuck"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

func validateShortenArgs(args ...string) (bool, error) {
	n := len(args)
	if n != 2 {
		return false, fmt.Errorf("wrong number of arguments to shorten: expected 1 `shorten <program>`, but got %v", n-1)
	}
	return true, nil
}

func shortenCommand(args ...string) (*dgo.MessageEmbed, error) {
	var err error
	var ok bool

	if ok, err = validateShortenArgs(args...); !ok {
		return &dgo.MessageEmbed{
			Title:       "Invalid number of arguments",
			Description: err.Error(),
			Color:       ErrorColor,
			Type:        dgo.EmbedTypeArticle,
		}, err
	}

	program := args[1]
	shortened := bf.Shorten(program)

	return &dgo.MessageEmbed{
		Color: SuccessColor,
		Fields: []*dgo.MessageEmbedField{
			{Name: "Original program", Value: program, Inline: false},
			{Name: "Short version", Value: shortened, Inline: false},
		},
		Type: dgo.EmbedTypeArticle,
	}, err

}
