package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func GetCharts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// read charts
		var charts []Chart
		if err := db.Preload("Song").Find(&charts).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read charts: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// reply
		rep := &GetChartsReply{
			Success: true,
			Charts:  charts,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func GetSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// read songs
		var songs []Song
		if err := db.Find(&songs).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read songs: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// reply
		rep := &GetSongsReply{
			Success: true,
			Songs:   songs,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func GetBPM(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id
		id := c.Param("id")

		// read song
		var song Song
		if err := db.First(&song, id).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read song: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// try read from server
		if song.BPM == 0 {
			reply, err := getBPM(&song)
			if err != nil {
				rep := &GeneralReply{
					Success: false,
					Msg:     fmt.Sprintf("failed to calculate BPM: %v", err),
				}
				c.JSON(http.StatusInternalServerError, rep)
				return
			}
			song.BPM = reply.BPM
			song.Offset = reply.Offset
			if err := db.Save(&song).Error; err != nil {
				rep := &GeneralReply{
					Success: false,
					Msg:     fmt.Sprintf("failed to save updated song: %v", err),
				}
				c.JSON(http.StatusInternalServerError, rep)
				return
			}
		}

		// reply
		rep := &GetBPMReply{
			Success: true,
			BPM:     song.BPM,
			Offset:  song.Offset,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func GetSongFile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id
		id := c.Param("id")

		// get song from db
		var song Song
		if err := db.First(&song, id).Error; err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to find song from db: %v", err),
			)
		}

		// read file
		file, err := os.ReadFile(song.Location)
		if err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to read file: %v", err),
			)
			return
		}

		// reply
		if _, err := c.Writer.Write(file); err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to write file: %v", err),
			)
			return
		}
		c.Writer.Flush()
	}
}

