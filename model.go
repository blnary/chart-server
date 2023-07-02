package main

type Song struct {
	ID       uint    `json:"id" gorm:"primarykey"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	BPM      float64 `json:"bpm"`
	Offset   float64 `json:"offset"`
	Charts   []Chart `json:"charts"`
}

type Chart struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Content string `json:"content"`
	SongID  uint   `json:"song_id"`
	Song    Song   `json:"song"`
}

type Level struct {
	Name           string  `json:"name"`
	BPM            float64 `json:"bpm"`
	ID             int     `json:"id"`
	Offset         int     `json:"offset"`
	StartPos       int     `json:"startpos"`
	HardStartPos   int     `json:"hardStartpos"`
	EndPos         int     `json:"endpos"`
	HardEndPos     int     `json:"hardEndpos"`
	AudioID        int     `json:"audioId"`
	Difficulty     float64 `json:"difficulty"`
	SampleRate     int     `json:"sampleRate"`
	DifficultyLine []Point `json:"difficultyLine"`
	Notes          []Note  `json:"notes"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Note struct {
	ID int `json:"id"`
	P  int `json:"p"`
	D  int `json:"d"`
	S  int `json:"s"`
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

type GetBPMReply struct {
	Success bool    `json:"success"`
	BPM     float64 `json:"bpm"`
	Offset  float64 `json:"offset"`
	Msg     string  `json:"msg"`
}

type InternalGetBPMReply struct {
	BPM    float64 `json:"bpm"`
	Offset float64 `json:"offset"`
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
