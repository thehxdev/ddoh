package resolver

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Resolver struct {
	*http.Client
}

func Init() *Resolver {
	return &Resolver{
		Client: initHttpClient(),
	}
}

func (r *Resolver) Resolve(ctx context.Context, conn *net.UDPConn, addr net.Addr, buff []byte) error {
	reqPacket, err := bytesToDNSPacket(buff)
	if err != nil {
		return err
	}
	qName := reqPacket.Questions[0].Name

	log.Printf("new query -> %s\n", string(qName))
	req, err := newHttpRequest(ctx, bytes.NewReader(buff))
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if stat := resp.StatusCode; stat != http.StatusOK {
		log.Printf("got %d status code\n", stat)
		return err
	}

	n, err := resp.Body.Read(buff)
	if err != nil && err != io.EOF {
		return err
	}

	conn.WriteTo(buff[:n], addr)
	return nil
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
