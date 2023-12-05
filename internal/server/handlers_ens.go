package server

import (
	"context"
	"fmt"

	"github.com/goverland-labs/helpers-ens-resolver/internal/models"
	"github.com/goverland-labs/helpers-ens-resolver/internal/server/forms"
	"github.com/goverland-labs/helpers-ens-resolver/proto"
)

type EnsClient interface {
	ResolveDomains(domain []string) ([]models.ResolvedModel, error)
	ResolveAddresses(address []string) ([]models.ResolvedModel, error)
}

type EnsHandler struct {
	proto.UnimplementedEnsServer

	client EnsClient
}

func NewEnsHandler(c EnsClient) *EnsHandler {
	return &EnsHandler{
		client: c,
	}
}

func (h *EnsHandler) ResolveAddresses(_ context.Context, req *proto.ResolveAddressesRequest) (*proto.ResolveResponse, error) {
	form, err := forms.NewResolveAddressesForm().ParseAndValidate(req)
	if err != nil {
		return nil, ResolveError(err)
	}

	list, err := h.client.ResolveDomains(form.Domains)
	if err != nil {
		return nil, fmt.Errorf("h.client.ResolveDomains: %w", err)
	}

	addresses := make([]*proto.Address, 0, len(form.Domains))
	for i := range list {
		addresses = append(addresses, &proto.Address{
			Address: list[i].Address,
			EnsName: list[i].Domain,
		})
	}

	return &proto.ResolveResponse{Addresses: addresses}, nil
}

func (h *EnsHandler) ResolveDomains(_ context.Context, req *proto.ResolveDomainsRequest) (*proto.ResolveResponse, error) {
	form, err := forms.NewResolveDomainsForm().ParseAndValidate(req)
	if err != nil {
		return nil, ResolveError(err)
	}

	list, err := h.client.ResolveAddresses(form.Addresses)
	if err != nil {
		return nil, fmt.Errorf("h.client.ResolveAddresses: %w", err)
	}

	addresses := make([]*proto.Address, 0, len(form.Addresses))
	for i := range list {
		addresses = append(addresses, &proto.Address{
			Address: list[i].Address,
			EnsName: list[i].Domain,
		})
	}

	return &proto.ResolveResponse{Addresses: addresses}, nil
}
