package infura

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"

	internalcache "github.com/goverland-labs/helpers-ens-resolver/internal/cache"
	"github.com/goverland-labs/helpers-ens-resolver/internal/config"
	"github.com/goverland-labs/helpers-ens-resolver/internal/models"
)

type Client struct {
	client *ethclient.Client
	cache  *internalcache.Cache
}

func NewClient(infura config.Infura) (*Client, error) {
	cl, err := ethclient.Dial(fmt.Sprintf("%s%s", infura.Endpoint, infura.Key))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: cl,
		cache:  internalcache.NewCache(),
	}, nil
}

func (c *Client) ResolveDomains([]string) ([]models.ResolvedModel, error) {
	return nil, errors.New("implement me")
}

func (c *Client) ResolveAddresses([]string) ([]models.ResolvedModel, error) {
	return nil, errors.New("implement me")
}

func (c *Client) ResolveDomain(domain string) (string, error) {
	addr, err := c.cache.Get(domain)
	if err == nil {
		return addr, nil
	}

	address, err := ens.Resolve(c.client, domain)
	if address == ens.UnknownAddress {
		return "", ErrUnknownAddress
	}

	if err != nil {
		return "", err
	}

	c.cache.Set(domain, address.String())

	return address.String(), nil
}

func (c *Client) ResolveAddress(address string) (string, error) {
	domain, err := c.cache.Get(address)
	if err == nil {
		return domain, nil
	}

	domain, err = ens.ReverseResolve(c.client, common.HexToAddress(address))
	if domain == "" {
		return "", ErrUnknownDomain
	}

	if err != nil {
		return "", err
	}

	c.cache.Set(address, domain)

	return domain, nil
}
