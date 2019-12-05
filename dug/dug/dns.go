package dug

import (
	"fmt"
	"golang.org/x/net/context"
	"net"
	"time"
)

type AddrResponse struct {
	addrs []string
	err   error
}

func TimeoutLookupHost(host string, timeout float64) (addrs []string, err error) {
	timeout_ns := timeout * 1e9
	timeout_d := time.Duration(int(timeout_ns))
	ctx, cancel := context.WithTimeout(context.Background(), timeout_d)
	defer cancel()

	ch := make(chan AddrResponse, 1)
	defer close(ch)

	go func() {
		select {
		default:
			addrs, err := net.LookupHost(host)
			ch <- AddrResponse{addrs: addrs, err: err}
		case <-ctx.Done():
			fmt.Println("Exiting gofunc")
			return
		}
	}()

	select {
	case out := <-ch:
		return out.addrs, out.err
	case <-ctx.Done():
		return nil, fmt.Errorf("Lookup '%s' timed out after %.3fs", host)
	}
}
