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

	if note.Archived {
		archivedNote, err := Archive(l.workDir, note)
		if err != nil {
			return note, err
		}

		return archivedNote, nil
	}
	dir := fmt.Sprintf("%s/%s/", l.workDir, note.User.Username)

	if Archived(dir, id) {
		activeNote, err := Unarchive(l.workDir, note)
		if err != nil {
			return note, err
		}

		return activeNote, nil
	}

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
	fileName, err := findFile(fmt.Sprintf("%s/%s/active/", l.workDir, user.Username), id)
	if err != nil {
		return err
	}

	err = os.RemoveAll(fmt.Sprintf("%s/%s/active/%s", l.workDir, user.Username, fileName))
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalFileSystem) ListActiveNotes(body io.ReadCloser) ([]models.Note, error) {
	var user models.User
	reqBody, _ := ioutil.ReadAll(body)
	json.Unmarshal(reqBody, &user)
	notes := []models.Note{}
	dir := fmt.Sprintf("%s/%s/active/", l.workDir, user.Username)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []models.Note{}, err
	}
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return []models.Note{}, err
		}

		note := models.Note{
			User:    user,
			Id:      strings.Split(strings.Split(file.Name(), "_")[1], ".")[0],
			Name:    strings.Split(file.Name(), "_")[0],
			Content: string(content),
		}
		err = validateNote(note)
		if err != nil {
			return []models.Note{}, err
		}

		notes = append(notes, note)
	}
	return notes, nil
}

func Archive(dir string, note models.Note) (models.Note, error) {
	oldPath := fmt.Sprintf("%s/%s/active/", dir, note.User.Username)
	newPath := fmt.Sprintf("%s/%s/archived/", dir, note.User.Username)

	archivedNote, err := moveNote(note, oldPath, newPath)
	if err != nil {
		return models.Note{}, err
	}

	return archivedNote, nil
}

func Unarchive(dir string, note models.Note) (models.Note, error) {
	oldPath := fmt.Sprintf("%s/%s/archived/", dir, note.User.Username)
	newPath := fmt.Sprintf("%s/%s/active/", dir, note.User.Username)

	activeNote, err := moveNote(note, oldPath, newPath)
	if err != nil {
		return models.Note{}, err
	}

	return activeNote, nil
}

func moveNote(note models.Note, from, to string) (models.Note, error) {
	fileName, err := findFile(from, note.Id)
	if err != nil {
		return models.Note{}, err
	}
	fromFile := fmt.Sprintf("%s%s", from, fileName)
	toFile := fmt.Sprintf("%s%s", to, fileName)

	err = os.MkdirAll(to, 0777)
	if err != nil {
		return models.Note{}, err
	}

	err = os.Rename(fromFile, toFile)
	if err != nil {
		return models.Note{}, err
	}

	note.Name = strings.Split(fileName, "_")[0]
	oldContent, err := os.ReadFile(toFile)
	if err != nil {
		return models.Note{}, err
	}

	note.Content = string(oldContent)

	return note, err
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
	var fileName string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		name := strings.Split(strings.Split(file.Name(), "_")[1], ".")[0]
		if name == id {
			fileName = file.Name()
		}
	}

	if fileName == "" {
		return "", errors.New("file does not exist")
	}

	return fileName, nil
}

func Archived(dir string, id string) bool {
	file, _ := findFile(dir+"archived/", id)
	return file != ""
}
