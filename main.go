package main

import (
	"encoding/json"
	"github.com/TwinProduction/go-github-wip/config"
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
	} else if pullRequestEvent.GetAction() != "edited" && pullRequestEvent.GetAction() != "opened" {
		// Ignore pull request events that don't modify the title
		return
	}

	// FIXME: If it's a new commit, remove old check run on old commit, and apply it on the new commit.
	// Or should the check runs be converted to a check suite instead?

	writer.WriteHeader(200)
	log.Printf(
		"[webhookHandler] Got a PR event from %s/%s#%d with title: %s\n",
		pullRequestEvent.GetRepo().GetOwner().GetLogin(),
		pullRequestEvent.GetRepo().GetName(),
		pullRequestEvent.GetPullRequest().GetNumber(),
		pullRequestEvent.GetPullRequest().GetTitle(),
	)
	// If title starts with "[WIP]", then set task to in progress
	if strings.HasPrefix(pullRequestEvent.GetPullRequest().GetTitle(), "[WIP]") {
		if config.Get().IsDebugging() {
			log.Printf("[webhookHandler] (SetAsWip) Body: %v\n", pullRequestEvent)
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetRepo().GetOwner().GetLogin())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetRepo().GetName())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetRef())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetSHA())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetInstallation().GetID())
		}
		pr := pullRequestEvent.GetPullRequest()
		println("[A]")
		util.SetAsWip(
			pullRequestEvent.GetRepo().GetOwner().GetLogin(),
			pullRequestEvent.GetRepo().GetName(),
			pr.GetHead().GetRef(),
			pr.GetHead().GetSHA(),
			pullRequestEvent.GetInstallation().GetID(),
		)
		println("[B]")
		util.ToggleWipLabelOnIssue(
			pullRequestEvent.GetRepo().GetOwner().GetLogin(),
			pullRequestEvent.GetRepo().GetName(),
			pr.GetNumber(),
			pullRequestEvent.GetInstallation().GetID(),
			true,
		)
	} else if strings.HasPrefix(*pullRequestEvent.GetChanges().Title.From, "[WIP]") {
		if config.Get().IsDebugging() {
			log.Printf("[webhookHandler] (ClearWip) Body: %v\n", pullRequestEvent)
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetRepo().GetOwner().GetLogin())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetRepo().GetName())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetRef())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetPullRequest().GetHead().GetSHA())
			log.Printf("[webhookHandler] (SetAsWip) %v\n", pullRequestEvent.GetInstallation().GetID())
		}
		pr := pullRequestEvent.GetPullRequest()
		println("[C]")
		util.ClearWip(
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
		println("[D]")
		util.ToggleWipLabelOnIssue(
			pullRequestEvent.GetRepo().GetOwner().GetLogin(),
			pullRequestEvent.GetRepo().GetName(),
			pr.GetNumber(),
			pullRequestEvent.GetInstallation().GetID(),
			false,
		)
	}
}
