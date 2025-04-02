package waf

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func Block(ipSetName string, addIP string, removeIP string, region string) {
	// Convert string to proper Scope type
	var scope types.Scope = types.ScopeRegional

	// Load AWS configuration - will use EC2 instance role credentials automatically
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// Create WAF client
	client := wafv2.NewFromConfig(cfg)
	var mylimit int32 = 100

	// List IP Sets to get the ID
	listInput := &wafv2.ListIPSetsInput{
		Scope: scope,
		Limit: &mylimit,
	}

	listResult, err := client.ListIPSets(context.TODO(), listInput)
	if err != nil {
		log.Fatalf("failed to list IP sets: %v", err)
	}

	// Find the target IP set ID
	var ipSetId string
	var ipSetARN string
	var lockToken string
	for _, ipSet := range listResult.IPSets {
		if *ipSet.Name == ipSetName {
			ipSetId = *ipSet.Id
			ipSetARN = *ipSet.ARN
			break
		}
	}

	// Get IP set details
	getInput := &wafv2.GetIPSetInput{
		Id:    &ipSetId,
		Name:  &ipSetName,
		Scope: scope,
	}

	getResult, err := client.GetIPSet(context.TODO(), getInput)
	if err != nil {
		log.Fatalf("failed to get IP set details: %v", err)
	}

	// Save the current lock token for updates
	lockToken = *getResult.LockToken

	// Create a copy of the current addresses
	addresses := make([]string, len(getResult.IPSet.Addresses))
	copy(addresses, getResult.IPSet.Addresses)

	// Check if we need to update the IP set
	needsUpdate := false

	// Add IP address if specified
	if addIP != "" {
		// Check if the IP is already in the set
		exists := false
		for _, addr := range addresses {
			if addr == addIP {
				exists = true
				break
			}
		}

		if !exists {
			addresses = append(addresses, addIP)
			fmt.Printf("Adding IP: %s\n", addIP)
			needsUpdate = true
		} else {
			fmt.Printf("IP %s already exists in the set\n", addIP)
		}
	}

	// Remove IP address if specified
	if removeIP != "" {
		for i, addr := range addresses {
			if addr == removeIP {
				// Remove the IP by replacing it with the last element and truncating
				addresses[i] = addresses[len(addresses)-1]
				addresses = addresses[:len(addresses)-1]
				fmt.Printf("Removing IP: %s\n", removeIP)
				needsUpdate = true
				break
			}
		}
	}

	// Update the IP set if needed
	if needsUpdate {
		updateInput := &wafv2.UpdateIPSetInput{
			Id:        &ipSetId,
			Name:      &ipSetName,
			Scope:     scope,
			Addresses: addresses,
			LockToken: &lockToken,
		}

		_, err = client.UpdateIPSet(context.TODO(), updateInput)
		if err != nil {
			log.Fatalf("failed to update IP set: %v", err)
		}
		fmt.Println("IP set updated successfully")

		// Refresh IP set details after update
		getResult, err = client.GetIPSet(context.TODO(), getInput)
		if err != nil {
			log.Fatalf("failed to get updated IP set details: %v", err)
		}
	}

	// Output results
	fmt.Printf("IP Set: %s\n", ipSetName)
	fmt.Printf("ID: %s\n", ipSetId)
	fmt.Printf("ARN: %s\n", ipSetARN)
	fmt.Printf("Description: %s\n", *getResult.IPSet.Description)
	fmt.Println("\nIP Addresses:")

	for _, address := range getResult.IPSet.Addresses {
		fmt.Println(address)
	}

	// Export to file if needed
	outputFile := ipSetName + "-ips.txt"
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer f.Close()

	for _, address := range getResult.IPSet.Addresses {
		f.WriteString(address + "\n")
	}

	fmt.Printf("\nIP addresses exported to %s\n", outputFile)
}
