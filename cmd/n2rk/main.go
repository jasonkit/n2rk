package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jasonkit/n2rk/api/nikeplus"
)

type ConfigJson struct {
	NikePlusToken  string `json:"nike_plus_token"`
	RunKeeperToken string `json:"runkeeper_token"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
}

type Config struct {
	*ConfigJson
	StartTime time.Time
	EndTime   time.Time
}

func ReadConfig() *Config {
	var (
		buf       []byte
		err       error
		startTime time.Time
		endTime   time.Time
	)

	configFile := flag.String("c", "config.json", "config file")
	flag.Parse()

	if buf, err = ioutil.ReadFile(*configFile); err != nil {
		fmt.Printf("Failed to open config file: %v\n", err)
		return nil
	}

	var configJson ConfigJson

	if err = json.Unmarshal(buf, &configJson); err != nil {
		fmt.Printf("Invalid config file: %v\n", err)
		return nil
	}

	if startTime, err = time.ParseInLocation("2006-01-02", configJson.StartTime, time.Now().Location()); err != nil {
		fmt.Printf("Invalid start_time value: %v\n", err)
		return nil
	}

	if endTime, err = time.ParseInLocation("2006-01-02", configJson.EndTime, time.Now().Location()); err != nil {
		endTime = time.Now()
	}

	return &Config{
		ConfigJson: &configJson,
		StartTime:  startTime,
		EndTime:    endTime,
	}
}

func main() {
	config := ReadConfig()

	if config == nil {
		return
	}

	/*
		np := nikeplus.New(config.NikePlusToken)
		activities := np.Activities(config.StartTime, config.EndTime)

		fmt.Printf("Exporting...\n")
		nikeplus.Export(activities, "./nikeplus.bin")
	*/
	fmt.Printf("Importing...\n")
	activities := nikeplus.Import("./nikeplus.bin")
	fmt.Printf("Loaded %d activities\n", len(activities))

	/*
		rk := runkeeper.New(config.RunKeeperToken)
		rk.UploadNikePlusActivities(activities[1:])
	*/
}
