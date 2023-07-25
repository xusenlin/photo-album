package models

import (
	"github.com/disintegration/imaging"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"photoAlbum/pkg/utils"
	"sort"
	"strings"
	"sync"
)

var PhotoAlbumList PhotoAlbums
var PhotoAlbumMap = make(map[string]PhotoAlbumPath) //用来保证相册 shortUrl 唯一和快速定位相册

func InitPhotoAlbum(root string) error {

	photoAlbum, err := readPhotoAlbum(root)

	if err != nil {
		return err
	}
	PhotoAlbumList = *photoAlbum

	sort.Sort(PhotoAlbumList)
	for i := len(PhotoAlbumList) - 1; i >= 0; i-- {
		//这里必须使用倒序的方式生成 PhotoAlbumMap,因为如果有相同的相册标题，
		// 倒序会将最老的相册优先生成shortUrl，保证和之前的 shortUrl一样
		PhotoAlbum := PhotoAlbumList[i]
		keyword := utils.GenerateShortUrl(PhotoAlbum.Title, func(url, keyword string) bool {
			//保证 keyword 唯一
			_, ok := PhotoAlbumMap[keyword]
			return !ok
		})
		PhotoAlbumList[i].ShortUrl = keyword
		PhotoAlbumMap[keyword] = PhotoAlbum.Path
	}
	return nil
}

func readPhotoAlbum(absolutePath string) (*PhotoAlbums, error) {

	resultCh := make(chan *PhotoAlbum, 10)

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
			pa, err := parserPhotoAlbum(absolutePath, path)
			if err == nil {
				resultCh <- &pa
			}
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

	var photoAlbums PhotoAlbums
	for pa := range resultCh {
		photoAlbums = append(photoAlbums, *pa)
	}

	return &photoAlbums, nil
}

func parserPhotoAlbum(root, path string) (PhotoAlbum, error) {
	pa, err := parserYaml(path)
	if err != nil {
		return pa, err
	}
	dir := filepath.Dir(path)
	relPath, err := filepath.Rel(root, dir)
	if err != nil {
		return pa, err
	}
	pa.Photos, err = parserPhotos(dir)
	if err != nil {
		return pa, err
	}
	pa.Path = PhotoAlbumPath(relPath)

	return pa, nil
}

func parserPhotos(path string) (Photos, error) {
	var photos Photos

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return photos, err
	}
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			// 检查文件是否为图片格式
			if ext == ".jpg" || ext == ".jpeg" {
				// 解析图片元数据
				photo, err := parsePhotoData(path, file)
				if err != nil {
					// 图片解析失败，跳过该图片
					continue
				}
				photos = append(photos, photo)
			}
		}
	}

	return photos, nil
}

func parserYaml(path string) (PhotoAlbum, error) {
	var photoAlbum PhotoAlbum
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
func parsePhotoData(path string, file fs.FileInfo) (Photo, error) {

	photo := Photo{}
	filePath := filepath.Join(path, file.Name())
	//fmt.Println(filePath)

	img, err := imaging.Open(filePath)
	if err != nil {
		return photo, err
	}
	cover := imaging.Resize(img, 0, 600, imaging.NearestNeighbor)

	err = imaging.Save(cover, filepath.Join(path, "cover_"+file.Name()))
	if err != nil {
		return photo, err
	}
	// 在这里进行图片元数据的解析
	// 使用标准库或第三方库来解析图片元数据，并将其存储在一个 Photo 对象中

	// 示例：创建一个空的 Photo 对象

	// 示例：设置一些示意性的数据
	//photo.Name = filepath.Base(filePath)
	//photo.Format = strings.TrimPrefix(filepath.Ext(filePath), ".")
	//photo.ShotTime = time.Now()

	return photo, nil
}
