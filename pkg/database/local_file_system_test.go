package database_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/m-rcd/notes/pkg/database"
	"github.com/m-rcd/notes/pkg/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalFileSystem", func() {
	var db database.Database
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "local_file_system_test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	Context("CREATE", func() {
		It("creates a new note file", func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Note1", Content: "Miawwww", User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			newNote, err := db.Create(r)
			Expect(err).NotTo(HaveOccurred())
			filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, newNote.User.Username, newNote.Name, newNote.Id)
			Expect(filepath).To(BeAnExistingFile())
		})

		Context("when error occurs", func() {
			Context("when name is not set", func() {
				It("does not create a note file and raises an error", func() {
					db = database.NewLocalFileSystem(tempDir)
					err := db.Open()
					Expect(err).NotTo(HaveOccurred())

					note := models.Note{Name: "", Content: "Miawwww", User: models.User{Username: "Casper"}}

					noteBytes, err := json.Marshal(note)
					Expect(err).NotTo(HaveOccurred())
					r := io.NopCloser(strings.NewReader(string(noteBytes)))

					_, err = db.Create(r)
					Expect(err).To(MatchError("name must be set"))
				})
			})

			Context("when User username is not set", func() {
				It("does not create a note file and raises an error", func() {
					db = database.NewLocalFileSystem(tempDir)
					err := db.Open()
					Expect(err).NotTo(HaveOccurred())

					note := models.Note{Name: "Note1", Content: "Miawwww", User: models.User{Username: ""}}

					noteBytes, err := json.Marshal(note)
					Expect(err).NotTo(HaveOccurred())
					r := io.NopCloser(strings.NewReader(string(noteBytes)))

					_, err = db.Create(r)
					Expect(err).To(MatchError("user must be set"))
				})
			})
		})
	})

	Context("UPDATE", func() {
		var existingNote models.Note
		BeforeEach(func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			existingNote, err = db.Create(r)
			Expect(err).NotTo(HaveOccurred())
		})

		It("updates a previously saved note", func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Note1", Content: "BOOOO", User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			_, err = db.Update(existingNote.Id, r)
			Expect(err).NotTo(HaveOccurred())
			filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
			content, err := os.ReadFile(filepath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("BOOOO"))
		})

		Context("when an error occurs", func() {
			Context("when file does not exist", func() {
				It("does not update the note and raises an error", func() {
					db = database.NewLocalFileSystem(tempDir)
					err := db.Open()
					Expect(err).NotTo(HaveOccurred())

					note := models.Note{Name: "Note2", Content: "BOOOO", User: models.User{Username: "Casper"}}

					noteBytes, err := json.Marshal(note)
					Expect(err).NotTo(HaveOccurred())
					r := io.NopCloser(strings.NewReader(string(noteBytes)))

					_, err = db.Update(existingNote.Id, r)
					Expect(err).To(MatchError("file does not exist"))
					filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
					content, err := os.ReadFile(filepath)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(content)).To(Equal("Miaaaww"))
				})
			})

			Context("when name is not set", func() {
				It("does not update the note and raises an error", func() {
					db = database.NewLocalFileSystem(tempDir)
					err := db.Open()
					Expect(err).NotTo(HaveOccurred())

					note := models.Note{Name: "", Content: "BOOOO", User: models.User{Username: "Casper"}}

					noteBytes, err := json.Marshal(note)
					Expect(err).NotTo(HaveOccurred())
					r := io.NopCloser(strings.NewReader(string(noteBytes)))

					_, err = db.Update(existingNote.Id, r)
					Expect(err).To(MatchError("name must be set"))
					filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
					content, err := os.ReadFile(filepath)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(content)).To(Equal("Miaaaww"))
				})
			})

			Context("when user is not set", func() {
				It("does not update the note and raises an error", func() {
					db = database.NewLocalFileSystem(tempDir)
					err := db.Open()
					Expect(err).NotTo(HaveOccurred())

					note := models.Note{Name: "Note1", Content: "BOOOO", User: models.User{Username: ""}}

					noteBytes, err := json.Marshal(note)
					Expect(err).NotTo(HaveOccurred())
					r := io.NopCloser(strings.NewReader(string(noteBytes)))

					_, err = db.Update(existingNote.Id, r)
					Expect(err).To(MatchError("user must be set"))
					filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
					content, err := os.ReadFile(filepath)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(content)).To(Equal("Miaaaww"))
				})
			})
		})
	})

	Context("Delete", func() {
		var existingNote models.Note
		BeforeEach(func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			existingNote, err = db.Create(r)
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes a note", func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			user := models.User{Username: "Casper"}

			userBytes, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(userBytes)))

			err = db.Delete(existingNote.Id, r)
			Expect(err).NotTo(HaveOccurred())
			filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)

			Expect(filepath).NotTo(BeAnExistingFile())
		})

		Context("when errors occur", func() {
			It("does not delete the file and raises an error", func() {
				db = database.NewLocalFileSystem(tempDir)
				err := db.Open()
				Expect(err).NotTo(HaveOccurred())

				user := models.User{Username: "Casper"}

				userBytes, err := json.Marshal(user)
				Expect(err).NotTo(HaveOccurred())
				r := io.NopCloser(strings.NewReader(string(userBytes)))

				err = db.Delete("123", r)
				Expect(err).To(MatchError("file does not exist"))
				filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)

				Expect(filepath).To(BeAnExistingFile())
			})
		})
	})

	Context("ARCHIVE", func() {
		var existingNote models.Note

		BeforeEach(func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			existingNote, err = db.Create(r)
			Expect(existingNote.Archived).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("archives a note", func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Archived: true, User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			updatedNote, err := db.Update(existingNote.Id, r)
			Expect(err).NotTo(HaveOccurred())
			activeFilepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
			archivedFilePath := fmt.Sprintf("%s/notes/%s/archived/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
			Expect(activeFilepath).NotTo(BeAnExistingFile())
			Expect(archivedFilePath).To(BeAnExistingFile())
			Expect(updatedNote.Archived).To(BeTrue())
			Expect(updatedNote.Name).To(Equal(existingNote.Name))
			Expect(updatedNote.Content).To(Equal(existingNote.Content))
		})

		Context("when content and name are passed as attributes", func() {
			It("archives the note but does not update name/content", func() {
				db = database.NewLocalFileSystem(tempDir)
				err := db.Open()
				Expect(err).NotTo(HaveOccurred())

				note := models.Note{Name: "Note2", Content: "NewContent", Archived: true, User: models.User{Username: "Casper"}}

				noteBytes, err := json.Marshal(note)
				Expect(err).NotTo(HaveOccurred())
				r := io.NopCloser(strings.NewReader(string(noteBytes)))

				updatedNote, err := db.Update(existingNote.Id, r)
				Expect(err).NotTo(HaveOccurred())
				activeFilepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
				archivedFilePath := fmt.Sprintf("%s/notes/%s/archived/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
				Expect(activeFilepath).NotTo(BeAnExistingFile())
				Expect(archivedFilePath).To(BeAnExistingFile())
				Expect(updatedNote.Archived).To(BeTrue())
				Expect(updatedNote.Name).To(Equal(existingNote.Name))
				Expect(updatedNote.Content).To(Equal(existingNote.Content))
			})
		})
	})

})
