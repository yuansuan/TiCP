package openstack

import (
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

var (
	ErrEmptyIdentityEndpoint            = errors.New("empty identity endpoint")
	ErrEmptyApplicationCredentialID     = errors.New("empty application credentialID")
	ErrEmptyApplicationCredentialSecret = errors.New("empty application credential secret")
)

func newProvider(authOpts gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	var err error
	if err = validateAuthOpts(authOpts); err != nil {
		return nil, fmt.Errorf("invalid auth options, %w", err)
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, fmt.Errorf("authenticate client failed, %w", err)
	}

	return provider, nil
}

func validateAuthOpts(authOpts gophercloud.AuthOptions) error {
	if authOpts.IdentityEndpoint == "" {
		return ErrEmptyIdentityEndpoint
	}

	if authOpts.ApplicationCredentialID == "" {
		return ErrEmptyApplicationCredentialID
	}

	if authOpts.ApplicationCredentialSecret == "" {
		return ErrEmptyApplicationCredentialSecret
	}

	return nil
}
