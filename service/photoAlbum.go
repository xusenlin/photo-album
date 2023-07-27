package service

import (
	"fmt"
	"github.com/disintegration/imaging"
	"gopkg.in/yaml.v3"
	"image"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"photoAlbum/global"
	"photoAlbum/models"
	"photoAlbum/pkg/utils"
	"sort"
	"strings"
	"sync"
)

func InitPhotoAlbum(root string) (*models.PhotoAlbums, error) {

	photoAlbums, err := readPhotoAlbum(root)

	if err != nil {
		return nil, err
	}
	sort.Sort(photoAlbums)

	return &photoAlbums, nil
}

func readPhotoAlbum(absolutePath string) (models.PhotoAlbums, error) {

	resultCh := make(chan *models.PhotoAlbum, 10)

	wg := sync.WaitGroup{}

	err := filepath.Walk(absolutePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !utils.IsYamlFile(info) {
			return nil
		}
		//因为有Yaml文件,当前被认为是一个相册集
		//如果同一个文件夹有多个Yaml，也会被当多个相册集
		wg.Add(1)
		go func() {
			defer wg.Done()
			pa := parserPhotoAlbum(absolutePath, path)
			resultCh <- &pa
		}()
		return nil
	})

	if err != nil {
		return nil, err
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var photoAlbums models.PhotoAlbums
	for pa := range resultCh {
		photoAlbums = append(photoAlbums, *pa)
	}

	return photoAlbums, nil
}

func parserPhotoAlbum(root, path string) models.PhotoAlbum {
	pa, err := parserYaml(path)
	if err != nil {
		pa.Error = err
		return pa
	}
	dir := filepath.Dir(path)
	relPath, err := filepath.Rel(root, dir)
	if err != nil {
		pa.Error = err
		return pa
	}
	pa.Photos, err = parserPhotos(dir)
	if err != nil {
		pa.Error = err
		return pa
	}
	pa.Path = models.PhotoAlbumPath(relPath)
	pa.Count = len(pa.Photos)

	return pa
}

func parserPhotos(dir string) (models.Photos, error) {
	var photos models.Photos

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return photos, err
	}
	for _, file := range files {
		if !file.IsDir() {
			if strings.Contains(file.Name(), "_COVER") {
				continue
			}
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" {
				// 解析图片元数据
				photo := parsePhotoData(dir, file)
				photos = append(photos, photo)
			}
		}
	}

	return photos, nil
}

func parserYaml(path string) (models.PhotoAlbum, error) {
	var photoAlbum models.PhotoAlbum
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return photoAlbum, err
	}
	err = yaml.Unmarshal(content, &photoAlbum)
	if err != nil {
		return photoAlbum, err
	}
	return photoAlbum, nil
}

// 解析图片元数据
func parsePhotoData(dir string, file fs.FileInfo) models.Photo {

	photo := models.Photo{}

	filePath := filepath.Join(dir, file.Name())
	coverPath := buildCoverPath(filePath)

	img, err := imaging.Open(filePath)
	if err != nil {
		photo.Error = err
		return photo
	}

	photo.Width = img.Bounds().Dx()
	photo.Height = img.Bounds().Dy()
	photo.Size = fmt.Sprintf("%.2f", float64(file.Size())/(1024*1024))
	photo.Format = filepath.Ext(file.Name())
	photo.Name = strings.TrimSuffix(file.Name(), photo.Format)

	photo.ParseExifByPath(filePath)

	err = buildPhotoCover(img, coverPath)
	if err != nil {
		photo.Error = err
		return photo
	}

	return photo
}

func buildPhotoCover(img image.Image, coverPath string) error {
	if utils.IsFile(coverPath) { //有封面了
		img, err := imaging.Open(coverPath)
		if err != nil {
			return err
		}
		if img.Bounds().Dy() == global.Config.CoverHeight {
			return nil
		}
	}
	cover := imaging.Resize(img, 0, global.Config.CoverHeight, imaging.NearestNeighbor)
	err := imaging.Save(cover, coverPath)
	if err != nil {
		return err
	}
	fmt.Println(coverPath)
	return nil
}

func buildCoverPath(path string) string {
	dir := filepath.Dir(path)
	fileName := filepath.Base(path)
	fileExt := filepath.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	coverFileName := fileNameWithoutExt + "_COVER" + fileExt
	return filepath.Join(dir, coverFileName)
}
