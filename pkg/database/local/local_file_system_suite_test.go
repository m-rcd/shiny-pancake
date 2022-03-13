package local_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLocalFileSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Local File System Suite")
}
