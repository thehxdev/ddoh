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

	var ip string
	if len(config.Global.DoHIP) == 0 {
		u, err := url.Parse(config.Global.DoHServer)
		if err != nil {
			log.Fatal(err)
		}

		dohIPs, err := net.LookupHost(u.Hostname())
		if err != nil {
			log.Fatal(err)
		}

		ip = dohIPs[0]
	} else {
		ip = config.Global.DoHIP
	}

	return &http.Client{
		// Request timeout
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// since we send requests to a constant URL, it's better to
				// resolve the host to it's IP addresses and use the IP address
				// directly.
				addr = net.JoinHostPort(ip, "443")
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}
}

func newHttpRequest(body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), "POST", config.Global.DoHServer, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")
	return req, nil
}
