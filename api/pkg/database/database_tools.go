package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func Copy_Schema(db *sql.DB, config Config, useIAM bool, NewSchema string) error {

	// Step 1: Dump the public schema structure
	dumpFile := "schema_dump.sql"
	err := dumpSchemaStructure(config, dumpFile, useIAM)
	if err != nil {
		fmt.Printf("Failed to dump schema structure: %v", err)
		return err
	}

	// Step 2: Create new schema and modify the dump file
	fmt.Println("Creating and modifying schema dump...")
	modifiedDump := "modified_schema.sql"
	err = createAndModifyDump(dumpFile, modifiedDump, NewSchema)
	if err != nil {
		fmt.Printf("Failed to create and modify dump: %v", err)
		return err
	}

	// Step 3: Apply the modified structure
	fmt.Println("Applying modified schema structure...")
	err = applyModifiedStructure(config, modifiedDump, useIAM)
	if err != nil {
		fmt.Printf("Failed to apply modified structure: %v", err)
		return err
	}

	// Step 4: Copy data from public to new schema
	fmt.Println("Copying data to new schema...")
	err = copyData(db, NewSchema)
	if err != nil {
		fmt.Printf("Failed to copy data: %v", err)
		return err
	}

	fmt.Println("Schema copying completed successfully!")

	return err
}

// Get IAM authentication token
func getIAMToken(config Config) (string, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Default region
	}

	port := strconv.Itoa(config.Port)

	cmd := exec.Command("aws", "rds", "generate-db-auth-token",
		"--hostname", config.Host,
		"--port", port,
		"--region", region,
		"--username", config.User)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute AWS CLI: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// Dump the structure of public schema using pg_dump
func dumpSchemaStructure(config Config, dumpFile string, useIAM bool) error {
	fmt.Println("Dumping public schema structure...")

	var cmd *exec.Cmd

	port := strconv.Itoa(config.Port)

	pgDumpArgs := []string{
		"-h", config.Host,
		"-p", port,
		"-U", config.User,
		"-d", config.DBName,
		"-n", "public",
		"-s", // schema only, no data
		"-f", dumpFile,
	}

	cmd = exec.Command("pg_dump", pgDumpArgs...)

	// Set password or token as environment variable
	cmd.Env = os.Environ()
	if useIAM {
		token, err := getIAMToken(config)
		if err != nil {
			return fmt.Errorf("failed to get IAM token: %w", err)
		}
		cmd.Env = append(cmd.Env, "PGPASSWORD="+token)
	} else {
		cmd.Env = append(cmd.Env, "PGPASSWORD="+config.Password)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Create new schema and modify the dumped SQL
func createAndModifyDump(dumpFile, modifiedDump, newSchema string) error {
	// Read the original dump file
	content, err := ioutil.ReadFile(dumpFile)
	if err != nil {
		return fmt.Errorf("failed to read dump file: %w", err)
	}

	// Modify content: replace schema references
	modified := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;\n\n", newSchema)
	modifiedContent := strings.ReplaceAll(string(content), "public.", newSchema+".")
	modifiedContent = strings.ReplaceAll(modifiedContent, "SET search_path = public", "SET search_path = "+newSchema)
	modified += modifiedContent

	// Write to the new file
	err = ioutil.WriteFile(modifiedDump, []byte(modified), 0644)
	if err != nil {
		return fmt.Errorf("failed to write modified dump: %w", err)
	}

	return nil
}

// Apply the modified structure to the database
func applyModifiedStructure(config Config, modifiedDump string, useIAM bool) error {
	var cmd *exec.Cmd

	port := strconv.Itoa(config.Port)

	psqlArgs := []string{
		"-h", config.Host,
		"-p", port,
		"-U", config.User,
		"-d", config.DBName,
		"-f", modifiedDump,
	}

	cmd = exec.Command("psql", psqlArgs...)

	// Set password or token as environment variable
	cmd.Env = os.Environ()
	if useIAM {
		token, err := getIAMToken(config)
		if err != nil {
			return fmt.Errorf("failed to get IAM token: %w", err)
		}
		cmd.Env = append(cmd.Env, "PGPASSWORD="+token)
	} else {
		cmd.Env = append(cmd.Env, "PGPASSWORD="+config.Password)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("psql failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Copy data from public schema to new schema
func copyData(db *sql.DB, newSchema string) error {
	// Get list of tables in public schema
	rows, err := db.Query(`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`)
	if err != nil {
		return fmt.Errorf("failed to get table list: %w", err)
	}
	defer rows.Close()

	// Copy data for each table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}

		fmt.Printf("Copying data for table: %s\n", tableName)

		// Copy data using INSERT INTO ... SELECT
		query := fmt.Sprintf("INSERT INTO %s.%s SELECT * FROM public.%s",
			newSchema, tableName, tableName)

		_, err := db.Exec(query)
		if err != nil {
			// Continue with other tables even if one fails
			fmt.Printf("Warning: Failed to copy data for table %s: %v\n", tableName, err)
			continue
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating through tables: %w", err)
	}

	return nil
}

// ValidateUUID checks if the provided string is a valid UUID
// Returns the validated UUID string and nil if valid
// Returns empty string and error if invalid
func ValidateUUID(uuid string) (string, error) {
	// Remove any whitespace
	uuid = strings.TrimSpace(uuid)

	// Check if empty
	if uuid == "" {
		return "", errors.New("uuid cannot be empty")
	}

	// Standard UUID format: 8-4-4-4-12 hexadecimal digits
	// Example: 550e8400-e29b-41d4-a716-446655440000
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	if !uuidPattern.MatchString(uuid) {
		return "", fmt.Errorf("invalid UUID format: %s", uuid)
	}

	return uuid, nil
}
