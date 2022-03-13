package database

import (
	"io"

	"github.com/m-rcd/notes/pkg/models"
)

type Database interface {
	Open() error
	Close() error
	Create(body io.ReadCloser) (models.Note, error)
}
