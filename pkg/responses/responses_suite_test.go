package responses_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResponses(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Responses Suite")
}
