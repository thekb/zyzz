package api

import (
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/zyzz/db/models"
	"gopkg.in/redis.v5"
	"time"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/jmoiron/sqlx"
)


type GetEvent struct {
	Common
}

type GetEvents struct {
	Common
}

type CreateEvent struct {
	Common
}

type UpdateEvent struct {
	Common
}

type GetEventStreams struct {
	Common
}

type GetCricketScores struct {
	Common
}

type EventScore struct {
	InningsRequirement string `json:"innings-requirement"`
	Team1              string `json:"team-1"`
	Team2              string `json:"team-2"`
	Score              string `json:"score"`
	Ttl                int `json:"ttl"`
	MatchType          string `json:"type"`
}

var EventChannels = make(map[string]chan bool)

func (ge *GetEvent) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(SHORT_ID)
	event, err := models.GetEventForShortId(ge.DB, shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, Response{Data:event})
	return
}

func (ce *CreateEvent) Serve(ctx *iris.Context) {
	var event models.Event
	err := ctx.ReadJSON(&event)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to decode request %s", err.Error())
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	event.ShortId = getNewShortId()
	id, err := models.CreateEvent(ce.DB, &event)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	EventChannels[event.ShortId] = make(chan bool)
	go scoreTicker(ce.R, EventChannels[event.ShortId], event.ShortId, event.MatchUrl)
	event, _ = models.GetEventForId(ce.DB, id)
	ctx.JSON(iris.StatusOK, Response{Data:event})
	return
}

func (ce *UpdateEvent) Serve(ctx *iris.Context) {
	//var event models.Event
	event_shortId := ctx.GetString(SHORT_ID)
	event, err := models.GetEventForShortId(ce.DB, event_shortId)
	if err != nil {
		ctx.Log(iris.ProdMode, "Not able to fetch an event %s", err.Error())
		ctx.JSON(iris.StatusNotFound, Response{Error:err.Error()})
		return
	}
	err = ctx.ReadJSON(&event)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to decode request %s", err.Error())
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	err = models.UpdateEvent(ce.DB, &event)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, Response{Error:err.Error()})
		return
	}
	fmt.Println("running now", event.RunningNow)
	close, ok := EventChannels[event.ShortId]
	if event.RunningNow == 0 {
		fmt.Println("in the if loop")
		close <- true
	} else {
		if !ok {
			EventChannels[event.ShortId] = make(chan bool)
			go scoreTicker(ce.R, close, event.ShortId, event.MatchUrl)
		}
	}
	//go scoreTicker(ce.CC, ce.DB)
	event, _ = models.GetEventForId(ce.DB, event.Id)
	ctx.JSON(iris.StatusOK, Response{Data:event})
	return
}

func (ge *GetEvents) Serve(ctx *iris.Context) {
	events, err := models.GetEvents(ge.DB)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, &Response{Data:events})
	return
}

func (ges *GetEventStreams) Serve(ctx *iris.Context) {
	event_shortId := ctx.GetString(SHORT_ID)
	streams, err := models.GetStreams(ges.DB, event_shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	for i := 0; i < len(streams); i++ {
		user_publish, err := models.GetUserForId(ges.DB, int64(streams[i].CreatorId))
		if err == nil {
			streams[i].User = user_publish
		}

	}
	ctx.JSON(iris.StatusOK, &Response{Data:streams})
	return
}

func (gcs *GetCricketScores) Serve(ctx *iris.Context) {
	event_shortId := ctx.GetString(SHORT_ID)
	redisEntries, err := gcs.R.Get(event_shortId).Result()
	if err != nil {
		fmt.Println("No data for event", event_shortId)
	}
	var MatchDetails MatchDetails
	err = json.Unmarshal([]byte(redisEntries), &MatchDetails)
	if err != nil {
		fmt.Println("err is ", err)
	}
	if err != nil {
		fmt.Println("unable to get data for event", err)
	}
	ctx.JSON(iris.StatusOK, &Response{Data:MatchDetails})
	return
}

func scoreTicker(r *redis.Client, close chan bool, eventId string, DataPath string) {
	cricURL := fmt.Sprintf("%s%s", DataPath, "commentary.xml")
	// initialize ticker
	ticker := time.Tick(time.Second * 2)

	Loop:
	for {
		select {
		case <- close:
			fmt.Println("i am in case close chan")
			break Loop
		case <- ticker:
			updateScore(r, cricURL, eventId)
		}
	}
}

func updateScore(r *redis.Client, cricURL, eventId string) {
	req, _ := http.NewRequest("GET", cricURL, bytes.NewBufferString(""))
	client := &http.Client{}
	resp, err := client.Do(req)
	// close response body
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Response from cric server failed:", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		MatchDetails, err := ReadMatchData(resp.Body)
		match, err := json.Marshal(MatchDetails)
		if err == nil {

			err = r.Set(eventId, match, 0).Err()
			if err != nil {
				fmt.Println("Error occured while setting in redis", err)
			}
		}
	}
}

func StartEventTickers(db *sqlx.DB, r *redis.Client) {
	events, err := models.GetEvents(db)
	if err != nil {
		fmt.Println("Error occured while getting events", err)
		return
	}
	for i := 0; i < len(events); i++ {
		event := events[i]
		if event.RunningNow == 1 {
			close, ok := EventChannels[event.ShortId]
			if !ok {
				EventChannels[event.ShortId] = make(chan bool)
				fmt.Println("calling the ticker")
				go scoreTicker(r, close, event.ShortId, event.MatchUrl)
			}
		}
	}
}