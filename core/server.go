package core

import (
	"fmt"
	"io"
	"log"
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

type TunnelConfig struct {
	CustomTransport *http.Transport
}

func (t *TunnelConfig) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ct := &transport{
		customTransport: t.CustomTransport,
	}
	if req.Method == http.MethodConnect {
		ct.proxyConnect(w, req)
	} else {
		ct.proxyHandler(w, req)
	}
}

func (t *transport) proxyConnect(w http.ResponseWriter, req *http.Request) {
	log.Printf("CONNECT requested to %v (from %v)", req.Host, req.RemoteAddr)
	targetConn, err := t.customTransport.DialContext(req.Context(), "tcp", req.Host)
	if err != nil {
		log.Println("failed to dial to target", req.Host)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Fatal("http server doesn't support hijacking connection")
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Fatal("http hijacking failed")
	}

	log.Println("tunnel established")
	go tunnelConn(targetConn, clientConn)
	go tunnelConn(clientConn, targetConn)
}

func tunnelConn(dst io.WriteCloser, src io.ReadCloser) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
}

func (t *transport) proxyHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("=== req ===")
	fmt.Println(req)
	fmt.Println("===========")

	director := func(target *http.Request) {
		target.URL = req.URL
		target.Host = req.Host
		fmt.Println("=== target ===")
		fmt.Println(target)
		fmt.Println("===========")
	}

	proxy := &httputil.ReverseProxy{
		Director:  director,
		Transport: t.customTransport,
	}

	proxy.ServeHTTP(w, req)
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

type transportOptions struct {
	CustomNetDialer *net.Dialer
}

func CreateTransportWithNic(nic net.Interface) *http.Transport {
	return CreateTransport(&transportOptions{
		CustomNetDialer: NewNetDialerWithNicBinding(nic),
	})
}

func CreateTransport(to *transportOptions) *http.Transport {
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
