package metrics

import (
	"net/http"
	"time"
)

const (
	HeaderAlias = "alias"
)

type RequestWatcher struct {
	name string
}

func NewRequestWatcher(name string) *RequestWatcher {
	return &RequestWatcher{
		name: name,
	}
}

func (m *RequestWatcher) RoundTrip(r *http.Request) (*http.Response, error) {
	var err error
	defer func(start time.Time) {
		CollectRequestsMetric(m.name, r.Header.Get(HeaderAlias), err, start)
	}(time.Now())

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
