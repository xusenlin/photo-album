package models

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"math"
	"os"
	"photoAlbum/pkg/dateTime"
	"photoAlbum/pkg/utils"
	"strings"
	"time"
)

type Photo struct {
	Size         string // 图片大小（以M为单位）
	Name         string // 图片名称
	Format       string // 图片格式（例如：JPEG、PNG等）
	Path         PhotoAlbumPath
	Height       int
	Width        int
	ShotTime     time.Time // 拍摄时间
	Camera       string
	CameraModel  string // 相机型号
	ExposureTime string // 曝光时间  快门速度
	Aperture     string // 光圈值
	ISO          string // ISO感光度
	FocalLength  string // 焦距毫米
	Error        error  //解析错误的信息
}

type Photos []Photo

type PhotoAlbumPath string //相册相对路径可能是(.)表示一级目录

// PhotoAlbum 当前目录存在yaml则视为一个相册集
type PhotoAlbum struct {
	Title        string            `yaml:"title"`
	Author       string            `yaml:"author"`
	CreatedAt    dateTime.DateTime `yaml:"createdAt"`
	Descriptions string            `yaml:"descriptions"`
	Path         PhotoAlbumPath
	Photos       Photos
	Count        int
	Error        error //解析错误的信息
}

type PhotoAlbums []PhotoAlbum

func (a PhotoAlbums) Pagination(pageNum, pageSize int) (PhotoAlbums, int, int, []int) {

	l := len(a)
	totalPages := int(math.Ceil(float64(l) / float64(pageSize)))

	if pageNum > totalPages {
		pageNum = totalPages
	}

	startIndex := (pageNum - 1) * pageSize
	endIndex := startIndex + pageSize

	if endIndex > l {
		endIndex = l
	}

	return a[startIndex:endIndex], pageNum, pageSize, utils.SpreadDigit(totalPages)
}

func (a PhotoAlbums) Len() int { return len(a) }

func (a PhotoAlbums) Less(i, j int) bool {
	return a[i].CreatedAt.Time.After(a[j].CreatedAt.Time)
}

func (a PhotoAlbums) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (p *Photo) ParseExifByPath(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		p.Error = err
		return
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		p.Error = err
		return
	}

	if t, err := x.DateTime(); err == nil {
		p.ShotTime = t
	}
	if cam, err := x.Get(exif.Make); err == nil {
		p.Camera = strings.Trim(cam.String(), "\"")
	}
	if camModel, err := x.Get(exif.Model); err == nil {
		p.CameraModel = strings.Trim(camModel.String(), "\"")
	}
	if exposure, err := x.Get(exif.ExposureTime); err == nil {
		numerator, denominator, _ := exposure.Rat2(0)
		p.ExposureTime = fmt.Sprintf("1/%v", denominator/numerator)

	}

	if aperture, err := x.Get(exif.FNumber); err == nil {
		numerator, denominator, _ := aperture.Rat2(0)
		p.Aperture = fmt.Sprintf("%.1f", float64(numerator)/float64(denominator))

	}

	if iso, err := x.Get(exif.ISOSpeedRatings); err == nil {
		p.ISO = iso.String()
	}

	if focalLength, err := x.Get(exif.FocalLength); err == nil {
		numerator, denominator, _ := focalLength.Rat2(0)
		p.FocalLength = fmt.Sprintf("%.0f", float64(numerator)/float64(denominator))
	}
	//fmt.Println(p)
}
