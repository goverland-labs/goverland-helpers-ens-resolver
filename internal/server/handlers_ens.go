package server

import (
	"context"
	"errors"

	"helpers-ens-resolver/internal/infura"
	"helpers-ens-resolver/internal/server/forms"
	"helpers-ens-resolver/proto"
)

type EnsClient interface {
	ResolveDomain(domain string) (string, error)
	ResolveAddress(address string) (string, error)
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

	addresses := make([]*proto.Address, 0, len(form.Domains))

	for _, domain := range form.Domains {
		addr := &proto.Address{
			EnsName: domain,
		}

		// TODO: Use cache
		resolved, err := h.client.ResolveDomain(domain)

		if err == nil {
			addr.Address = resolved
		}

		if err != nil && !errors.Is(err, infura.ErrUnknownAddress) {
			return nil, ResolveError(err)
		}

		addresses = append(addresses, addr)
	}

	return &proto.ResolveResponse{Addresses: addresses}, nil
}

func (h *EnsHandler) ResolveDomains(_ context.Context, req *proto.ResolveDomainsRequest) (*proto.ResolveResponse, error) {
	form, err := forms.NewResolveDomainsForm().ParseAndValidate(req)
	if err != nil {
		return nil, ResolveError(err)
	}

	addresses := make([]*proto.Address, 0, len(form.Addresses))

	for _, addr := range form.Addresses {
		res := &proto.Address{
			Address: addr,
		}

		// TODO: Use cache
		resolved, err := h.client.ResolveAddress(addr)

		if err == nil {
			res.EnsName = resolved
		}

		if err != nil && !errors.Is(err, infura.ErrUnknownDomain) {
			return nil, ResolveError(err)
		}

		addresses = append(addresses, res)
	}

	return &proto.ResolveResponse{Addresses: addresses}, nil
}
