package resolver

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"
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
	req, err := newDohPostRequest(ctx, bytes.NewReader(buff))
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

	n, err := resp.Body.Read(buff[:cap(buff)])
	if err != nil && err != io.EOF {
		return err
	}

	conn.WriteTo(buff[:n], addr)
	return nil
}
