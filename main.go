package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/m-rcd/notes/pkg/database"
	"github.com/m-rcd/notes/pkg/handler"

	"github.com/gorilla/mux"
)

var (
	db  database.Database
	err error
	h   handler.Handler
)

func main() {
	fmt.Println("Listening on port 10000")
	var workDir string

	flag.StringVar(&workDir, "directory", "/tmp", "notes location when `--db` set to `local` (default: /tmp)")
	flag.Parse()

	db = database.NewLocalFileSystem(workDir)
	err = db.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	handleRequests(db)
}

func handleRequests(db database.Database) {
	h = handler.New(db)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", h.HomePage)
	myRouter.HandleFunc("/note", h.CreateNewNote).Methods("POST")
	myRouter.HandleFunc("/note/{id}", h.UpdateNote).Methods("PATCH")
	myRouter.HandleFunc("/note/{id}", h.DeleteNote).Methods("DELETE")
	myRouter.HandleFunc("/notes/active", h.ListActiveNotes).Methods("GET")
	myRouter.HandleFunc("/notes/archived", h.ListArchivedNotes).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
