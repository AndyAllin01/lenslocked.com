package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	//	"strings"
)

//this is not stored in the database
type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
	//return "/" + i.RelativePath()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)

}

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}
	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

//commented out temporarily due to repeating image url
/*func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	fmt.Println("################## 1", path)
	imgStrings, err := filepath.Glob(filepath.Join(path, "*"))
	fmt.Println("################## 2", imgStrings)
	if err != nil {
		return nil, err
	}
	ret := make([]Image, len(imgStrings))
	for i := range imgStrings {
		fmt.Println("################## 3", imgStrings[i])

		imgStrings[i] = filepath.ToSlash("/" + imgStrings[i])
		fmt.Println("################## 4", imgStrings[i])
		imgStrings[i] = strings.Replace(imgStrings[i], path, "", 1)
		fmt.Println("################## 5", imgStrings[i])
		ret[i] = Image{
			Filename:  imgStrings[i],
			GalleryID: galleryID,
		}
		fmt.Println("################## 6", ret[i])

	}
	return ret, nil
}*/

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	// Setup the Image slice we are returning
	ret := make([]Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = Image{
			Filename:  filepath.Base(imgStr),
			GalleryID: galleryID,
		}
	}
	return ret, nil
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

// Going to need this when we know it is already made
func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries", fmt.Sprintf("%v", galleryID))
}

// Use the imagePath method we just made
func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
