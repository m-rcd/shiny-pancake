package database

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/m-rcd/notes/pkg/models"
)

type LocalFileSystem struct {
	workDir string
}

func NewLocalFileSystem(workDir string) *LocalFileSystem {
	return &LocalFileSystem{
		workDir: workDir + "/notes",
	}
}

func (l *LocalFileSystem) Open() error {
	err := os.MkdirAll(l.workDir, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalFileSystem) Close() error {
	return nil
}

var note models.Note

func (l *LocalFileSystem) Create(body io.ReadCloser) (models.Note, error) {
	reqBody, _ := ioutil.ReadAll(body)
	json.Unmarshal(reqBody, &note)

	folder := fmt.Sprintf("%s/%s/active/", l.workDir, note.User.Username)
	err := os.MkdirAll(folder, 0777)
	if err != nil {
		return models.Note{}, err
	}
	path := fmt.Sprintf("%s%s.txt", folder, note.Name)
	err = ioutil.WriteFile(path, []byte(note.Content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
	return note, nil
}
