package main

import (
	"encoding/json"
	"net/http"
	"sort"

	"time"

	"github.com/gin-gonic/gin"

	"strconv"
)

type Event struct {
	Deadline time.Time
	Title    string
	Memo     string
}

var table map[int]Event
var lastid int

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

	event := Event{t, body.Title, body.Memo}
	lastid++
	table[lastid] = event

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "registered",
		"id":      lastid,
	})
}

func getAllEvents(ctx *gin.Context) {
	type Response struct {
		ID       int    `json:"id"`
		Deadline string `json:"deadline"`
		Title    string `json:"title"`
		Memo     string `json:"memo"`
	}

	keys := []int{}
	for id := range table {
		keys = append(keys, id)
	}
	sort.Ints(keys)
	entries := []Response{}
	for _, k := range keys {
		y := table[k]
		e := Response{k, y.Deadline.Format(time.RFC3339), y.Title, y.Memo}
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
	if event, ok := table[id]; ok {
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

func setupRouter() *gin.Engine {
	table = make(map[int]Event)
	lastid = 0

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
