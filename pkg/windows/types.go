package windows

type FileVersionInfo struct {
	FilePath       string `json:"file_path"`
	ProductName    string `json:"product_name"`
	CompanyName    string `json:"company_name"`
	FileVersion    string `json:"file_version"`
	ProductVersion string `json:"product_version"`
}
