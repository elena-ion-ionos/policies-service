package config

import (
	"net"
	"net/http"
	"time"

	"github.com/ionos-cloud/go-paaskit/pkg/paasclient/transport"

	"github.com/spf13/cobra"
)

type HttpClient struct {
	MaxIdleConnsPerHost   int
	Timeout               time.Duration
	DialTimeout           time.Duration
	KeepAlive             time.Duration
	IdleConnTimeout       time.Duration
	TlsHandShakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	ResponseHeaderTimeout time.Duration
}

// http client and dialer defaults
const (
	defaultMaxIdleConnsPerHost   = 100             // limits the number of connections that can be open at a time for each of the host from your pod
	defaultResponseHeaderTimeout = 1 * time.Second // expect server to reply with a response header within 1 second or the client will retry.
	defaultTimeout               = 4 * time.Second // default request context time out after 4 seconds so the client can retry the default amount of retries which is 3
	defaultDialTimeout           = 4 * time.Second // limits the time spent establishing a TCP connection (if a new one is needed).
	defaultKeepAlive             = 30 * time.Second
	defaultIdleConnTimeout       = 90 * time.Second
	defaultTlsHandShakeTimeout   = 10 * time.Second // limits the time spent performing the TLS handshake.
	defaultExpectContinueTimeout = 1 * time.Second  // limits the time the client will wait between sending the request headers when including an Expect: 100-continue and receiving the go-ahead to send the body
)

func (o *HttpClient) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&o.MaxIdleConnsPerHost, "HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST", defaultMaxIdleConnsPerHost, "")
	cmd.Flags().DurationVar(&o.ResponseHeaderTimeout, "HTTP_CLIENT_RESPONSE_HEADER_TIMEOUT", defaultResponseHeaderTimeout, "")
	// NOTE: Make sure to always set this to 4 * RetryTimeout
	cmd.Flags().DurationVar(&o.Timeout, "HTTP_CLIENT_TIMEOUT", defaultTimeout, "")
	cmd.Flags().DurationVar(&o.DialTimeout, "HTTP_CLIENT_DIAL_TIMEOUT", defaultDialTimeout, "")
	cmd.Flags().DurationVar(&o.KeepAlive, "HTTP_CLIENT_KEEP_ALIVE", defaultKeepAlive, "")
	cmd.Flags().DurationVar(&o.IdleConnTimeout, "HTTP_CLIENT_IDLE_CONN_TIMEOUT", defaultIdleConnTimeout, "")
	cmd.Flags().DurationVar(&o.TlsHandShakeTimeout, "HTTP_CLIENT_TLS_HANDSHAKE_TIMEOUT", defaultTlsHandShakeTimeout, "")
	cmd.Flags().DurationVar(&o.ExpectContinueTimeout, "HTTP_CLIENT_EXPECT_CONTINUE_TIMEOUT", defaultExpectContinueTimeout, "")
}

func (o *HttpClient) applyToTransport(transport *http.Transport) {
	transport.MaxIdleConnsPerHost = o.MaxIdleConnsPerHost
	transport.ResponseHeaderTimeout = o.ResponseHeaderTimeout
	transport.IdleConnTimeout = o.IdleConnTimeout
	transport.TLSHandshakeTimeout = o.TlsHandShakeTimeout
	transport.ExpectContinueTimeout = o.ExpectContinueTimeout
	transport.DialContext = (&net.Dialer{
		Timeout:   o.DialTimeout,
		KeepAlive: o.KeepAlive,
	}).DialContext
}

func (o *HttpClient) Configure(client *http.Client) {
	setClientTransport(client, o.applyToTransport)

	client.Timeout = o.Timeout
}

func setClientTransport(client *http.Client, setTransport func(transport *http.Transport)) {
	if client.Transport == nil {
		client.Transport = &http.Transport{}
	}
	setRoundTripper(client.Transport, setTransport)
}

func setRoundTripper(rt http.RoundTripper, setTransport func(transport *http.Transport)) {
	if val, ok := rt.(transport.ChainableRoundTripper); ok {
		if val.GetNext() == nil {
			t := http.Transport{}
			setTransport(&t)
			val.SetNextRoundTripper(&t)
			return
		}
		setRoundTripper(val.GetNext(), setTransport)
		return
	} else if val, ok := rt.(*http.Transport); ok {
		setTransport(val)
		return
	} else {
		t := &http.Transport{}
		setTransport(t)
		return
	}
}
