package upload

type WebPUpload struct {
	OriginalFileName string   `json:"originalFileName"`
	FileName         string   `json:"fileName"`
	Variants         []string `json:"variants"`
	FullPaths        []string `json:"fullPaths"`
}

type Variant = string

const (
	// full quality, low compression
	FullQuality Variant = "x1"

	// half quality, medium compression
	HalfQuality Variant = "x2"

	// third quality, high compression
	ThirdQuality Variant = "x3"
)
