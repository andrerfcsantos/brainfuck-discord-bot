package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

func helpCommand(args ...string) (*dgo.MessageEmbed, error) {
	return &dgo.MessageEmbed{
		Title: "Brainfuck Bot Help",
		Fields: []*dgo.MessageEmbedField{
			{Name: "Usage", Value: "`!bf <command> [arguments]`", Inline: false},
			{
				Name: "Available commands",
				Value: "`!bf help` - Prints this message\n" +
					"`!bf exec [input] <program>` - Executes a brainfuck program\n" +
					"`!bf encode <target_output>` - Creates a Brainfuck program that outputs the characters in the target output\n" +
					"`!bf shorten <program>` - Creates a shorter version of the program. Aliases: `short`",
				Inline: false,
			},
		},
		Color: InfoColor,
		Type:  dgo.EmbedTypeArticle,
	}, nil
}
