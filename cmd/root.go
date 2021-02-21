package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const (
	RedditHost = "strapi.reddit.com"
	// VIDEO_ID    = "llyqyx"
	// TOKEN = "-Sy9mdZh4u6IM38lCKvuYYDuuTlo"
)

var token string

var rootCmd = &cobra.Command{
	Use:   "redditstats [list of video ids space-separated]",
	Short: "RedditStats is a tool to pull post stats from Reddit.",
	Long:  "RedditStats is a tool to pull post stats from Reddit.",
	Run: func(cmd *cobra.Command, args []string) {
		for _, videoID := range args {
			streamStats, err := getStreamStats(token, videoID)
			if err != nil {
				log.Println("[ERROR] video id:", videoID, "error:", err.Error())
			} else {
				statsB, err := json.MarshalIndent(streamStats, "", "  ")
				if err != nil {
					log.Println("[ERROR] video id:", videoID, "error:", err.Error())
				} else {
					fmt.Println(string(statsB))
				}
			}
		}
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "Authorization bearer token")
	rootCmd.MarkFlagRequired("token")
}

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

func getStreamStats(token, videoID string) (StreamStats, error) {
	log.Println("Attempting to retreive stats for video ID:", videoID)

	var streamStats StreamStats

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s/videos/t3_%s", RedditHost, videoID),
		nil,
	)
	if err != nil {
		log.Println(err.Error())
		return streamStats, err
	}

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return streamStats, err
	}

	log.Println("Response code:", res.StatusCode)
	if res.StatusCode != 200 {
		return streamStats, errors.New("reddit API responded with an error code")
	}

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
