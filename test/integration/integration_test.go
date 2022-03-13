package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

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
			Expect("/tmp/notes/Pantalaimon/active/note1.txt").To(BeAnExistingFile())
		})
	})

	Context("PATCH request", func() {
		BeforeEach(func() {
			c := http.Client{}
			postData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am a new note!","user":{"username":"Pantalaimon"}}`))
			resp, err := c.Post("http://localhost:10000/note", "application/json", postData)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect("/tmp/notes/Pantalaimon/active/note1.txt").To(BeAnExistingFile())
		})

		AfterEach(func() {
			os.RemoveAll("/tmp/notes")
		})

		It("updates a previously saved note", func() {
			c := http.Client{}
			patchData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am updated!","user":{"username":"Pantalaimon"}}`))
			req, err := http.NewRequest("PATCH", "http://localhost:10000/note/note1", patchData)
			Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			dat, err := os.ReadFile("/tmp/notes/Pantalaimon/active/note1.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(string(dat)).To(Equal("I am updated!"))
		})
	})

	Context("DELETE request", func() {
		BeforeEach(func() {
			c := http.Client{}
			postData := bytes.NewBuffer([]byte(`{"name":"note1","content":"I am a new note!","user":{"username":"Pantalaimon"}}`))
			resp, err := c.Post("http://localhost:10000/note", "application/json", postData)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect("/tmp/notes/Pantalaimon/active/note1.txt").To(BeAnExistingFile())
		})

		AfterEach(func() {
			os.RemoveAll("/tmp/notes")
		})

		It("deletes a note", func() {
			c := http.Client{}
			data := bytes.NewBuffer([]byte(`{"username":"Pantalaimon"}`))
			req, err := http.NewRequest("DELETE", "http://localhost:10000/note/note1", data)
			Expect(err).NotTo(HaveOccurred())
			resp, _ := c.Do(req)
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer req.Body.Close()

			Expect("/tmp/notes/Pantalaimon/active/note1.txt").NotTo(BeAnExistingFile())
		})
	})
})
