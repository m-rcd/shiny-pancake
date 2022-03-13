package responses

import "github.com/m-rcd/notes/pkg/models"

type Response interface {
	Failure(message string) JsonNoteResponse
	Success(data []models.Note, message string) JsonNoteResponse
}
