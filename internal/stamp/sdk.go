package stamp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/models"
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
	err = s.sendRequest(req, &result)
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

func (s *SDK) sendRequest(req *http.Request, result interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; charset=utf-8")

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

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

func (s *SDK) parseError(status int, body []byte) error {
	switch status {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusTooManyRequests:
		var resp map[string]interface{}

		if err := json.Unmarshal(body, &resp); err != nil {
			return NewTooManyRequestsError(0)
		}

		seconds, ok := resp["retry_after"].(int)
		if !ok {
			return NewTooManyRequestsError(0)
		}

		return NewTooManyRequestsError(time.Duration(seconds) * time.Second)
	case http.StatusBadRequest:
		var errors map[string]interface{}

		if err := json.Unmarshal(body, &errors); err != nil {
			return NewValidationError(err.Error(), nil)
		}

		return NewValidationError("validation error", errors)
	}
	return ErrInternalServer
}
