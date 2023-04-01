package core

import (
	"net"
	"net/http"
	"net/http/httputil"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type transport struct {
	customTransport *http.Transport
}

func SetServer(customTransport *http.Transport) *http.ServeMux {
	server := http.NewServeMux()

	ct := &transport{
		customTransport: customTransport,
	}

	server.HandleFunc("/", ct.proxyHandler)

	return server
}

func (t *transport) proxyHandler(w http.ResponseWriter, req *http.Request) {
	director := func(target *http.Request) {
		target.URL = req.URL
		target.Host = req.Host
	}

	proxy := &httputil.ReverseProxy{
		Director:  director,
		Transport: t.customTransport,
	}

	proxy.ServeHTTP(w, req)
}

// func newTransportWithNic(nic net.Interface) *http.Transport {
// 	return &http.Transport{
// 		DialContext: (&net.Dialer{
// 			Control: func(network, address string, conn syscall.RawConn) error {
// 				var operr error
// 				if err := conn.Control(func(fd uintptr) {
// 					operr = unix.BindToDevice(int(fd), nic.Name)
// 				}); err != nil {
// 					return err
// 				}
// 				return operr
// 			},
// 		}).DialContext,
// 		MaxIdleConns:          100,
// 		IdleConnTimeout:       90 * time.Second,
// 		TLSHandshakeTimeout:   10 * time.Second,
// 		ExpectContinueTimeout: 1 * time.Second,
// 	}
// }

func NewDialControlWithBindNic(nic net.Interface) func(network, address string, conn syscall.RawConn) error {
	return func(network, address string, conn syscall.RawConn) error {
		var operr error
		if err := conn.Control(func(fd uintptr) {
			operr = unix.BindToDevice(int(fd), nic.Name)
		}); err != nil {
			return err
		}
		return operr
	}
}

func NewNetDialerWithNicBinding(nic net.Interface) *net.Dialer {
	return &net.Dialer{
		Control: func(network, address string, conn syscall.RawConn) error {
			var operr error
			if err := conn.Control(func(fd uintptr) {
				operr = unix.BindToDevice(int(fd), nic.Name)
			}); err != nil {
				return err
			}
			return operr
		},
	}
}

type TransportOptions struct {
	CustomNetDialer *net.Dialer
}

func CreateTransportWithNic(nic net.Interface) *http.Transport {
	return CreateTransport(&TransportOptions{
		CustomNetDialer: NewNetDialerWithNicBinding(nic),
	})
}

func CreateTransport(to *TransportOptions) *http.Transport {
	transporter := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if to != nil {
		if to.CustomNetDialer != nil {
			transporter.DialContext = (to.CustomNetDialer).DialContext
		}
	}

	return transporter
}
