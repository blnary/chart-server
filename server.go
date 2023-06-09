package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	connStr    = flag.String("conn_str", "iloveyourmum", "Connect string of MySQL")
	listenAddr = flag.String("listen_addr", "0.0.0.0:51121", "Listen address")
	credential = flag.String("credential", "speciallady", "Credential of admin auth")
)

const BPM_SERVER_URL = "http://bpm-server:8000/"

func main() {

	// init database
	flag.Parse()
	db, err := gorm.Open(mysql.Open(*connStr), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1;").AutoMigrate(&Song{}, &Chart{})
	if err != nil {
		log.Fatalln(err)
	}

	// init server
	server := gin.Default()
	public := server.Group("")
	public.GET("songs/:id", GetSongFile(db))
	public.GET("charts", GetCharts(db))
	public.GET("songs", GetSongs(db))
	public.GET("bpm/:id", GetBPM(db))
	public.GET("test/bpm", TestBPM(db))
	public.POST("bpm/:id", SetBPM(db))
	public.POST("cal/:id", CalBPM(db))
	public.POST("sync/:id", SyncBPM(db))
	public.POST("charts", PostCharts(db))
	public.POST("songs", PostSongs(db))
	admin := server.Group("")
	admin.Use(AuthMiddleware())
	admin.POST("charts/:id", PostChart(db))
	admin.DELETE("charts/:id", DeleteChart(db))
	admin.DELETE("songs/:id", DeleteSong(db))

	// run server
	log.Printf("listening at %s", *listenAddr)
	if err := server.Run(*listenAddr); err != nil {
		log.Fatalln(err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCredential := c.GetHeader("Credential")
		cond := userCredential == *credential
		if !cond {
			rep := &GeneralReply{
				Success: false,
				Msg:     "credential incorrect",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, rep)
			return
		}
		c.Next()
	}
}
