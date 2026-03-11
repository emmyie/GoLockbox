package vault

type FileEntry struct {
	Name string `json:"name"`
}

type metadata struct {
	Version int         `json:"version"`
	Entries []FileEntry `json:"entries"`
}
