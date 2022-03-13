package database_test

import (
	"encoding/json"
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

			_, err = db.Create(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(tempDir + "/notes/Casper/active/note1.txt").To(BeAnExistingFile())
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
		BeforeEach(func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			_, err = db.Create(r)
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

			_, err = db.Update("Note1", r)
			Expect(err).NotTo(HaveOccurred())
			dat, err := os.ReadFile(tempDir + "/notes/Casper/active/note1.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(dat)).To(Equal("BOOOO"))
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

					_, err = db.Update("Note2", r)
					Expect(err).To(MatchError("file does not exist"))
					dat, err := os.ReadFile(tempDir + "/notes/Casper/active/note1.txt")
					Expect(err).NotTo(HaveOccurred())
					Expect(string(dat)).To(Equal("Miaaaww"))
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

					_, err = db.Update("Note1", r)
					Expect(err).To(MatchError("name must be set"))
					dat, err := os.ReadFile(tempDir + "/notes/Casper/active/note1.txt")
					Expect(err).NotTo(HaveOccurred())
					Expect(string(dat)).To(Equal("Miaaaww"))
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

					_, err = db.Update("Note1", r)
					Expect(err).To(MatchError("user must be set"))
					dat, err := os.ReadFile(tempDir + "/notes/Casper/active/note1.txt")
					Expect(err).NotTo(HaveOccurred())
					Expect(string(dat)).To(Equal("Miaaaww"))
				})
			})
		})
	})

	Context("Delete", func() {
		BeforeEach(func() {
			db = database.NewLocalFileSystem(tempDir)
			err := db.Open()
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Note1", Content: "Miaaaww", User: models.User{Username: "Casper"}}

			noteBytes, err := json.Marshal(note)
			Expect(err).NotTo(HaveOccurred())
			r := io.NopCloser(strings.NewReader(string(noteBytes)))

			_, err = db.Create(r)
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

			err = db.Delete("Note1", r)
			Expect(err).NotTo(HaveOccurred())
			Expect(tempDir + "/notes/Casper/active/note1.txt").NotTo(BeAnExistingFile())
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

				err = db.Delete("Note2", r)
				Expect(err).To(MatchError("file does not exist"))
				Expect(tempDir + "/notes/Casper/active/note1.txt").To(BeAnExistingFile())
			})
		})
	})
})
