package sql_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSqlSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "sql Suite")
}
