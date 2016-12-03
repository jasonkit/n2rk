package nikeplus

import (
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

type Activity struct {
	Id            string
	Type          string
	StartTime     time.Time
	TimeZone      string
	Status        string
	DeviceType    string
	MetricSummary *MetricSummary
	Tags          []*Tag
	GPS           *GPS
}

type MetricSummary struct {
	Calories int
	Fuel     int
	Distance float64
	Steps    int
	Duration time.Duration
}

type Tag struct {
	Type  string
	Value string
}

type GPS struct {
	ElevationLoss  float64
	ElevationGain  float64
	ElevationMax   float64
	ElevationMin   float64
	IntervalMetric int
	IntervalUnit   string
	Waypoints      []*Waypoint
}

type Waypoint struct {
	Latitude  float64
	Longitude float64
	Elevation float64
}

func NewGPS(json []byte) *GPS {
	eloss, _ := jsonparser.GetFloat(json, "elevationLoss")
	egain, _ := jsonparser.GetFloat(json, "elevationGain")
	emax, _ := jsonparser.GetFloat(json, "elevationMax")
	emin, _ := jsonparser.GetFloat(json, "elevationMin")
	intervalMetric, _ := jsonparser.GetInt(json, "intervalMetric")
	intervalUnit, _ := jsonparser.GetString(json, "intervalUnit")

	var waypoints []*Waypoint
	jsonparser.ArrayEach(json, func(v []byte, t jsonparser.ValueType, idx int, err error) {
		waypoints = append(waypoints, NewWaypoint(v))
	}, "waypoints")

	return &GPS{
		ElevationLoss:  eloss,
		ElevationGain:  egain,
		ElevationMax:   emax,
		ElevationMin:   emin,
		IntervalMetric: int(intervalMetric),
		IntervalUnit:   intervalUnit,
		Waypoints:      waypoints,
	}
}

func NewTag(json []byte) *Tag {
	t, _ := jsonparser.GetString(json, "tagType")
	v, _ := jsonparser.GetString(json, "tagValue")
	return &Tag{t, v}
}

func NewWaypoint(json []byte) *Waypoint {
	lat, _ := jsonparser.GetFloat(json, "latitude")
	lng, _ := jsonparser.GetFloat(json, "longitude")
	elevation, _ := jsonparser.GetFloat(json, "elevation")
	return &Waypoint{lat, lng, elevation}
}

func NewMetricSummary(json []byte) *MetricSummary {
	calories, _ := jsonparser.GetString(json, "calories")
	fuel, _ := jsonparser.GetString(json, "fuel")
	distance, _ := jsonparser.GetString(json, "distance")
	steps, _ := jsonparser.GetString(json, "steps")
	duration, _ := jsonparser.GetString(json, "duration")

	caloriesVal, _ := strconv.ParseInt(calories, 10, 64)
	fuelVal, _ := strconv.ParseInt(fuel, 10, 64)
	stepsVal, _ := strconv.ParseInt(steps, 10, 64)
	distanceVal, _ := strconv.ParseFloat(distance, 64)
	durationVal, _ := time.Parse("15:04:05.000", duration)

	return &MetricSummary{
		Calories: int(caloriesVal),
		Fuel:     int(fuelVal),
		Distance: distanceVal,
		Steps:    int(stepsVal),
		Duration: time.Duration(durationVal.Hour())*time.Hour +
			time.Duration(durationVal.Minute())*time.Minute +
			time.Duration(durationVal.Second())*time.Second,
	}
}

func NewActivity(json []byte) *Activity {
	id, _ := jsonparser.GetString(json, "activityId")
	t, _ := jsonparser.GetString(json, "activityType")
	st, _ := jsonparser.GetString(json, "startTime")
	tz, _ := jsonparser.GetString(json, "activityTimeZone")
	status, _ := jsonparser.GetString(json, "status")
	deviceType, _ := jsonparser.GetString(json, "deviceType")

	var tags []*Tag
	jsonparser.ArrayEach(json, func(v []byte, t jsonparser.ValueType, idx int, err error) {
		tags = append(tags, NewTag(v))
	}, "tags")

	startTime, _ := time.Parse("2006-01-02T15:04:05Z", st)
	timezone, _ := time.LoadLocation(tz)

	return &Activity{
		Id:         id,
		Type:       t,
		StartTime:  startTime.In(timezone),
		TimeZone:   tz,
		Status:     status,
		DeviceType: deviceType,
		Tags:       tags,
	}
}

type ByStartTime []*Activity

func (a ByStartTime) Len() int           { return len(a) }
func (a ByStartTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartTime) Less(i, j int) bool { return a[i].StartTime.UnixNano() < a[j].StartTime.UnixNano() }
