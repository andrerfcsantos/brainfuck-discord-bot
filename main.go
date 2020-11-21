package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	dgo "github.com/bwmarrin/discordgo"
)

// Session handler
var session *dgo.Session

// Log file handler
var logFileHandler *os.File

// Configurations
var (
	bot_prefix string
	bot_token  string
)

func main() {
	var err error

	// Read configs
	err = readConfig()
	if err != nil {
		fmt.Printf("error reading config: %v", err)
		return
	}

	// Setup logger
	err = setupLogger()
	defer func() {
		if logFileHandler != nil {
			logFileHandler.Close()
		}
	}()

	if err != nil {
		fmt.Printf("error setting up logger: %v", err)
		return
	}

	// Try to connect and authenticate to discord
	bot_token = viper.GetString("bot_token")
	bot_prefix = viper.GetString("bot_prefix")

	session, err = dgo.New("Bot " + bot_token)
	if err != nil {
		log.Printf("error getting new session: %v", err)
		return
	}

	err = session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		return
	}
	defer session.Close()

	log.Infof("Started brainfuck bot. Using prefix %v with token starting in %v", bot_prefix, bot_token[:4])

	session.UpdateStatus(0, bot_prefix+" help")
	session.AddHandler(newMessageHandler)

	// Wait for a CTRL-C or other control signal to terminate
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	s := <-sc
	log.Infof("Stopping brainfuck bot due to signal %v received", s.String())
}

// newMessageHandler handles new messages received
func newMessageHandler(s *dgo.Session, m *dgo.MessageCreate) {

	if !strings.HasPrefix(m.Content, bot_prefix) {
		return
	}

	args := ParseCommand(m.Content)
	fmt.Printf("args: %#v\n", args)

	// Check if the message is intended for this bot
	if len(args) == 0 || args[0] != bot_prefix {
		return
	}

	if m.Content == bot_prefix {
		// User called the bot but didn't specify a command,
		// assume help command
		args = []string{bot_prefix, "help"}
	}

	var outMessage *dgo.MessageEmbed
	var err, sendErr error

	switch args[1] {
	case "help":
		outMessage, err = helpCommand(args[1:]...)
	case "exec":
		outMessage, err = execCommand(args[1:]...)
	default:
		err = fmt.Errorf("Command **%v** does not exist: type `%v help` to see the list of available commands", args[1], bot_prefix)
		outMessage = &dgo.MessageEmbed{
			Title:       fmt.Sprintf("Command **%v** does not exist", args[1]),
			Description: err.Error(),
			Color:       ErrorColor,
			Type:        dgo.EmbedTypeArticle,
		}
	}

	_, sendErr = s.ChannelMessageSendEmbed(m.ChannelID, outMessage)

	log.WithFields(log.Fields{
		"guild":           m.GuildID,
		"author_id":       m.Author.ID,
		"author_username": m.Author.Username,
		"raw_command":     m.Content,
		"process_error":   err,
		"send_error":      sendErr,
		"out_title":       outMessage.Title,
		"out_description": outMessage.Description,
		"out_fields":      outMessage.Fields,
	}).Info("command received")
}

func setupLogger() error {
	var err error

	log.SetFormatter(&log.JSONFormatter{})

	if runtime.GOOS == "windows" {
		logFileHandler, err = os.OpenFile("./requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		err = os.MkdirAll("/var/log/brainfuck-bot", 0755)
		if err != nil {
			return err
		}

		logFileHandler, err = os.OpenFile("/var/log/brainfuck-bot/requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	if err != nil {
		return err
	}

	log.SetOutput(logFileHandler)

	return nil
}
