package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/m-rcd/notes/pkg/models"
	"github.com/m-rcd/notes/pkg/utils"
)

type SQL struct {
	Db       *sql.DB
	username string
	password string
	address  string
	port     string
}

func NewSQL(username, password, address, port string) *SQL {
	return &SQL{
		username: username,
		password: password,
		address:  address,
		port:     port,
	}
}

func (s *SQL) Open() error {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", s.username, s.password, s.address, s.port, "notes")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return err
	}

	s.Db = db

	if _, err := s.Db.Exec(CreateNoteTable); err != nil {
		return err
	}

	return nil
}

func (s *SQL) Close() error {
	return s.Db.Close()
}

func (s *SQL) Create(body io.ReadCloser) (models.Note, error) {
	var note models.Note

	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return models.Note{}, err
	}

	if err := json.Unmarshal(reqBody, &note); err != nil {
		return models.Note{}, err
	}

	sql := fmt.Sprintf("INSERT INTO notes(name, content, username, archived) VALUES ('%s', '%s', '%s', '%v')", note.Name, note.Content, note.User.Username, 0)
	savedNote, err := s.Db.Exec(sql)
	if err != nil {
		return models.Note{}, err
	}

	id, err := savedNote.LastInsertId()
	if err != nil {
		return models.Note{}, err
	}

	note.Id = strconv.FormatInt(id, 10)

	return note, nil
}

func (s *SQL) Update(id string, body io.ReadCloser) (models.Note, error) {
	var existingNote models.Note

	result := s.Db.QueryRow("SELECT id, name, content, archived FROM notes WHERE id=" + id)
	if err := result.Scan(&existingNote.Id, &existingNote.Name, &existingNote.Content, &existingNote.Archived); err != nil {
		return models.Note{}, err
	}

	var note models.Note
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return models.Note{}, err
	}

	if err := json.Unmarshal(reqBody, &note); err != nil {
		return models.Note{}, err
	}

	if !utils.IsSet(note.Name) {
		note.Name = existingNote.Name
	}

	if !utils.IsSet(note.Content) {
		note.Content = existingNote.Content
	}

	if !utils.IsSet(note.User.Username) {
		note.User = existingNote.User
	}

	note.Id = id

	if note.Archived {
		if err := archive(s, note); err != nil {
			return models.Note{}, err
		}

		return note, nil
	}

	if !note.Archived && existingNote.Archived {
		if err := unarchive(s, note); err != nil {
			return models.Note{}, err
		}

		return note, nil
	}

	if _, err := s.Db.Exec("UPDATE notes set name=?, content=?, archived=? where id=?", note.Name, note.Content, 0, id); err != nil {
		return models.Note{}, err
	}

	return note, nil
}

func (s *SQL) Delete(id string, body io.ReadCloser) error {
	if _, err := s.Db.Exec("DELETE FROM notes WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

func (s *SQL) ListActiveNotes(body io.ReadCloser) ([]models.Note, error) {
	result, err := s.Db.Query("SELECT * FROM notes WHERE archived=0")
	if err != nil {
		return []models.Note{}, err
	}
	defer result.Close()

	notes, err := listNotes(result)
	if err != nil {
		return []models.Note{}, err
	}

	return notes, nil
}

func (s *SQL) ListArchivedNotes(body io.ReadCloser) ([]models.Note, error) {
	result, err := s.Db.Query("SELECT * FROM notes WHERE archived=1")
	if err != nil {
		return []models.Note{}, err
	}
	defer result.Close()

	notes, err := listNotes(result)
	if err != nil {
		return []models.Note{}, err
	}

	return notes, nil
}

func archive(s *SQL, note models.Note) error {
	if _, err := s.Db.Exec("UPDATE notes set archived=? where id=?", 1, note.Id); err != nil {
		return err
	}

	return nil
}

func unarchive(s *SQL, note models.Note) error {
	if _, err := s.Db.Exec("UPDATE notes set archived=? where id=?", 0, note.Id); err != nil {
		return err
	}

	return nil
}

func listNotes(result *sql.Rows) ([]models.Note, error) {
	var (
		notes []models.Note
		note  models.Note
	)

	for result.Next() {
		if err := result.Scan(&note.Id, &note.Name, &note.Content, &note.Archived, &note.User.Username); err != nil {
			return []models.Note{}, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}
