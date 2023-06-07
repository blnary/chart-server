package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func PostSong(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// read song
		file, err := c.FormFile("song")
		if err != nil {
			rep := &PostSongReply{
				Success: false,
				Msg:     "failed to read song",
			}
			c.JSON(http.StatusBadRequest, rep)
		}

		// save song file to disk
		base := filepath.Base(file.Filename)
		ext := filepath.Ext(file.Filename)
		dir := fmt.Sprintf("./song/%s%s", uuid.NewV4().String(), ext)
		if err := c.SaveUploadedFile(file, dir); err != nil {
			rep := &PostSongReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save song to %s", dir),
			}
			c.JSON(http.StatusInternalServerError, rep)
		}

		// save song info to db
		song := &Song{
			Name:     base,
			Location: dir,
		}
		if err := db.Create(song).Error; err != nil {
			rep := &PostSongReply{
				Success: false,
				Msg:     fmt.Sprintf("failed to save song to database"),
			}
			c.JSON(http.StatusInternalServerError, rep)
		}

		// reply
		rep := &PostSongReply{
			Success:  true,
			Location: dir,
		}
		c.JSON(http.StatusOK, rep)
	}
}
