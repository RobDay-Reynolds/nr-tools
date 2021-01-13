package kafkaquota_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/monkeyherder/nr-tools/kafka-quota"
)

var _ = Describe("Quotas", func() {
	Describe("GenerateQuotasForClients", func() {
		var (
			clients           []KafkaClient
			quotas            []KafkaQuota
			averageThroughput int
			maxThroughput     int
		)

		BeforeEach(func() {
			averageThroughput = 10
			maxThroughput = 10
		})

		JustBeforeEach(func() {
			clients = []KafkaClient{
				KafkaClient{
					ID:                "my-client",
					AverageThroughput: averageThroughput,
					MaxThroughput:     maxThroughput,
				},
			}

			quotas = GenerateQuotasForClients(clients, 5)
		})

		It("calculates by multiplying the current max throughput by the given multiplier", func() {
			Expect(quotas).To(HaveLen(1))
			Expect(quotas[0].ClientID).To(Equal("my-client"))
			Expect(quotas[0].Quota).To(Equal(50))
		})

		Context("when the clients have throughput above 1 kB/s but below 1 MB/s", func() {
			BeforeEach(func() {
				averageThroughput = 1200
				maxThroughput = 1200
			})

			It("rounds up quotas to the next kB", func() {
				Expect(quotas).To(HaveLen(1))
				Expect(quotas[0].ClientID).To(Equal("my-client"))
				Expect(quotas[0].Quota).To(Equal(6144))
			})
		})

		Context("when the clients have throughput above 1 MB/s", func() {
			BeforeEach(func() {
				averageThroughput = 1200000
				maxThroughput = 1200000
			})

			It("rounds up quotas to the next MB", func() {
				Expect(quotas).To(HaveLen(1))
				Expect(quotas[0].ClientID).To(Equal("my-client"))
				Expect(quotas[0].Quota).To(Equal(6291456))
			})
		})
	})
})
