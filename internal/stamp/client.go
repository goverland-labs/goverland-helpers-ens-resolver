package stamp

import (
	"errors"
	"fmt"

	internalcache "github.com/goverland-labs/helpers-ens-resolver/internal/cache"
	"github.com/goverland-labs/helpers-ens-resolver/internal/models"
)

const (
	maxLookupSize = 50
)

type Client struct {
	sdk   *SDK
	cache *internalcache.Cache
}

func NewClient(sdk *SDK) (*Client, error) {
	return &Client{
		sdk:   sdk,
		cache: internalcache.NewCache(),
	}, nil
}

func (c *Client) ResolveDomains(_ []string) ([]models.ResolvedModel, error) {
	return nil, errors.New("implement me")
}

func (c *Client) ResolveAddresses(addresses []string) ([]models.ResolvedModel, error) {
	resp := make([]models.ResolvedModel, 0, len(addresses))
	missed := make([]string, 0, len(addresses))
	for i := range addresses {
		domain, err := c.cache.Get(addresses[i])
		if err != nil {
			missed = append(missed, addresses[i])
			continue
		}

		resp = append(resp, models.ResolvedModel{
			Address: addresses[i],
			Domain:  domain,
		})
	}

	if len(missed) == 0 {
		return resp, nil
	}

	for _, chunk := range chunkSlice(missed, maxLookupSize) {
		list, err := c.sdk.ResolveAddresses(chunk)
		if err != nil {
			return nil, fmt.Errorf("c.sdk.ResolveAddresses: %w", err)
		}

		for i := range list {
			c.cache.Set(list[i].Address, list[i].Domain)
			resp = append(resp, list[i])
		}
	}

	return resp, nil
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
