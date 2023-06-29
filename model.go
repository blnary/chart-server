package main

type Song struct {
	ID       uint    `json:"id" gorm:"primarykey"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Charts   []Chart `json:"charts"`
}

type Chart struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Content string `json:"content"`
	SongID  uint   `json:"song_id"`
	Song    Song   `json:"song"`
}

type GeneralReply struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type GetSongsReply struct {
	Success bool   `json:"success"`
	Songs   []Song `json:"songs"`
	Msg     string `json:"msg"`
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
