package imagestore

import (
	"errors"
	"strings"
)

func validateSegment(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("path segment cannot be empty")
	}
	if strings.Contains(value, "..") {
		return errors.New("path segment cannot contain '..'")
	}
	if strings.ContainsAny(value, "/\\") {
		return errors.New("path segment cannot contain path separators")
	}
	return nil
}

// validatePath validates a relative path (e.g. module/date/name).
func validatePath(relPath string) error {
	if strings.TrimSpace(relPath) == "" {
		return errors.New("path cannot be empty")
	}
	if strings.HasPrefix(relPath, "/") || strings.HasPrefix(relPath, "\\") {
		return errors.New("path must be relative")
	}
	if strings.Contains(relPath, "..") {
		return errors.New("path cannot contain '..'")
	}
	for _, seg := range strings.Split(relPath, "/") {
		if seg == "" {
			continue
		}
		if err := validateSegment(seg); err != nil {
			return err
		}
	}
	return nil
}

// buildDirectURL returns the direct download URL for a file via the /d/ path.
// endpoint is the OpenList endpoint (e.g. https://example.com), pathPrefix is the store path prefix,
// relPath is the relative path (e.g. module/date/filename.jpg).
func buildDirectURL(endpoint, pathPrefix, relPath string) string {
	base := strings.TrimRight(endpoint, "/") + "/d/" + strings.Trim(pathPrefix, "/")
	return base + "/" + strings.TrimLeft(relPath, "/")
}
