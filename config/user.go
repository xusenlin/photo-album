package config

type UserConfig struct {
	ListenPort      string `yaml:"listenPort"`
	CoverHeight     int    `yaml:"coverHeight"`
	PhotoAlbumPath  string `yaml:"photoAlbumPath"`
	SiteName        string `yaml:"siteName"`
	Author          string `yaml:"author"`
	HtmlKeywords    string `yaml:"htmlKeywords"`
	HtmlDescription string `yaml:"htmlDescription"`
}
