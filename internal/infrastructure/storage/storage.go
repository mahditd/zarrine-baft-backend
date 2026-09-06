package storage

import (
	"mime/multipart"
)

type Storage interface {
	Save(file *multipart.FileHeader) (string, error)
	Delete(path string) error
}