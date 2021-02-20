package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	REDDIT_HOST = "strapi.reddit.com"
	VIDEO_ID    = "llyqyx"
	TOKEN       = "-Sy9mdZh4u6IM38lCKvuYYDuuTlo"
)

// StreamStatsResponseBody ...
type StreamStatsResponseBody struct {
	Data StreamStats `json:"data"`
}

// StreamStats ...
type StreamStats struct {
	BroadcastTime           float64         `json:"broadcast_time"`
	Upvotes                 int             `json:"upvotes"`
	Downvotes               int             `json:"downvotes"`
	ShareLink               string          `json:"share_link"`
	TotalContinuousWatchers int             `json:"total_continuous_watchers"`
	TotalStreams            int             `json:"total_streams"`
	UniqueWatchers          int             `json:"unique_watchers"`
	Post                    StreamStatsPost `json:"post"`
}

// StreamStatsPost ...
type StreamStatsPost struct {
	CommentCount float64 `json:"commentCount"`
}

func getStreamStats(videoID string) (StreamStats, error) {
	log.Println("Attempting to retreive stats for video ID:", VIDEO_ID)

	var streamStats StreamStats

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s/videos/t3_%s", REDDIT_HOST, VIDEO_ID),
		nil,
	)
	if err != nil {
		log.Println(err.Error())
		return streamStats, err
	}

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", TOKEN))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return streamStats, err
	}

	log.Println("Response code:", res.StatusCode)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return streamStats, err
	}

	var response StreamStatsResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err.Error())
		return streamStats, err
	}

	streamStats = response.Data

	return streamStats, nil
}

func main() {
	fmt.Println("hello reddit-stats")

	streamStats, err := getStreamStats(VIDEO_ID)

	dataB, err := json.MarshalIndent(streamStats, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dataB))
}
