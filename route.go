package main

import (
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
			rep := &GetChartsReply{
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
		dir := fmt.Sprintf("./songs/%s%s", uuid.NewV4().String(), ext)
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
