package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/m-rcd/notes/pkg/models"
	"github.com/m-rcd/notes/pkg/responses"
	"github.com/m-rcd/notes/pkg/utils"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("When running the note server", func() {
	var (
		tempDir string
	)

	BeforeEach(func() {
		Expect(godotenv.Load("./../../.env")).To(Succeed())

		var err error
		tempDir, err = ioutil.TempDir("", "local_integration_test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	localArgsBuilder := func() []string {
		return []string{"--db", "local", "--directory", tempDir}
	}

	sqlArgsBuilder := func() []string {
		return []string{"--db", "sql"}
	}

	table.DescribeTable("the user can manipulate notes", func(getArgs func() []string) {
		var (
			err     error
			session *gexec.Session
			note1   models.Note
			note2   models.Note

			c    = http.Client{}
			args = getArgs()
		)

		if databaseNotRunning(args[1]) {
			Skip("skipped because SQL database not set and running")
		}

		command := exec.Command(cliBin, args...)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		defer func() {
			session.Terminate().Wait()
		}()

		By("creating a note")
		Eventually(func(g Gomega) error {
			postData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am a new note!","user":{"username":"Pantalaimon"}}`))
			resp, err := c.Post("http://localhost:10000/note", "application/json", postData)
			g.Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			g.Expect(response.Type).To(Equal("success"))
			g.Expect(response.StatusCode).To(Equal(200))
			g.Expect(response.Message).To(Equal("The note was successfully created"))
			note1 = response.Data[0]
			g.Expect(note1.Name).To(Equal("note1"))

			return nil
		}, "20s").Should(Succeed())

		By("updating a note")
		Eventually(func(g Gomega) error {
			patchData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am updated!","user":{"username":"Pantalaimon"}}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note1.Id, patchData)
			g.Expect(err).NotTo(HaveOccurred())
			resp, err := c.Do(req)
			g.Expect(err).NotTo(HaveOccurred())
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			g.Expect(response.Type).To(Equal("success"))
			g.Expect(response.StatusCode).To(Equal(200))
			g.Expect(response.Message).To(Equal("The note was successfully updated"))
			note1 = response.Data[0]
			g.Expect(note1.Content).To(Equal("I am updated!"))

			return nil
		}, "20s").Should(Succeed())

		By("creating a second note")
		Eventually(func(g Gomega) error {
			postData := bytes.NewBuffer([]byte(`{"name":"note2","content":"I am a second note!","user":{"username":"Pantalaimon"}}`))
			resp, err := c.Post("http://localhost:10000/note", "application/json", postData)
			g.Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			g.Expect(response.Type).To(Equal("success"))
			g.Expect(response.StatusCode).To(Equal(200))
			g.Expect(response.Message).To(Equal("The note was successfully created"))
			note2 = response.Data[0]
			g.Expect(note2.Name).To(Equal("note2"))

			return nil
		}, "20s").Should(Succeed())

		By("listing active notes")
		Eventually(func(g Gomega) error {
			data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
			req, err := http.NewRequest("GET", "http://localhost:10000/notes/active", data)
			g.Expect(err).NotTo(HaveOccurred())
			resp, err := c.Do(req)
			g.Expect(err).NotTo(HaveOccurred())
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			var list []models.Note
			json.Unmarshal(body, &list)
			g.Expect(len(list)).To(Equal(2))
			g.Expect(list[0]).To(Equal(note1))
			g.Expect(list[1]).To(Equal(note2))
			return nil

		}, "20s").Should(Succeed())

		By("archiving a note")
		Eventually(func(g Gomega) error {
			patchData := bytes.NewBuffer([]byte(`{"archived":true,"user":{"username":"Pantalaimon"}}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note1.Id, patchData)
			g.Expect(err).NotTo(HaveOccurred())
			resp, err := c.Do(req)
			g.Expect(err).NotTo(HaveOccurred())
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()
			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			g.Expect(response.Type).To(Equal("success"))
			g.Expect(response.StatusCode).To(Equal(200))
			g.Expect(response.Message).To(Equal("The note was successfully updated"))
			note1 = response.Data[0]
			g.Expect(note1.Archived).To(BeTrue())
			return nil

		}, "20s").Should(Succeed())

		By("listing archived notes")
		Eventually(func(g Gomega) error {
			data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
			req, err := http.NewRequest("GET", "http://localhost:10000/notes/archived", data)
			g.Expect(err).NotTo(HaveOccurred())
			resp, err := c.Do(req)
			g.Expect(err).NotTo(HaveOccurred())
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			var list []models.Note
			json.Unmarshal(body, &list)
			g.Expect(len(list)).To(Equal(1))
			g.Expect(list[0]).To(Equal(note1))
			return nil

		}, "20s").Should(Succeed())

		By("unarchiving a note")
		Eventually(func(g Gomega) error {
			patchData := bytes.NewBuffer([]byte(`{"archived":false,"user":{"username":"Pantalaimon"}}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note1.Id, patchData)
			g.Expect(err).NotTo(HaveOccurred())
			resp, err := c.Do(req)
			g.Expect(err).NotTo(HaveOccurred())
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()
			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			g.Expect(response.Type).To(Equal("success"))
			g.Expect(response.StatusCode).To(Equal(200))
			g.Expect(response.Message).To(Equal("The note was successfully updated"))
			note1 = response.Data[0]
			g.Expect(note1.Archived).To(BeFalse())
			return nil

		}, "20s").Should(Succeed())

		By("deleting the first note")
		Eventually(func(g Gomega) error {
			data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
			req, err := http.NewRequest("DELETE", "http://localhost:10000/note/"+note1.Id, data)
			g.Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			g.Expect(response.Type).To(Equal("success"))
			g.Expect(response.StatusCode).To(Equal(200))
			g.Expect(response.Message).To(Equal("The note was successfully deleted"))

			return nil

		}, "20s").Should(Succeed())

		By("deleting the second note")
		Eventually(func(g Gomega) error {
			data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
			req, err := http.NewRequest("DELETE", "http://localhost:10000/note/"+note2.Id, data)
			g.Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			g.Expect(response.Type).To(Equal("success"))
			g.Expect(response.StatusCode).To(Equal(200))
			g.Expect(response.Message).To(Equal("The note was successfully deleted"))

			return nil

		}, "20s").Should(Succeed())
	},
		table.Entry("local", localArgsBuilder),
		table.Entry("sql", sqlArgsBuilder),
	)
})

func databaseNotRunning(storage string) bool {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	return storage == "sql" && (!utils.IsSet(username) || !utils.IsSet(password))
}
