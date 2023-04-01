package main

// import (
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"io"
// 	"net"
// 	"net/http"
// 	"strconv"
// )

// type cfg struct {
// 	httpPort int
// }

// func main() {
// 	fmt.Println("+++ Halaproxy starting +++")

// 	var config cfg
// 	var httpPort string
// 	var bindNic string

// 	flag.StringVar(&httpPort, "http-port", "", "cli command")
// 	flag.StringVar(&bindNic, "bind-nic", "", "cli command")
// 	flag.Parse()

// 	validHttpPort, err := strconv.Atoi(httpPort)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if validHttpPort > 0 {
// 		config.httpPort = validHttpPort
// 	}

// 	fmt.Println(bindNic)

// 	// clients := []*http.Client{}

// 	nics, err := net.Interfaces()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(nics)

// 	// ports := []int{33001, 33002}
// 	// found := 0

// 	// fmt.Println("reading device nic")
// 	// for _, nic := range nics {
// 	// 	fmt.Printf("%+v\n", nic)

// 	// 	if strings.HasPrefix(nic.Name, "enx") {
// 	// 		clients = append(clients, &http.Client{
// 	// 			Transport: newTransport(nic),
// 	// 		})

// 	// 		go newServer(ports[found], nic)
// 	// 		found++
// 	// 	}
// 	// }

// 	// if len(clients) == 0 {
// 	// 	fmt.Println("no clients found")
// 	// 	os.Exit(0)
// 	// }

// 	// fmt.Println("Connected to", len(clients), "clients")

// 	// netCheck(clients)

// 	// http.ListenAndServe(":8800", nil)

// 	// readSignal := make(chan os.Signal, 1)

// 	// signal.Notify(
// 	// 	readSignal,
// 	// 	syscall.SIGTERM,
// 	// 	syscall.SIGINT,
// 	// )

// 	// defaultServer := newDefaultServer(config.httpPort)

// 	// go func() {}()

// 	// <-readSignal

// }

// func intToPort(port int) string {
// 	return fmt.Sprintf(":%d", port)
// }

// // func serverSetup(port int) *http.ServerMux {
// // 	server := http.NewServeMux()

// // 	ct := customTransport{}

// // 	server.HandleFunc("/", ct.TransparentHttpProxy)

// // 	return server
// // }

// // func newServer(port int, nic net.Interface) {
// // 	server := http.NewServeMux()
// // 	ct := customTransport{Nic: nic}

// // 	server.HandleFunc("/", ct.TransparentHttpProxy)

// // 	fmt.Println("serving", port, "with nic", nic.Name)

// // 	http.ListenAndServe(fmt.Sprintf(":%d", port), server)
// // }

// // type customTransport struct {
// // 	Nic net.Interface
// // }

// // func (ct customTransport) TransparentHttpProxy(w http.ResponseWriter, r *http.Request) {
// // 	// log.Println("r url", r.URL)
// // 	// director := func(target *http.Request) {
// // 	// 	target
// // 	// }
// // 	// proxy := &httputil.ReverseProxy{
// // 	// 	Director: director,
// // 	// 	Transport: &http.Transport{
// // 	// 		DialTLS: dialTLS,
// // 	// 	},
// // 	// 	// Transport: newTransport(ct.Nic),
// // 	// }
// // 	// proxy.ServeHTTP(w, r)

// // 	proxy := httputil.NewSingleHostReverseProxy(r.URL)

// // 	proxy.Transport = &http.Transport{
// // 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// // 		Dial:            dialTLS,
// // 	}

// // 	proxy.Director = func(req *http.Request) {
// // 		req.URL = r.URL
// // 		req.Host = r.Host
// // 	}

// // 	proxy.ServeHTTP(w, r)
// // }

// // func dialTLS(network, addr string) (net.Conn, error) {
// // 	log.Println(network, addr)

// // 	conn, err := net.Dial(network, addr)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	host, _, err := net.SplitHostPort(addr)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	cfg := &tls.Config{ServerName: host}

// // 	tlsConn := tls.Client(conn, cfg)
// // 	if err := tlsConn.Handshake(); err != nil {
// // 		conn.Close()
// // 		return nil, err
// // 	}

// // 	cs := tlsConn.ConnectionState()
// // 	cert := cs.PeerCertificates[0]

// // 	// Verify here
// // 	cert.VerifyHostname(host)
// // 	log.Println(cert.Subject)

// // 	return tlsConn, nil
// // }

// // func newTransport(nic net.Interface) *http.Transport {
// // 	// var dnsResolverIP = "8.8.8.8:53" // Google DNS resolver.
// // 	// var dnsResolverProto = "udp"     // Protocol to use for the DNS resolver
// // 	// var dnsResolverTimeoutMs = 5000  // Timeout (ms) for the DNS
// // 	// var certPath = ""
// // 	// var caCertPath = ""
// // 	// var keyPath = ""
// // 	// var caCertPool *x509.CertPool
// // 	// // Create a CA certificate pool and add cert.pem to it

