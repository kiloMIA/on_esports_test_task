package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandlePoll(s *discordgo.Session, m *discordgo.MessageCreate) {
    content := strings.TrimPrefix(m.Content, "!poll ")
    if content == "" || !strings.Contains(content, "|") {
        s.ChannelMessageSend(m.ChannelID, "Usage: !poll Question | Option 1 | Option 2 | Option 3...")
        return
    }

    parts := strings.Split(content, "|")
    if len(parts) < 3 { 
        s.ChannelMessageSend(m.ChannelID, "You must include a question and at least two options.")
        return
    }

    question := strings.TrimSpace(parts[0])
    options := parts[1:]
    for i, option := range options {
        options[i] = strings.TrimSpace(option)
    }

    pollMessage := fmt.Sprintf("**%s**\n", question)
    emojis := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣"} 
    if len(options) > len(emojis) {
        s.ChannelMessageSend(m.ChannelID, "Currently, only up to 9 options are supported.")
        return
    }

    for i, option := range options {
        pollMessage += fmt.Sprintf("%s %s\n", emojis[i], option)
    }

    sentMessage, err := s.ChannelMessageSend(m.ChannelID, pollMessage)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to send poll message.")
        return
    }

    for i := 0; i < len(options) && i < len(emojis); i++ {
        s.MessageReactionAdd(m.ChannelID, sentMessage.ID, emojis[i])
    }
}