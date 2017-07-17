package file

type Owner struct {
	User    string
	Pattern string
}

type OwnersFile struct {
	Flags    map[string]bool
	Ignored  []string
	Owners   []*Owner
	Watchers []*Owner
}
