package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Comcast/webpa-common/xhttp"
	"github.com/Comcast/webpa-common/xmetrics"
)

// This can be used to further test prometheus metrics. Add your prometheus metrics here.
func NewPrometheusMockRegistry() []xmetrics.Metric {
	return []xmetrics.Metric{
		{
			Name:    OutboundRequestDuration,
			Help:    "The time for outbound request to get a response",
			Type:    "histogram",
			Buckets: []float64{0.10, 0.20, 0.50, 0.100, 0.200, 0.500, 1.00, 2.00, 5.00},
		},
	}
}

// NewRegistryMock creates a NewRegistryMock
func NewRegistryMock(m xmetrics.Module) (xmetrics.Registry, error) {
	return xmetrics.NewRegistry(nil, m)
}

/*
// outboundSenderMock is a mock CaduceusOutboundSender for testing OutboundMeasures.  It's fields are
// the minimal needed to fulfill the test below. More fields may be required for future test.
type outboundSenderMock struct {
	logger           log.Logger
	transport        *http.Transport
	outboundMeasures OutboundMeasures
	deliverUntil     time.Time
	deliveryRetries  int
	deliveryInterval time.Duration
	registry         xmetrics.Registry
}
*/

// NewOutboundSender creates a new outboundSenderMock for testing.
func NewOutboundSenderMock(m xmetrics.Module) *CaduceusOutboundSender {
	reg, _ := NewRegistryMock(m)
	return &CaduceusOutboundSender{
		logger:           getLogger(),
		transport:        &http.Transport{},
		outboundMeasures: NewOutboundMeasures(reg), // need new registry here.
		deliveryRetries:  1,
		//	deliveryInterval: time.Duration(),
	}
}

/*
// Fullfil the required methods to abide to the OutboundSender interface so OutboundRequestDuration can be tested
func (o outboundSenderMock) Update(webhook.W) error  { return errors.New("test") }
func (o outboundSenderMock) Shutdown(bool)           {}
func (o outboundSenderMock) RetiredSince() time.Time { return o.deliverUntil }
func (o outboundSenderMock) Queue(*wrp.Message)      {}
*/

// TestOutboundRequestDuration tests if OutboundRequestDuration is working correctly.
func TestOutboundRequestDuration(t *testing.T) {
	var (
		m            = NewPrometheusMockRegistry
		trans        = &transport{}
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
	resp, err := roundTripper.RoundTrip(req)
	if err != nil {
		fmt.Print(err)
	}
}
