package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type WeatherApiResponse struct {
    Location struct {
        Name string `json:"name"`
    } `json:"location"`
    Current struct {
        TempC float64 `json:"temp_c"`
        Condition struct {
            Text string `json:"text"`
        } `json:"condition"`
    } `json:"current"`
}

type Poll struct {
    ID       string 
    Question string 
    Options  []string 
    Votes    map[string]int 
}

type Reminder struct {
    UserID    string    
    ChannelID string    
    Message   string    
    Time      time.Time 
}

type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate)

var Commands = map[string]CommandHandler{
    "ping":    handlePing,
    "pong":    handlePong,
    "help":    handleHelp,
    "poll":    handlePoll, 
	"weather": handleWeather,
	"translate": handleTranslate,
	"remindme": handleRemindMe,
}

var polls = make(map[string]*Poll)
var reminders []Reminder

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

func handlePoll(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Split(m.Content, "\"") 
    if len(args) < 3 {
        s.ChannelMessageSend(m.ChannelID, "Usage: !poll \"Question\" \"Option1\" \"Option2\" ...")
        return
    }

    question := args[1]
    options := args[2:]

    poll := &Poll{
        ID:       generatePollID(), 
        Question: question,
        Options:  options,
        Votes:    make(map[string]int),
    }
    polls[poll.ID] = poll

    pollMessage := fmt.Sprintf("Poll ID: %s\nQuestion: %s\nOptions:\n", poll.ID, poll.Question)
    for i, option := range poll.Options {
        pollMessage += fmt.Sprintf("%d: %s\n", i+1, option)
    }
    s.ChannelMessageSend(m.ChannelID, pollMessage)
}

func generatePollID() string {
    return fmt.Sprintf("%d", len(polls)+1)
}

func handleWeather(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Usage: !weather <Location>")
        return
    }

    location := strings.Join(args[1:], " ")
    apiKey := os.Getenv("WEATHER_API_TOKEN")
    apiUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, url.QueryEscape(location))

    resp, err := http.Get(apiUrl)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to get weather information.")
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to read weather information.")
        return
    }

    var weatherResponse WeatherApiResponse
    if err := json.Unmarshal(body, &weatherResponse); err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to parse weather information.")
        return
    }

    reply := fmt.Sprintf("Weather in %s: %s, %f°C", weatherResponse.Location.Name, weatherResponse.Current.Condition.Text, weatherResponse.Current.TempC)
    s.ChannelMessageSend(m.ChannelID, reply)
}

func handleTranslate(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 3 {
        s.ChannelMessageSend(m.ChannelID, "Usage: !translate <LanguageCode> <Text>")
        return
    }

    languageCode := args[1]
    text := strings.Join(args[2:], " ")


    requestBody, err := json.Marshal(map[string]interface{}{
        "q":      text,
        "source": "auto", 
        "target": languageCode,
    })
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error preparing translation request.")
        return
    }

    apiUrl := "https://libretranslate.com/translate"
    req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBody))
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to create translation request.")
        return
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to get translation.")
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to read translation response.")
        return
    }

    var translateResponse struct {
        TranslatedText string `json:"translatedText"`
    }
    if err := json.Unmarshal(body, &translateResponse); err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to parse translation response.")
        return
    }

    if translateResponse.TranslatedText != "" {
        s.ChannelMessageSend(m.ChannelID, translateResponse.TranslatedText)
    } else {
        s.ChannelMessageSend(m.ChannelID, "Translation failed or no translations found.")
    }
}

func handleRemindMe(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 3 {
        s.ChannelMessageSend(m.ChannelID, "Usage: !remindme <Time> <Message>")
        return
    }

    duration, err := time.ParseDuration(args[1])
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Invalid time format. Please use formats like 10s, 2m, 1h.")
        return
    }

    message := strings.Join(args[2:], " ")
    reminderTime := time.Now().Add(duration)

    reminder := Reminder{
        UserID:    m.Author.ID,
        ChannelID: m.ChannelID,
        Message:   message,
        Time:      reminderTime,
    }
    reminders = append(reminders, reminder)

    go func(rem Reminder) {
        time.Sleep(duration)
        s.ChannelMessageSend(rem.ChannelID, fmt.Sprintf("<@%s>, Reminder: %s", rem.UserID, rem.Message))
    }(reminder)

    s.ChannelMessageSend(m.ChannelID, "Reminder set!")
}