package reviewers

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func allFilesCovered(fileToOwners map[string][]string, reviewers []string) bool {
	var reviewerSet = make(map[string]bool)
	for _, reviewer := range reviewers {
		reviewerSet[reviewer] = true
	}

FileLoop:
	for _, owners := range fileToOwners {
		if len(owners) == 0 {
			continue
		}

		for _, owner := range owners {
			if reviewerSet[owner] {
				continue FileLoop
			}
		}

		return false
	}

	return true
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomName(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestSetCoverAlgorithm(t *testing.T) {
	owners := []string{"alice", "bob", "carol", "dave", "eve", "fred"}
	distribution := []int{0, 1, 1, 2, 2, 2, 3}

	for rounds := 0; rounds < 100; rounds++ {
		fileToOwners := make(map[string][]string)

		for i := 0; i < 50; i++ {
			name := randomName(5)
			numReviewers := distribution[rand.Intn(len(distribution))]
			fileToOwners[name] = []string{}

			for j := 0; j < numReviewers; j++ {
				fileToOwners[name] = append(fileToOwners[name], owners[rand.Intn(len(owners))])
			}
		}

		reviewers := SuggestReviewers(fileToOwners)
		assert.True(t, allFilesCovered(fileToOwners, reviewers))
	}
}
