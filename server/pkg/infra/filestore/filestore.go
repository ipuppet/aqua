package filestore

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

// Client provides file store operations with an explicit base path.
type Client struct {
	BasePath string
}

// NewClient returns a Client with the given base path.
func NewClient(basePath string) *Client {
	return &Client{BasePath: basePath}
}

func (c *Client) path(app, fileName string) string {
	base := c.BasePath
	if base == "" {
		base, _ = os.Getwd()
		base = filepath.Join(base, "storage")
	}
	return filepath.Join(base, app, fileName)
}

// GetContent reads and returns the string content of a stored file.
func (c *Client) GetContent(app, fileName string) (string, error) {
	content, err := os.ReadFile(c.path(app, fileName))
	return string(content), err
}

// GetJSON reads a stored JSON file and unmarshals it into v.
func (c *Client) GetJSON(app, fileName string, v any) error {
	content, err := os.ReadFile(c.path(app, fileName))
	if err != nil {
		return err
	}
	return json.Unmarshal(content, v)
}

// SetContent writes the given content to a stored file, overwriting existing data.
func (c *Client) SetContent(app, fileName, content string) error {
	return os.WriteFile(c.path(app, fileName), []byte(content), 0666)
}

// AppendContent appends the given content as a new line to a stored file.
func (c *Client) AppendContent(app, fileName, content string) error {
	file, err := os.OpenFile(c.path(app, fileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return writer.Flush()
}
