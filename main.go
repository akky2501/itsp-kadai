package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Event struct {
	ID       int `gorm:"primary_key"`
	Deadline time.Time
	Title    string
	Memo     string
}

var db *gorm.DB = nil

func registerEvent(ctx *gin.Context) {
	var body struct {
		Deadline string `json:"deadline"`
		Title    string `json:"title"`
		Memo     string `json:"memo"`
	}

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "invalid request format",
		})
		return
	}

	t, err := time.Parse(time.RFC3339, body.Deadline)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "invalid date format",
		})
		return
	}

	event := Event{}
	event.Deadline = t
	event.Title = body.Title
	event.Memo = body.Memo

	db.Create(&event)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "registered",
		"id":      event.ID,
	})
}

func getAllEvents(ctx *gin.Context) {
	type Response struct {
		ID       int    `json:"id"`
		Deadline string `json:"deadline"`
		Title    string `json:"title"`
		Memo     string `json:"memo"`
	}

	events := []Event{}
	db.Find(&events)
	sort.Slice(events, func(i, j int) bool { return events[i].ID < events[j].ID })
	entries := []Response{}
	for _, y := range events {
		e := Response{y.ID, y.Deadline.Format(time.RFC3339), y.Title, y.Memo}
		entries = append(entries, e)
	}

	json, err := json.Marshal(entries)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failure",
			"message": "internal server error",
		})
		return
	}

	ctx.String(http.StatusOK, string(json))
}

func getEvent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	event := Event{}
	event.ID = id
	if !db.First(&event).RecordNotFound() {
		ctx.JSON(http.StatusOK, gin.H{
			"id":       id,
			"deadline": event.Deadline.Format(time.RFC3339),
			"title":    event.Title,
			"memo":     event.Memo,
		})
	} else {
		ctx.String(http.StatusNotFound, "")
	}
}

func deleteAllEventsFromDB() {
	e := Event{}
	db.Delete(e)
}

func setupRouter() *gin.Engine {
	if db == nil {
		var err error
		db, err = gorm.Open("mysql", "root:@tcp(localhost:3306)/itspkadai?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			panic("failed to connect database")
		}

		db.AutoMigrate(&Event{})
	}

	r := gin.Default()
	r.POST("/api/v1/event", registerEvent)
	r.GET("/api/v1/event", getAllEvents)
	r.GET("/api/v1/event/:id", getEvent)

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
