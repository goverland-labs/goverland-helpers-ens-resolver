package alchemy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/metrics"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/pkg/sdk"
)

const (
	baseURL = "https://eth-mainnet.g.alchemy.com/nft/v2/%s/getNFTs?owner=%s&contractAddresses[]=0x57f1887a8BF19b14fC0dF6Fd9B2acc9Af147eA85"
)

type (
	Client struct {
		apiKey string
		client *http.Client
	}

	NFT struct {
		Title    string `json:"title"`
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		TokenID string `json:"tokenId"`
	}

	NFTResponse struct {
		OwnedNfts []NFT  `json:"ownedNfts"`
		PageKey   string `json:"pageKey,omitempty"`
	}
)

func NewClient(apiKey string, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		apiKey: apiKey,
		client: client,
	}
}

func (s *Client) GetENSNames(owner string) ([]string, error) {
	url := fmt.Sprintf(baseURL, s.apiKey, owner)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set(metrics.HeaderAlias, "getNFTs")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("s.client.Get: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	log.Debug().Fields(map[string]interface{}{
		"status_code": resp.StatusCode,
		"url":         req.URL.String(),
	}).Msg(string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, s.parseError(resp.StatusCode, body)
	}

	var nftResp NFTResponse
	if err = json.Unmarshal(body, &nftResp); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	names := make([]string, 0, len(nftResp.OwnedNfts))
	for _, nft := range nftResp.OwnedNfts {
		switch {
		case nft.Title != "":
			names = append(names, nft.Title)
		case nft.Metadata.Name != "":
			names = append(names, nft.Metadata.Name)
		}
	}

	return names, nil
}

func (s *Client) parseError(status int, body []byte) error {
	switch status {
	case http.StatusNotFound:
		return sdk.ErrNotFound
	case http.StatusUnauthorized:
		return sdk.ErrUnauthorized
	case http.StatusForbidden:
		return sdk.ErrForbidden
	case http.StatusTooManyRequests:
		return sdk.NewTooManyRequestsError(0)
	case http.StatusBadRequest:
		var errors map[string]interface{}

		if err := json.Unmarshal(body, &errors); err != nil {
			return sdk.NewValidationError(err.Error(), nil)
		}

		return sdk.NewValidationError("validation error", errors)
	}

	return sdk.ErrInternalServer
}
