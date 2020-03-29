package auth

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"

	"github.com/desteves/realm/pkg/options"
)

var _ = Describe("Auth", func() {
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

					expectedToken := oauth2.Token{}
					Expect(client.HTTPClient).Should(BeNil())
					Expect(client.Token).Should(BeEquivalentTo(&expectedToken))
				})

				// this is private
				It("should have valid oauth2 endpoint", func() {

					client, err := NewClient(&opts)
					Expect(err).ShouldNot(HaveOccurred())

					expectedEndpoint := oauth2.Endpoint{
						AuthURL:  "https://stitch.mongodb.com/api/client/v2.0/app/graphqlserver-lrnqt/auth/providers/anon-user/login",
						TokenURL: "https://stitch.mongodb.com/api/client/v2.0/auth/session",
					}
					Expect(client.HTTPClient).Should(BeNil())
					Expect(client.oauth.Endpoint).Should(BeEquivalentTo(expectedEndpoint))
				})
			})
			// Context("that are invalid", func() {
			// 	// missing appid
			// 	// unsupported provider
			// 	// missing provider-specific creds

			// })
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
				FIt("should return a valid token", func() {
					nc, err := NewClient(&opts)
					Expect(err).ShouldNot(HaveOccurred())

					err = nc.Connect()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(nc.Token.AccessToken).ShouldNot(BeEmpty())
					Expect(nc.Token.RefreshToken).ShouldNot(BeEmpty())
					Expect(nc.Token.Extra("user_id")).ShouldNot(BeNil())
					Expect(nc.Token.Extra("device_id")).ShouldNot(BeNil())
				})
			})
			Context("using local-userpass", func() {

			})
			Context("using oauth2-google", func() {

			})
			Context("using key", func() {

			})
			Context("using custom-token", func() {

			})

		})
		// should fail to retrieveFirstToken()
		Context("with invalid credentials", func() {

			Context("using invalid appid", func() {

			})

			Context("using invalid prodiver", func() {

			})

		})
	})
})