func SetBPM(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the song ID from the request parameters
		id := c.Param("id")

		// Read the song from the database
		var song Song
		if err := db.First(&song, id).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read song: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// Parse the BPM and offset values from the request body
		var requestBody struct {
			BPM    float64 `json:"bpm"`
			Offset float64 `json:"offset"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to parse request body: %v", err),
			}
			c.JSON(http.StatusBadRequest, rep)
			return
		}

		// Update the song's BPM and offset
		song.BPM = requestBody.BPM
		song.Offset = requestBody.Offset

		// Save the updated song to the database
		if err := db.Save(&song).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save updated song: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// Reply with a success message
		rep := &GeneralReply{
			Success: true,
			Msg:     "BPM and offset successfully updated.",
		}
		c.JSON(http.StatusOK, rep)
	}
}

func SyncBPM(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the song ID from the request parameters
		id := c.Param("id")

		// Read the song from the database
		var song Song
		if err := db.First(&song, id).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read song: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// Check if song has charts
		if len(song.Charts) == 0 {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("song has no charts"),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// Read BPM from song's chart
		chart := song.Charts[0]
		var level Level
		if err := json.Unmarshal([]byte(chart.Content), &level); err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to get chart of BPM: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// Update the song's BPM and offset
		song.BPM = level.BPM
		song.Offset = float64(level.Offset) / 44.1

		// Save the updated song to the database
		if err := db.Save(&song).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save updated song: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// Reply with a success message
		rep := &GeneralReply{
			Success: true,
			Msg:     "BPM and offset successfully updated.",
		}
		c.JSON(http.StatusOK, rep)
	}
}

func PostCharts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// read chart
		var req PostChartsRequest
		if err := c.BindJSON(&req); err != nil {
			return
		}

		// save chart
		chart := &Chart{
			Content: req.Content,
			SongID:  req.SongID,
		}
		if err := db.Create(chart).Error; err != nil {
			rep := &PostChartsReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save chart to db: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// reply
		rep := &PostChartsReply{
			Success: true,
			ChartID: chart.ID,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func PostChart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// read chart
		var req PostChartRequest
		if err := c.BindJSON(&req); err != nil {
			return
		}

		// get param
		id := c.Param("id")

		// load chart from db
		var chart Chart
		if err := db.First(&chart, id).Error; err != nil {
			rep := &PostChartReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read chart from db: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// save chart
		chart.Content = req.Content
		if err := db.Save(chart).Error; err != nil {
			rep := &PostChartReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save chart to db: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// reply
		rep := &PostChartReply{
			Success: true,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func PostSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// read song
		file, err := c.FormFile("song")
		if err != nil {
			rep := &PostSongsReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read song: %v", err),
			}
			c.JSON(http.StatusBadRequest, rep)
			return
		}

		// save song file to disk
		base := filepath.Base(file.Filename)
		ext := filepath.Ext(file.Filename)
		dir := fmt.Sprintf("/songs/%s%s", uuid.NewV4().String(), ext)
		if err := c.SaveUploadedFile(file, dir); err != nil {
			rep := &PostSongsReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save song to %s: %v", dir, err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// save song info to db
		song := &Song{
			Name:     base,
			BPM:      0,
			Offset:   0,
			Location: dir,
		}
		if err := db.Create(song).Error; err != nil {
			rep := &PostSongsReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save song to database: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// reply
		rep := &PostSongsReply{
			Success: true,
			SongID:  song.ID,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func DeleteChart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id
		id := c.Param("id")

		// delete chart from db
		if err := db.Delete(id).Error; err != nil {
			rep := &DeleteChartReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to delete chart from db: %v", err),
			}
			c.JSON(http.StatusBadRequest, rep)
			return
		}

		// reply
		rep := &DeleteChartReply{
			Success: true,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func DeleteSong(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id
		id := c.Param("id")

		// delete related charts
		if err := db.Where("song_id = ?", id).Delete(&Chart{}).Error; err != nil {
			rep := &DeleteSongReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to delete related charts from db: %v", err),
			}
			c.JSON(http.StatusBadRequest, rep)
			return
		}

		// delete song from db
		if err := db.Delete(&Song{}, id).Error; err != nil {
			rep := &DeleteSongReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to delete song from db: %v", err),
			}
			c.JSON(http.StatusBadRequest, rep)
			return
		}

		// reply
		rep := &DeleteSongReply{
			Success: true,
		}
		c.JSON(http.StatusOK, rep)
	}
}

func TestBPM(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// read song
		var songs []Song
		if err := db.Preload("Charts").Find(&songs).Error; err != nil {
			rep := &GeneralReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to read songs: %v", err),
			}
			c.JSON(http.StatusInternalServerError, rep)
			return
		}

		// start testing each songs
		bpmCorrect := 0
		testTotal := 0
		offsetCorrect := 0
		for _, song := range songs {
			if len(song.Charts) == 0 {
				continue
			}
			chart := song.Charts[0]
			var level Level
			if err := json.Unmarshal([]byte(chart.Content), &level); err != nil {
				continue
			}
			bpm := level.BPM
			spb := 60.0 / bpm
			offset := float64(level.Offset) / float64(level.SampleRate)
			if spb > 0 {
				for offset > spb {
					offset -= spb
				}
			}
			offset *= 1000
			reply, err := getBPM(&song)
			if err != nil {
				continue
			}
			bpmErr := eqErr(reply.BPM, bpm)
			if bpmErr < 0.5 {
				bpmCorrect++
			}
			testTotal++
			offsetErr := relErr(reply.Offset, offset, spb*1000)
			if offsetErr < 30 {
				offsetCorrect++
			}
			fmt.Printf(
				"%d %.2f %.1f | BPM %.2f - %.2f | Offset %.1f - %.1f\n",
				song.ID,
				bpmErr,
				offsetErr,
				bpm,
				reply.BPM,
				offset,
				reply.Offset,
			)
		}

		// reply
		rep := &GeneralReply{
			Success: true,
			Msg:     fmt.Sprintf("BPM: %d / %d correct; Offset: %d / %d correct", bpmCorrect, testTotal, offsetCorrect, testTotal),
		}
		c.JSON(http.StatusOK, rep)
	}
}
