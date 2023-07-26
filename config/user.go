package config

type UserConfig struct {
	ListenPort      string `yaml:"listenPort"`
	CoverHeight     int    `yaml:"coverHeight"`
	PhotoAlbumPath  string `yaml:"photoAlbumPath"`
	SiteName        string `yaml:"siteName"`
	HtmlKeywords    string `yaml:"htmlKeywords"`
	HtmlDescription string `yaml:"htmlDescription"`
}
