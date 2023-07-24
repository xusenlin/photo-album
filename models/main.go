package models

import (
	"photoAlbum/pkg/dateTime"
)

type Photo struct {
	Size         int               // 图片大小（以字节为单位）
	Name         string            // 图片名称
	Format       string            // 图片格式（例如：JPEG、PNG等）
	ShotTime     dateTime.DateTime // 拍摄时间
	Location     string            // 拍摄地点
	CameraModel  string            // 相机型号
	ExposureTime string            // 曝光时间
	Aperture     string            // 光圈值
	ISO          int               // ISO感光度
	FocalLength  string            // 焦距
	ShutterSpeed string            // 快门速度
}

type Photos []Photo

type PhotoAlbumPath string //相册相对路径

// PhotoAlbum 主目录下一个目录视为一个相册
type PhotoAlbum struct {
	Title        string
	Author       string
	CreatedAt    dateTime.DateTime
	Descriptions string
	ShortUrl     string //快捷URL
	Path         PhotoAlbumPath
	Photos       Photos
}

type PhotoAlbums []PhotoAlbum

func (a PhotoAlbums) Len() int { return len(a) }

func (a PhotoAlbums) Less(i, j int) bool {
	return a[i].CreatedAt.Time.After(a[j].CreatedAt.Time)
}

func (a PhotoAlbums) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
