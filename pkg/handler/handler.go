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
	id := mux.Vars(r)["id"]

	note, err := h.db.Update(id, r.Body)
	if err != nil {
		response = noteResponse.Failure(err.Error())
	} else {
		response = noteResponse.Success([]models.Note{note}, "The note was successfully updated")
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.db.Delete(id, r.Body)
	if err != nil {
		response = noteResponse.Failure(err.Error())
	} else {
		response = noteResponse.Success([]models.Note{}, "The note was successfully deleted")
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ListActiveNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.db.ListActiveNotes(r.Body)
	if err != nil {
		response = noteResponse.Failure(err.Error())
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(notes)
	}

}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Note!")
}
