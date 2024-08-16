package store

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageStore struct {
	dir   string
	index int
}

func NewImageStore(dir string) *ImageStore {
	os.MkdirAll(dir, os.ModePerm)
	return &ImageStore{dir, -1}
}

func (is *ImageStore) ListFilenames() ([]string, error) {
	files, err := os.ReadDir(is.dir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jpg") {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func (is *ImageStore) Save(image image.Image) (string, error) {
	filename := fmt.Sprintf("%v.jpg", time.Now().Unix())
	file, err := os.Create(filepath.Join(is.dir, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := jpeg.Encode(file, image, nil); err != nil {
		return "", err
	}

	return filename, nil
}

func (is *ImageStore) Load(filename string) (image.Image, error) {
	f, err := os.Open(filepath.Join(is.dir, filename))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return img, nil
}

func (is *ImageStore) LoadNext() (image.Image, string, error) {
	filenames, _ := is.ListFilenames()
	if len(filenames) == 0 {
		return nil, "", errors.New("image folder is empty")
	}
	filename := filenames[(is.index+len(filenames))%len(filenames)]
	img, err := is.Load(filename)
	is.index++
	if is.index >= len(filenames) {
		is.index = 0
	}
	return img, filename, err
}
