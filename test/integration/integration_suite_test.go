package integration_test

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {

	RegisterFailHandler(Fail)
	var session *gexec.Session

	BeforeSuite(func() {
		cliBin, err := gexec.Build("github.com/m-rcd/notes")
		Expect(err).NotTo(HaveOccurred())
		command := exec.Command(cliBin)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func(g Gomega) error {
			c := http.Client{}
			resp, err := c.Get("http://localhost:10000/")
			g.Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(body).To(ContainSubstring("Welcome to Note"))
			return nil
		}, "2s").Should(Succeed())
	})

	AfterSuite(func() {
		session.Terminate().Wait()
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Integration Suite")
}
