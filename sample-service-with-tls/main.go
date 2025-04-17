package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, mutual TLS world!")
}

func main() {
	// Get cert path from ENV or default
	certPath := "/home/ballerina"

	// Construct full paths
	serverCert := filepath.Join(certPath, "server.crt")
	serverKey := filepath.Join(certPath, "server.key")
	caCertPath := filepath.Join(certPath, "ca.pem")

	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatalf("failed to load server cert/key: %v", err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("failed to read ca.pem: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("failed to add CA cert to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":8443",
		Handler:   http.HandlerFunc(helloHandler),
		TLSConfig: tlsConfig,
	}

	fmt.Println("🚀 Server started at https://localhost:8443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
