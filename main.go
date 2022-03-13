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
	"github.com/m-rcd/notes/pkg/utils"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Listening on port 10000")

	var (
		storage string
		workDir string
	)

	flag.StringVar(&storage, "db", "local", "store notes on the local filesystem or in an SQL database (default: local)")
	flag.StringVar(&workDir, "directory", "/tmp", "notes location when `--db` set to `local` (default: /tmp)")
	flag.Parse()

	db := getDb(storage, workDir)

	if err := db.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	handleRequests(db)
}

func handleRequests(db database.Database) {
	h := handler.New(db)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", h.HomePage)
	myRouter.HandleFunc("/note", h.CreateNewNote).Methods("POST")
	myRouter.HandleFunc("/note/{id}", h.UpdateNote).Methods("PATCH")
	myRouter.HandleFunc("/note/{id}", h.DeleteNote).Methods("DELETE")
	myRouter.HandleFunc("/notes/active", h.ListActiveNotes).Methods("GET")
	myRouter.HandleFunc("/notes/archived", h.ListArchivedNotes).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func getDb(storage, workDir string) database.Database {
	var db database.Database

	switch storage {
	case "sql":
		username := os.Getenv("DB_USERNAME")
		if !utils.IsSet(username) {
			fmt.Println("DB_USERNAME must be set")
			os.Exit(1)
		}
		password := os.Getenv("DB_PASSWORD")
		if !utils.IsSet(password) {
			fmt.Println("DB_PASSWORD must be set")
			os.Exit(1)
		}
		db = sql.NewSQL(username, password, "127.0.0.1", "3306")
	default:
		db = local.NewLocalFileSystem(workDir)
	}

	return db
}
