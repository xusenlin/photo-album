package models

import (
	"errors"
	"os"
	"photoAlbum/pkg/utils"
	"sort"
)

var PhotoAlbumList PhotoAlbums
var PhotoAlbumMap = make(map[string]PhotoAlbumPath) //用来保证相册 shortUrl 唯一和快速定位相册

func InitPhotoAlbum(root string) error {
	var err error
	PhotoAlbumList, err = RecursiveReadPhotoAlbum(root, "/")
	if err != nil {
		return err
	}
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

func RecursiveReadPhotoAlbum(root, path string) (PhotoAlbums, error) {
	var photoAlbums PhotoAlbums
	dirInfo, err := os.Stat(path)

	if err != nil {
		return photoAlbums, err
	}

	if !dirInfo.IsDir() {
		return photoAlbums, errors.New("目标不是一个目录")
	}

	return photoAlbums, nil
}
