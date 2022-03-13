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
	var (
		id       = "1"
		name     = "Note1"
		content  = "Miawww"
		username = "Casper"
		archived = false
	)

	Context("Create", func() {
		It("creates a new note", func() {
			s := sql.NewSQL("username", "password", "127.0.0.1", "3306")
			db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			s.Db = db
			defer db.Close()

			note := models.Note{Name: name, Content: content, Archived: archived, User: models.User{Username: username}}
			bytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			reader := io.NopCloser(strings.NewReader(string(bytes)))
			mock.ExpectExec("INSERT INTO notes").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			newNote, err := s.Create(reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(newNote.Name).To(Equal(name))
		})
	})

	Context("Update", func() {
		It("updates a previously saved note", func() {
			s := sql.NewSQL("username", "password", "127.0.0.1", "3306")
			db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			s.Db = db
			defer db.Close()
			existingNote := models.Note{Id: id, Name: name, Content: content, Archived: archived, User: models.User{Username: username}}

			requestData := models.Note{Name: name, Content: "updated", Archived: archived, User: models.User{Username: username}}
			bytes, err := json.Marshal(requestData)
			Expect(err).NotTo(HaveOccurred())
			reader := io.NopCloser(strings.NewReader(string(bytes)))

			rows := sqlmock.NewRows([]string{"id", "name", "content", "archived"}).
				AddRow(existingNote.Id, existingNote.Name, existingNote.Content, existingNote.Archived)
			mock.ExpectQuery("SELECT id, name, content, archived FROM notes WHERE id=" + id).WillReturnRows(rows)

			mock.ExpectExec("UPDATE notes").WithArgs(requestData.Name, requestData.Content, 0, existingNote.Id).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			updatedNote, err := s.Update(existingNote.Id, reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedNote.Content).To(Equal(requestData.Content))
		})
	})

	Context("Delete", func() {
		It("deletes a note", func() {
			s := sql.NewSQL("username", "password", "127.0.0.1", "3306")
			db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			s.Db = db
			defer db.Close()
			existingNote := models.Note{Id: id, Name: name, Content: content, Archived: archived, User: models.User{Username: username}}
			mock.ExpectExec("DELETE FROM notes WHERE id = ?").WithArgs(existingNote.Id).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			reader := io.NopCloser(strings.NewReader(""))

			err = s.Delete(existingNote.Id, reader)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Archive", func() {
		It("archives a note", func() {
			s := sql.NewSQL("username", "password", "127.0.0.1", "3306")
			db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			s.Db = db
			defer db.Close()
			existingNote := models.Note{Id: id, Name: name, Content: content, Archived: archived, User: models.User{Username: username}}

			requestData := models.Note{Archived: true, User: models.User{Username: username}}
			bytes, err := json.Marshal(requestData)
			Expect(err).NotTo(HaveOccurred())
			reader := io.NopCloser(strings.NewReader(string(bytes)))

			rows := sqlmock.NewRows([]string{"id", "name", "content", "archived"}).
				AddRow(existingNote.Id, existingNote.Name, existingNote.Content, existingNote.Archived)
			mock.ExpectQuery("SELECT id, name, content, archived FROM notes WHERE id=1").WillReturnRows(rows)

			mock.ExpectExec("UPDATE notes").WithArgs(1, existingNote.Id).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			updatedNote, err := s.Update(existingNote.Id, reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedNote.Archived).To(Equal(requestData.Archived))
		})
	})

	Context("Unarchive", func() {
		It("unarchives a note", func() {
			s := sql.NewSQL("username", "password", "127.0.0.1", "3306")
			db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			s.Db = db
			defer db.Close()
			existingNote := models.Note{Id: id, Name: name, Content: content, Archived: true, User: models.User{Username: username}}

			requestData := models.Note{Archived: false, User: models.User{Username: username}}
			bytes, err := json.Marshal(requestData)
			Expect(err).NotTo(HaveOccurred())
			reader := io.NopCloser(strings.NewReader(string(bytes)))

			rows := sqlmock.NewRows([]string{"id", "name", "content", "archived"}).
				AddRow(existingNote.Id, existingNote.Name, existingNote.Content, existingNote.Archived)
			mock.ExpectQuery("SELECT id, name, content, archived FROM notes WHERE id=" + id).WillReturnRows(rows)

			mock.ExpectExec("UPDATE notes").WithArgs(0, existingNote.Id).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			updatedNote, err := s.Update(existingNote.Id, reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedNote.Archived).To(Equal(requestData.Archived))
		})
	})
})
