/**
 * Copyright 2017 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package main

import (
	"net/http"
	"time"

	"github.com/Comcast/webpa-common/health"
	"github.com/Comcast/webpa-common/webhook"
	"github.com/Comcast/webpa-common/wrp"
	"github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

// mockHandler only needs to mock the `HandleRequest` method
type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) HandleRequest(workerID int, msg *wrp.Message) {
	m.Called(workerID, msg)
}

// mockHealthTracker needs to mock things from both the `HealthTracker`
// interface as well as the `health.Monitor` interface
type mockHealthTracker struct {
	mock.Mock
}

func (m *mockHealthTracker) SendEvent(healthFunc health.HealthFunc) {
	m.Called(healthFunc)
}

func (m *mockHealthTracker) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	m.Called(response, request)
}

func (m *mockHealthTracker) IncrementBucket(inSize int) {
	m.Called(inSize)
}

// mockOutboundSender needs to mock things that the `OutboundSender` does
type mockOutboundSender struct {
	mock.Mock
}

func (m *mockOutboundSender) Extend(until time.Time) {
	m.Called(until)
}

func (m *mockOutboundSender) Shutdown(gentle bool) {
	m.Called(gentle)
}

func (m *mockOutboundSender) RetiredSince() time.Time {
	arguments := m.Called()
	return arguments.Get(0).(time.Time)
}

func (m *mockOutboundSender) Queue(msg *wrp.Message) {
	m.Called(msg)
}

// mockSenderWrapper needs to mock things that the `SenderWrapper` does
type mockSenderWrapper struct {
	mock.Mock
}

func (m *mockSenderWrapper) Update(list []webhook.W) {
	m.Called(list)
}

func (m *mockSenderWrapper) Queue(msg *wrp.Message) {
	m.Called(msg)
}

func (m *mockSenderWrapper) Shutdown(gentle bool) {
	m.Called(gentle)
}

// mockCounter provides the mock implementation of the metrics.Counter object
type mockCounter struct {
	mock.Mock
}

func (m *mockCounter) Add(delta float64) {
	m.Called(delta)
}

func (m *mockCounter) With(labelValues ...string) metrics.Counter {
	args := m.Called(labelValues)
	return args.Get(0).(metrics.Counter)
}

// mockGauge provides the mock implementation of the metrics.Counter object
type mockGauge struct {
	mock.Mock
}

func (m *mockGauge) Add(delta float64) {
	m.Called(delta)
}

func (m *mockGauge) Set(value float64) {
	m.Called(value)
}

func (m *mockGauge) With(labelValues ...string) metrics.Gauge {
	args := m.Called(labelValues)
	return args.Get(0).(metrics.Gauge)
}

// mockHistogram provides the mock implementation of the metrics.Histogram object
type mockHistogram struct {
	mock.Mock
}

// number of observations, sum of observed values, ---> allowing the calculation of the observed values.
// is inherently a counter
//
// rate(http_request_duration_seconds_sum[5m])
// rate(http_request_duration_seconds_count[5m])

// mockCaduceusMetricsRegistry provides the mock implementation of the
// CaduceusMetricsRegistry  object
type mockCaduceusMetricsRegistry struct {
	mock.Mock
}

func (m *mockCaduceusMetricsRegistry) NewCounter(name string) metrics.Counter {
	args := m.Called(name)
	return args.Get(0).(metrics.Counter)
}

func (m *mockCaduceusMetricsRegistry) NewCounterVec(name string) *prometheus.CounterVec {
	args := m.Called(name)
	return args.Get(0).(*prometheus.CounterVec)
}

func (m *mockCaduceusMetricsRegistry) NewGauge(name string) metrics.Gauge {
	args := m.Called(name)
	return args.Get(0).(metrics.Gauge)
}

func (m *mockCaduceusMetricsRegistry) NewHistogram(name string) *prometheus.Histogram {
	args := m.Called(name)
	return args.Get(0).(*prometheus.Histogram)
}

func (m *mockCaduceusMetricsRegistry) NewHistogramVec(name string) *prometheus.HistogramVec {
	args := m.Called(name)
	return args.Get(0).(*prometheus.HistogramVec)
}
func (m *mockCaduceusMetricsRegistry) NewHistogramVecEx(namespace, subsystem, name string) *prometheus.HistogramVec {
	args := m.Called(name)
	return args.Get(0).(*prometheus.HistogramVec)
}
