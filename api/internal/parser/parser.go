package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// ExtractUniqueIPsFromHandshakeErrors parses a log file and returns a slice of unique IP addresses
// found in log lines containing "handshake error"
func ExtractUniqueIPsFromHandshakeErrors(logFilePath string) ([]string, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Regular expression to match IPv4 addresses
	ipv4Regex := regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)

	// Regular expression to match IPv6 addresses
	ipv6Regex := regexp.MustCompile(`\b(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,7}:\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,5}(?::[0-9a-fA-F]{1,4}){1,2}\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,4}(?::[0-9a-fA-F]{1,4}){1,3}\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,3}(?::[0-9a-fA-F]{1,4}){1,4}\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,2}(?::[0-9a-fA-F]{1,4}){1,5}\b|` +
		`\b[0-9a-fA-F]{1,4}:(?::[0-9a-fA-F]{1,4}){1,6}\b|` +
		`\b:(?::[0-9a-fA-F]{1,4}){1,7}\b|` +
		`\b::(?:[fF]{4}:)?(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)

	// Set to track unique IP addresses
	uniqueIPs := make(map[string]struct{})

	scanner := bufio.NewScanner(file)

	// Scan the log file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Only process lines that contain "handshake error"
		if strings.Contains(line, "handshake error") {
			// Find all IPv4 addresses in the line
			ipv4Matches := ipv4Regex.FindAllString(line, -1)
			for _, ip := range ipv4Matches {
				uniqueIPs[ip] = struct{}{}
			}

			// Find all IPv6 addresses in the line
			ipv6Matches := ipv6Regex.FindAllString(line, -1)
			for _, ip := range ipv6Matches {
				uniqueIPs[ip] = struct{}{}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Convert map keys to slice
	result := make([]string, 0, len(uniqueIPs))
	for ip := range uniqueIPs {
		result = append(result, ip)
	}

	return result, nil
}
