package runkeeper

import (
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/jasonkit/n2rk/api/nikeplus"
	"github.com/parnurzeal/gorequest"
)

type RunKeeper struct {
	req   *gorequest.SuperAgent
	token string
}

func New(token string) *RunKeeper {
	return &RunKeeper{
		req:   gorequest.New(),
		token: token,
	}
}

const (
	host                      = "https://api.runkeeper.com"
	fitnessActivitiesEndPoint = host + "/fitnessActivities"
)

func (rk *RunKeeper) addAuth(req *gorequest.SuperAgent) *gorequest.SuperAgent {
	return req.Set("Authorization", "Bearer "+rk.token)
}

func (rk *RunKeeper) getWithAuth(url string) *gorequest.SuperAgent {
	return rk.addAuth(rk.req.Get(url))
}

func (rk *RunKeeper) postWithAuth(url string) *gorequest.SuperAgent {
	return rk.addAuth(rk.req.Post(url))
}

func (rk *RunKeeper) UploadFitnessActivities(activities []*FitnessActivity) {
	for _, v := range activities {
		rk.UploadFitnessActivity(v)
	}
}

func (rk *RunKeeper) UploadFitnessActivity(activity *FitnessActivity) {
	fmt.Printf("       Uploading %s", activity.Notes)
	_, _, errors := rk.postWithAuth(fitnessActivitiesEndPoint).
		Set("Content-Type", "application/vnd.com.runkeeper.NewFitnessActivity+json").
		Send(activity.Json()).
		End()
	if errors != nil {
		fmt.Printf("\r[Fail]\r\n")
		fmt.Printf("Failed in uploading %s: %v\n", activity.Notes, errors)
	} else {
		fmt.Printf("\r[Done]\r\n")
	}
}

func (rk *RunKeeper) UploadNikePlusActivities(activities []*nikeplus.Activity) {
	for _, v := range activities {
		rk.UploadNikePlusActivity(v)
	}
}

func (rk *RunKeeper) UploadNikePlusActivity(activity *nikeplus.Activity) {
	rk.UploadFitnessActivity(NewFitnessActivity(activity))
}

func (rk *RunKeeper) LastRunningTime() time.Time {
	var (
		body string
		err  []error
	)

	url := fitnessActivitiesEndPoint
	for {
		for {
			_, body, err = rk.getWithAuth(url).
				Set("Accept", "application/vnd.com.runkeeper.FitnessActivityFeed+json").
				End()
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				break
			}
		}

		st := ""
		cnt := 0
		dur := int64(0)
		jsonparser.ArrayEach([]byte(body), func(v []byte, t jsonparser.ValueType, idx int, err error) {
			actType, _ := jsonparser.GetString(v, "type")
			if actType == "Running" && st == "" {
				st, _ = jsonparser.GetString(v, "start_time")
				dur, _ = jsonparser.GetInt(v, "duration")
			}
			cnt += 1
		}, "items")

		if st != "" {
			startTime, _ := time.ParseInLocation("Mon, 2 Jan 2006 15:05:05", st, time.Now().Location())
			return startTime.Add(time.Duration(dur) * time.Second)
		} else {
			var err error
			url, err = jsonparser.GetString([]byte(body), "next")
			if err != nil || cnt == 0 {
				return time.Now()
			} else {
				url = host + url
			}
		}
	}
}
