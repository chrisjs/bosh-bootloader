package manifests

type ReleaseManifestBuilder struct{}

func NewReleaseManifestBuilder() ReleaseManifestBuilder {
	return ReleaseManifestBuilder{}
}

func (r ReleaseManifestBuilder) Build() []Release {
	return []Release{
		{
			Name: "bosh",
			URL:  "https://s3.amazonaws.com/bbl-precompiled-bosh-releases/release-bosh-257-on-ubuntu-trusty-stemcell-3262.tgz",
			SHA1: "fee9a89b044879ea5c9c17d239ee62c606c84a60",
		},
		{
			Name: "bosh-aws-cpi",
			URL:  "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=53",
			SHA1: "3a5988bd2b6e951995fe030c75b07c5b922e2d59",
		},
	}
}
