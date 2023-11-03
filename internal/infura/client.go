package infura

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"

	"helpers-ens-resolver/internal/config"
)

type Client struct {
	client *ethclient.Client
	cache  *Cache
}

func NewClient(infura config.Infura) (*Client, error) {
	cl, err := ethclient.Dial(fmt.Sprintf("%s%s", infura.Endpoint, infura.Key))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: cl,
		cache:  NewCache(),
	}, nil
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
