package handlers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kiloMIA/on_esports_test_task/internal/models"
)

var (
	reminders []models.Reminder
	mu        sync.Mutex 
)

func cleanupReminders() {
	mu.Lock()
	defer mu.Unlock()

	currentTime := time.Now()
	var activeReminders []models.Reminder
	for _, reminder := range reminders {
		if reminder.Time.After(currentTime) {
			activeReminders = append(activeReminders, reminder)
		}
	}
	reminders = activeReminders
}

func HandleRemindMe(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	reminder := models.Reminder{
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
		Message:   message,
		Time:      reminderTime,
	}

	mu.Lock()
	reminders = append(reminders, reminder)
	mu.Unlock()

	time.AfterFunc(duration, func() {
		s.ChannelMessageSend(reminder.ChannelID, fmt.Sprintf("<@%s>, Reminder: %s", reminder.UserID, reminder.Message))
		cleanupReminders()
	})

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Reminder set for %s! I'll remind you to: %s", reminderTime.Format("Mon Jan 2 15:04:05 MST 2006"), message))
}