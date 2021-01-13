package kafkaquota_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKafkaQuota(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KafkaQuota Suite")
}
