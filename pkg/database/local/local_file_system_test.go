package local_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/m-rcd/notes/pkg/database"
	"github.com/m-rcd/notes/pkg/database/local"
	"github.com/m-rcd/notes/pkg/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalFileSystem", func() {
	var (
		db      database.Database
		tempDir string
		err     error
	)

	BeforeEach(func() {
		tempDir, err = ioutil.TempDir("", "local_file_system_test")
		Expect(err).NotTo(HaveOccurred())

		db = local.NewLocalFileSystem(tempDir)
		err = db.Open()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	Context("CREATE", func() {
		It("creates a new note file", func() {
			note := models.Note{Name: "Note1", Content: "Miawwww", User: models.User{Username: "Casper"}}
			r := buildReader(note)

			newNote, err := db.Create(r)
			Expect(err).NotTo(HaveOccurred())
			filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, newNote.User.Username, newNote.Name, newNote.Id)
			Expect(filepath).To(BeAnExistingFile())
		})

		Context("when error occurs", func() {
			Context("when name is not set", func() {
				It("does not create a note file and raises an error", func() {
					note := models.Note{Name: "", Content: "Miawwww", User: models.User{Username: "Casper"}}
					r := buildReader(note)

					_, err = db.Create(r)
					Expect(err).To(MatchError("name must be set"))
				})
			})

			Context("when User username is not set", func() {
				It("does not create a note file and raises an error", func() {
					note := models.Note{Name: "Note1", Content: "Miawwww", User: models.User{Username: ""}}
					r := buildReader(note)

					_, err = db.Create(r)
					Expect(err).To(MatchError("user must be set"))
				})
			})
		})
	})

	Context("UPDATE", func() {
		var existingNote models.Note

		BeforeEach(func() {
			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}
			existingNote = createNote(note, db)
		})

		It("updates a previously saved note", func() {
			note := models.Note{Name: "Note1", Content: "BOOOO", User: models.User{Username: "Casper"}}
			r := buildReader(note)

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
					note := models.Note{Name: "Note2", Content: "BOOOO", User: models.User{Username: "Casper"}}
					r := buildReader(note)

					_, err = db.Update(existingNote.Id, r)
					Expect(err).To(MatchError(ContainSubstring("file does not exist")))
					filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)
					content, err := os.ReadFile(filepath)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(content)).To(Equal("Miaaaww"))
				})
			})

			Context("when name is not set", func() {
				It("does not update the note and raises an error", func() {
					note := models.Note{Name: "", Content: "BOOOO", User: models.User{Username: "Casper"}}
					r := buildReader(note)

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
					note := models.Note{Name: "Note1", Content: "BOOOO", User: models.User{Username: ""}}
					r := buildReader(note)

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
			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}
			existingNote = createNote(note, db)
		})

		It("deletes a note", func() {
			user := models.User{Username: "Casper"}
			r := buildReader(user)

			err = db.Delete(existingNote.Id, r)
			Expect(err).NotTo(HaveOccurred())
			filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)

			Expect(filepath).NotTo(BeAnExistingFile())
		})

		Context("when errors occur", func() {
			It("does not delete the file and raises an error", func() {
				user := models.User{Username: "Casper"}
				r := buildReader(user)

				err = db.Delete("123", r)
				Expect(err).To(MatchError(ContainSubstring("file does not exist")))
				filepath := fmt.Sprintf("%s/notes/%s/active/%s_%s.txt", tempDir, existingNote.User.Username, existingNote.Name, existingNote.Id)

				Expect(filepath).To(BeAnExistingFile())
			})
		})
	})

	Context("ARCHIVE", func() {
		var existingNote models.Note

		BeforeEach(func() {
			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}
			existingNote = createNote(note, db)
		})

		It("archives a note", func() {
			note := models.Note{Archived: true, User: models.User{Username: "Casper"}}
			r := buildReader(note)

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
				note := models.Note{Name: "Note2", Content: "NewContent", Archived: true, User: models.User{Username: "Casper"}}
				r := buildReader(note)

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
				note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}
				archivedNote = createNote(note, db)

				updatedNote := models.Note{Name: "Note2", Content: "NewContent", Archived: true, User: models.User{Username: "Casper"}}
				rr := buildReader(updatedNote)

				archivedNote, err = db.Update(existingNote.Id, rr)
				Expect(err).NotTo(HaveOccurred())
			})

			It("unarchives a note", func() {
				note := models.Note{Archived: false, User: models.User{Username: "Casper"}}
				r := buildReader(note)

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

	Context("LIST active notes", func() {
		var note1 models.Note
		var note2 models.Note

		BeforeEach(func() {
			noteData1 := models.Note{Name: "Note1", Content: "Kirjava", User: models.User{Username: "Lyra"}}
			note1 = createNote(noteData1, db)

			noteData2 := models.Note{Name: "Note2", Content: "Pantalaimon", User: models.User{Username: "Lyra"}}
			note2 = createNote(noteData2, db)
		})

		It("returns a list of active notes", func() {
			user := models.User{Username: "Lyra"}
			r := buildReader(user)

			notes, err := db.ListActiveNotes(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(notes)).To(Equal(2))
			Expect(notes[0]).To(Equal(note1))
			Expect(notes[1]).To(Equal(note2))
		})
	})

	Context("LIST archived notes", func() {
		var note1 models.Note
		var note2 models.Note
		var archivedNote1 models.Note
		var archivedNote2 models.Note

		BeforeEach(func() {
			noteData1 := models.Note{Name: "Note1", Content: "Kirjava", User: models.User{Username: "Lyra"}}
			note1 = createNote(noteData1, db)

			noteData2 := models.Note{Name: "Note2", Content: "Pantalaimon", User: models.User{Username: "Lyra"}}
			note2 = createNote(noteData2, db)

			updatedNote1 := models.Note{Archived: true, User: models.User{Username: "Lyra"}}

			rr1 := buildReader(updatedNote1)

			archivedNote1, err = db.Update(note1.Id, rr1)
			Expect(err).NotTo(HaveOccurred())

			updatedNote2 := models.Note{Archived: true, User: models.User{Username: "Lyra"}}
			rr2 := buildReader(updatedNote2)

			archivedNote2, err = db.Update(note2.Id, rr2)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a list of archived notes", func() {
			user := models.User{Username: "Lyra"}
			r := buildReader(user)

			notes, err := db.ListArchivedNotes(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(notes)).To(Equal(2))
			Expect(notes[0]).To(Equal(archivedNote1))
			Expect(notes[1]).To(Equal(archivedNote2))
		})
	})

})

func buildReader(data interface{}) io.ReadCloser {
	bytes, err := json.Marshal(data)
	Expect(err).NotTo(HaveOccurred())
	reader := io.NopCloser(strings.NewReader(string(bytes)))
	return reader
}

func createNote(note models.Note, db database.Database) models.Note {
	r := buildReader(note)
	note, err := db.Create(r)
	Expect(err).NotTo(HaveOccurred())
	return note
}
