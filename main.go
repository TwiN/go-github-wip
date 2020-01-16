package main

import (
	"encoding/json"
	"github.com/TwinProduction/go-github-wip/config"
	"github.com/TwinProduction/go-github-wip/util"
	"github.com/google/go-github/v25/github"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	config.Validate()
	http.HandleFunc("/", webhookHandler)
	log.Println("[main] Listening to port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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
		return
	} else if pullRequestEvent.GetAction() != "edited" && pullRequestEvent.GetAction() != "opened" && pullRequestEvent.GetAction() != "synchronize" {
		if config.Get().IsDebugging() {
			if isActionCompletelyIgnored(pullRequestEvent.GetAction()) {
				log.Printf("[webhookHandler] Got a non-edit event: %v\n", string(bodyData))
			}
		}
		// Ignore pull request events that don't modify the title
		return
	}

	writer.WriteHeader(200)
	log.Printf(
		"[webhookHandler] Got a PR event from %s/%s#%d with title: %s\n",
		pullRequestEvent.GetRepo().GetOwner().GetLogin(),
		pullRequestEvent.GetRepo().GetName(),
		pullRequestEvent.GetPullRequest().GetNumber(),
		pullRequestEvent.GetPullRequest().GetTitle(),
	)
	// If title starts with "[WIP]", then set task to in progress
	if config.Get().HasWipPrefix(pullRequestEvent.GetPullRequest().GetTitle()) {
		if config.Get().IsDebugging() {
			log.Printf("[webhookHandler] (SetAsWip) Body: %v\n", pullRequestEvent)
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetRepo().GetOwner().GetLogin())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetRepo().GetName())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetRef())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetSHA())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetInstallation().GetID())
		}
		pr := pullRequestEvent.GetPullRequest()
		go util.SetAsWip(
			pullRequestEvent.GetRepo().GetOwner().GetLogin(),
			pullRequestEvent.GetRepo().GetName(),
			pr.GetHead().GetRef(),
			pr.GetHead().GetSHA(),
			pullRequestEvent.GetInstallation().GetID(),
		)
		go util.ToggleWipLabelOnIssue(
			pullRequestEvent.GetRepo().GetOwner().GetLogin(),
			pullRequestEvent.GetRepo().GetName(),
			pr.GetNumber(),
			pullRequestEvent.GetInstallation().GetID(),
			true,
		)
	} else if config.Get().HasWipPrefix(*pullRequestEvent.GetChanges().Title.From) {
		if config.Get().IsDebugging() {
			log.Printf("[webhookHandler] (ClearWip) Body: %v\n", pullRequestEvent)
			log.Printf("[webhookHandler] (ClearWip) %v\n", pullRequestEvent.GetRepo().GetOwner().GetLogin())
			log.Printf("[webhookHandler] (ClearWip) %v\n", pullRequestEvent.GetRepo().GetName())
			log.Printf("[webhookHandler] (ClearWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetRef())
			log.Printf("[webhookHandler] (ClearWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetSHA())
			log.Printf("[webhookHandler] (ClearWip) %v\n", pullRequestEvent.GetInstallation().GetID())
		}
		pr := pullRequestEvent.GetPullRequest()
		go util.ClearWip(
			pullRequestEvent.GetRepo().GetOwner().GetLogin(),
			pullRequestEvent.GetRepo().GetName(),
			pr.GetHead().GetRef(),
			pr.GetHead().GetSHA(),
			pullRequestEvent.GetInstallation().GetID(),
			util.GetCheckRunId(
				pullRequestEvent.GetRepo().GetOwner().GetLogin(),
				pullRequestEvent.GetRepo().GetName(),
				pr.GetHead().GetRef(),
				pullRequestEvent.GetInstallation().GetID(),
			),
		)
		go util.ToggleWipLabelOnIssue(
			pullRequestEvent.GetRepo().GetOwner().GetLogin(),
			pullRequestEvent.GetRepo().GetName(),
			pr.GetNumber(),
			pullRequestEvent.GetInstallation().GetID(),
			false,
		)
	}
}

// Events that we really don't care about can just be plainly ignored
func isActionCompletelyIgnored(action string) bool {
	return action == "labeled"
}
