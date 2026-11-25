package imagestore

import "time"

// ImageInfo contains detailed information about an uploaded image.
type ImageInfo struct {
	Path    string    `json:"path"`    // Relative path (e.g. module/date/name)
	Name    string    `json:"name"`    // File name
	URL     string    `json:"url"`     // Full URL path
	Size    int64     `json:"size"`    // File size in bytes
	ModTime time.Time `json:"modTime"` // Modification time
}
