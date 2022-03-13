package responses

import "github.com/m-rcd/notes/pkg/models"

type JsonNoteResponse struct {
	Type       string        `json:"type"`
	StatusCode int           `json:status_code`
	Data       []models.Note `json:"data"`
	Message    string        `json:"message"`
}

type NoteResponse struct {
	Response JsonNoteResponse
}

func NewNoteResponse() *NoteResponse {
	return &NoteResponse{}
}

func (b *NoteResponse) Failure(message string) JsonNoteResponse {
	b.Response = JsonNoteResponse{Type: "failed", StatusCode: 500, Data: []models.Note{}, Message: message}
	return b.Response
}

func (b *NoteResponse) Success(data []models.Note, message string) JsonNoteResponse {
	b.Response = JsonNoteResponse{Type: "success", StatusCode: 200, Data: data, Message: message}
	return b.Response
}
