package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwinProduction/go-github-wip/util"
	"github.com/google/go-github/v25/github"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
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
	_ = json.Unmarshal(bodyData, &pullRequestEvent)
	log.Println("Title:" + pullRequestEvent.GetPullRequest().GetTitle())
	// If title starts with "[WIP]", then set task to in progress (see https://github.com/wip/app)
	if strings.HasPrefix(pullRequestEvent.GetPullRequest().GetTitle(), "[WIP]") {
		pr := pullRequestEvent.GetPullRequest()
		util.SetAsWip(
			pullRequestEvent.Repo.Owner.GetName(),
			pullRequestEvent.Repo.GetName(),
			pr.Head.GetRef(),
			pr.Head.GetSHA(),
			pullRequestEvent.Installation.GetAppID(),
			pullRequestEvent.Installation.GetID(),
		)
	} else if strings.HasPrefix(*pullRequestEvent.GetChanges().Title.From, "[WIP]") {
		pr := pullRequestEvent.GetPullRequest()
		util.ClearWip(
			pullRequestEvent.Repo.Owner.GetName(),
			pullRequestEvent.Repo.GetName(),
			pr.Head.GetRef(),
			pr.Head.GetSHA(),
			pullRequestEvent.Installation.GetAppID(),
			pullRequestEvent.Installation.GetID(),
			util.GetCheckRunId(
				pullRequestEvent.Repo.Owner.GetName(),
				pullRequestEvent.Repo.GetName(),
				pr.Head.GetRef(),
				pullRequestEvent.Installation.GetAppID(),
				pullRequestEvent.Installation.GetID(),
			),
		)
	}
	fmt.Fprint(writer, "Title:"+pullRequestEvent.GetPullRequest().GetTitle())
}
