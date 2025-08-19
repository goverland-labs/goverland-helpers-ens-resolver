package stamp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/metrics"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/models"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/pkg/sdk"
)

type SDK struct {
	baseUrl string
	client  *http.Client
}

type resolveAddressesParams struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type resolveAddressResponse struct {
	Result map[string]string `json:"result"`
}

func NewSDK(baseURL string, client *http.Client) *SDK {
	if client == nil {
		client = http.DefaultClient
	}

	return &SDK{
		baseUrl: baseURL,
		client:  client,
	}
}

func (s *SDK) ResolveAddresses(addresses []string) ([]models.ResolvedModel, error) {
	params := resolveAddressesParams{
		Method: "lookup_addresses",
		Params: addresses,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json.Marshalt: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.baseUrl, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	var result resolveAddressResponse
	err = s.sendRequest(req, params.Method, &result)
	if err != nil {
		return nil, fmt.Errorf("s.sendRequest: %w", err)
	}

	resp := make([]models.ResolvedModel, 0, len(result.Result))
	for address, domain := range result.Result {
		resp = append(resp, models.ResolvedModel{
			Address: address,
			Domain:  domain,
		})
	}

	return resp, nil
}

func (s *SDK) sendRequest(req *http.Request, alias string, result interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set(metrics.HeaderAlias, alias)

	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("s.client.Do: %w", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	log.Debug().Fields(map[string]interface{}{
		"status_code": res.StatusCode,
		"url":         req.URL.String(),
	}).Msg(string(body))

	if res.StatusCode != http.StatusOK {
		return s.parseError(res.StatusCode, body)
	}

	if result == nil {
		return nil
	}

	if err = json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

func (s *SDK) parseError(status int, body []byte) error {
	switch status {
	case http.StatusNotFound:
		return sdk.ErrNotFound
	case http.StatusUnauthorized:
		return sdk.ErrUnauthorized
	case http.StatusForbidden:
		return sdk.ErrForbidden
	case http.StatusTooManyRequests:
		var resp map[string]interface{}

		if err := json.Unmarshal(body, &resp); err != nil {
			return sdk.NewTooManyRequestsError(0)
		}

		seconds, ok := resp["retry_after"].(int)
		if !ok {
			return sdk.NewTooManyRequestsError(0)
		}

		return sdk.NewTooManyRequestsError(time.Duration(seconds) * time.Second)
	case http.StatusBadRequest:
		var errors map[string]interface{}

		if err := json.Unmarshal(body, &errors); err != nil {
			return sdk.NewValidationError(err.Error(), nil)
		}

		return sdk.NewValidationError("validation error", errors)
	}
	return sdk.ErrInternalServer
}
