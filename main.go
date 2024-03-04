package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"
)

type IfconfigInfo struct {
	IP         string `json:"ip"`
	IPDecimal  int    `json:"ip_decimal"`
	Country    string `json:"country"`
	CountryIso string `json:"country_iso"`
	CountryEu  bool   `json:"country_eu"`
	Asn        string `json:"asn"`
	AsnOrg     string `json:"asn_org"`
	Hostname   string `json:"hostname"`
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	hosts := []string{"hoto.moe", "jp1.hotomoe.net", "us1.hotomoe.net", "jp1-cf.hotomoe.net", "us1-cf.hotomoe.net"}

	info, err := getIfconfigInfo("")
	if err == nil {
		println("Client IP:", info.IP)
		println("Client Country:", info.Country)
		println("Client ASN:", info.Asn)
	} else {
		panic("Failed to get client IP")
	}

	println()

	for _, host := range hosts {
		println("Host:", host)

		ip := getIpAddress(host)
		println("IP:", ip)

		info, err := getIfconfigInfo(ip)
		if err == nil {
			println("ASN:", info.Asn)
		}
		latency := getLatency(host)
		println("Latency:", latency, "ms")

		println()
	}

	println("Press Enter to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getLatency(host string) int64 {
	currentTime := time.Now().UnixMilli()

	resp, err := http.Get("https://" + host)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	latency := time.Now().UnixMilli() - currentTime

	return latency
}

func getIpAddress(dnsName string) string {
	ips, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip4", dnsName)
	if err != nil {
		return "unknown"
	}
	return ips[0].String()
}

func getIfconfigInfo(ip string) (IfconfigInfo, error) {
	var info IfconfigInfo

	resp, err := http.Get("https://ifconfig.co/json?ip=" + ip)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&info)

	return info, nil
}