// // 	// if caCertPath != "" {
// // 	// 	caCert, err := ioutil.ReadFile(caCertPath)
// // 	// 	if err != nil {
// // 	// 		log.Fatalf("[ERROR] [proxy,nomad] [message: failed to read Nomad CA Cert]")
// // 	// 	}
// // 	// 	caCertPool = x509.NewCertPool()
// // 	// 	caCertPool.AppendCertsFromPEM(caCert)
// // 	// }

// // 	return &http.Transport{
// // 		// TLSClientConfig: &tls.Config{
// // 		// 	InsecureSkipVerify: true,
// // 		// },
// // 		// DialTLS: newDialTls,
// // 		DialContext: (&net.Dialer{
// // 			Control: setNewDialControl(nic),
// // 		}).DialContext,
// // 		MaxIdleConns:          100,
// // 		IdleConnTimeout:       90 * time.Second,
// // 		TLSHandshakeTimeout:   10 * time.Second,
// // 		ExpectContinueTimeout: 1 * time.Second,
// // 	}
// // }

// // func setNewDialControl(nic net.Interface) func(network, address string, conn syscall.RawConn) error {
// // 	// fmt.Println("no custom config use nic.name as network binding", nic.Name)

// // 	return func(network, address string, conn syscall.RawConn) error {
// // 		var operr error
// // 		if err := conn.Control(func(fd uintptr) {
// // 			operr = unix.BindToDevice(int(fd), nic.Name)
// // 		}); err != nil {
// // 			return err
// // 		}
// // 		return operr
// // 	}
// // }

// func netCheck(clients []*http.Client) {
// 	type response struct {
// 		Status      string  `json:"status"`
// 		Country     string  `json:"country"`
// 		CountryCode string  `json:"countryCode"`
// 		Region      string  `json:"region"`
// 		RegionName  string  `json:"regionName"`
// 		City        string  `json:"city"`
// 		Zip         string  `json:"zip"`
// 		Lat         float64 `json:"lat"`
// 		Lon         float64 `json:"lon"`
// 		Timezone    string  `json:"timezone"`
// 		Isp         string  `json:"isp"`
// 		Org         string  `json:"org"`
// 		As          string  `json:"as"`
// 		Query       string  `json:"query"`
// 	}

// 	checkIPReq, err := http.NewRequest(http.MethodGet, "http://ip-api.com/json", nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, client := range clients {
// 		resp, err := client.Do(checkIPReq)
// 		if err != nil {
// 			fmt.Println("failed to checkIPReq :", err)
// 		}

// 		defer resp.Body.Close()

// 		result, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println("failed to read result :", err)
// 		}

// 		var responseParser response

// 		err = json.Unmarshal(result, &responseParser)
// 		if err != nil {
// 			fmt.Println("failed to read result :", err)
// 		}

// 		fmt.Println("Got IP Address:", responseParser.Query, "from provider:", responseParser.Isp)
// 	}
// }

// // // func NewProxy(addr, caCertPath, certPath, keyPath string) http.Handler {
// // // 	remoteURL, _ := url.Parse(addr)

// // // 	proxy := httputil.NewSingleHostReverseProxy(remoteURL)

// // // 	if certPath != "" && keyPath != "" {
// // // 		var caCertPool *x509.CertPool
// // // 		// Create a CA certificate pool and add cert.pem to it
// // // 		if caCertPath != "" {
// // // 			caCert, err := ioutil.ReadFile(caCertPath)
// // // 			if err != nil {
// // // 				log.Fatalf("[ERROR] [proxy,nomad] [message: failed to read Nomad CA Cert]")
// // // 			}
// // // 			caCertPool = x509.NewCertPool()
// // // 			caCertPool.AppendCertsFromPEM(caCert)
// // // 		}

// // // 		// Create an HTTPS client and supply the created CA pool and certificate
// // // 		proxy.Transport = &http.Transport{
// // // 			TLSClientConfig: &tls.Config{
// // // 				RootCAs:    caCertPool,
// // // 				MinVersion: tls.VersionTLS13,
// // // 				MaxVersion: tls.VersionTLS13,
// // // 				GetClientCertificate: func(chi *tls.CertificateRequestInfo) (*tls.Certificate, error) {
// // // 					cert, err := tls.LoadX509KeyPair(certPath, keyPath)
// // // 					if err != nil {
// // // 						return nil, err
// // // 					}

// // // 					return &cert, nil
// // // 				},
// // // 			},
// // // 		}
// // // 	}

// // // 	return proxy
// // // }
