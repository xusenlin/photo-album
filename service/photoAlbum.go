package service

import (
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

func InitPhotoAlbum(root string) (*models.PhotoAlbums, *models.PhotoAlbumMap, error) {

	var photoAlbumMap = make(models.PhotoAlbumMap)

	photoAlbums, err := readPhotoAlbum(root)

	if err != nil {
		return nil, nil, err
	}
	sort.Sort(photoAlbums)

	for i := len(photoAlbums) - 1; i >= 0; i-- {
		//这里必须使用倒序的方式生成 PhotoAlbumMap,因为如果有相同的相册标题，
		// 倒序会将最老的相册优先生成shortUrl，保证和之前的 shortUrl一样
		photoAlbum := photoAlbums[i]
		keyword := utils.GenerateShortUrl(photoAlbum.Title, func(url, keyword string) bool {
			//保证 keyword 唯一
			_, ok := photoAlbumMap[keyword]
			return !ok
		})
		photoAlbums[i].ShortUrl = keyword
		photoAlbumMap[keyword] = photoAlbum.Path
	}
	return &photoAlbums, &photoAlbumMap, nil
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
			//比较耗时
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

	return pa
}

func parserPhotos(path string) (models.Photos, error) {
	var photos models.Photos

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return photos, err
	}
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			// 检查文件是否为图片格式,并且排除封面
			if strings.HasPrefix(file.Name(), "cover_") {
				continue
			}
			if ext == ".jpg" || ext == ".jpeg" {
				// 解析图片元数据
				photo := parsePhotoData(path, file)
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
func parsePhotoData(path string, file fs.FileInfo) models.Photo {

	photo := models.Photo{}
	coverName := "cover_" + file.Name()
	filePath := filepath.Join(path, file.Name())

	img, err := imaging.Open(filePath)
	if err != nil {
		photo.Error = err
		return photo
	}
	coverFile := filepath.Join(path, coverName)

	err = buildPhotoCover(img, coverFile)
	if err != nil {
		photo.Error = err
		return photo
	}

	// 在这里进行图片元数据的解析
	// 使用标准库或第三方库来解析图片元数据，并将其存储在一个 Photo 对象中

	// 示例：创建一个空的 Photo 对象

	// 示例：设置一些示意性的数据
	//photo.Name = filepath.Base(filePath)
	//photo.Format = strings.TrimPrefix(filepath.Ext(filePath), ".")
	//photo.ShotTime = time.Now()

	return photo
}

func buildPhotoCover(img image.Image, coverFile string) error {
	if utils.IsFile(coverFile) { //有封面了
		img, err := imaging.Open(coverFile)
		if err != nil {
			return err
		}
		if img.Bounds().Dy() == global.Config.CoverHeight {
			return nil
		}
	}
	cover := imaging.Resize(img, 0, global.Config.CoverHeight, imaging.NearestNeighbor)
	err := imaging.Save(cover, coverFile)
	if err != nil {
		return err
	}
	return nil
}
