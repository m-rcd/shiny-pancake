package responses

import "github.com/m-rcd/notes/pkg/models"

type JsonNoteResponse struct {
	Type       string        `json:"type"`
	StatusCode int           `json:"status_code"`
	Data       []models.Note `json:"data"`
	Message    string        `json:"message"`
}

func Failure(message string) JsonNoteResponse {
	return JsonNoteResponse{Type: "failed", StatusCode: 500, Data: []models.Note{}, Message: message}
}

func Success(data []models.Note, message string) JsonNoteResponse {
	return JsonNoteResponse{Type: "success", StatusCode: 200, Data: data, Message: message}
}
