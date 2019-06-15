package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v25/github"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//client, ctx := util.GetGithubClient()
	//events, response, err := client.Issues.ListEvents(ctx, &github.ListOptions{})
	//fmt.Println(events)
	//fmt.Println(response)
	//fmt.Println(err)
	http.HandleFunc("/", webhookHandler)
	log.Println("[main] Listening to port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

// See https://github.com/gdperkins/gondle/blob/master/consumer.go
func webhookHandler(writer http.ResponseWriter, request *http.Request) {
	bodyData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("[webhookHandler] Unable to read body: %s\n", err.Error())
		writer.WriteHeader(500)
		return
	}
	writer.WriteHeader(200)
	pullRequestEvent := github.PullRequestEvent{}
	json.Unmarshal(bodyData, &pullRequestEvent)
	log.Println("Title:" + pullRequestEvent.GetPullRequest().GetTitle())
	// If title starts with "WIP" or "[WIP]", then set task to in progress (see https://github.com/wip/app)
	fmt.Fprint(writer, "Title:"+pullRequestEvent.GetPullRequest().GetTitle())
}
