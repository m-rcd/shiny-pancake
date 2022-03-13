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

	_, err = s.Db.Exec(CreateNoteTable)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQL) Close() error {
	return s.Db.Close()
}

func (s *SQL) Create(body io.ReadCloser) (models.Note, error) {
	stmt, err := s.Db.Prepare("INSERT INTO notes(name, content, username, archived) VALUES(?, ?, ?, ?)")
	if err != nil {
		return models.Note{}, err
	}
	var note models.Note
	reqBody, _ := ioutil.ReadAll(body)
	json.Unmarshal(reqBody, &note)
	savedNote, err := stmt.Exec(note.Name, note.Content, note.User.Username, false)
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
	return models.Note{}, nil
}

func (s *SQL) Delete(id string, body io.ReadCloser) error {
	return nil
}

func (s *SQL) ListActiveNotes(body io.ReadCloser) ([]models.Note, error) {
	return []models.Note{}, nil
}

func (s *SQL) ListArchivedNotes(body io.ReadCloser) ([]models.Note, error) {
	return []models.Note{}, nil
}
