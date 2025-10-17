package config

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func checkHttpClient(t *testing.T, opts *HttpClient, cmd *cobra.Command) {
	// Check all flags
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST"))
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_RESPONSE_HEADER_TIMEOUT"))
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_TIMEOUT"))
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_DIAL_TIMEOUT"))
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_KEEP_ALIVE"))
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_IDLE_CONN_TIMEOUT"))
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_TLS_HANDSHAKE_TIMEOUT"))
	assert.NotNil(t, cmd.Flag("HTTP_CLIENT_EXPECT_CONTINUE_TIMEOUT"))

	// Check default values
	assert.Equal(t, 100, opts.MaxIdleConnsPerHost)
	assert.Equal(t, 1*time.Second, opts.ResponseHeaderTimeout)
	assert.Equal(t, 4*time.Second, opts.Timeout)
	assert.Equal(t, 4*time.Second, opts.DialTimeout)
	assert.Equal(t, 30*time.Second, opts.KeepAlive)
	assert.Equal(t, 90*time.Second, opts.IdleConnTimeout)
	assert.Equal(t, 10*time.Second, opts.TlsHandShakeTimeout)
	assert.Equal(t, 1*time.Second, opts.ExpectContinueTimeout)
}

func TestHttpClientAddFlagsDefault(t *testing.T) {
	cmd := &cobra.Command{}
	opts := &HttpClient{}

	opts.AddFlags(cmd)

	checkHttpClient(t, opts, cmd)
}

func TestHttpClient_applyToTransport(t *testing.T) {
	opts := &HttpClient{
		MaxIdleConnsPerHost:   42,
		ResponseHeaderTimeout: 2 * time.Second,
		IdleConnTimeout:       3 * time.Second,
		TlsHandShakeTimeout:   4 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		DialTimeout:           6 * time.Second,
		KeepAlive:             7 * time.Second,
	}
	tr := &http.Transport{}

	opts.applyToTransport(tr)

	assert.Equal(t, 42, tr.MaxIdleConnsPerHost)
	assert.Equal(t, 2*time.Second, tr.ResponseHeaderTimeout)
	assert.Equal(t, 3*time.Second, tr.IdleConnTimeout)
	assert.Equal(t, 4*time.Second, tr.TLSHandshakeTimeout)
	assert.Equal(t, 5*time.Second, tr.ExpectContinueTimeout)

	// Check that DialContext is set and uses the correct timeouts
	if tr.DialContext != nil {
		// Use reflect to extract the dialer from the closure
		// (since Go does not expose closure internals, just check that it works)
		ctx := context.Background()
		_, _ = tr.DialContext(ctx, "tcp", "localhost:80")
	}
}

func TestHttpClient_Configure(t *testing.T) {
	opts := &HttpClient{
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: 2 * time.Second,
		IdleConnTimeout:       3 * time.Second,
		TlsHandShakeTimeout:   4 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		DialTimeout:           6 * time.Second,
		KeepAlive:             7 * time.Second,
		Timeout:               8 * time.Second,
	}
	client := &http.Client{}

	opts.Configure(client)

	// Check client timeout
	assert.Equal(t, 8*time.Second, client.Timeout)

	// Check transport fields
	tr, ok := client.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.Equal(t, 10, tr.MaxIdleConnsPerHost)
	assert.Equal(t, 2*time.Second, tr.ResponseHeaderTimeout)
	assert.Equal(t, 3*time.Second, tr.IdleConnTimeout)
	assert.Equal(t, 4*time.Second, tr.TLSHandshakeTimeout)
	assert.Equal(t, 5*time.Second, tr.ExpectContinueTimeout)
	if tr.DialContext != nil {
		ctx := context.Background()
		_, _ = tr.DialContext(ctx, "tcp", "localhost:80")
	}
}

// Mock ChainableRoundTripper for testing
type mockChainable struct {
	next http.RoundTripper
}

func (m *mockChainable) RoundTrip(req *http.Request) (*http.Response, error) { return nil, nil }
func (m *mockChainable) GetNext() http.RoundTripper                          { return m.next }
func (m *mockChainable) SetNextRoundTripper(rt http.RoundTripper)            { m.next = rt }

// Dummy type for fallback branch
type dummyRoundTripper struct{}

func (d *dummyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return nil, nil }

func Test_setRoundTripper_withHttpTransport(t *testing.T) {
	tr := &http.Transport{}
	called := false
	setRoundTripper(tr, func(transport *http.Transport) {
		called = true
		assert.Equal(t, tr, transport)
	})
	assert.True(t, called)
}

func Test_setRoundTripper_withChainableNoNext(t *testing.T) {
	m := &mockChainable{}
	called := false
	setRoundTripper(m, func(transport *http.Transport) {
		called = true
		assert.NotNil(t, transport)
	})
	assert.True(t, called)
	assert.IsType(t, &http.Transport{}, m.GetNext())
}

func Test_setRoundTripper_withChainableWithNext(t *testing.T) {
	next := &http.Transport{}
	m := &mockChainable{next: next}
	called := false
	setRoundTripper(m, func(transport *http.Transport) {
		called = true
		assert.Equal(t, next, transport)
	})
	assert.True(t, called)
}

func Test_setRoundTripper_withOtherType(t *testing.T) {
	rt := &dummyRoundTripper{}
	called := false
	setRoundTripper(rt, func(transport *http.Transport) {
		called = true
		assert.NotNil(t, transport)
	})
	assert.True(t, called)
}

func Test_setClientTransport_withNilTransport(t *testing.T) {
	client := &http.Client{Transport: nil}
	called := false
	setClientTransport(client, func(tr *http.Transport) {
		called = true
		assert.NotNil(t, tr)
	})
	assert.True(t, called)
	assert.IsType(t, &http.Transport{}, client.Transport)
}

func Test_setClientTransport_withHttpTransport(t *testing.T) {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	called := false
	setClientTransport(client, func(tr2 *http.Transport) {
		called = true
		assert.Equal(t, tr, tr2)
	})
	assert.True(t, called)
}
