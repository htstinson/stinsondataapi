package schema

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Schema struct {
	DB             *sql.DB
	FromSchemaName string
	ToSchemaName   string
}

// CopySchema creates a new schema with the specified name and copies all tables,
// sequences, functions, and data from the source schema to the new schema.
func (schema *Schema) CopySchema(ctx context.Context) error {
	fmt.Println("schema CopySchema")

	// Step 1: Create the new schema
	fmt.Printf("Creating schema: %s\n", schema.ToSchemaName)
	_, err := schema.DB.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema.ToSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Step 2: Get list of all tables in source schema
	fmt.Println("Getting list of tables in source schema")
	q := fmt.Sprintf("SELECT tablename FROM pg_tables WHERE schemaname = '%s'", schema.FromSchemaName)
	rows, err := schema.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating tables: %w", err)
	}

	// Step 3: Get sequences
	fmt.Println("Getting list of sequences in source schema")
	q = fmt.Sprintf("SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = '%s'", schema.FromSchemaName)
	seqRows, err := schema.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to get sequences: %w", err)
	}
	defer seqRows.Close()

	var sequences []string
	for seqRows.Next() {
		var seqName string
		if err := seqRows.Scan(&seqName); err != nil {
			return fmt.Errorf("failed to scan sequence name: %w", err)
		}
		sequences = append(sequences, seqName)
	}
	if err = seqRows.Err(); err != nil {
		return fmt.Errorf("error iterating sequences: %w", err)
	}

	// Step 4: Get views
	fmt.Println("Getting list of views in source schema")
	q = fmt.Sprintf("SELECT table_name FROM information_schema.views WHERE table_schema = '%s'", schema.FromSchemaName)
	viewRows, err := schema.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to get views: %w", err)
	}
	defer viewRows.Close()

	var views []string
	for viewRows.Next() {
		var viewName string
		if err := viewRows.Scan(&viewName); err != nil {
			return fmt.Errorf("failed to scan view name: %w", err)
		}
		views = append(views, viewName)
	}
	if err = viewRows.Err(); err != nil {
		return fmt.Errorf("error iterating views: %w", err)
	}

	// Step 5: Get table structures (DDL) and create them in the new schema
	fmt.Println("Creating tables in new schema")

	// Start a transaction for schema creation
	tx, err := schema.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create sequences first (since tables may depend on them)
	for _, seqName := range sequences {
		fmt.Printf("Creating sequence: %s\n", seqName)

		// Get the current value of the sequence
		var lastVal int64
		var isCalled bool
		err := schema.DB.QueryRowContext(ctx, fmt.Sprintf("SELECT last_value, is_called FROM %s.%s", schema.FromSchemaName, seqName)).Scan(&lastVal, &isCalled)
		if err != nil {
			return fmt.Errorf("failed to get sequence value: %w", err)
		}

		// Create the sequence in the new schema
		_, err = tx.ExecContext(ctx, fmt.Sprintf("CREATE SEQUENCE %s.%s", schema.ToSchemaName, seqName))
		if err != nil {
			return fmt.Errorf("failed to create sequence: %w", err)
		}

		// Set the sequence to match the original
		if isCalled {
			_, err = tx.ExecContext(ctx, fmt.Sprintf("SELECT setval('%s.%s', %d, true)", schema.ToSchemaName, seqName, lastVal))
		} else {
			_, err = tx.ExecContext(ctx, fmt.Sprintf("SELECT setval('%s.%s', %d, false)", schema.ToSchemaName, seqName, lastVal))
		}
		if err != nil {
			return fmt.Errorf("failed to set sequence value: %w", err)
		}
	}

	// For each table, get its structure and recreate it in the new schema
	for _, tableName := range tables {
		fmt.Printf("Processing table: %s\n", tableName)

		// Get table creation SQL
		var tableSQL string
		err := schema.DB.QueryRowContext(ctx, `
			SELECT 
				'CREATE TABLE IF NOT EXISTS ' || $1 || '.' || c.relname || ' (' || 
				string_agg(
					column_name || ' ' || data_type || 
					CASE 
						WHEN character_maximum_length IS NOT NULL THEN '(' || character_maximum_length || ')' 
						ELSE '' 
					END || 
					CASE 
						WHEN is_nullable = 'NO' THEN ' NOT NULL' 
						ELSE '' 
					END || 
					CASE 
						WHEN column_default IS NOT NULL THEN ' DEFAULT ' || column_default 
						ELSE '' 
					END,
					', '
				) || ');'
			FROM 
				information_schema.columns
			JOIN 
				pg_class c ON c.relname = table_name
			JOIN 
				pg_namespace n ON n.oid = c.relnamespace AND n.nspname = table_schema
			WHERE 
				table_schema = $2 AND table_name = $3
			GROUP BY 
				c.relname;
		`, schema.ToSchemaName, schema.FromSchemaName, tableName).Scan(&tableSQL)

		if err != nil {
			return fmt.Errorf("failed to get table definition for %s: %w", tableName, err)
		}

		// Create the table
		_, err = tx.ExecContext(ctx, tableSQL)
		if err != nil {
			return fmt.Errorf("failed to create table %s: %w", tableName, err)
		}

		// Get primary key constraints
		rows, err := schema.DB.QueryContext(ctx, `
			SELECT
				'ALTER TABLE ' || $1 || '.' || tc.table_name || 
				' ADD CONSTRAINT ' || tc.constraint_name || ' PRIMARY KEY (' ||
				string_agg(kcu.column_name, ', ') || ');'
			FROM
				information_schema.table_constraints tc
			JOIN
				information_schema.key_column_usage kcu ON kcu.constraint_name = tc.constraint_name
				AND kcu.table_schema = tc.table_schema
				AND kcu.table_name = tc.table_name
			WHERE
				tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $2 AND tc.table_name = $3
			GROUP BY
				tc.table_name, tc.constraint_name;
		`, schema.ToSchemaName, schema.FromSchemaName, tableName)

		if err != nil {
			return fmt.Errorf("failed to get primary key constraints: %w", err)
		}

		// Apply primary key constraints
		for rows.Next() {
			var pkSQL string
			if err := rows.Scan(&pkSQL); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan primary key SQL: %w", err)
			}

			_, err = tx.ExecContext(ctx, pkSQL)
			if err != nil {
				rows.Close()
				return fmt.Errorf("failed to apply primary key constraint: %w", err)
			}
		}
		rows.Close()
	}

	// Commit the transaction for schema creation
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Step 6: Now copy data table by table
	fmt.Println("Copying data to new schema")
	for _, tableName := range tables {
		fmt.Printf("Copying data for table: %s\n", tableName)

		// Get columns for this table
		colRows, err := schema.DB.QueryContext(ctx, `
			SELECT column_name
			FROM information_schema.columns
			WHERE table_schema = $1 AND table_name = $2
			ORDER BY ordinal_position
		`, schema.FromSchemaName, tableName)

		if err != nil {
			return fmt.Errorf("failed to get columns for table %s: %w", tableName, err)
		}

		var columns []string
		for colRows.Next() {
			var col string
			if err := colRows.Scan(&col); err != nil {
				colRows.Close()
				return fmt.Errorf("failed to scan column name: %w", err)
			}
			columns = append(columns, col)
		}
		colRows.Close()

		// Build column list string
		columnList := strings.Join(columns, ", ")

		// Copy data using INSERT INTO ... SELECT
		_, err = schema.DB.ExecContext(ctx, fmt.Sprintf(`
			INSERT INTO %s.%s (%s)
			SELECT %s FROM %s.%s
		`, schema.ToSchemaName, tableName, columnList, columnList, schema.FromSchemaName, tableName))

		if err != nil {
			// Log error but continue with other tables
			fmt.Printf("Warning: Failed to copy data for table %s: %v\n", tableName, err)
		}
	}

	// Step 7: Create foreign key constraints
	fmt.Println("Creating foreign key constraints")
	fkRows, err := schema.DB.QueryContext(ctx, `
		SELECT
			'ALTER TABLE ' || $1 || '.' || tc.table_name || 
			' ADD CONSTRAINT ' || tc.constraint_name || ' FOREIGN KEY (' ||
			string_agg(kcu.column_name, ', ') || ') REFERENCES ' || $1 || '.' || 
			ccu.table_name || ' (' || string_agg(ccu.column_name, ', ') || ');'
		FROM
			information_schema.table_constraints tc
		JOIN
			information_schema.key_column_usage kcu ON kcu.constraint_name = tc.constraint_name
			AND kcu.table_schema = tc.table_schema
			AND kcu.table_name = tc.table_name
		JOIN
			information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		WHERE
			tc.constraint_type = 'FOREIGN KEY' AND tc.table_schema = $2
		GROUP BY
			tc.table_name, tc.constraint_name, ccu.table_name;
	`, schema.ToSchemaName, schema.FromSchemaName)

	if err != nil {
		return fmt.Errorf("failed to get foreign key constraints: %w", err)
	}
	defer fkRows.Close()

	// Apply foreign key constraints
	for fkRows.Next() {
		var fkSQL string
		if err := fkRows.Scan(&fkSQL); err != nil {
			return fmt.Errorf("failed to scan foreign key SQL: %w", err)
		}

		_, err = schema.DB.ExecContext(ctx, fkSQL)
		if err != nil {
			// Log error but continue with other constraints
			fmt.Printf("Warning: Failed to apply foreign key constraint: %v\n", err)
		}
	}

	// Step 8: Create indexes
	log.Println("Creating indexes")
	idxRows, err := schema.DB.QueryContext(ctx, `
    SELECT
        'CREATE INDEX ' || indexname || ' ON ' || $1 || '.' || tablename || ' USING ' || 
        substring(indexdef from position(' USING ' in indexdef) for 999)
    FROM
        pg_indexes
    WHERE
        schemaname = $2 AND indexname NOT IN (
            SELECT conname FROM pg_constraint WHERE contype = 'p'
        );
	`, schema.ToSchemaName, schema.FromSchemaName)

	if err != nil {
		return fmt.Errorf("failed to get indexes: %w", err)
	}
	defer idxRows.Close()

	// Create indexes
	for idxRows.Next() {
		var idxSQL string
		if err := idxRows.Scan(&idxSQL); err != nil {
			return fmt.Errorf("failed to scan index SQL: %w", err)
		}

		_, err = schema.DB.ExecContext(ctx, idxSQL)
		if err != nil {
			// Log error but continue with other indexes
			fmt.Printf("Warning: Failed to create index: %v\n", err)
		}
	}

	// Step 9: Create views
	fmt.Println("Creating views")
	for _, viewName := range views {
		// Get view definition
		var viewDef string
		err := schema.DB.QueryRowContext(ctx, `
			SELECT 'CREATE VIEW ' || $1 || '.' || table_name || ' AS ' || view_definition
			FROM information_schema.views
			WHERE table_schema = $2 AND table_name = $3
		`, schema.ToSchemaName, schema.FromSchemaName, viewName).Scan(&viewDef)

		if err != nil {
			return fmt.Errorf("failed to get view definition for %s: %w", viewName, err)
		}

		// Rewrite the view definition to reference the new schema
		viewDef = strings.Replace(viewDef, schema.FromSchemaName+".", schema.ToSchemaName+".", -1)

		// Create the view
		_, err = schema.DB.ExecContext(ctx, viewDef)
		if err != nil {
			fmt.Printf("Warning: Failed to create view %s: %v\n", viewName, err)
		}
	}

	fmt.Printf("Schema '%s' successfully created and populated\n", schema.ToSchemaName)
	return nil
}
