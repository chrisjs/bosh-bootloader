package state_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pivotal-cf-experimental/bosh-bootloader/state"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	var (
		store   state.Store
		tempDir string
	)

	BeforeEach(func() {
		store = state.NewStore()

		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		state.ResetEncode()
	})

	Describe("Set", func() {
		It("stores the aws credentials", func() {
			err := store.Set(tempDir, state.State{
				AWSAccessKeyID:     "some-aws-access-key-id",
				AWSSecretAccessKey: "some-aws-secret-access-key",
				AWSRegion:          "some-region",
			})
			Expect(err).NotTo(HaveOccurred())

			data, err := ioutil.ReadFile(filepath.Join(tempDir, "state.json"))
			Expect(err).NotTo(HaveOccurred())
			Expect(data).To(MatchJSON(`{
				"AWSAccessKeyID": "some-aws-access-key-id",
				"AWSSecretAccessKey": "some-aws-secret-access-key",
				"AWSRegion": "some-region"
			}`))
		})

		Context("failure cases", func() {
			It("fails to open the state.json file", func() {
				err := os.Chmod(tempDir, 0000)
				Expect(err).NotTo(HaveOccurred())

				err = store.Set(tempDir, state.State{})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})

			It("fails to write the state.json file", func() {
				state.SetEncode(func(io.Writer, interface{}) error {
					return errors.New("failed to encode")
				})

				err := store.Set(tempDir, state.State{})
				Expect(err).To(MatchError("failed to encode"))
			})
		})
	})

	Describe("Get", func() {
		It("gets the aws credentials", func() {
			err := ioutil.WriteFile(filepath.Join(tempDir, "state.json"), []byte(`{
				"AWSAccessKeyID": "some-aws-access-key-id",
				"AWSSecretAccessKey": "some-aws-secret-access-key",
				"AWSRegion": "some-aws-region"
			}`), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			s, err := store.Get(tempDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(s).To(Equal(state.State{
				AWSAccessKeyID:     "some-aws-access-key-id",
				AWSSecretAccessKey: "some-aws-secret-access-key",
				AWSRegion:          "some-aws-region",
			}))
		})

		Context("when the state.json file doesn't exist", func() {
			It("returns an empty state object", func() {
				s, err := store.Get(tempDir)
				Expect(err).NotTo(HaveOccurred())

				Expect(s).To(Equal(state.State{}))
			})
		})

		Context("failure cases", func() {
			It("fails to open the state.json file", func() {
				err := os.Chmod(tempDir, 0000)
				Expect(err).NotTo(HaveOccurred())

				_, err = store.Get(tempDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})

			It("fails to decode the state.json file", func() {
				err := ioutil.WriteFile(filepath.Join(tempDir, "state.json"), []byte(`%%%%`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				_, err = store.Get(tempDir)
				Expect(err).To(MatchError(ContainSubstring("invalid character")))
			})
		})
	})
})