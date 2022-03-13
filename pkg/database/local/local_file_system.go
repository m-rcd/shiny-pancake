package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/m-rcd/notes/pkg/models"
	"github.com/m-rcd/notes/pkg/utils"
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
	return os.MkdirAll(l.workDir, 0777)
}

func (l *LocalFileSystem) Close() error {
	return nil
}

func (l *LocalFileSystem) Create(body io.ReadCloser) (models.Note, error) {
	var note models.Note
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return note, err
	}

	if err := json.Unmarshal(reqBody, &note); err != nil {
		return note, err
	}

	if err := validateNote(note); err != nil {
		return note, err
	}

	note.Id = newId()

	activeDir := fmt.Sprintf("%s/%s/active/", l.workDir, note.User.Username)
	if err := os.MkdirAll(activeDir, 0777); err != nil {
		return note, err
	}

	fileName := fmt.Sprintf("%s_%s.txt", note.Name, note.Id)
	filePath := fmt.Sprintf("%s%s", activeDir, fileName)
	if err := ioutil.WriteFile(filePath, []byte(note.Content), 0777); err != nil {
		return note, err
	}

	return note, nil
}

func (l *LocalFileSystem) Update(id string, body io.ReadCloser) (models.Note, error) {
	var note models.Note
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return note, err
	}

	if err := json.Unmarshal(reqBody, &note); err != nil {
		return note, err
	}

	note.Id = id

	if note.Archived {
		archivedNote, err := Archive(l.workDir, note)
		if err != nil {
			return note, err
		}

		return archivedNote, nil
	}

	dir := fmt.Sprintf("%s/%s/", l.workDir, note.User.Username)
	if archived(dir, id) {
		activeNote, err := Unarchive(l.workDir, note)
		if err != nil {
			return note, err
		}

		return activeNote, nil
	}

	if err := validateNote(note); err != nil {
		return models.Note{}, err
	}

	fileName := fmt.Sprintf("%s_%s.txt", note.Name, note.Id)
	filePath := fmt.Sprintf("%s/%s/active/%s", l.workDir, note.User.Username, fileName)
	if err := validateFileExists(filePath); err != nil {
		return models.Note{}, err
	}

	if err := ioutil.WriteFile(filePath, []byte(note.Content), 0777); err != nil {
		return models.Note{}, err
	}

	return note, nil
}

func (l *LocalFileSystem) Delete(id string, body io.ReadCloser) error {
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	var user models.User
	if err := json.Unmarshal(reqBody, &user); err != nil {
		return err
	}

	fileName, err := findFile(fmt.Sprintf("%s/%s/active/", l.workDir, user.Username), id)
	if err != nil {
		return err
	}

	return os.RemoveAll(fmt.Sprintf("%s/%s/active/%s", l.workDir, user.Username, fileName))
}

func (l *LocalFileSystem) ListActiveNotes(body io.ReadCloser) ([]models.Note, error) {
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return []models.Note{}, err
	}

	var user models.User
	if err := json.Unmarshal(reqBody, &user); err != nil {
		return []models.Note{}, err
	}

	dir := fmt.Sprintf("%s/%s/active/", l.workDir, user.Username)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []models.Note{}, err
	}

	notes, err := ListNotes(dir, files, user, false)
	if err != nil {
		return []models.Note{}, err
	}

	return notes, nil
}

func (l *LocalFileSystem) ListArchivedNotes(body io.ReadCloser) ([]models.Note, error) {
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return []models.Note{}, err
	}

	var user models.User
	if err := json.Unmarshal(reqBody, &user); err != nil {
		return []models.Note{}, err
	}

	dir := fmt.Sprintf("%s/%s/archived/", l.workDir, user.Username)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []models.Note{}, err
	}

	notes, err := ListNotes(dir, files, user, true)
	if err != nil {
		return []models.Note{}, err
	}

	return notes, nil
}

func ListNotes(dir string, files []fs.FileInfo, user models.User, archived bool) ([]models.Note, error) {
	notes := []models.Note{}
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return []models.Note{}, err
		}

		note := models.Note{
			User:     user,
			Id:       strings.Split(strings.Split(file.Name(), "_")[1], ".")[0],
			Name:     strings.Split(file.Name(), "_")[0],
			Content:  string(content),
			Archived: archived,
		}

		if err := validateNote(note); err != nil {
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

	if err := os.MkdirAll(to, 0777); err != nil {
		return models.Note{}, err
	}

	if err := os.Rename(fromFile, toFile); err != nil {
		return models.Note{}, err
	}

	note.Name = strings.Split(fileName, "_")[0]
	oldContent, err := os.ReadFile(toFile)
	if err != nil {
		return models.Note{}, err
	}

	note.Content = string(oldContent)

	return note, nil
}

func validateNote(note models.Note) error {
	if !utils.IsSet(note.Name) {
		return errors.New("name must be set")
	}

	if !utils.IsSet(note.User.Username) {
		return errors.New("user must be set")
	}

	return nil
}

func validateFileExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file does not exist: %s", err)
	}

	return nil
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

	if !utils.IsSet(fileName) {
		return "", errors.New("file does not exist")
	}

	return fileName, nil
}

func archived(dir string, id string) bool {
	file, err := findFile(dir+"archived/", id)
	if err != nil {
		return false
	}

	return utils.IsSet(file)
}
