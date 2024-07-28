package resolver

import (
	"bytes"
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
