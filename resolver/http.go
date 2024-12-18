package resolver

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/thehxdev/ddoh/config"
)

const DNS_MESSAGE_HEADER = "application/dns-message"

func initHttpClient() *http.Client {
	dialer := &net.Dialer{
		// Lookup timeout
		Timeout: time.Second * 10,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", net.JoinHostPort(config.Global.LocalResolver, "53"))
			},
		},
	}
	net.DefaultResolver = dialer.Resolver

	u, err := url.Parse(config.Global.DoHServer)
	if err != nil {
		log.Fatal(err)
	}

	port := "443"
	if p := u.Port(); p != "" {
		port = p
	}

	ip := config.Global.DoHIP
	if len(config.Global.DoHIP) == 0 {
		dohIPs, err := net.LookupHost(u.Hostname())
		if err != nil {
			log.Fatal(err)
		}
		ip = dohIPs[0]
	}

	return &http.Client{
		// Request timeout
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// since we send requests to a constant URL, it's better to
				// resolve the host to it's IP addresses and use the IP address
				// directly.
				addr = net.JoinHostPort(ip, port)
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}
}

func newDohPostRequest(ctx context.Context, body *bytes.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", config.Global.DoHServer, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", DNS_MESSAGE_HEADER)
	req.Header.Set("Accept", DNS_MESSAGE_HEADER)
	return req, nil
}
