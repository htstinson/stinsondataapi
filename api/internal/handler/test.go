package handler

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	searcher "github.com/htstinson/business_searcher"
	"github.com/joho/godotenv"
)

func (h *Handler) test() {
	fmt.Println("test")
	configFile := flag.String("config", "", "Path to a single YAML configuration file (optional)")
	configDir := flag.String("config-dir", "./search_definitions", "Path to directory containing YAML configuration files")
	outputDir := flag.String("output-dir", "./search_results", "Path to directory for output files")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found, using system environment variables")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("ERROR: GOOGLE_API_KEY environment variable is not set")
	}

	// Ensure output directory exists
	if err := searcher.EnsureDirectory(*outputDir); err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// Determine which config files to process
	var configFiles []string
	var err error

	if *configFile != "" {
		// Single file mode (backward compatibility)
		// If the config file doesn't have a path separator, check in config-dir
		configPath := *configFile
		if !strings.Contains(configPath, string(filepath.Separator)) && !filepath.IsAbs(configPath) {
			// Just a filename, look in config-dir
			configPath = filepath.Join(*configDir, configPath)
			fmt.Printf("Looking for config file in: %s\n", configPath)
		}
		configFiles = []string{configPath}
		fmt.Printf("Processing single config file: %s\n", configPath)
	} else {
		// Directory mode (new default)
		fmt.Printf("Looking for config files in: %s\n", *configDir)
		configFiles, err = searcher.LoadConfigsFromDirectory(*configDir)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		fmt.Printf("Found %d config file(s)\n", len(configFiles))
		for i, cf := range configFiles {
			fmt.Printf("  [%d] %s\n", i+1, cf)
		}
	}

	// Process each config file
	totalSearches := 0
	totalSuccessful := 0
	totalFailed := 0

	for _, cfgFile := range configFiles {
		fmt.Printf("\n--- Processing: %s ---\n", cfgFile)

		config, err := searcher.LoadConfig(cfgFile)
		if err != nil {
			log.Printf("ERROR: Failed to load configuration from '%s': %v", cfgFile, err)
			continue
		}

		client, err := searcher.NewSearchClient(apiKey, config)
		if err != nil {
			log.Printf("ERROR: Failed to create search client for '%s': %v", cfgFile, err)
			continue
		}

		// Generate output filename
		timestamp := time.Now().Format("20060102_150405")
		configName := strings.TrimSuffix(filepath.Base(cfgFile), filepath.Ext(cfgFile))
		outputFilename := filepath.Join(*outputDir, fmt.Sprintf("search_results_%s_%s.json", configName, timestamp))

		fmt.Printf("Executing searches...\n")

		// Build output structure
		output := searcher.OutputResult{
			Timestamp:     time.Now().Format(time.RFC3339),
			ConfigFile:    cfgFile,
			Configuration: client.BuildConfigurationOutput(),
			Searches:      client.ExecuteAllSearches(),
		}

		// Write JSON to file
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			log.Printf("ERROR: Failed to create output file '%s': %v", outputFilename, err)
			continue
		}

		encoder := json.NewEncoder(outputFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			outputFile.Close()
			log.Printf("ERROR: Failed to write JSON to '%s': %v", outputFilename, err)
			continue
		}
		outputFile.Close()

		// Calculate statistics
		successCount := 0
		for _, search := range output.Searches {
			if search.Error == "" {
				successCount++
			}
		}

		fmt.Printf("✓ Search results saved to: %s\n", outputFilename)
		fmt.Printf("  Total searches: %d\n", len(output.Searches))
		fmt.Printf("  Successful: %d\n", successCount)
		if successCount < len(output.Searches) {
			fmt.Printf("  Failed: %d\n", len(output.Searches)-successCount)
		}

		totalSearches += len(output.Searches)
		totalSuccessful += successCount
		totalFailed += (len(output.Searches) - successCount)
	}

	// Print summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Config files processed: %d\n", len(configFiles))
	fmt.Printf("Total searches executed: %d\n", totalSearches)
	fmt.Printf("Total successful: %d\n", totalSuccessful)
	if totalFailed > 0 {
		fmt.Printf("Total failed: %d\n", totalFailed)
	}
	fmt.Printf("Results directory: %s\n", *outputDir)
}
