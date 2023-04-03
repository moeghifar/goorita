package core

import (
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"
)

type transport struct {
	customTransport *http.Transport
}

func SetServer(customTransport *http.Transport) *TunnelConfig {
	return &TunnelConfig{
		CustomTransport: customTransport,
	}
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
	targetConn, err := t.customTransport.DialContext(req.Context(), "tcp", req.Host)
	if err != nil {
		log.Error().Msgf("failed on DialContext to %v", req.Host)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Panic().Msg("http server doesn't support hijacking connection")
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Panic().Msg("http hijacking failed")
	}

	go connectTunnel(targetConn, clientConn)
	go connectTunnel(clientConn, targetConn)
	log.Info().Msgf("tunnel established with CONNECT request to %v from %v", req.Host, req.RemoteAddr)
}

func connectTunnel(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
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
		DialContext:           (&net.Dialer{}).DialContext, // this is a must to define default net.Dialer.dialContext
	}

	if to != nil {
		if to.CustomNetDialer != nil {
			// overwrite DialContext on behalf custom nic binding
			transporter.DialContext = (to.CustomNetDialer).DialContext
		}
	}

	return transporter
}
