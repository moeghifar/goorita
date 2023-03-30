package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("+++ Go Multiclientproxy +++")

	clients := []*http.Client{}

	nics, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	ports := []int{33001, 33002}
	found := 0

	for _, nic := range nics {
		// fmt.Printf("%+v\n", nic)

		if strings.HasPrefix(nic.Name, "enx") {
			clients = append(clients, &http.Client{
				Transport: newTransport(nic),
			})

			go newServer(ports[found], nic)
			found++
		}
	}

	if len(clients) == 0 {
		fmt.Println("no clients found")
		os.Exit(0)
	}

	fmt.Println("Connected to", len(clients), "clients")

	netCheck(clients)

	http.ListenAndServe(":8800", nil)
}

func newServer(port int, nic net.Interface) {
	server := http.NewServeMux()
	ct := customTransport{Nic: nic}

	server.HandleFunc("/", ct.TransparentHttpProxy)

	fmt.Println("serving", port, "with nic", nic.Name)

	http.ListenAndServe(fmt.Sprintf(":%d", port), server)
}

type customTransport struct {
	Nic net.Interface
}

func (ct customTransport) TransparentHttpProxy(w http.ResponseWriter, r *http.Request) {
	director := func(target *http.Request) {
		target.URL.Scheme = r.URL.Scheme
		target.URL.Path = r.URL.Path
		target.Header.Set("Pass-Via-Go-Proxy", "1")
	}
	proxy := &httputil.ReverseProxy{Director: director, Transport: newTransport(ct.Nic)}
	proxy.ServeHTTP(w, r)
}

func newTransport(nic net.Interface) *http.Transport {
	var dnsResolverIP = "8.8.8.8:53" // Google DNS resolver.
	var dnsResolverProto = "udp"     // Protocol to use for the DNS resolver
	var dnsResolverTimeoutMs = 5000  // Timeout (ms) for the DNS

	return &http.Transport{
		DialContext: (&net.Dialer{
			Resolver: &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
					}
					return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
				},
			},
			Control: setNewDialControl(nic),
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func setNewDialControl(nic net.Interface) func(network, address string, conn syscall.RawConn) error {
	// fmt.Println("no custom config use nic.name as network binding", nic.Name)

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

func netCheck(clients []*http.Client) {
	type response struct {
		Status      string  `json:"status"`
		Country     string  `json:"country"`
		CountryCode string  `json:"countryCode"`
		Region      string  `json:"region"`
		RegionName  string  `json:"regionName"`
		City        string  `json:"city"`
		Zip         string  `json:"zip"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		Timezone    string  `json:"timezone"`
		Isp         string  `json:"isp"`
		Org         string  `json:"org"`
		As          string  `json:"as"`
		Query       string  `json:"query"`
	}

	checkIPReq, err := http.NewRequest(http.MethodGet, "http://ip-api.com/json", nil)
	if err != nil {
		panic(err)
	}

	for _, client := range clients {
		resp, err := client.Do(checkIPReq)
		if err != nil {
			fmt.Println("failed to checkIPReq :", err)
		}

		defer resp.Body.Close()

		result, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read result :", err)
		}

		var responseParser response

		err = json.Unmarshal(result, &responseParser)
		if err != nil {
			fmt.Println("failed to read result :", err)
		}

		fmt.Println("Got IP Address:", responseParser.Query, "from provider:", responseParser.Isp)
	}
}
