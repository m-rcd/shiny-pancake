package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m-rcd/notes/pkg/database"
	"github.com/m-rcd/notes/pkg/models"
	"github.com/m-rcd/notes/pkg/responses"
)

type Handler struct {
	db database.Database
}

var noteResponse = responses.NewNoteResponse()
var response responses.JsonNoteResponse

func New(db database.Database) Handler {
	return Handler{db: db}
}

func (h *Handler) CreateNewNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newNote, err := h.db.Create(r.Body)

	if err != nil {
		response = noteResponse.Failure(err.Error())
	} else {
		response = noteResponse.Success([]models.Note{newNote}, "The note was successfully created")
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	note, err := h.db.Update(name, r.Body)
	if err != nil {
		response = noteResponse.Failure(err.Error())
	} else {
		response = noteResponse.Success([]models.Note{note}, "The note was successfully updated")
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Note!")
}
