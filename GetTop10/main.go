package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "YOUR DEVELOPER KEY"
const targetkeyword = "Apex Legends"

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
	nowUTC := time.Now().UTC()
	jst := time.FixedZone("JST", +9*60*60)
	nowJST := nowUTC.In(jst)
	add3day := nowJST.AddDate(0, 0, 3)

	fmt.Printf("from %v to %v search...\n\n", nowJST.Format(time.RFC3339), add3day.Format(time.RFC3339))

	call := service.Search.List([]string{"id, snippet"}).
		Q(targetkeyword).
		MaxResults(10).
		Type("video").
		RegionCode("JP").
		PublishedAfter(nowJST.Format(time.RFC3339)).
		PublishedBefore(add3day.Format(time.RFC3339)).
		Order("rating")

	response, err := call.Do()
	handleError(err, "")

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		}
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
