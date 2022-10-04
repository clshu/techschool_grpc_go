package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

// ImageStore is an interface for storing images.
type ImageStore interface {
	Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore is an implementation of ImageStore that saves images to disk.
type DiskImageStore struct {
	mutex sync.RWMutex
	imageFolder string
	images map[string]*ImageInfo
}

// ImageInfo contains information about an image.
type ImageInfo struct {
	LaptopID string
	Type string
	Path string
}

// NewDiskImageStore creates a new DiskImageStore.
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images: make(map[string]*ImageInfo),
	}
}

// Save saves an image to disk.
func (store *DiskImageStore) Save(
		laptopID string,
		imageType string,
		imageData bytes.Buffer,
	) (string, error) {
		
		imageID, err := uuid.NewRandom()
		if err != nil {
			return "", fmt.Errorf("cannot generate image ID: %w", err)
		}

		createDirIfNotExist(store.imageFolder)

		imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)

		file, err := os.Create(imagePath)
		if err != nil {
			return "", fmt.Errorf("cannot create image file: %w", err)
		}      

		_, err = imageData.WriteTo(file)
		if err != nil {
			return "", fmt.Errorf("cannot write image data to file: %w", err)
		}

		store.mutex.Lock()
		defer store.mutex.Unlock()

		store.images[imageID.String()] = &ImageInfo{
			LaptopID: laptopID,
			Type: imageType,
			Path: imagePath,
		}

		return imageID.String(), nil
}

// saveImageToFile saves an image to a file.
func saveImageToFile(folder string, laptopID string, imageType string, imageData []byte) (string, error) {
	return "", nil
}

func createDirIfNotExist(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}
}