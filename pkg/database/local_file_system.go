package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/m-rcd/notes/pkg/models"
	uuid "github.com/nu7hatch/gouuid"
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

func (l *LocalFileSystem) Create(body io.ReadCloser) (models.Note, error) {
	var note models.Note
	reqBody, _ := ioutil.ReadAll(body)
	json.Unmarshal(reqBody, &note)
	note.Id = newId()
	err := validateNote(note)
	if err != nil {
		return models.Note{}, err
	}

	folder := fmt.Sprintf("%s/%s/active/", l.workDir, note.User.Username)
	err = os.MkdirAll(folder, 0777)
	if err != nil {
		return models.Note{}, err
	}
	fileName := fmt.Sprintf("%s_%s.txt", note.Name, note.Id)
	filePath := fmt.Sprintf("%s%s", folder, fileName)
	err = ioutil.WriteFile(filePath, []byte(note.Content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
	return note, nil
}

func (l *LocalFileSystem) Update(id string, body io.ReadCloser) (models.Note, error) {
	var note models.Note
	reqBody, _ := ioutil.ReadAll(body)
	json.Unmarshal(reqBody, &note)
	note.Id = id

	err := validateNote(note)
	if err != nil {
		return models.Note{}, err
	}
	fileName := fmt.Sprintf("%s_%s.txt", note.Name, note.Id)

	filePath := fmt.Sprintf("%s/%s/active/%s", l.workDir, note.User.Username, fileName)
	err = validateFileExists(filePath)

	if err != nil {
		return models.Note{}, err
	}

	err = ioutil.WriteFile(filePath, []byte(note.Content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
	return note, nil
}

func (l *LocalFileSystem) Delete(id string, body io.ReadCloser) error {
	var user models.User
	reqBody, _ := ioutil.ReadAll(body)
	json.Unmarshal(reqBody, &user)
	filePath, err := findFile(fmt.Sprintf("%s/%s/active/", l.workDir, user.Username), id)
	if err != nil {
		return err
	}

	err = os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return nil
}

func validateNote(note models.Note) error {
	if !isSet(note.Name) {
		return errors.New("name must be set")
	}

	if !isSet(note.User.Username) {
		return errors.New("user must be set")
	}
	return nil
}

func validateFileExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		return errors.New("file does not exist")
	}
	return nil
}

func isSet(attr string) bool {
	return attr != ""
}

func newId() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func findFile(dir string, id string) (string, error) {
	var f string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		name := strings.Split(strings.Split(file.Name(), "_")[1], ".")[0]
		if name == id {
			f = dir + file.Name()
		}
	}

	if f == "" {
		return "", errors.New("file does not exist")
	}

	return f, nil
}
