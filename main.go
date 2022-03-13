package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/m-rcd/notes/pkg/database"
	"github.com/m-rcd/notes/pkg/database/local"
	"github.com/m-rcd/notes/pkg/database/sql"
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

	var storage string
	var workDir string
	flag.StringVar(&storage, "db", "local", "store notes on the local filesystem or in an SQL database (default: local)")

	flag.StringVar(&workDir, "directory", "/tmp", "notes location when `--db` set to `local` (default: /tmp)")
	flag.Parse()

	switch storage {
	case "sql":
		username := os.Getenv("DB_USERNAME")
		if username == "" {
			fmt.Println("DB_USERNAME must be set")
			os.Exit(1)
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			fmt.Println("DB_PASSWORD must be set")
			os.Exit(1)
		}
		db = sql.NewSQL(username, password, "127.0.0.1", "3306")
	default:
		db = local.NewLocalFileSystem(workDir)
	}

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
