package upload

type WebPUpload struct {
	OriginalFileName string   `json:"originalFileName"`
	FileName         string   `json:"fileName"`
	Variants         []string `json:"variants"`
	FullPaths        []string `json:"fullPaths"`
}
