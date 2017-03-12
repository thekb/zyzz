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
	ticker := time.NewTicker(2 * time.Second).C
	go scoreTicker(EventChannels[event.ShortId], ticker, event.ShortId, event.MatchUrl)
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
	if event.RunningNow == 0 {
		fmt.Println("in the if loop")
		close_chan := EventChannels[event.ShortId]
		//ticker := make(chan time.Time)
		close_chan <- true
		ticker := make(chan time.Time)
		go scoreTicker(close_chan, ticker, event.ShortId, event.MatchUrl)
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
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0, // use default DB
	})
	pong, err := r.Ping().Result()
	fmt.Println(pong, err)
	redisEntries, err := r.Get(event_shortId).Result()
	r.Close()
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

func scoreTicker(close_chan chan bool, ticker <- chan time.Time, eventId string, DataPath string) {
	//cric_url := "http://synd.cricbuzz.com/j2me/1.0/livematches.xml"
	cric_url := fmt.Sprintf("%s%s", DataPath, "commentary.xml")
	ForLoop:
	for {
		select {
		case <-ticker:
		//data := url.Values{}
		//data.Set("apikey", "bRGI79IqC3S5fh52QlsgxNnJbu72")
		//data.Add("unique_id", strconv.Itoa(matchId))
			req, _ := http.NewRequest("GET", cric_url, bytes.NewBufferString(""))
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Response from cric server failed ", err)
				continue
			}
			defer resp.Body.Close()
			status := resp.StatusCode
		//var eventScore EventScore
			if status == 200 {
				//body, err := ioutil.ReadAll(resp.Body)
				MatchDetails, err := ReadMatchData(resp.Body)
				match, err := json.Marshal(MatchDetails)
				if err == nil {
					r := redis.NewClient(&redis.Options{
						Addr:     "localhost:6379",
						Password: "", // no password set
						DB:      0, // use default DB
					})
					pong, err := r.Ping().Result()
					fmt.Println(pong, err)
					err = r.Set(eventId, match, 0).Err()
					if err != nil {
						fmt.Println("Error occured while setting in redis", err)
					}
					r.Close()
				}
			}
		case <-close_chan:
			break ForLoop
		default:
			continue
		}
	}
}
