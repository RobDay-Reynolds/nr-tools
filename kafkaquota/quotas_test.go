package kafkaquota_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/monkeyherder/nr-tools/kafkaquota"
)

var _ = Describe("Quotas", func() {
	Describe("GenerateQuotasForClients", func() {
		var (
			clients           []KafkaClient
			quotas            []KafkaQuota
			averageThroughput float64
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

			quotas = GenerateQuotasForClients(clients, 4)
		})

		It("calculates by multiplying the current average throughput by the given multiplier", func() {
			Expect(quotas).To(HaveLen(1))
			Expect(quotas[0].ClientID).To(Equal("my-client"))
			Expect(quotas[0].Quota).To(Equal(40))
			Expect(quotas[0].HumanReadableQuota).To(Equal("40 B/sec"))
		})

		Context("when the clients have throughput above 1 kB/sec but below 1 MB/sec", func() {
			BeforeEach(func() {
				averageThroughput = 2400
				maxThroughput = 2400
			})

			It("rounds up quotas to the next kB", func() {
				Expect(quotas).To(HaveLen(1))
				Expect(quotas[0].ClientID).To(Equal("my-client"))
				Expect(quotas[0].HumanReadableQuota).To(Equal("10 kB/sec"))
			})
		})

		Context("when the clients have throughput above 1 MB/sec", func() {
			BeforeEach(func() {
				averageThroughput = 2400000
				maxThroughput = 2400000
			})

			It("rounds up quotas to the next MB", func() {
				Expect(quotas).To(HaveLen(1))
				Expect(quotas[0].ClientID).To(Equal("my-client"))
				Expect(quotas[0].HumanReadableQuota).To(Equal("10 MB/sec"))
			})
		})

		Context("when the multiplier would return a value higher than 60 MB/sec", func() {
			BeforeEach(func() {
				averageThroughput = 50000000
				maxThroughput = 50000000
			})

			It("returns a quota of 60 MB/sec", func() {
				Expect(quotas).To(HaveLen(1))
				Expect(quotas[0].ClientID).To(Equal("my-client"))
				Expect(quotas[0].Quota).To(Equal(62914560))
				Expect(quotas[0].HumanReadableQuota).To(Equal("60 MB/sec"))
			})
		})

		Context("when the max throughput is higher than the average times the multiplier", func() {
			BeforeEach(func() {
				averageThroughput = 100
				maxThroughput = 1300
			})

			It("returns a quota equal to half of the multiplier times the max", func() {
				Expect(quotas).To(HaveLen(1))
				Expect(quotas[0].ClientID).To(Equal("my-client"))
				Expect(quotas[0].Quota).To(Equal(3072))
				Expect(quotas[0].HumanReadableQuota).To(Equal("3 kB/sec"))
			})

			Context("when the quota would otherwise be above 60 MB/sec", func() {
				BeforeEach(func() {
					averageThroughput = 10
					maxThroughput = 115242648
				})

				It("returns a quota of 60 MB/sec", func() {
					Expect(quotas).To(HaveLen(1))
					Expect(quotas[0].ClientID).To(Equal("my-client"))
					Expect(quotas[0].Quota).To(Equal(62914560))
					Expect(quotas[0].HumanReadableQuota).To(Equal("60 MB/sec"))
				})
			})
		})

		Context("when there are multiple clients", func() {
			JustBeforeEach(func() {
				clients = []KafkaClient{
					KafkaClient{
						ID:                "my-client",
						AverageThroughput: 1234567,
						MaxThroughput:     33837956,
					},
					KafkaClient{
						ID:                "my-client-2",
						AverageThroughput: 1234567,
						MaxThroughput:     31322453,
					},
					KafkaClient{
						ID:                "my-client-3",
						AverageThroughput: 5478000,
						MaxThroughput:     115242648,
					},
				}

				quotas = GenerateQuotasForClients(clients, 4)
			})

			It("returns quotas for each of them", func() {
				Expect(quotas).To(HaveLen(3))

				Expect(quotas[0].ClientID).To(Equal("my-client"))
				Expect(quotas[0].Quota).To(Equal(62914560))
				Expect(quotas[0].HumanReadableQuota).To(Equal("60 MB/sec"))

				Expect(quotas[1].ClientID).To(Equal("my-client-2"))
				Expect(quotas[1].Quota).To(Equal(62914560))
				Expect(quotas[1].HumanReadableQuota).To(Equal("60 MB/sec"))

				Expect(quotas[2].ClientID).To(Equal("my-client-3"))
				Expect(quotas[2].Quota).To(Equal(62914560))
				Expect(quotas[2].HumanReadableQuota).To(Equal("60 MB/sec"))
			})
		})
	})
})
