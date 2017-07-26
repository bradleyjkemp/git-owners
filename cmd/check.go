package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

var githubPR bool

type GithubReviews []struct {
	User struct {
		Login string
	}
	State string
}

// checkCmd represents the check command
var CheckCmd = &cobra.Command{
	Use:   "check reviewer1 reviewer2",
	Short: "Check a PR has been accepted by all necessary owners",
	Long: `Given a list of reviewers this checks that every file modified on this branch has at least one owner who has accepted the PR
	
If the --github flag is used then the $TRAVIS_REPO_SLUG and $TRAVIS_PULL_REQUEST
environment variables will be used to automatically retrieve usernames of approved reviewers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var approvers []string
		if githubPR {
			approvers, err := getGithubReviewers()
			fmt.Println(approvers)
			return err
		} else {
			approvers = args
		}

		filesToOwners, err := getFilesToOwnersForPR()
		if err != nil {
			return errors.Wrap(err, "failed to get owners for changed files")
		}

		if allFilesCovered(filesToOwners, approvers) {
			fmt.Println("All files approved by at least one owner")
			return nil
		}

		fmt.Println("At least one file not approved by an owner")
		os.Exit(1)
		return nil
	},
}

func init() {
	CheckCmd.Flags().BoolVarP(&githubPR, "github", "g", false, "Automatically retrieve acceptances from github")
	RootCmd.AddCommand(CheckCmd)
}

func allFilesCovered(fileToOwners map[string][]string, approvers []string) bool {
	var approverSet = make(map[string]bool)
	for _, approver := range approvers {
		approverSet[approver] = true
	}

FileLoop:
	for _, owners := range fileToOwners {
		if len(owners) == 0 {
			continue
		}

		for _, owner := range owners {
			if approverSet[owner] {
				continue FileLoop
			}
		}

		return false
	}

	return true
}

func getGithubReviewers() ([]string, error) {
	repoSlug := os.Getenv("TRAVIS_REPO_SLUG")
	pullRequest := os.Getenv("TRAVIS_PULL_REQUEST")
	if repoSlug == "" || pullRequest == "" {
		return nil, errors.New("Repo slug or pull request id environment variable not set")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s/reviews", repoSlug, pullRequest)
	fmt.Println("getting", url)
	r, err := (&http.Client{Timeout: 10 * time.Second}).Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var reviews GithubReviews
	err = json.NewDecoder(r.Body).Decode(&reviews)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode github response")
	}
	fmt.Println(reviews)

	approvals := make(map[string]bool)
	for _, review := range reviews {
		approvals[review.User.Login] = review.State == "APPROVED"
	}

	var approvers []string
	for approver := range approvals {
		approvers = append(approvers, approver)
	}

	return approvers, nil
}
