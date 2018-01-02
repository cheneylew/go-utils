package utils

import (
	"os"
	"image/jpeg"
	"github.com/nfnt/resize"
	"path"
	"strings"
	"image"
	"image/png"
	"image/gif"
	"path/filepath"
)

const (
	ImageType_JPG = 1 + iota
	ImageType_PNG
	ImageType_GIF
)

/*图片缩略图*/
func ImageThumbnail(orgImagePath string, width uint) (thumbPath string, err error) {
	dir := path.Dir(orgImagePath)
	imageName := strings.Replace(orgImagePath, dir+"/", "", -1)
	thumbnailName := "thumbnail_" + imageName
	thumbnailPath := path.Join(dir, thumbnailName)

	ext := strings.Replace(strings.ToLower(path.Ext(orgImagePath)), ".", "", -1)
	imageType := ImageType_JPG
	switch ext {
	case "jpg", "jpeg":
		imageType = ImageType_JPG
	case "png":
		imageType = ImageType_PNG
	case "gif":
		imageType = ImageType_GIF
	}

	// open "test.jpg"
	file, err := os.Open(orgImagePath)
	if err != nil {
		return "", err
	}

	// decode jpeg into image.Image
	var img image.Image = nil
	var decodeErr error
	switch imageType {
	case ImageType_JPG:
		img, decodeErr = jpeg.Decode(file)
	case ImageType_PNG:
		img, decodeErr = png.Decode(file)
	case ImageType_GIF:
		img, decodeErr = gif.Decode(file)

	}
	if err != nil {
		return "", decodeErr
	}
	file.Close()


	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(width, 0, img, resize.Lanczos3)

	out, createErr := os.Create(thumbnailPath)
	if err != nil {
		return "", createErr
	}
	defer out.Close()

	// write new image to file
	switch imageType {
	case ImageType_JPG:
		jpeg.Encode(out, m, nil)
	case ImageType_PNG:
		png.Encode(out, m)
	case ImageType_GIF:
		gif.Encode(out, m, nil)
	}

	return thumbnailPath, nil
}

func ThumbnailPath(rawPath string) string {
	dir := Dir(rawPath)
	name := filepath.Base(rawPath)
	return dir+"/"+"thumbnail_"+name
}