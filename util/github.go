package util

import (
	"context"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v25/github"
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
	itr, err := ghinstallation.NewKeyFromFile(transport, int(appId), int(installationId), os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	return github.NewClient(&http.Client{Transport: itr})
}

func SetAsWip(userName, repositoryName, branch, commit string, appId, installationId int64) *github.CheckRun {
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

	log.Printf("[SetAsWip] %v\n", checkRun)

	if err != nil {
		panic(err.Error())
	}

	log.Printf("[SetAsWip] Response: %s\n", resp.Status)
	return checkRun
}

func ClearWip(userName, repositoryName, branch, commit string, appId, installationId int64, checkRunId int64) {
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

	log.Printf("[ClearWip] %v\n", checkRun)

	if err != nil {
		panic(err.Error())
	}

	log.Printf("[ClearWip] Response: %s\n", resp.Status)
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

	log.Printf("[GetCheckRunId] %v\n", checkRun)

	if err != nil {
		panic(err.Error())
	}

	log.Printf("[GetCheckRunId] Response: %s\n", resp.Status)
	return *checkRun.CheckRuns[0].ID
}
