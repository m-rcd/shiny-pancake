package sql_test

import (
	"encoding/json"
	"io"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/m-rcd/notes/pkg/database/sql"
	"github.com/m-rcd/notes/pkg/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/DATA-DOG/go-sqlmock"
)

var _ = Describe("Sql", func() {
	Context("Create", func() {
		It("creates a new note", func() {
			s := sql.NewSQL("username", "password", "127.0.0.1", "3306")
			db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			s.Db = db
			defer db.Close()

			note := models.Note{Name: "Note1", Content: "Miawwww", Archived: false, User: models.User{Username: "Casper"}}
			bytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			reader := io.NopCloser(strings.NewReader(string(bytes)))
			mock.ExpectExec("INSERT INTO notes").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			newNote, err := s.Create(reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(newNote.Name).To(Equal("Note1"))
		})
	})

	Context("Update", func() {
		It("updates a previously saved note", func() {
			s := sql.NewSQL("username", "password", "127.0.0.1", "3306")
			db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			s.Db = db
			defer db.Close()
			existingNote := models.Note{Id: "1", Name: "Note1", Content: "Miawwww", Archived: false, User: models.User{Username: "Casper"}}

			requestData := models.Note{Name: "Note1", Content: "updated", Archived: false, User: models.User{Username: "Casper"}}
			bytes, err := json.Marshal(requestData)
			Expect(err).NotTo(HaveOccurred())
			reader := io.NopCloser(strings.NewReader(string(bytes)))

			rows := sqlmock.NewRows([]string{"id", "name", "content", "archived"}).
				AddRow(existingNote.Id, existingNote.Name, existingNote.Content, existingNote.Archived)
			mock.ExpectQuery("SELECT id, name, content, archived FROM notes WHERE id=1").WillReturnRows(rows)

			mock.ExpectExec("UPDATE notes").WithArgs(requestData.Name, requestData.Content, 0, existingNote.Id).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			updatedNote, err := s.Update(existingNote.Id, reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedNote.Content).To(Equal(requestData.Content))
		})
	})
})
