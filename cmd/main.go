package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kiloMIA/on_esports_test_task/internal/handlers"
)

func main() {
    token := os.Getenv("DISCORD_BOT_TOKEN")
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        fmt.Println("Error creating Discord session,", err)
        return
    }

    dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
        if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "!") {
            return
        }

        command := strings.TrimSpace(m.Content[1:])
        switch command {
        case "help":
            handlers.HandleHelp(s, m)
        case "poll":
            handlers.HandlePoll(s, m)
        case "weather":
            handlers.HandleWeather(s, m)
        case "translate":
            handlers.HandleTranslate(s, m)
        case "remindme":
            handlers.HandleRemindMe(s, m)
        default:
            s.ChannelMessageSend(m.ChannelID, "Unknown command. Use !help to list all commands.")
        }
    })

    err = dg.Open()
    if err != nil {
        fmt.Println("Error opening connection,", err)
        return
    }

    fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    dg.Close()
}