package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	piondtls "github.com/pion/dtls/v3"
	coapdtls "github.com/plgd-dev/go-coap/v3/dtls"
	coapmessage "github.com/plgd-dev/go-coap/v3/message"
	coapmessagepool "github.com/plgd-dev/go-coap/v3/message/pool"
)

func sendCoap(host string, config *piondtls.Config, body string, path string) {
	co, err := coapdtls.Dial(host, config)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var resp *coapmessagepool.Message
	var coapErr error
	log.Printf("=> sending message via coap/dtls to %v", host)
	if body == "" {
		resp, coapErr = co.Get(ctx, path)
	} else {
		resp, coapErr = co.Post(ctx, path, coapmessage.TextPlain, bytes.NewReader([]byte(body)))
	}

	if coapErr != nil {
		log.Fatalf("Error sending request: %v", coapErr)
	}
	log.Printf("<= response: %+v", resp)
	responseBody, err := resp.ReadBody()
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	if len(responseBody) > 0 {
		log.Printf("<= body: %+v", string(responseBody))
	}
	log.Printf("<= done")
}

func sendUdp(host string, config *piondtls.Config, body string) {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		panic(err)
	}

	conn, err := piondtls.Dial("udp", addr, config)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	log.Printf("=> sending message via udp/dtls to %v", host)
	_, err = conn.Write([]byte(body))
	if err != nil {
		log.Fatalf("Error writing: %v", err)
	}
	log.Printf("<= done")
}

func sendHttp(host string, config *tls.Config, body string, path string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var req *http.Request
	var err error
	if body == "" {
		req, _ = http.NewRequestWithContext(ctx, "GET", host+path, nil)
	} else {
		req, _ = http.NewRequestWithContext(ctx, "POST", host+path, bytes.NewReader([]byte(body)))
	}

	transport := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: transport}
	log.Printf("=> Sending message via http to %v", host)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("request failed: %e", err)
	}
	log.Printf("<= response: %+v", resp)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	if len(responseBody) > 0 {
		log.Printf("<= body: %+v", string(responseBody))
	}
	log.Printf("<= done")
}

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	body := ""
	if fi.Mode()&os.ModeNamedPipe != 0 {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		body = string(input)
	}

	caCertFile := "../ca.pem"
	_, err = os.Stat(caCertFile)
	if err != nil {
		caCertFile = ""
		log.Printf("[WARNING] unable to load ca.pem; proceeding without checking server certificate")
	}

	certFile := os.Args[1]
	host := "receivers.iot-exchange.io"
	if len(os.Args) > 2 {
		host = os.Args[2]
	}
	path := "/uplink"
	if len(os.Args) > 3 {
		path = os.Args[3]
	}

	certificate, err := tls.LoadX509KeyPair(certFile, certFile)
	if err != nil {
		log.Fatalf("Errors loading keypair: %v", err)
	}

	caCertPool := x509.NewCertPool()
	hasCaCert := false
	if caCertFile != "" {
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			log.Fatalf("Failed to read CA certificate: %v", err)
		}

		caCertPool.AppendCertsFromPEM(caCert)
		hasCaCert = true
	}

	config := &piondtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: piondtls.RequireExtendedMasterSecret,
		InsecureSkipVerify:   !hasCaCert,
		RootCAs:              caCertPool,
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		RootCAs:            caCertPool,
		InsecureSkipVerify: !hasCaCert,
	}

	log.Printf("--------------------------------------------------------------------------------------")
	sendCoap(host+":5684", config, body, path)
	log.Printf("--------------------------------------------------------------------------------------")
	sendUdp(host+":4433", config, body)
	log.Printf("--------------------------------------------------------------------------------------")
	sendHttp("https://"+host, tlsConfig, body, path)
	log.Printf("--------------------------------------------------------------------------------------")
}
