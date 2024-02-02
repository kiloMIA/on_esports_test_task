package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kiloMIA/on_esports_test_task/internal/models"
	"github.com/bwmarrin/discordgo"
)

func HandleWeather(s *discordgo.Session, m *discordgo.MessageCreate) {
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

    var weatherResponse models.WeatherApiResponse
    if err := json.Unmarshal(body, &weatherResponse); err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to parse weather information.")
        return
    }

    reply := fmt.Sprintf("Weather in %s: %s, %f°C", weatherResponse.Location.Name, weatherResponse.Current.Condition.Text, weatherResponse.Current.TempC)
    s.ChannelMessageSend(m.ChannelID, reply)
}