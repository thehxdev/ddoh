package resolver

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/thehxdev/ddoh/config"
)

type Resolver struct {
	*http.Client
}

func Init() *Resolver {
	return &Resolver{
		Client: initHttpClient(),
	}
}

func (r *Resolver) Resolve(conn *net.UDPConn, addr net.Addr, buff []byte) {
	reqPacket, err := bytesToDNSPacket(buff)
	if err != nil {
		log.Println(err)
		return
	}
	qName := reqPacket.Questions[0].Name

	log.Printf("new query -> %s\n", string(qName))
	bodyBuff := bytesToBuffer(buff)
	req, err := newHttpRequest(bodyBuff)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if stat := resp.StatusCode; stat != http.StatusOK {
		log.Printf("got %d status code\n", stat)
		return
	}

	bodyBuff.Reset()
	_, err = bodyBuff.ReadFrom(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	conn.WriteTo(bodyBuff.Bytes(), addr)
}

func bytesToDNSPacket(body []byte) (*layers.DNS, error) {
	dns := &layers.DNS{}
	if err := dns.DecodeFromBytes(body, nil); err != nil {
		return nil, err
	}
	return dns, nil
}

func dnsPacketToBytes(dns *layers.DNS) []byte {
	buff := gopacket.NewSerializeBuffer()
	err := dns.SerializeTo(buff, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: false})
	if err != nil {
		return nil
	}
	return buff.Bytes()
}

func bytesToBuffer(b []byte) *bytes.Buffer {
	buff := &bytes.Buffer{}
	if b != nil {
		buff.Write(b)
	}
	return buff
}

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
