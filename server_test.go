package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      *int   `json:"id,omitempty"`
}

type RequestEvent struct {
	ID       *int   `json:"id,omitempty"`
	Deadline string `json:"deadline"`
	Title    string `json:"title"`
	Memo     string `json:"memo"`
}

func EqualP(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if assert.Equal(t, expected, actual, msgAndArgs) == false {
		assert.FailNow(t, "FailNow!", msgAndArgs)
	}
}

func EqualEvent(t *testing.T, expectedID int, expected, actual RequestEvent) {
	EqualP(t, expected.Deadline, actual.Deadline)
	EqualP(t, expected.Title, actual.Title)
	EqualP(t, expected.Memo, actual.Memo)
	EqualP(t, expectedID, *actual.ID)
}

func PostEvent(t *testing.T, r *gin.Engine, ev RequestEvent) int {
	w := httptest.NewRecorder()
	d, _ := json.Marshal(ev)
	req, _ := http.NewRequest("POST", "/api/v1/event", bytes.NewReader(d))
	r.ServeHTTP(w, req)

	EqualP(t, 200, w.Code)
	body := Response{}
	json.Unmarshal(w.Body.Bytes(), &body)
	EqualP(t, "success", body.Status)
	EqualP(t, "registered", body.Message)
	t.Logf("assigned event id: %v\n", *body.ID)
	return *body.ID
}

func GetEvent(t *testing.T, r *gin.Engine, id int) RequestEvent {
	w := httptest.NewRecorder()
	t.Logf("request to %v\n", "/api/v1/event/"+strconv.Itoa(id))
	req, _ := http.NewRequest("GET", "/api/v1/event/"+strconv.Itoa(id), nil)
	r.ServeHTTP(w, req)

	EqualP(t, 200, w.Code)
	body := RequestEvent{}
	json.Unmarshal(w.Body.Bytes(), &body)
	return body
}

func GetEvents(t *testing.T, r *gin.Engine) []RequestEvent {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/event", nil)
	r.ServeHTTP(w, req)

	EqualP(t, 200, w.Code)
	body := []RequestEvent{}
	json.Unmarshal(w.Body.Bytes(), &body)

	t.Logf("GET all events: %v\n", w.Body.String())
	return body
}

func Test1(t *testing.T) {
	r := setupRouter()
	deleteAllEventsFromDB()

	w := httptest.NewRecorder()
	d, _ := json.Marshal(RequestEvent{nil, "2019-06-11T14:00:00+09:00", "レポート提出", ""})
	req, _ := http.NewRequest("POST", "/api/v1/event", bytes.NewReader(d))
	r.ServeHTTP(w, req)

	EqualP(t, 200, w.Code)

	body := Response{}
	json.Unmarshal(w.Body.Bytes(), &body)
	EqualP(t, "success", body.Status)
	EqualP(t, "registered", body.Message)
	t.Logf("assigned event id: %v\n", *body.ID)
}

func Test2(t *testing.T) {
	r := setupRouter()
	deleteAllEventsFromDB()

	w := httptest.NewRecorder()
	d, _ := json.Marshal(RequestEvent{nil, "2019/06/11T14:00:00+09:00", "レポート提出", ""})
	req, _ := http.NewRequest("POST", "/api/v1/event", bytes.NewReader(d))
	r.ServeHTTP(w, req)

	EqualP(t, 400, w.Code)

	body := Response{}
	json.Unmarshal(w.Body.Bytes(), &body)
	EqualP(t, "failure", body.Status)
	EqualP(t, "invalid date format", body.Message)
}

func Test3(t *testing.T) {
	r := setupRouter()
	deleteAllEventsFromDB()

	requestEvent := RequestEvent{nil, "2019-06-11T14:00:00+09:00", "レポート提出", "memomemo"}
	id := PostEvent(t, r, requestEvent)

	body2 := GetEvent(t, r, id)
	EqualEvent(t, id, requestEvent, body2)
}

func Test4(t *testing.T) {
	r := setupRouter()
	deleteAllEventsFromDB()

	ev := RequestEvent{nil, "2019-06-11T14:00:00+09:00", "レポート提出", "memomemo"}
	id := PostEvent(t, r, ev)

	ev2 := RequestEvent{nil, "2019-06-12T14:00:00+09:00", "レポート提出2", "memomemo2"}
	id2 := PostEvent(t, r, ev2)

	evR := GetEvent(t, r, id)
	EqualEvent(t, id, ev, evR)

	ev2R := GetEvent(t, r, id2)
	EqualEvent(t, id2, ev2, ev2R)
}

func Test5(t *testing.T) {
	r := setupRouter()
	deleteAllEventsFromDB()

	ev := RequestEvent{nil, "2019-06-11T14:00:00+09:00", "レポート提出", "memomemo"}
	id := PostEvent(t, r, ev)

	ev2 := RequestEvent{nil, "2019-06-12T14:00:00+09:00", "レポート提出2", "memomemo2"}
	id2 := PostEvent(t, r, ev2)

	ev3 := RequestEvent{nil, "2019-06-13T14:00:00+09:00", "レポート提出3", "memomemo3"}
	id3 := PostEvent(t, r, ev3)

	ev4 := RequestEvent{nil, "2019-06-14T14:00:00+09:00", "レポート提出4", "memomemo4"}
	id4 := PostEvent(t, r, ev4)

	evs := GetEvents(t, r)
	EqualP(t, 4, len(evs))
	EqualEvent(t, id, ev, evs[0])
	EqualEvent(t, id2, ev2, evs[1])
	EqualEvent(t, id3, ev3, evs[2])
	EqualEvent(t, id4, ev4, evs[3])
}
