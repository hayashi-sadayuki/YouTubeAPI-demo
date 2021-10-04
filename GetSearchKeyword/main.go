package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "YOUR DEVELOPER KEY"
const targetVideoTitle = "SHOWROOM"

func main() {
	videos := make(map[string]string)

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	callApi(service, videos)

	printIDs(videos)
}

func callApi(service *youtube.Service, videos map[string]string) {
	nextPageToken := ""
	for i := 0; ; i++ {

		call := service.Search.List([]string{"id, snippet"}).
			Q(targetVideoTitle).
			MaxResults(50).
			Type("video").
			Order("date").
			PageToken(nextPageToken)
		response, err := call.Do()
		handleError(err, "")

		for _, item := range response.Items {
			switch item.Id.Kind {
			case "youtube#video":
				if strings.Contains(item.Snippet.Title, targetVideoTitle) {
					videos[item.Id.VideoId] = item.Snippet.Title
				}
			}
			if len(videos) >= 100 {
				return
			}
		}

		if response.NextPageToken == "" {
			break
		}

		nextPageToken = response.NextPageToken
	}
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func printIDs(matches map[string]string) {
	countNum := 1
	for id, title := range matches {
		fmt.Printf("%v: %v [https://www.youtube.com/watch?v=%v]\n", strconv.Itoa(countNum), title, id)
		countNum++
	}
	fmt.Printf("\n\n")
}
