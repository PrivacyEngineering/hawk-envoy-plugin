package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var dnsCache = make(map[string]string)

func findHostByIp(ip string) (string, error) {
	if host, has := dnsCache[ip]; has {
		return host, nil
	}
	addr, err := net.LookupAddr(ip)
	if err != nil {
		log.Printf("Error finding host by ip (ip=%v),: %s", ip, err)
		return ip, nil
	}
	log.Printf("Hosts found by IP (%s): %v", ip, addr)
	dnsCache[ip] = addr[0]
	return addr[0], nil
}

func fix(data []byte, usage map[string]interface{}) map[string]interface{} {
	fixTimestamp(data, usage)
	fixEndpoint(data, usage)
	fixInitiator(data, usage)
	return usage
}

func fixTimestamp(data []byte, usage map[string]interface{}) {
	metadata, ok := usage["metadata"].(map[string]interface{})
	if !ok {
		log.Printf("Unable to transform metadata to map")
		toDlq(data)
		return
	}

	unixTS, ok := metadata["timestamp"].(string)
	if !ok {
		log.Printf("Unable to transform metadata.timestamp to string. Metadata: %v", metadata)
		toDlq(data)
		return
	}

	i, err := strconv.ParseFloat(unixTS, 64)
	if err != nil {
		log.Printf("Unable to transform unix timestamp: %s", err)
		toDlq(data)
		return
	}
	i = i / 1e9

	tm := time.Unix(int64(i), 0).Format(time.RFC3339)
	metadata["timestamp"] = tm
}

func fixEndpoint(data []byte, usage map[string]interface{}) {
	endpoint, ok := usage["endpoint"].(map[string]interface{})
	if !ok {
		log.Printf("Unable to transform endpoint to map")
		toDlq(data)
		return
	}

	ip, ok := endpoint["host"].(string)
	if !ok {
		log.Printf("Unable to transform endpoint.host to string. Metadata: %v", endpoint)
		toDlq(data)
		return
	}

	host, err := findHostByIp(ip)
	if err != nil {
		log.Printf("Unable to transform find host by ip. Metadata: %v", endpoint)
		toDlq(data)
		return
	}

	endpoint["host"] = host

	method, ok := endpoint["method"].(string)
	if !ok {
		log.Printf("Unable to transform endpoint.method to string. Metadata: %v", endpoint)
		toDlq(data)
		return
	}

	path, ok := endpoint["path"].(string)
	if !ok {
		log.Printf("Unable to transform endpoint.path to string. Metadata: %v", endpoint)
		toDlq(data)
		return
	}
	parts := strings.Split(path, "?")
	endpoint["path"] = parts[0]

	protocol, ok := endpoint["protocol"].(string)
	if !ok {
		log.Printf("Unable to transform endpoint.protocol to string. Metadata: %v", endpoint)
		toDlq(data)
		return
	}

	if strings.Contains(strings.ToLower(protocol), "http") {
		protocol = "http"
	}
	// http:<method>:<host>:<path>
	id := protocol + ":" + method + ":" + host + ":" + path
	endpoint["id"] = id
}

func fixInitiator(data []byte, usage map[string]interface{}) {
	initiator, ok := usage["initiator"].(map[string]interface{})
	if !ok {
		log.Printf("Unable to transform initiator to map")
		toDlq(data)
		return
	}

	host, ok := initiator["host"].(string)
	if !ok {
		log.Printf("Unable to transform initiator.host to string. Metadata: %v", initiator)
		toDlq(data)
		return
	}

	parts := strings.Split(host, ":")
	initiator["host"] = parts[0]
}
