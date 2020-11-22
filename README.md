# Brainfuck Discord Bot

Discord bot that can execute [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) programs.

## Use the bot in your Discord server

The easiest way to start using the bot is to invite it to your discord server via this link:

- [Invite Brainfuck Bot](https://discord.com/oauth2/authorize?client_id=779135765031813130&permissions=125952&scope=bot)

## Usage

`!bf <command> [arguments]`

## Available commands

* `help` - Prints a help message

* `exec [input] <program>` - Executes a brainfuck program

* `encode <target_output>` - Creates a Brainfuck program that outputs the characters in the target output

* `shorten <program>` - Creates a shorter version of the program. Aliases: `short`


## Examples

### Simple "Hello World"

This is a short Hello World program.

![Brainfuck bot executing the a Hello World program](https://media.discordapp.net/attachments/246378961603526666/779766384975544360/discord_bot.png)

### Passing inputs

In this example, the numbers 67 and 68 are given as input to the program. The program stores these values and then outputs them. 

![Brainfuck bot executing a program with input](https://media.discordapp.net/attachments/737687180331319459/779767672690704394/discord_bot_input.png)

### Encoding outputs

You can also generate programs that produce a desired output.

![Brainfuck bot encoding output into a Brainfuck program](https://media.discordapp.net/attachments/246378961603526666/779818344877391872/unknown.png)


### Shorten programs

Get a shorter version of programs. For some programs, this shorter version can help readability.

![Brainfuck bot encoding output into a Brainfuck program](https://media.discordapp.net/attachments/737687180331319459/780206837353545728/unknown.png)
