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
})
