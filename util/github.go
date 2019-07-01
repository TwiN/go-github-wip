package util

import (
	"context"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v25/github"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	CheckRunName = "Work in Progress"
)

func GetGithubClient(appId, installationId int64) (*github.Client, context.Context) {
	return createGithubClient(appId, installationId), context.Background()
}

func createGithubClient(appId, installationId int64) *github.Client {
	transport := http.DefaultTransport
	itr, err := ghinstallation.NewKeyFromFile(transport, int(appId), int(installationId), os.Getenv("GO_GITHUB_WIP_APP_PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	return github.NewClient(&http.Client{Transport: itr})
}

func SetAsWip(userName, repositoryName, branch, commit string, appId, installationId int64) *github.CheckRun {
	log.Printf("[SetAsWip] Creating WIP CheckRun on branch %s from repository %s/%s to WIP", branch, userName, repositoryName)
	client, ctx := GetGithubClient(appId, installationId)
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

	log.Printf("[SetAsWip] Response status: %s\n", resp.Status)
	if checkRun != nil {
		log.Printf("[SetAsWip] Response body: %v\n", checkRun)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[SetAsWip] Response body: %s\n", body)
	}

	if err != nil {
		panic(err.Error())
	}

	return checkRun
}

func ClearWip(userName, repositoryName, branch, commit string, appId, installationId int64, checkRunId int64) {
	log.Printf("[ClearWip] Clearing WIP CheckRun on branch %s from repository %s/%s to WIP", branch, userName, repositoryName)
	client, ctx := GetGithubClient(appId, installationId)
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

	log.Printf("[ClearWip] Response status: %s\n", resp.Status)
	if checkRun != nil {
		log.Printf("[ClearWip] Response body: %v\n", checkRun)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[ClearWip] Response body: %s\n", body)
	}

	if err != nil {
		panic(err.Error())
	}
}

func GetCheckRunId(owner, repository, branch string, appId, installationId int64) int64 {
	client, ctx := GetGithubClient(appId, installationId)

	checkRunName := CheckRunName
	checkRun, resp, err := client.Checks.ListCheckRunsForRef(
		ctx,
		owner,
		repository,
		branch,
		&github.ListCheckRunsOptions{
			CheckName: &checkRunName,
		},
	)

	log.Printf("[GetCheckRunId] Response status: %s\n", resp.Status)
	if checkRun != nil {
		log.Printf("[GetCheckRunId] Response body: %v\n", checkRun)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[GetCheckRunId] Response body: %s\n", body)
	}

	if err != nil {
		panic(err.Error())
	}

	return *checkRun.CheckRuns[0].ID
}
