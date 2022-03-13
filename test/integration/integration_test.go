package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/m-rcd/notes/pkg/models"
	"github.com/m-rcd/notes/pkg/responses"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	Context("POST request", func() {
		AfterEach(func() {
			os.RemoveAll("/tmp/notes")
		})

		It("creates a new note", func() {
			c := http.Client{}
			postData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am a new note!","user":{"username":"Pantalaimon"}}`))
			resp, err := c.Post("http://localhost:10000/note", "application/json", postData)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			Expect(response.Type).To(Equal("success"))
			Expect(response.StatusCode).To(Equal(200))
			path := fmt.Sprintf("/tmp/notes/Pantalaimon/active/note1_%s.txt", response.Data[0].Id)
			Expect(path).To(BeAnExistingFile())
		})
	})

	Context("PATCH request", func() {
		var path string
		var note models.Note
		BeforeEach(func() {
			c := http.Client{}
			postData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am a new note!","user":{"username":"Pantalaimon"}}`))
			resp, err := c.Post("http://localhost:10000/note", "application/json", postData)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			note = response.Data[0]
			path = fmt.Sprintf("/tmp/notes/Pantalaimon/active/note1_%s.txt", note.Id)
			Expect(path).To(BeAnExistingFile())
		})

		AfterEach(func() {
			os.RemoveAll("/tmp/notes")
		})

		It("updates a previously saved note", func() {
			c := http.Client{}
			patchData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am updated!","user":{"username":"Pantalaimon"}}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note.Id, patchData)
			Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			dat, err := os.ReadFile(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(dat)).To(Equal("I am updated!"))
		})

		It("archives a previously saved note", func() {
			c := http.Client{}
			patchData := bytes.NewBuffer([]byte(`{"archived":true,"user":{"username":"Pantalaimon"}}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note.Id, patchData)
			Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()
			Expect(path).NotTo(BeAnExistingFile())
			archivedPath := fmt.Sprintf("/tmp/notes/Pantalaimon/archived/note1_%s.txt", note.Id)
			Expect(archivedPath).To(BeAnExistingFile())
		})
	})

	Context("DELETE request", func() {
		var path string
		var note models.Note
		BeforeEach(func() {
			c := http.Client{}
			postData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am a new note!","user":{"username":"Pantalaimon"}}`))
			resp, err := c.Post("http://localhost:10000/note", "application/json", postData)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			var response responses.JsonNoteResponse
			json.Unmarshal(body, &response)
			note = response.Data[0]
			path = fmt.Sprintf("/tmp/notes/Pantalaimon/active/note1_%s.txt", note.Id)
			Expect(path).To(BeAnExistingFile())
		})
		AfterEach(func() {
			os.RemoveAll("/tmp/notes")
		})

		It("deletes a note", func() {
			c := http.Client{}
			data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
			req, err := http.NewRequest("DELETE", "http://localhost:10000/note/"+note.Id, data)
			Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			Expect("/tmp/notes/Pantalaimon/active/note1.txt").NotTo(BeAnExistingFile())
		})
	})
})
