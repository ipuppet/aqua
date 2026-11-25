package openlist

// FsObject represents a file or directory in OpenList list response.
type FsObject struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	Modified string `json:"modified"`
	Type     int    `json:"type"` // 0=Unknown, 1=Folder, 2=Video, 3=Audio, 4=Text, 5=Image
}

// ListRequest is the request body for the list API.
type ListRequest struct {
	Path     string `json:"path"`
	Password string `json:"password,omitempty"`
	Refresh  bool   `json:"refresh,omitempty"`
	Page     int    `json:"page,omitempty"`
	PerPage  int    `json:"per_page,omitempty"`
}

// ListResponse is the response from List API.
type ListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Content []FsObject `json:"content"`
		Total   int        `json:"total"`
	} `json:"data"`
}

// RemoveRequest is the request body for Remove.
type RemoveRequest struct {
	Dir   string   `json:"dir"`
	Names []string `json:"names"`
}
