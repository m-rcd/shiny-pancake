package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/m-rcd/notes/pkg/database/databasefakes"
	"github.com/m-rcd/notes/pkg/handler"
	"github.com/m-rcd/notes/pkg/models"
	"github.com/m-rcd/notes/pkg/responses"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {

	Context("#CreateNewNote", func() {
		It("Creates a new note", func() {
			fake_db := new(databasefakes.FakeDatabase)

			h := handler.New(fake_db)
			r := httptest.NewRecorder()
			postData := bytes.NewBuffer([]byte(`{"name":"Vampires","content":"I SLAY","user":{"username":"Buffy"}}`))
			req, err := http.NewRequest("POST", "http://localhost:10000/note", postData)
			Expect(err).NotTo(HaveOccurred())

			note := models.Note{Name: "Vampires", Content: "I SLAY", User: models.User{Username: "Buffy"}}
			fake_db.CreateReturns(note, nil)
			h.CreateNewNote(r, req)
			Expect(fake_db.CreateCallCount()).To(Equal(1))
			var response responses.JsonNoteResponse

			json.Unmarshal(r.Body.Bytes(), &response)
			Expect(response.Type).To(Equal("success"))
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Data[0].Name).To(Equal(note.Name))
			Expect(response.Data[0].Content).To(Equal(note.Content))
			Expect(response.Data[0].User).To(Equal(note.User))
			Expect(response.Message).To(Equal("The note was successfully created"))
		})

		Context("when an error occurs", func() {
			It("does not create a note", func() {
				fake_db := new(databasefakes.FakeDatabase)

				h := handler.New(fake_db)
				r := httptest.NewRecorder()
				postData := bytes.NewBuffer([]byte(`{"name":"","content":"I SLAY","user":{"username":"Buffy"}}`))
				req, err := http.NewRequest("POST", "http://localhost:10000/note", postData)
				Expect(err).NotTo(HaveOccurred())

				fake_db.CreateReturns(models.Note{}, errors.New("Not created"))
				h.CreateNewNote(r, req)
				Expect(fake_db.CreateCallCount()).To(Equal(1))
				var response responses.JsonNoteResponse

				json.Unmarshal(r.Body.Bytes(), &response)
				Expect(response.Type).To(Equal("failed"))
				Expect(response.StatusCode).To(Equal(500))
				Expect(response.Message).To(Equal("Not created"))
			})
		})
	})

	Context("#UpdateNote", func() {
		It("handles PATCH request", func() {
			fake_db := new(databasefakes.FakeDatabase)

			data := bytes.NewBuffer([]byte(`{"name":"Vampires","content":"I SLAY A LOT","user":{"username":"Buffy"}}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/1", data)
			Expect(err).NotTo(HaveOccurred())
			r := httptest.NewRecorder()
			h := handler.New(fake_db)

			note := models.Note{Id: "1", Name: "Vampires", Content: "I SLAY A LOT", User: models.User{Username: "Buffy"}}
			fake_db.UpdateReturns(note, nil)
			h.UpdateNote(r, req)
			Expect(fake_db.UpdateCallCount()).To(Equal(1))
			var response responses.JsonNoteResponse

			json.Unmarshal(r.Body.Bytes(), &response)
			Expect(response.Type).To(Equal("success"))
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Data[0].Content).To(Equal("I SLAY A LOT"))
			Expect(response.Message).To(Equal("The note was successfully updated"))
		})

		Context("when an error occurs", func() {
			It("does not update the note", func() {
				fake_db := new(databasefakes.FakeDatabase)

				h := handler.New(fake_db)
				r := httptest.NewRecorder()
				patchData := bytes.NewBuffer([]byte(`{"name":"","content":"I SLAY","user":{"username":"Buffy"}}`))
				req, err := http.NewRequest("POST", "http://localhost:10000/note/1", patchData)
				Expect(err).NotTo(HaveOccurred())

				fake_db.UpdateReturns(models.Note{}, errors.New("Not updated"))
				h.UpdateNote(r, req)
				Expect(fake_db.UpdateCallCount()).To(Equal(1))
				var response responses.JsonNoteResponse

				json.Unmarshal(r.Body.Bytes(), &response)
				Expect(response.Type).To(Equal("failed"))
				Expect(response.StatusCode).To(Equal(500))
				Expect(response.Message).To(Equal("Not updated"))
			})
		})
	})

	Context("#DeleteNote", func() {
		It("handles DELETE request", func() {
			fake_db := new(databasefakes.FakeDatabase)

			data := bytes.NewBuffer([]byte(`{"username":"Buffy"}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/1", data)
			Expect(err).NotTo(HaveOccurred())
			r := httptest.NewRecorder()
			h := handler.New(fake_db)

			fake_db.DeleteReturns(nil)
			h.DeleteNote(r, req)
			Expect(fake_db.DeleteCallCount()).To(Equal(1))
			var response responses.JsonNoteResponse

			json.Unmarshal(r.Body.Bytes(), &response)
			Expect(response.Type).To(Equal("success"))
			Expect(response.StatusCode).To(Equal(200))
			Expect(len(response.Data)).To(Equal(0))
			Expect(response.Message).To(Equal("The note was successfully deleted"))
		})

		Context("when an error occurs", func() {
			It("does not delete the note", func() {
				fake_db := new(databasefakes.FakeDatabase)

				h := handler.New(fake_db)
				r := httptest.NewRecorder()
				data := bytes.NewBuffer([]byte(`{"username":"Buffy"}`))
				req, err := http.NewRequest("POST", "http://localhost:10000/note/1", data)
				Expect(err).NotTo(HaveOccurred())

				fake_db.DeleteReturns(errors.New("Not deleted"))
				h.DeleteNote(r, req)
				Expect(fake_db.DeleteCallCount()).To(Equal(1))
				var response responses.JsonNoteResponse

				json.Unmarshal(r.Body.Bytes(), &response)
				Expect(response.Type).To(Equal("failed"))
				Expect(response.StatusCode).To(Equal(500))
				Expect(response.Message).To(Equal("Not deleted"))
			})
		})
	})

	Context("#ListActiveNotes", func() {
		It("handles GET request", func() {
			fake_db := new(databasefakes.FakeDatabase)

			data := bytes.NewBuffer([]byte(`{"username":"Buffy"}`))
			req, err := http.NewRequest("GET", "http://localhost:10000/notes/active", data)
			Expect(err).NotTo(HaveOccurred())
			r := httptest.NewRecorder()
			h := handler.New(fake_db)

			note1 := models.Note{Id: "1", Name: "Vampires", Content: "I SLAY A LOT", User: models.User{Username: "Buffy"}}
			note2 := models.Note{Id: "2", Name: "Monsters", Content: "I EAT THEM", User: models.User{Username: "Buffy"}}

			fake_db.ListActiveNotesReturns([]models.Note{note1, note2}, nil)
			h.ListActiveNotes(r, req)
			Expect(fake_db.ListActiveNotesCallCount()).To(Equal(1))
			var list []models.Note
			json.Unmarshal(r.Body.Bytes(), &list)
			Expect(len(list)).To(Equal(2))
			Expect(list[0]).To(Equal(note1))
			Expect(list[1]).To(Equal(note2))
		})
	})

	Context("#ListArchivedNotes", func() {
		It("handles GET request", func() {
			fake_db := new(databasefakes.FakeDatabase)

			data := bytes.NewBuffer([]byte(`{"username":"Buffy"}`))
			req, err := http.NewRequest("GET", "http://localhost:10000/notes/archived", data)
			Expect(err).NotTo(HaveOccurred())
			r := httptest.NewRecorder()
			h := handler.New(fake_db)

			note1 := models.Note{Id: "1", Name: "Vampires", Content: "I SLAY A LOT", Archived: true, User: models.User{Username: "Buffy"}}
			note2 := models.Note{Id: "2", Name: "Monsters", Content: "I EAT THEM", Archived: true, User: models.User{Username: "Buffy"}}

			fake_db.ListArchivedNotesReturns([]models.Note{note1, note2}, nil)
			h.ListArchivedNotes(r, req)
			Expect(fake_db.ListArchivedNotesCallCount()).To(Equal(1))
			var list []models.Note
			json.Unmarshal(r.Body.Bytes(), &list)
			Expect(len(list)).To(Equal(2))
			Expect(list[0]).To(Equal(note1))
			Expect(list[1]).To(Equal(note2))
		})
	})
})
