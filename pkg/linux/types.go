package linux

type FileVersionInfo struct {
	FilePath     string   `json:"file_path"`
	FileType     string   `json:"file_type"`
	SONAME       string   `json:"soname,omitempty"`
	BuildID      string   `json:"build_id,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	Version      string   `json:"version,omitempty"`
}
