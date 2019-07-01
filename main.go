package main

import (
	"encoding/json"
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
	pullRequestEvent := github.PullRequestEvent{}
	err = json.Unmarshal(bodyData, &pullRequestEvent)
	if err != nil {
		log.Println("[webhookHandler] Ignoring request, because its body couldn't be unmarshalled to a PullRequestEvent")
		// This isn't a pull request event, ignore the request.
		writer.WriteHeader(400)
		return
	}
	writer.WriteHeader(200)
	log.Printf("[webhookHandler] Title: %s\n", pullRequestEvent.GetPullRequest().GetTitle())
	// If title starts with "[WIP]", then set task to in progress (see https://github.com/wip/app)
	if strings.HasPrefix(pullRequestEvent.GetPullRequest().GetTitle(), "[WIP]") {
		//log.Printf("[webhookHandler] Body: %v\n", pullRequestEvent)
		pr := pullRequestEvent.GetPullRequest()
		util.SetAsWip(
			pullRequestEvent.Repo.Owner.GetLogin(),
			pullRequestEvent.Repo.GetName(),
			pr.Head.GetRef(),
			pr.Head.GetSHA(),
			pullRequestEvent.Installation.GetAppID(),
			pullRequestEvent.Installation.GetID(),
		)
	} else if strings.HasPrefix(*pullRequestEvent.GetChanges().Title.From, "[WIP]") {
		//log.Printf("[webhookHandler] Body: %v\n", pullRequestEvent)
		pr := pullRequestEvent.GetPullRequest()
		util.ClearWip(
			pullRequestEvent.Repo.Owner.GetLogin(),
			pullRequestEvent.Repo.GetName(),
			pr.Head.GetRef(),
			pr.Head.GetSHA(),
			pullRequestEvent.Installation.GetAppID(),
			pullRequestEvent.Installation.GetID(),
			util.GetCheckRunId(
				pullRequestEvent.Repo.Owner.GetLogin(),
				pullRequestEvent.Repo.GetName(),
				pr.Head.GetRef(),
				pullRequestEvent.Installation.GetAppID(),
				pullRequestEvent.Installation.GetID(),
			),
		)
	}
}
