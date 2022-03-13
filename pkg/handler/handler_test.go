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

})
