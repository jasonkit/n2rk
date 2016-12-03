package runkeeper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jasonkit/n2rk/api/nikeplus"
)

type FitnessActivity struct {
	Type          string
	Equipment     string
	StartTime     time.Time
	TotalDistance float64
	Duration      time.Duration
	TotalCalories int
	Notes         string
	Path          []Path
}

type Path struct {
	Timestamp float64 `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Type      string  `json:"type"`
}

type FitnessActivityJson struct {
	Type          string  `json:"type"`
	Equipment     string  `json:"equipment"`
	StartTime     string  `json:"start_time"`
	TotalDistance float64 `json:"total_distance"`
	Duration      float64 `json:"duration"`
	TotalCalories float64 `json:"total_calories"`
	Notes         string  `json:"notes"`
	Path          []Path  `json:"path"`
}

func (fa *FitnessActivity) Json() string {
	faJson := FitnessActivityJson{
		Type:          fa.Type,
		Equipment:     fa.Equipment,
		StartTime:     fa.StartTime.Format("Mon, 2 Jan 2006 15:05:05"),
		TotalDistance: fa.TotalDistance,
		Duration:      fa.Duration.Seconds(),
		TotalCalories: float64(fa.TotalCalories),
		Notes:         fa.Notes,
		Path:          fa.Path,
	}

	buf, err := json.Marshal(faJson)
	if err != nil {
		fmt.Printf("Failed to marshal RunKeeper FitnessActivity data: %v\n", err)
		return ""
	}

	return string(buf)
}

func NewFitnessActivity(activity *nikeplus.Activity) *FitnessActivity {
	return &FitnessActivity{
		Type:          "Running",
		Equipment:     "None",
		StartTime:     activity.StartTime,
		TotalDistance: activity.MetricSummary.Distance * 1000.0,
		Duration:      activity.MetricSummary.Duration,
		TotalCalories: activity.MetricSummary.Calories,
		Notes:         NotesFromNikePlusActivity(activity),
		Path:          PathsFromNikePlusGPS(activity.GPS),
	}
}

func PathsFromNikePlusGPS(gps *nikeplus.GPS) []Path {
	var path []Path

	for i, w := range gps.Waypoints {
		t := "gps"

		if i == 0 {
			t = "start"
		} else if i == len(gps.Waypoints) {
			t = "end"
		}

		p := Path{
			Timestamp: (time.Duration(i*4106) * time.Millisecond).Seconds(),
			Latitude:  w.Latitude,
			Longitude: w.Longitude,
			Altitude:  w.Elevation,
			Type:      t,
		}

		path = append(path, p)
	}

	return path
}

func NotesFromNikePlusActivity(activity *nikeplus.Activity) string {
	weekday := activity.StartTime.Weekday().String()
	period := ""
	hr := activity.StartTime.Hour()
	switch {
	case (hr >= 0 && hr <= 4) || (hr >= 21 && hr <= 23):
		period = "Night "
	case hr >= 5 && hr <= 11:
		period = "Morning "
	case hr >= 12 && hr <= 13:
		period = "Noon "
	case hr >= 14 && hr <= 17:
		period = "Afternoon "
	case hr >= 18 && hr <= 20:
		period = "Evening "
	}

	return fmt.Sprintf("%s %sRun on %s", weekday, period, activity.StartTime.Format("2 Jan 2006"))
}
