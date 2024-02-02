package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleTranslate(s *discordgo.Session, m *discordgo.MessageCreate) {
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