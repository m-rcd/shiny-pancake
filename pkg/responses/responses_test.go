package responses_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/m-rcd/notes/pkg/models"
	"github.com/m-rcd/notes/pkg/responses"
)

var _ = Describe("Responses", func() {
	responder := responses.NewNoteResponder()

	Context("success", func() {
		It("returns a json response", func() {
			note := models.Note{Name: "Note", Content: "I am successful!", User: models.User{Username: "Sabriel"}}
			message := "Note created successfully"
			data := []models.Note{note}

			expectedResponse := responses.JsonNoteResponse{Type: "success", StatusCode: 200, Data: data, Message: message}
			Expect(responder.Success(data, message)).To(Equal(expectedResponse))
		})
	})

	Context("failure", func() {
		It("returns a json response", func() {
			message := "Note not created"

			expectedResponse := responses.JsonNoteResponse{Type: "failed", StatusCode: 500, Data: []models.Note{}, Message: message}
			Expect(responder.Failure(message)).To(Equal(expectedResponse))
		})
	})
})
