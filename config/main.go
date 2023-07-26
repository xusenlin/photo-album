package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	UserConfig
	SystemConfig
}

func New(path string) (*Config, error) {
	var userConfig UserConfig
	var systemConfig SystemConfig

	y, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(y, &userConfig)
	if err != nil {
		return nil, err
	}

	photoAlbumAbsolutePath, err := filepath.Abs(userConfig.PhotoAlbumPath)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	systemConfig.PhotoAlbumAbsolutePath = photoAlbumAbsolutePath

	return &Config{userConfig, systemConfig}, nil

}

//func (c *Config) name()  {
//
//}
