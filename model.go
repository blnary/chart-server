package main

import (
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Charts   []Chart `json:"charts"`
}

type Chart struct {
	gorm.Model
	Content string `json:"content"`
	SongID  uint   `json:"song_id"`
	Song    Song   `json:"song"`
}

type GetChartsReply struct {
	Success bool    `json:"success"`
	Charts  []Chart `json:"charts"`
	Msg     string  `json:"msg"`
}

type PostSongsReply struct {
	Success bool   `json:"success"`
	SongID  uint   `json:"song_id"`
	Msg     string `json:"msg"`
}

type PostChartsRequest struct {
	Content string `json:"content"`
	SongID  uint   `json:"song_id"`
}

type PostChartsReply struct {
	Success bool   `json:"success"`
	ChartID uint   `json:"chart_id"`
	Msg     string `json:"msg"`
}

type PostChartRequest struct {
	Content string `json:"content"`
}

type PostChartReply struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type DeleteChartReply struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type DeleteSongReply struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}
