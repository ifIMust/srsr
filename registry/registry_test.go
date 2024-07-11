package registry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ifIMust/srsr/registry"
)

var _ = Describe("Registry", func() {
	var reg registry.Registry

	BeforeEach(func() {
		reg = registry.NewServiceRegistry()
	})

	Describe("Register", func() {
		var reg_name string
		var reg_address string

		BeforeEach(func() {
			reg_name = "flardmaster"
			reg_address = "127.721.217.555:4343"
		})

		Context("when empty", func() {
			var id string
			var err error

			BeforeEach(func() {
				id, err = reg.Register(reg_name, reg_address)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should return non-empty ID", func() {
				Expect(id).NotTo(BeEmpty())
			})
		})
	})

	Describe("Lookup", func() {
		var lookup_name string
		var lookup_address string

		BeforeEach(func() {
			lookup_name = "flardmaster"
		})

		Context("when empty", func() {
			BeforeEach(func() {
				lookup_address = reg.Lookup(lookup_name)
			})
			It("should return empty address", func() {
				Expect(lookup_address).To(BeEmpty())
			})
		})
	})
})
