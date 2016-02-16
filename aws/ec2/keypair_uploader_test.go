package ec2_test

import (
	"errors"

	goaws "github.com/aws/aws-sdk-go/aws"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pivotal-cf-experimental/bosh-bootloader/aws/ec2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeEC2Client struct {
	ImportKeyPairCall struct {
		Receives struct {
			Input *awsec2.ImportKeyPairInput
		}
		Returns struct {
			Error error
		}
	}
}

func (c *FakeEC2Client) ImportKeyPair(input *awsec2.ImportKeyPairInput) (*awsec2.ImportKeyPairOutput, error) {
	c.ImportKeyPairCall.Receives.Input = input

	return nil, c.ImportKeyPairCall.Returns.Error
}

var _ = Describe("KeypairUploader", func() {
	var (
		ec2Client *FakeEC2Client
		uploader  ec2.KeypairUploader
	)

	BeforeEach(func() {
		ec2Client = &FakeEC2Client{}
		uploader = ec2.NewKeypairUploader()
	})

	Describe("Upload", func() {
		It("uploads the keypair to AWS", func() {
			err := uploader.Upload(ec2Client, ec2.Keypair{
				Name: "some-keypair",
				Key:  []byte("some-key"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(ec2Client.ImportKeyPairCall.Receives.Input).To(Equal(&awsec2.ImportKeyPairInput{
				KeyName:           goaws.String("some-keypair"),
				PublicKeyMaterial: []byte("some-key"),
			}))
		})

		Context("failure cases", func() {
			Context("when the import fails", func() {
				It("returns an error", func() {
					ec2Client.ImportKeyPairCall.Returns.Error = errors.New("failed to import keypair")

					err := uploader.Upload(ec2Client, ec2.Keypair{})
					Expect(err).To(MatchError(errors.New("failed to import keypair")))
				})
			})
		})
	})
})