package graphql

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/desteves/realm/pkg/options"
)

var _ = Describe("Graphql", func() {
	Describe("NewClient", func() {
		Context("with Options", func() {
			var opts options.ClientOptions
			Context("that are valid", func() {
				BeforeEach(func() {
					appid := "graphqlserver-lrnqt"
					auth := "anon-user"
					opts.AppID = &appid
					opts.AuthMechanism = &auth
				})
				It("should return no error", func() {
					_, err := NewClient(&opts)
					Expect(err).ShouldNot(HaveOccurred())
				})
				It("should return a client", func() {
					client, err := NewClient(&opts)
					Expect(err).ShouldNot(HaveOccurred())

					// this are private
					expectedURI := "https://stitch.mongodb.com/api/client/v2.0/app/graphqlserver-lrnqt/graphql"
					Expect(client.uri).Should(BeEquivalentTo(&expectedURI))
					Expect(client.client).ShouldNot(BeNil())
				})

			})
		})
	})
	Describe("Health", func() {
		Context("with valid client options", func() {
			var opts options.ClientOptions
			BeforeEach(func() {
				appid := "graphqlserver-lrnqt"
				auth := "anon-user"
				opts.AppID = &appid
				opts.AuthMechanism = &auth
			})

			Context("and client connected", func() {
				var nc *Client
				var err error
				BeforeEach(func() {
					nc, err = NewClient(&opts)
					Expect(err).ShouldNot(HaveOccurred())

					err = nc.Connect()
					Expect(err).ShouldNot(HaveOccurred())
				})
				It("should return no error", func() {

					var response Response
					err := nc.Health(&response)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

		})
	})
	Describe("Connect", func() {
		Context("with valid client options", func() {
			var opts options.ClientOptions
			Context("using anon-user", func() {
				BeforeEach(func() {
					appid := "graphqlserver-lrnqt"
					auth := "anon-user"
					opts.AppID = &appid
					opts.AuthMechanism = &auth
				})
				It("should return no error", func() {
					nc, err := NewClient(&opts)
					Expect(err).ShouldNot(HaveOccurred())

					err = nc.Connect()
					Expect(err).ShouldNot(HaveOccurred())

				})
			})

		})
	})
	Describe("Query", func() {
		Context("with valid client options", func() {
			var opts options.ClientOptions
			BeforeEach(func() {
				appid := "graphqlserver-lrnqt"
				auth := "anon-user"
				opts.AppID = &appid
				opts.AuthMechanism = &auth
			})
			Context("and client connected", func() {
				var nc *Client
				var err error
				BeforeEach(func() {
					nc, err = NewClient(&opts)
					Expect(err).ShouldNot(HaveOccurred())

					err = nc.Connect()
					Expect(err).ShouldNot(HaveOccurred())
				})
				Context("and valid arguments", func() {

					var ctx context.Context
					var query interface{}
					var variables map[string]interface{}
					var response Response

					BeforeEach(func() {
						ctx = context.TODO()
						// query = interface{}
						// variables = map[string]interface{}
						response = Response{}
					})
					It("should return no error", func() {
						// query = {}
						err := nc.Query(ctx, query, variables, &response)
						Expect(err).ShouldNot(HaveOccurred())
					})

				})

			})

		})
	})
	Describe("Mutate", func() {

	})
})
