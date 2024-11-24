package main

import "time"

type message struct {
	Name      string
	Text      string
	AvatarURL string
	Timestamp time.Time
}
