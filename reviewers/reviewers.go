package reviewers

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
	var largestOwner string

	for owner, set := range ownerSets {
		if len(set) > maxSize {
			maxSize = len(set)
			largestOwner = owner
		}
	}

	return largestOwner
}

func removeFromOwnership(ownerSets map[string]map[string]bool, files []string) {
	for owner, _ := range ownerSets {
		for _, file := range files {
			delete(ownerSets[owner], file)
		}
	}
}

func allFilesCovered(ownerSets map[string]map[string]bool) bool {
	for _, set := range ownerSets {
		if len(set) > 0 {
			return false
		}
	}

	return true
}

func suggestReviewers(fileToOwners map[string][]string) []string {
	ownersToFiles := createOwnerFilesList(fileToOwners)
	ownersToOwnership := createOwnersSets(fileToOwners)

	var reviewers []string

	for !allFilesCovered(ownersToOwnership) {
		reviewer := largestOwnership(ownersToOwnership)
		reviewers = append(reviewers, reviewer)

		removeFromOwnership(ownersToOwnership, ownersToFiles[reviewer])
	}

	return reviewers
}
