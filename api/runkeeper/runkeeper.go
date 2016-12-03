package runkeeper

import (
	"fmt"

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
