package forms

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/protocol/enspb"
)

var domainRE = regexp.MustCompile(`[a-zA-Z0-9-.]+\.eth`)

var (
	ErrMissedDomains = errors.New("missed domains")
	ErrWrongDomain   = errors.New("wrong domain")
)

type ResolveAddressesForm struct {
	Domains []string
}

func NewResolveAddressesForm() *ResolveAddressesForm {
	return &ResolveAddressesForm{}
}

func (f *ResolveAddressesForm) ParseAndValidate(req *enspb.ResolveAddressesRequest) (*ResolveAddressesForm, error) {
	domains := make([]string, 0, len(req.GetDomains()))

	for _, domain := range req.GetDomains() {
		cleaned := strings.TrimSpace(domain)
		if cleaned == "" {
			continue
		}

		if !domainRE.MatchString(cleaned) {
			return nil, fmt.Errorf("%w: %s", ErrWrongDomain, cleaned)
		}

		domains = append(domains, cleaned)
	}

	if len(domains) == 0 {
		return nil, ErrMissedDomains
	}

	f.Domains = domains

	return f, nil
}

func (f *ResolveAddressesForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"domains": f.Domains,
	}
}
