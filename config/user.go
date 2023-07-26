package config

type UserConfig struct {
	ListenPort     string `yaml:"listenPort"`
	CoverHeight    int    `yaml:"coverHeight"`
	PhotoAlbumPath string `yaml:"photoAlbumPath"`
}
