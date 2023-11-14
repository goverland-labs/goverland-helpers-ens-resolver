package forms

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/goverland-labs/helpers-ens-resolver/proto"
)

var addressRE = regexp.MustCompile(`0x[a-fA-F0-9]{40}`)

var (
	ErrMissedAddresses = errors.New("missed addresses")
	ErrWrongAddress    = errors.New("wrong address")
)

type ResolveDomainsForm struct {
	Addresses []string
}

func NewResolveDomainsForm() *ResolveDomainsForm {
	return &ResolveDomainsForm{}
}

func (f *ResolveDomainsForm) ParseAndValidate(req *proto.ResolveDomainsRequest) (*ResolveDomainsForm, error) {
	addresses := make([]string, 0, len(req.GetAddresses()))

	for _, addr := range req.GetAddresses() {
		cleaned := strings.TrimSpace(addr)
		if cleaned == "" {
			continue
		}

		if !addressRE.MatchString(cleaned) {
			return nil, fmt.Errorf("%w: %s", ErrWrongAddress, cleaned)
		}

		addresses = append(addresses, cleaned)
	}

	if len(addresses) == 0 {
		return nil, ErrMissedAddresses
	}

	f.Addresses = addresses

	return f, nil
}

func (f *ResolveDomainsForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"addresses": f.Addresses,
	}
}
