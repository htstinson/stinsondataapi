package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Get the direct TCP/IP connection address (Layer 3)
func getTCPAddr(r *http.Request) string {
	fmt.Printf("[%v] getTCPAddr.\n", time.Now().Format(time.RFC3339))
	// RemoteAddr contains the actual TCP connection address (IP:port)
	// This is the most reliable source of the client's direct IP
	// but will be the proxy's IP if the client is behind a proxy
	addr := r.RemoteAddr

	// RemoteAddr includes both IP and port (e.g., 192.168.1.1:12345)
	// Extract just the IP part
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		// If there's an error splitting, just return the whole thing
		fmt.Printf("[%v] getTCPAddr err. %v\n", time.Now().Format(time.RFC3339), err.Error())
		return addr
	}

	return ip
}

// Log middleware that captures the TCP address
func IpLoggingMiddleware(next http.Handler) http.Handler {
	//fmt.Printf("[%v] ipLoggingMiddleware.\n", time.Now().Format(time.RFC3339))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the TCP address
		//ipAddr := getTCPAddr(r)

		// Log the connection information
		fmt.Printf("[%v] [IpLoggingMiddleware] Method: %s, Path: %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		fmt.Printf("[%v] [IpLoggingMiddleware] X-Forwarded-For: %s\n", time.Now().Format(time.RFC3339), r.Header.Get("X-Forwarded-For"))

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// Handler function that displays the IP information
func ipInfoHandler(w http.ResponseWriter, r *http.Request) {
	ipAddr := getTCPAddr(r)

	// Is this a private IP address?
	ip := net.ParseIP(ipAddr)
	isPrivate := isPrivateIP(ip)

	fmt.Fprintf(w, "Connection Information\n\n")
	fmt.Fprintf(w, "Your TCP/IP address: %s\n", ipAddr)
	fmt.Fprintf(w, "Is private address: %t\n", isPrivate)

	// For educational purposes, also show what the headers claim
	// (but we're not using these for our actual IP detection)
	fmt.Fprintf(w, "\nHTTP Headers (NOT TRUSTED):\n")
	fmt.Fprintf(w, "X-Forwarded-For: %s\n", r.Header.Get("X-Forwarded-For"))
	fmt.Fprintf(w, "X-Real-IP: %s\n", r.Header.Get("X-Real-IP"))
}

// Check if an IP is a private address
func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	// Private IPv4 ranges
	privateRanges := []struct {
		start net.IP
		end   net.IP
	}{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
		{net.ParseIP("127.0.0.0"), net.ParseIP("127.255.255.255")},
	}

	for _, r := range privateRanges {
		if bytes4(ip) >= bytes4(r.start) && bytes4(ip) <= bytes4(r.end) {
			return true
		}
	}

	// Check for IPv6 private addresses
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	return false
}

// Helper function to convert IPv4 to uint32 for range comparison
func bytes4(ip net.IP) uint32 {
	if len(ip) == 16 {
		// Convert IPv4-mapped IPv6 to IPv4
		ip = ip[12:16]
	}
	if len(ip) != 4 {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
