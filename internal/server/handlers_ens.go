package server

import (
	"context"
	"fmt"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/protocol/enspb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/models"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/server/forms"
)

type PrimaryEnsResolver interface {
	ResolveDomains(domain []string) ([]models.ResolvedModel, error)
	ResolveAddresses(address []string) ([]models.ResolvedModel, error)
}

type AllEnsResolver interface {
	GetENSNames(owner string) ([]string, error)
}

type EnsHandler struct {
	enspb.UnimplementedEnsServer

	prEnsResolver  PrimaryEnsResolver
	allEnsResolver AllEnsResolver
}

func NewEnsHandler(per PrimaryEnsResolver, aer AllEnsResolver) *EnsHandler {
	return &EnsHandler{
		prEnsResolver:  per,
		allEnsResolver: aer,
	}
}

func (h *EnsHandler) ResolveAddresses(_ context.Context, req *enspb.ResolveAddressesRequest) (*enspb.ResolveResponse, error) {
	form, err := forms.NewResolveAddressesForm().ParseAndValidate(req)
	if err != nil {
		return nil, ResolveError(err)
	}

	list, err := h.prEnsResolver.ResolveDomains(form.Domains)
	if err != nil {
		return nil, fmt.Errorf("h.prEnsResolver.ResolveDomains: %w", err)
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

	list, err := h.prEnsResolver.ResolveAddresses(form.Addresses)
	if err != nil {
		return nil, fmt.Errorf("h.prEnsResolver.ResolveAddresses: %w", err)
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

func (h *EnsHandler) ResolveAllDomains(_ context.Context, req *enspb.ResolveAllDomainsRequest) (*enspb.ResolveAllDomainsResponse, error) {
	names, err := h.allEnsResolver.GetENSNames(req.GetAddress())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "h.allEnsResolver.GetENSNames: %s", err.Error())
	}

	return &enspb.ResolveAllDomainsResponse{
		Domains: names,
	}, nil
}
