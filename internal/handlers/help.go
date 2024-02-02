package handlers

import "github.com/bwmarrin/discordgo"

func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
    helpMessage := "Available commands:\n" +
        "!help - Lists all available commands.\n" +
        "!poll \"Question\" | \"Option1\" | \"Option2\" ... - Creates a new poll with specified question and options. \n Example Usage: !poll What's your favorite color? | Red | Blue | Green \n" +
        "!weather <Location> - Shows current weather for the specified location.\n Example Usage: !weather New York \n" +
        "!remindme <Time> <Message> - Sets a reminder with the specified message after the given time duration.\n Example Usage: !remindme 1m Take a break!"
    s.ChannelMessageSend(m.ChannelID, helpMessage)
}