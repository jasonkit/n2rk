package nikeplus

import (
	"sort"
	"time"

	"github.com/buger/jsonparser"
	"github.com/parnurzeal/gorequest"
)

type NikePlus struct {
	req   *gorequest.SuperAgent
	token string
}

func New(token string) *NikePlus {
	return &NikePlus{
		req:   gorequest.New(),
		token: token,
	}
}

const (
	host               = "https://api.nike.com"
	prefix             = host + "/v1/me/sport/"
	activitiesEndPoint = prefix + "activities"
)

func activityDetailEndPoint(id string) string {
	return activitiesEndPoint + "/" + id
}

func activityGPSEndPoint(id string) string {
	return activityDetailEndPoint(id) + "/gps"
}

func (np *NikePlus) reqWithToken(url string) *gorequest.SuperAgent {
	return np.req.Get(url).Param("access_token", np.token)
}

func (np *NikePlus) Activities(start, end time.Time) []*Activity {

	activities := make(map[string]*Activity)

	_, body, _ := np.reqWithToken(activitiesEndPoint).
		Param("startDate", start.Format("2006-01-02")).
		Param("endDate", end.Format("2006-01-02")).
		Param("count", "20").
		End()

	for {
		repeated := false
		jsonparser.ArrayEach([]byte(body), func(v []byte, t jsonparser.ValueType, idx int, err error) {
			activity := NewActivity(v)
			if activity.Status == "COMPLETE" && activity.Type == "RUN" {
				if activities[activity.Id] == nil {
					activities[activity.Id] = activity
				} else {
					repeated = true
				}
			}
		}, "data")

		if repeated {
			break
		} else {
			next, nextType, _, _ := jsonparser.Get([]byte(body), "paging", "next")

			if nextType == jsonparser.String {
				_, body, _ = np.req.Get(host + string(next)).End()
			} else {
				break
			}
		}
	}

	var activityList []*Activity
	for _, v := range activities {
		np.FillMetricSummary(v)
		np.FillGPS(v)
		activityList = append(activityList, v)
	}

	sort.Sort(ByStartTime(activityList))

	return activityList
}

func (np *NikePlus) FillMetricSummary(a *Activity) {
	_, body, _ := np.reqWithToken(activityDetailEndPoint(a.Id)).End()
	metricSummaryJson, _, _, _ := jsonparser.Get([]byte(body), "metricSummary")
	a.MetricSummary = NewMetricSummary(metricSummaryJson)
}

func (np *NikePlus) FillGPS(a *Activity) {
	_, body, _ := np.reqWithToken(activityGPSEndPoint(a.Id)).End()
	a.GPS = NewGPS([]byte(body))
}
