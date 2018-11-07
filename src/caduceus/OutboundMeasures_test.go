package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Comcast/webpa-common/xhttp"
	"github.com/Comcast/webpa-common/xmetrics"
)

// This can be used to test futre prometheus metrics. Add your prometheus metrics here.
func NewPrometheusMockRegistry() []xmetrics.Metric {
	return []xmetrics.Metric{
		{
			Name:    OutboundRequestDuration,
			Help:    "The time for outbound request to get a response",
			Type:    "histogram",
			Buckets: []float64{0.10, 0.20, 0.50, 1.00, 2.00, 5.00},
		},
	}
}

// NewRegistryMock creates a NewRegistryMock
func NewRegistryMock(m xmetrics.Module) (xmetrics.Registry, error) {
	return xmetrics.NewRegistry(nil, m)
}

// NewOutboundSender creates a new outboundSenderMock for testing.
func NewOutboundSenderMock(m xmetrics.Module) *CaduceusOutboundSender {
	reg, _ := NewRegistryMock(m)
	return &CaduceusOutboundSender{
		logger:           getLogger(),
		transport:        &http.Transport{},
		outboundMeasures: NewOutboundMeasures(reg),
		deliveryRetries:  1,
		//	deliveryInterval: time.Duration(),
	}
}

// TestOutboundRequestDuration tests if OutboundRequestDuration is working properly.
func TestOutboundRequestDuration(t *testing.T) {
	var (
		m            = NewPrometheusMockRegistry
		req          = httptest.NewRequest("GET", "http://example.com/foo", nil)
		obs          = NewOutboundSenderMock(m)
		retryOptions = xhttp.RetryOptions{
			Logger:   obs.logger,
			Retries:  obs.deliveryRetries,
			Interval: obs.deliveryInterval,
			Counter:  &SimpleCounter{},
			// Always retry on failures up to the max count.
			ShouldRetry: func(error) bool { return true },
		}
	)

	roundTripper := NewOutboundRoundTripper(retryOptions, obs)
	_, err := roundTripper.RoundTrip(req)
	if err != nil {
		fmt.Print(err)
	}
}
