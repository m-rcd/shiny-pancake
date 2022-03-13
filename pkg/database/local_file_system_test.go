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

	Context("DELETE", func() {
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

		Context("UNARCHIVE", func() {
			var archivedNote models.Note

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
				updatedNote := models.Note{Name: "Note2", Content: "NewContent", Archived: true, User: models.User{Username: "Casper"}}

				updatedNoteBytes, err := json.Marshal(updatedNote)
				Expect(err).NotTo(HaveOccurred())
				rr := io.NopCloser(strings.NewReader(string(updatedNoteBytes)))

				archivedNote, err = db.Update(existingNote.Id, rr)
				Expect(err).NotTo(HaveOccurred())
			})

			It("unarchives a note", func() {
				db = database.NewLocalFileSystem(tempDir)
				err := db.Open()
				Expect(err).NotTo(HaveOccurred())

				note := models.Note{Archived: false, User: models.User{Username: "Casper"}}

				noteBytes, err := json.Marshal(note)
				Expect(err).NotTo(HaveOccurred())
				r := io.NopCloser(strings.NewReader(string(noteBytes)))

				updatedNote, err := db.Update(archivedNote.Id, r)
				Expect(err).NotTo(HaveOccurred())
				activeFilepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, archivedNote.User.Username, archivedNote.Name, archivedNote.Id)
				archivedFilePath := fmt.Sprintf("%s/notes/%s/archived/%s_%s.txt", tempDir, archivedNote.User.Username, archivedNote.Name, archivedNote.Id)
				Expect(activeFilepath).To(BeAnExistingFile())
				Expect(archivedFilePath).NotTo(BeAnExistingFile())
				Expect(updatedNote.Archived).To(BeFalse())
				Expect(updatedNote.Name).To(Equal(archivedNote.Name))
				Expect(updatedNote.Content).To(Equal(archivedNote.Content))
			})

		})
	})

	Context("LISTACTIVENOTES", func() {
		var note1 models.Note
		var note2 models.Note

		BeforeEach(func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			noteData1 := models.Note{Name: "Note1", Content: "Kirjava", User: models.User{Username: "Lyra"}}

			noteBytes1, err := json.Marshal(noteData1)
			Expect(err).NotTo(HaveOccurred())
			r1 := io.NopCloser(strings.NewReader(string(noteBytes1)))

			note1, err = db.Create(r1)
			Expect(err).NotTo(HaveOccurred())

			noteData2 := models.Note{Name: "Note2", Content: "Pantalaimon", User: models.User{Username: "Lyra"}}

			noteBytes2, err := json.Marshal(noteData2)
			Expect(err).NotTo(HaveOccurred())
			r2 := io.NopCloser(strings.NewReader(string(noteBytes2)))

			note2, err = db.Create(r2)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a list of active notes", func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			user := models.User{Username: "Lyra"}

			userBytes, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(userBytes)))

			notes, err := db.ListActiveNotes(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(notes)).To(Equal(2))
			Expect(notes[0]).To(Equal(note1))
			Expect(notes[1]).To(Equal(note2))
		})
	})

	Context("LISTARCHIVEDNOTES", func() {
		var note1 models.Note
		var note2 models.Note
		var archivedNote1 models.Note
		var archivedNote2 models.Note

		BeforeEach(func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			noteData1 := models.Note{Name: "Note1", Content: "Kirjava", User: models.User{Username: "Lyra"}}

			noteBytes1, err := json.Marshal(noteData1)
			Expect(err).NotTo(HaveOccurred())
			r1 := io.NopCloser(strings.NewReader(string(noteBytes1)))

			note1, err = db.Create(r1)
			Expect(err).NotTo(HaveOccurred())

			noteData2 := models.Note{Name: "Note2", Content: "Pantalaimon", User: models.User{Username: "Lyra"}}

			noteBytes2, err := json.Marshal(noteData2)
			Expect(err).NotTo(HaveOccurred())
			r2 := io.NopCloser(strings.NewReader(string(noteBytes2)))

			note2, err = db.Create(r2)
			Expect(err).NotTo(HaveOccurred())

			updatedNote1 := models.Note{Archived: true, User: models.User{Username: "Lyra"}}

			updatedNote1Bytes, err := json.Marshal(updatedNote1)
			Expect(err).NotTo(HaveOccurred())
			rr1 := io.NopCloser(strings.NewReader(string(updatedNote1Bytes)))

			archivedNote1, err = db.Update(note1.Id, rr1)
			Expect(err).NotTo(HaveOccurred())

			updatedNote2 := models.Note{Archived: true, User: models.User{Username: "Lyra"}}

			updatedNote2Bytes, err := json.Marshal(updatedNote2)
			Expect(err).NotTo(HaveOccurred())
			rr2 := io.NopCloser(strings.NewReader(string(updatedNote2Bytes)))

			archivedNote2, err = db.Update(note2.Id, rr2)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a list of archived notes", func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			user := models.User{Username: "Lyra"}

			userBytes, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(userBytes)))

			notes, err := db.ListArchivedNotes(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(notes)).To(Equal(2))
			Expect(notes[0]).To(Equal(archivedNote1))
			Expect(notes[1]).To(Equal(archivedNote2))
		})
	})

})
