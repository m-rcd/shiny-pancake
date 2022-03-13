package sql_test

import (
	_ "github.com/go-sql-driver/mysql"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("Sql", func() {
	// var s *database.SQL
	// BeforeSuite(func() {
	// 	err := godotenv.Load("./../../.env")
	// 	Expect(err).NotTo(HaveOccurred())
	// 	s = database.NewSQL(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), "127.0.0.1", "3306")
	// 	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), "127.0.0.1", "3306", "test_notes")
	// 	db, err := sql.Open("mysql", connString)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	s.Db = db

	// 	_, err = s.Db.Exec(database.CreateNoteTestTable)
	// 	Expect(err).NotTo(HaveOccurred())
	// })

	// AfterSuite(func() {
	// 	_, err := s.Db.Exec(database.DropNoteTestTable)
	// 	Expect(err).NotTo(HaveOccurred())
	// })

	// Context("CREATE", func() {
	// 	It("creates a new note", func() {
	// 		note := models.Note{Name: "Note1", Content: "Miawwww", User: models.User{Username: "Casper"}}
	// 		r := buildReader(note)
	// 		newNote, err := s.Create(r)
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(newNote.Name).To(Equal(note.Name))
	// 		Expect(newNote.Content).To(Equal(note.Content))
	// 	})
	// })
})
