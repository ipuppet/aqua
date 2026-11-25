package types

type FileStoreConfig struct {
	URLPrefix string          `mapstructure:"url_prefix" json:"url_prefix"`
	SavePath  string          `mapstructure:"save_path" json:"save_path"`
	Key       string          `mapstructure:"key" json:"key"`
	OpenList  *OpenListConfig `mapstructure:"openlist" json:"openlist,omitempty"`
}

// OpenListConfig holds credentials and endpoint for uploading to OpenList.
type OpenListConfig struct {
	Enable     bool   `mapstructure:"enable" json:"enable"`
	Username   string `mapstructure:"username" json:"username"`
	Password   string `mapstructure:"password" json:"password"`
	Endpoint   string `mapstructure:"endpoint" json:"endpoint"`
	PathPrefix string `mapstructure:"path_prefix" json:"path_prefix"`
}
