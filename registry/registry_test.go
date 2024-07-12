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

	Describe("Deregister", func() {
		var err error
		var id string
		var reg_name string
		var reg_address string

		BeforeEach(func() {
			id = "234343242342"
			reg_name = "callisto"
			reg_address = "234.234.222.111:3232"
		})
		Context("when empty", func() {
			BeforeEach(func() {
				err = reg.Deregister(id)
			})
			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("when non-empty, but ID is not a match", func() {
			BeforeEach(func() {
				reg.Register(reg_name, reg_address)
				err = reg.Deregister(id)
			})
			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("with matching ID", func() {
			BeforeEach(func() {
				id, err = reg.Register(reg_name, reg_address)
				err = reg.Deregister(id)
			})
			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			Context("after Deregister", func() {
				When("calling Lookup on deregistered ID", func() {
					It("should return an empty address", func() {
						Expect(reg.Lookup(reg_name)).To(BeEmpty())
					})
				})
				When("calling Deregister again", func() {
					It("should return an error", func() {
						err = reg.Deregister(id)
						Expect(err).NotTo(BeNil())
					})
				})
			})
		})
	})

	Describe("Lookup", func() {
		var reg_address string
		var lookup_name string
		var lookup_address string

		BeforeEach(func() {
			lookup_name = "flardmaster"
			reg_address = "128.128.128.128:128"
		})

		Context("when empty", func() {
			BeforeEach(func() {
				lookup_address = reg.Lookup(lookup_name)
			})
			It("should return empty address", func() {
				Expect(lookup_address).To(BeEmpty())
			})
		})

		Context("when name exists", func() {
			BeforeEach(func() {
				reg.Register(lookup_name, reg_address)
				lookup_address = reg.Lookup(lookup_name)
			})
			It("should return stored address", func() {
				Expect(lookup_address).To(Equal(reg_address))
			})
		})
	})
})
