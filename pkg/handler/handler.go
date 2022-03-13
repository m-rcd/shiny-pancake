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

func New(db database.Database) Handler {
	return Handler{
		db: db,
	}
}

func (h *Handler) CreateNewNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response responses.JsonNoteResponse
	newNote, err := h.db.Create(r.Body)
	if err != nil {
		response = responses.Failure(err.Error())
	} else {
		response = responses.Success([]models.Note{newNote}, "The note was successfully created")
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var response responses.JsonNoteResponse
	note, err := h.db.Update(id, r.Body)
	if err != nil {
		response = responses.Failure(err.Error())
	} else {
		response = responses.Success([]models.Note{note}, "The note was successfully updated")
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var response responses.JsonNoteResponse
	err := h.db.Delete(id, r.Body)
	if err != nil {
		response = responses.Failure(err.Error())
	} else {
		response = responses.Success([]models.Note{}, "The note was successfully deleted")
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ListActiveNotes(w http.ResponseWriter, r *http.Request) {
	var response responses.JsonNoteResponse

	notes, err := h.db.ListActiveNotes(r.Body)
	if err != nil {
		response = responses.Failure(err.Error())
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(notes)
	}
}

func (h *Handler) ListArchivedNotes(w http.ResponseWriter, r *http.Request) {
	var response responses.JsonNoteResponse

	notes, err := h.db.ListArchivedNotes(r.Body)
	if err != nil {
		response = responses.Failure(err.Error())
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(notes)
	}
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Note!")
}
