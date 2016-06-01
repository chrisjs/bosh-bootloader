package commands

import (
	"github.com/pivotal-cf-experimental/bosh-bootloader/bosh"
	"github.com/pivotal-cf-experimental/bosh-bootloader/storage"
)

type DeleteLBs struct {
	awsCredentialValidator    awsCredentialValidator
	availabilityZoneRetriever availabilityZoneRetriever
	certificateManager        certificateManager
	infrastructureManager     infrastructureManager
	logger                    logger
	boshCloudConfigurator     boshCloudConfigurator
	cloudConfigManager        cloudConfigManager
	boshClientProvider        boshClientProvider
}

type cloudConfigManager interface {
	Update(cloudConfigInput bosh.CloudConfigInput, boshClient bosh.Client) error
}

func NewDeleteLBs(awsCredentialValidator awsCredentialValidator, availabilityZoneRetriever availabilityZoneRetriever,
	certificateManager certificateManager, infrastructureManager infrastructureManager, logger logger,
	boshCloudConfigurator boshCloudConfigurator, cloudConfigManager cloudConfigManager,
	boshClientProvider boshClientProvider,
) DeleteLBs {
	return DeleteLBs{
		awsCredentialValidator:    awsCredentialValidator,
		availabilityZoneRetriever: availabilityZoneRetriever,
		certificateManager:        certificateManager,
		infrastructureManager:     infrastructureManager,
		logger:                    logger,
		boshCloudConfigurator:     boshCloudConfigurator,
		cloudConfigManager:        cloudConfigManager,
		boshClientProvider:        boshClientProvider,
	}
}

func (c DeleteLBs) Execute(subcommandFlags []string, state storage.State) (storage.State, error) {
	err := c.awsCredentialValidator.Validate()
	if err != nil {
		return state, err
	}

	if err := checkBBLAndLB(state, c.boshClientProvider, c.infrastructureManager); err != nil {
		return state, err
	}

	azs, err := c.availabilityZoneRetriever.Retrieve(state.AWS.Region)
	if err != nil {
		return state, err
	}

	stack, err := c.infrastructureManager.Describe(state.Stack.Name)
	if err != nil {
		return state, err
	}

	cloudConfigInput := c.boshCloudConfigurator.Configure(stack, azs)
	cloudConfigInput.LBs = nil

	boshClient := c.boshClientProvider.Client(state.BOSH.DirectorAddress, state.BOSH.DirectorUsername, state.BOSH.DirectorPassword)

	err = c.cloudConfigManager.Update(cloudConfigInput, boshClient)
	if err != nil {
		return state, err
	}

	_, err = c.infrastructureManager.Update(state.KeyPair.Name, len(azs), state.Stack.Name, "", "")
	if err != nil {
		return state, err
	}

	err = c.certificateManager.Delete(state.Stack.CertificateName)
	if err != nil {
		return state, err
	}

	state.Stack.LBType = ""
	state.Stack.CertificateName = ""

	return state, nil
}
