package server

import (
	"context"
	"fmt"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/models"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/server/forms"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/protocol/enspb"
)

type EnsClient interface {
	ResolveDomains(domain []string) ([]models.ResolvedModel, error)
	ResolveAddresses(address []string) ([]models.ResolvedModel, error)
}

type EnsHandler struct {
	enspb.UnimplementedEnsServer

	client EnsClient
}

func NewEnsHandler(c EnsClient) *EnsHandler {
	return &EnsHandler{
		client: c,
	}
}

func (h *EnsHandler) ResolveAddresses(_ context.Context, req *enspb.ResolveAddressesRequest) (*enspb.ResolveResponse, error) {
	form, err := forms.NewResolveAddressesForm().ParseAndValidate(req)
	if err != nil {
		return nil, ResolveError(err)
	}

	list, err := h.client.ResolveDomains(form.Domains)
	if err != nil {
		return nil, fmt.Errorf("h.client.ResolveDomains: %w", err)
	}

	addresses := make([]*enspb.Address, 0, len(form.Domains))
	for i := range list {
		addresses = append(addresses, &enspb.Address{
			Address: list[i].Address,
			EnsName: list[i].Domain,
		})
	}

	return &enspb.ResolveResponse{Addresses: addresses}, nil
}

func (h *EnsHandler) ResolveDomains(_ context.Context, req *enspb.ResolveDomainsRequest) (*enspb.ResolveResponse, error) {
	form, err := forms.NewResolveDomainsForm().ParseAndValidate(req)
	if err != nil {
		return nil, ResolveError(err)
	}

	list, err := h.client.ResolveAddresses(form.Addresses)
	if err != nil {
		return nil, fmt.Errorf("h.client.ResolveAddresses: %w", err)
	}

	addresses := make([]*enspb.Address, 0, len(form.Addresses))
	for i := range list {
		addresses = append(addresses, &enspb.Address{
			Address: list[i].Address,
			EnsName: list[i].Domain,
		})
	}

	return &enspb.ResolveResponse{Addresses: addresses}, nil
}
