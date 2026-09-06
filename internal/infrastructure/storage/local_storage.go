package storage

import (
	"bytes"
	"image/jpeg"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
	}
}

func (s *LocalStorage) Save(file *multipart.FileHeader) (string, error) {

	err := os.MkdirAll(
		s.basePath,
		os.ModePerm,
	)

	if err != nil {
		return "", err
	}

	src, err := file.Open()

	if err != nil {
		return "", err
	}

	defer src.Close()

	img, err := imaging.Decode(src)

	if err != nil {
		return "", err
	}

	img = imaging.Fit(
		img,
		1200,
		1200,
		imaging.Lanczos,
	)

	filename := uuid.New().String() + ".jpg"

	path := filepath.Join(
		s.basePath,
		filename,
	)

	buffer := bytes.NewBuffer(nil)

	err = jpeg.Encode(
		buffer,
		img,
		&jpeg.Options{
			Quality: 85,
		},
	)

	if err != nil {
		return "", err
	}

	err = os.WriteFile(
		path,
		buffer.Bytes(),
		0644,
	)

	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *LocalStorage) Delete(path string) error {

	return os.Remove(path)
}

func saveUploadedFile(
	file *multipart.FileHeader,
	path string,
) error {

	src, err := file.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	dst, err := os.Create(path)

	if err != nil {
		return err
	}

	defer dst.Close()

	buffer := make([]byte, 32*1024)

	for {

		n, err := src.Read(buffer)

		if n > 0 {
			_, writeErr := dst.Write(buffer[:n])

			if writeErr != nil {
				return writeErr
			}
		}

		if err != nil {
			break
		}
	}

	return nil
}
