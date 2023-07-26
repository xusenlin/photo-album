package global

import "photoAlbum/models"

var PhotoAlbumList *models.PhotoAlbums
var PhotoAlbumMap *models.PhotoAlbumMap //用来保证相册 shortUrl 唯一和快速定位相册
