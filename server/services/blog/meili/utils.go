package meili

import "regexp"

// stripMarkdown removes basic Markdown symbols to avoid interfering with full-text search tokenization.
func stripMarkdown(md string) string {
	re := regexp.MustCompile(`[#*_\-\[\]()>]`)
	return re.ReplaceAllString(md, " ")
}
