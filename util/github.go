package util

import (
	"context"
	"github.com/TwinProduction/go-github-wip/config"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v25/github"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	CheckRunName        = "Work in Progress"
	WipLabelName        = "wip"
	WipLabelColor       = "FBCA04"
	WipLabelDescription = "Work in progress"
)

func GetGithubClient(installationId int64) (*github.Client, context.Context) {
	// TODO: Cache client based on installationId
	return createGithubClient(installationId), context.Background()
}

func createGithubClient(installationId int64) *github.Client {
	log.Printf("[createGithubClient] Creating client for appId=%d and installationId=%d\n", int(config.Get().GetAppId()), int(installationId))
	transport := http.DefaultTransport
	itr, err := ghinstallation.NewKeyFromFile(transport, int(config.Get().GetAppId()), int(installationId), config.Get().GetPrivateKeyFileName())
	if err != nil {
		log.Printf("[createGithubClient] Failed create Github client: %v\n", err)
		log.Fatal(err)
	}
	if config.Get().IsDebugging() {
		_, err := itr.Token()
		if err != nil {
			log.Printf("[createGithubClient] Failed to get token: %v\n", err)
		}
	}
	return github.NewClient(&http.Client{Transport: itr})
}

func SetAsWip(userName, repositoryName, branch, commit string, installationId int64) *github.CheckRun {
	log.Printf("[SetAsWip] Creating WIP CheckRun on branch %s from repository %s/%s to WIP", branch, userName, repositoryName)
	client, ctx := GetGithubClient(installationId)
	status := "in_progress"
	outputTitle := "Do not merge!"
	outputSummary := "Do not merge!"
	checkRun, resp, err := client.Checks.CreateCheckRun(
		ctx,
		userName,
		repositoryName,
		github.CreateCheckRunOptions{
			Name:       CheckRunName,
			HeadBranch: branch,
			HeadSHA:    commit,
			Status:     &status,
			Output: &github.CheckRunOutput{
				Title:   &outputTitle,
				Summary: &outputSummary,
			},
		},
	)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("[SetAsWip] Response status: %s\n", resp.Status)
	if config.Get().IsDebugging() {
		if checkRun != nil {
			log.Printf("[SetAsWip] Response body: %v\n", checkRun)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err.Error())
			}
			log.Printf("[SetAsWip] Response body: %s\n", body)
		}
	}
	return checkRun
}

func ClearWip(userName, repositoryName, branch, commit string, installationId int64, checkRunId int64) {
	log.Printf("[ClearWip] Clearing WIP CheckRun on branch %s from repository %s/%s to WIP", branch, userName, repositoryName)
	client, ctx := GetGithubClient(installationId)
	status := "completed"
	conclusion := "success"
	checkRun, resp, err := client.Checks.UpdateCheckRun(
		ctx,
		userName,
		repositoryName,
		checkRunId,
		github.UpdateCheckRunOptions{
			Name:       CheckRunName,
			HeadBranch: &branch,
			HeadSHA:    &commit,
			Status:     &status,
			Conclusion: &conclusion,
			CompletedAt: &github.Timestamp{
				Time: time.Now(),
			},
		},
	)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("[ClearWip] Response status: %s\n", resp.Status)
	if config.Get().IsDebugging() {
		if checkRun != nil {
			log.Printf("[ClearWip] Response body: %v\n", checkRun)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err.Error())
			}
			log.Printf("[ClearWip] Response body: %s\n", body)
		}
	}
}

func GetCheckRunId(owner, repository, branch string, installationId int64) int64 {
	client, ctx := GetGithubClient(installationId)
	checkRunName := CheckRunName
	listCheckRuns, resp, err := client.Checks.ListCheckRunsForRef(
		ctx,
		owner,
		repository,
		branch,
		&github.ListCheckRunsOptions{
			CheckName: &checkRunName,
		},
	)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("[GetCheckRunId] Response status: %s\n", resp.Status)
	if config.Get().IsDebugging() {
		if listCheckRuns != nil {
			log.Printf("[GetCheckRunId] Response body: %v\n", listCheckRuns)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err.Error())
			}
			log.Printf("[GetCheckRunId] Response body: %s\n", body)
		}
	}

	return listCheckRuns.CheckRuns[0].GetID()
}

func ToggleWipLabelOnIssue(userName, repositoryName string, issueNumber int, installationId int64, addLabel bool) {
	var verb string
	if addLabel {
		verb = "Adding"
	} else {
		verb = "Removing"
	}
	log.Printf("[toggleWipLabelOnIssue] %s WIP label on %s/%s#%d\n", verb, userName, repositoryName, issueNumber)
	client, ctx := GetGithubClient(installationId)

	createWipLabelIfNotExist(userName, repositoryName, installationId)

	var labels []*github.Label
	var err error

	if addLabel {
		labels, _, err = client.Issues.AddLabelsToIssue(
			ctx,
			userName,
			repositoryName,
			issueNumber,
			[]string{WipLabelName},
		)
	} else {
		_, err = client.Issues.RemoveLabelForIssue(
			ctx,
			userName,
			repositoryName,
			issueNumber,
			WipLabelName,
		)
	}
	if err != nil {
		panic(err.Error())
	}
	if config.Get().IsDebugging() && labels != nil {
		log.Printf("[toggleWipLabelOnIssue] Response body: %v\n", labels)
	}
}

func createWipLabelIfNotExist(userName, repositoryName string, installationId int64) {
	// TODO: Create cache with key "userName/repositoryName" and value representing whether the label has already been created
	client, ctx := GetGithubClient(installationId)
	_, resp, err := client.Issues.GetLabel(
		ctx,
		userName,
		repositoryName,
		WipLabelName,
	)
	if err != nil {
		if resp.StatusCode == 404 {
			err = nil
			labelName := WipLabelName
			labelColor := WipLabelColor
			labelDescription := WipLabelDescription

			// Create label
			label, _, err := client.Issues.CreateLabel(
				ctx,
				userName,
				repositoryName,
				&github.Label{
					Name:        &labelName,
					Color:       &labelColor,
					Description: &labelDescription,
				},
			)
			if err != nil {
				panic(err.Error())
			}
			log.Printf("[createWipLabelIfNotExist] Successfully created label '%s' on %s/%s\n", WipLabelName, userName, repositoryName)
			if config.Get().IsDebugging() {
				if label != nil {
					log.Printf("[createWipLabelIfNotExist] Response body: %v\n", label)
				} else {
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						panic(err.Error())
					}
					log.Printf("[createWipLabelIfNotExist] Response body: %s\n", body)
				}
			}
		} else {
			panic(err)
		}
	}
}
