package manifests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/bosh-bootloader/boshinit/manifests"
)

var _ = Describe("ResourcePoolsManifestBuilder", func() {
	var resourcePoolsManifestBuilder manifests.ResourcePoolsManifestBuilder

	BeforeEach(func() {
		resourcePoolsManifestBuilder = manifests.NewResourcePoolsManifestBuilder()
	})

	Describe("ResourcePools", func() {
		It("returns all resource pools for manifest", func() {
			resourcePools := resourcePoolsManifestBuilder.Build(manifests.ManifestProperties{AvailabilityZone: "some-az"})

			Expect(resourcePools).To(HaveLen(1))
			Expect(resourcePools).To(ConsistOf([]manifests.ResourcePool{
				{
					Name:    "vms",
					Network: "private",
					Stemcell: manifests.Stemcell{
						URL:  "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=3262.12",
						SHA1: "90e9825b814da801e1aff7b02508fdada8e155cb",
					},
					CloudProperties: manifests.ResourcePoolCloudProperties{
						InstanceType: "m3.xlarge",
						EphemeralDisk: manifests.EphemeralDisk{
							Size: 25000,
							Type: "gp2",
						},
						AvailabilityZone: "some-az",
					},
				},
			}))
		})
	})
})
