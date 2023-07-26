package models

import (
	"photoAlbum/pkg/dateTime"
)

type Photo struct {
	Size         int    // 图片大小（以字节为单位）
	Name         string // 图片名称
	Format       string // 图片格式（例如：JPEG、PNG等）
	Height       int
	Width        int
	ShotTime     dateTime.DateTime // 拍摄时间
	CameraModel  string            // 相机型号
	ExposureTime string            // 曝光时间  快门速度
	Aperture     string            // 光圈值
	ISO          int               // ISO感光度
	FocalLength  string            // 焦距
	Error        error             //解析错误的信息
}

type Photos []Photo

type PhotoAlbumPath string //相册相对路径可能是(.)表示一级目录

// PhotoAlbum 当前目录存在yaml则视为一个相册集
type PhotoAlbum struct {
	Title        string            `yaml:"title"`
	Author       string            `yaml:"author"`
	CreatedAt    dateTime.DateTime `yaml:"createdAt"`
	Descriptions string            `yaml:"descriptions"`
	ShortUrl     string            //快捷URL
	Path         PhotoAlbumPath
	Photos       Photos
	Error        error //解析错误的信息
}

type PhotoAlbums []PhotoAlbum

func (a PhotoAlbums) Len() int { return len(a) }

func (a PhotoAlbums) Less(i, j int) bool {
	return a[i].CreatedAt.Time.After(a[j].CreatedAt.Time)
}

func (a PhotoAlbums) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
