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
			s.Db = db
			Expect(err).NotTo(HaveOccurred())

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
})
