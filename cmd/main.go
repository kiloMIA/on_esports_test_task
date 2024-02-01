package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate)

var Commands = map[string]CommandHandler{
	"ping":    handlePing,
	"pong":    handlePong,
	"help":    handleHelp,
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    if !strings.HasPrefix(m.Content, "!") {
        return
    }

    command := strings.TrimSpace(m.Content[1:])
    if handler, found := Commands[command]; found {
        go handler(s, m)
    } else {
        s.ChannelMessageSend(m.ChannelID, "Unknown command. Use !help to list all commands.")
    }
}

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func handlePong(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Ping!")
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := "Available commands:\n" +
		"!ping - Responds with 'Pong!'\n" +
		"!pong - Responds with 'Ping!'\n" +
		"!help - Lists all available commands."
	s.ChannelMessageSend(m.ChannelID, helpMessage)
}
