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
			resp, err := c.Do(req)
			Expect(err).NotTo(HaveOccurred())
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
			resp, err := c.Do(req)
			Expect(err).NotTo(HaveOccurred())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()
			Expect(path).NotTo(BeAnExistingFile())
			archivedPath := fmt.Sprintf("/tmp/notes/Pantalaimon/archived/note1_%s.txt", note.Id)
			Expect(archivedPath).To(BeAnExistingFile())
		})

		Context("when a note is archived", func() {
			BeforeEach(func() {
				c := http.Client{}
				patchData := bytes.NewBuffer([]byte(`{"archived":true,"user":{"username":"Pantalaimon"}}`))
				req, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note.Id, patchData)
				Expect(err).NotTo(HaveOccurred())
				_, err = c.Do(req)
				Expect(err).NotTo(HaveOccurred())
			})

			It("unarchives a previously archived note", func() {
				c := http.Client{}
				patchData := bytes.NewBuffer([]byte(`{"archived":false,"user":{"username":"Pantalaimon"}}`))
				req, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note.Id, patchData)
				Expect(err).NotTo(HaveOccurred())
				resp, err := c.Do(req)
				Expect(err).NotTo(HaveOccurred())
				_, err = ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				defer req.Body.Close()
				archivedPath := fmt.Sprintf("/tmp/notes/Pantalaimon/archived/note1_%s.txt", note.Id)
				Expect(archivedPath).NotTo(BeAnExistingFile())
				Expect(path).To(BeAnExistingFile())
			})
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

	Context("GET request", func() {
		var note1 models.Note
		var note2 models.Note

		BeforeEach(func() {
			c := http.Client{}
			postData1 := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am a new note!","user":{"username":"Pantalaimon"}}`))
			resp1, err := c.Post("http://localhost:10000/note", "application/json", postData1)
			Expect(err).NotTo(HaveOccurred())
			defer resp1.Body.Close()
			body1, err := ioutil.ReadAll(resp1.Body)
			Expect(err).NotTo(HaveOccurred())
			var response1 responses.JsonNoteResponse
			json.Unmarshal(body1, &response1)
			note1 = response1.Data[0]

			postData2 := bytes.NewBuffer([]byte(`{"name":"note2","content":"I am another note!","user":{"username":"Pantalaimon"}}`))
			resp2, err := c.Post("http://localhost:10000/note", "application/json", postData2)
			Expect(err).NotTo(HaveOccurred())
			defer resp2.Body.Close()
			body2, err := ioutil.ReadAll(resp2.Body)
			Expect(err).NotTo(HaveOccurred())
			var response2 responses.JsonNoteResponse
			json.Unmarshal(body2, &response2)
			note2 = response2.Data[0]
		})

		It("lists active notes for user", func() {
			c := http.Client{}
			data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
			req, err := http.NewRequest("GET", "http://localhost:10000/notes/active", data)
			Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			var list []models.Note
			json.Unmarshal(body, &list)
			Expect(len(list)).To(Equal(2))
			Expect(list[0]).To(Equal(note1))
			Expect(list[1]).To(Equal(note2))
		})

		Context("when notes are archived", func() {
			var archivedNote1 models.Note
			var archivedNote2 models.Note

			BeforeEach(func() {
				c := http.Client{}
				patchData1 := bytes.NewBuffer([]byte(`{"archived":true,"user":{"username":"Pantalaimon"}}`))
				req1, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note1.Id, patchData1)
				Expect(err).NotTo(HaveOccurred())
				resp1, err := c.Do(req1)
				Expect(err).NotTo(HaveOccurred())
				body1, err := ioutil.ReadAll(resp1.Body)
				Expect(err).NotTo(HaveOccurred())
				defer req1.Body.Close()
				var response1 responses.JsonNoteResponse
				json.Unmarshal(body1, &response1)
				archivedNote1 = response1.Data[0]

				patchData2 := bytes.NewBuffer([]byte(`{"archived":true,"user":{"username":"Pantalaimon"}}`))
				req2, err := http.NewRequest("PATCH", "http://localhost:10000/note/"+note2.Id, patchData2)
				Expect(err).NotTo(HaveOccurred())
				resp2, err := c.Do(req2)
				Expect(err).NotTo(HaveOccurred())
				body2, err := ioutil.ReadAll(resp2.Body)
				Expect(err).NotTo(HaveOccurred())
				defer req2.Body.Close()
				var response2 responses.JsonNoteResponse
				json.Unmarshal(body2, &response2)
				archivedNote2 = response2.Data[0]
			})

			It("lists archived notes for user", func() {
				c := http.Client{}
				data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
				req, err := http.NewRequest("GET", "http://localhost:10000/notes/archived", data)
				Expect(err).NotTo(HaveOccurred())
				resp, _ := c.Do(req)
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				defer req.Body.Close()

				var list []models.Note
				json.Unmarshal(body, &list)
				Expect(len(list)).To(Equal(2))
				Expect(list[0]).To(Equal(archivedNote1))
				Expect(list[1]).To(Equal(archivedNote2))
			})
		})
	})
})
