package reviewers

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func createOwnersSets(fileToOwners map[string][]string) map[string]map[string]bool {
	ownersToOwnership := make(map[string]map[string]bool)
	for file, owners := range fileToOwners {
		for _, owner := range owners {
			if ownersToOwnership[owner] == nil {
				ownersToOwnership[owner] = make(map[string]bool)
			}

			ownersToOwnership[owner][file] = true
		}
	}

	return ownersToOwnership
}

func createOwnerFilesList(fileToOwners map[string][]string) map[string][]string {
	ownersToFiles := make(map[string][]string)
	for file, owners := range fileToOwners {
		for _, owner := range owners {
			ownersToFiles[owner] = append(ownersToFiles[owner], file)
		}
	}

	return ownersToFiles
}

func largestOwnership(ownerSets map[string]map[string]bool) string {
	var maxSize int
	var largestOwners []string

	for owner, set := range ownerSets {
		if len(set) > maxSize {
			maxSize = len(set)
			largestOwners = []string{owner}
		} else if len(set) == maxSize {
			// To properly randomly pick between owners of the same size
			// we need to keep all of them and select a random one at the end
			largestOwners = append(largestOwners, owner)
		}
	}

	return largestOwners[r.Intn(len(largestOwners))]
}

func removeFromOwnership(ownerSets map[string]map[string]bool, files []string) {
	for owner, _ := range ownerSets {
		for _, file := range files {
			delete(ownerSets[owner], file)
		}
	}
}

func noFilesLeft(ownerSets map[string]map[string]bool) bool {
	for _, set := range ownerSets {
		if len(set) > 0 {
			return false
		}
	}

	return true
}

func SuggestReviewers(fileToOwners map[string][]string) []string {
	ownersToFiles := createOwnerFilesList(fileToOwners)
	ownersToOwnership := createOwnersSets(fileToOwners)

	var reviewers []string

	for !noFilesLeft(ownersToOwnership) {
		reviewer := largestOwnership(ownersToOwnership)
		reviewers = append(reviewers, reviewer)

		removeFromOwnership(ownersToOwnership, ownersToFiles[reviewer])
	}

	return reviewers
}
