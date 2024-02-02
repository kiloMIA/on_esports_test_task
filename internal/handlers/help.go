package handlers

import "github.com/bwmarrin/discordgo"

func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
    helpMessage := "Available commands:\n" +
        "!ping - Responds with 'Pong!'\n" +
        "!pong - Responds with 'Ping!'\n" +
        "!help - Lists all available commands.\n" +
        "!poll \"Question\" \"Option1\" \"Option2\" ... - Creates a new poll with specified question and options.\n" +
        "!weather <Location> - Shows current weather for the specified location.\n" +
        "!translate <LanguageCode> <Text> - Translates the given text into the specified language.\n" +
        "!remindme <Time> <Message> - Sets a reminder with the specified message after the given time duration."
    s.ChannelMessageSend(m.ChannelID, helpMessage)
}