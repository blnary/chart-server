package main

import (
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	Name     string
	Location string
	Charts   []Chart
}

type Chart struct {
	gorm.Model
	Content string
	SongID  uint
}

type PostSongReply struct {
	Success  bool
	Location string
	Msg      string
}
