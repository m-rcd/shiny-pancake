package database

import (
	"io"

	"github.com/m-rcd/notes/pkg/models"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Database
type Database interface {
	Open() error
	Close() error
	Create(body io.ReadCloser) (models.Note, error)
}
