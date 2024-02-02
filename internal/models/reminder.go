package models

import "time"

type Reminder struct {
    UserID    string    
    ChannelID string    
    Message   string    
    Time      time.Time 
}