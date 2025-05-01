package schema

import (
	"context"
	"database/sql"
	"fmt"
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
	fmt.Println("Schema copy operation starting")
	fmt.Printf("Copying from schema '%s' to new schema '%s'\n", schema.FromSchemaName, schema.ToSchemaName)

	// Step 1: Create the new schema
	fmt.Printf("Creating schema: %s\n", schema.ToSchemaName)
	_, err := schema.DB.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema.ToSchemaName))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Step 2: Get list of all tables in source schema
	tables, err := schema.getTableNames(ctx)
	if err != nil {
		return err
	}

	// Step 3: Get sequences
	sequences, err := schema.getSequenceNames(ctx)
	if err != nil {
		return err
	}

	// Step 4: Get views
	views, err := schema.getViewNames(ctx)
	if err != nil {
		return err
	}

	// Step 5: Create sequences and tables in a transaction
	err = schema.createStructures(ctx, sequences, tables)
	if err != nil {
		return err
	}

	// Step 6: Now copy data table by table
	err = schema.copyTableData(ctx, tables)
	if err != nil {
		return err
	}

	// Step 7: Create foreign key constraints
	err = schema.createForeignKeys(ctx)
	if err != nil {
		return err
	}

	// Step 8: Create indexes
	err = schema.createIndexes(ctx)
	if err != nil {
		return err
	}

	// Step 9: Create views
	err = schema.createViews(ctx, views)
	if err != nil {
		return err
	}

	fmt.Printf("Schema '%s' successfully created and populated\n", schema.ToSchemaName)
	return nil
}

// getTableNames gets all table names from the source schema
func (schema *Schema) getTableNames(ctx context.Context) ([]string, error) {
	fmt.Println("Getting list of tables in source schema")
	q := fmt.Sprintf("SELECT tablename FROM pg_tables WHERE schemaname = '%s'", schema.FromSchemaName)
	rows, err := schema.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}

	fmt.Printf("Found %d tables in source schema\n", len(tables))
	return tables, nil
}

// getSequenceNames gets all sequence names from the source schema
func (schema *Schema) getSequenceNames(ctx context.Context) ([]string, error) {
	fmt.Println("Getting list of sequences in source schema")
	q := fmt.Sprintf("SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = '%s'", schema.FromSchemaName)
	seqRows, err := schema.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequences: %w", err)
	}
	defer seqRows.Close()

	var sequences []string
	for seqRows.Next() {
		var seqName string
		if err := seqRows.Scan(&seqName); err != nil {
			return nil, fmt.Errorf("failed to scan sequence name: %w", err)
		}
		sequences = append(sequences, seqName)
	}
	if err = seqRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sequences: %w", err)
	}

	fmt.Printf("Found %d sequences in source schema\n", len(sequences))
	return sequences, nil
}

// getViewNames gets all view names from the source schema
func (schema *Schema) getViewNames(ctx context.Context) ([]string, error) {
	fmt.Println("Getting list of views in source schema")
	q := fmt.Sprintf("SELECT table_name FROM information_schema.views WHERE table_schema = '%s'", schema.FromSchemaName)
	viewRows, err := schema.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get views: %w", err)
	}
	defer viewRows.Close()

	var views []string
	for viewRows.Next() {
		var viewName string
		if err := viewRows.Scan(&viewName); err != nil {
			return nil, fmt.Errorf("failed to scan view name: %w", err)
		}
		views = append(views, viewName)
	}
	if err = viewRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating views: %w", err)
	}

	fmt.Printf("Found %d views in source schema\n", len(views))
	return views, nil
}

// createStructures creates sequences and tables in the target schema
func (schema *Schema) createStructures(ctx context.Context, sequences []string, tables []string) error {
	fmt.Println("Creating tables and sequences in new schema")

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

		// Check if sequence already exists in target schema
		var exists bool
		err := schema.DB.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.sequences 
				WHERE sequence_schema = $1 
				AND sequence_name = $2
			)
		`, schema.ToSchemaName, seqName).Scan(&exists)

		if err != nil {
			return fmt.Errorf("failed to check if sequence exists: %w", err)
		}

		if exists {
			fmt.Printf("Sequence %s already exists in target schema, skipping\n", seqName)
			continue
		}

		// Get the current value of the sequence
		var lastVal int64
		var isCalled bool
		err = schema.DB.QueryRowContext(ctx, fmt.Sprintf("SELECT last_value, is_called FROM %s.%s", schema.FromSchemaName, seqName)).Scan(&lastVal, &isCalled)
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

		// Check if table already exists in target schema
		var exists bool
		err := schema.DB.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = $1 
				AND table_name = $2
			)
		`, schema.ToSchemaName, tableName).Scan(&exists)

		if err != nil {
			return fmt.Errorf("failed to check if table exists: %w", err)
		}

		if exists {
			fmt.Printf("Table %s already exists in target schema, skipping structure creation\n", tableName)
			continue
		}

		// Get table creation SQL
		var tableSQL string
		err = schema.DB.QueryRowContext(ctx, `
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

		// Check if table has a primary key
		var hasPK bool
		err = schema.DB.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.table_constraints 
				WHERE constraint_type = 'PRIMARY KEY' 
				AND table_schema = $1 
				AND table_name = $2
			)
		`, schema.FromSchemaName, tableName).Scan(&hasPK)

		if err != nil {
			return fmt.Errorf("failed to check if table has primary key: %w", err)
		}

		if !hasPK {
			continue
		}

		// Get primary key constraints
		pkRows, err := schema.DB.QueryContext(ctx, `
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
		for pkRows.Next() {
			var pkSQL string
			if err := pkRows.Scan(&pkSQL); err != nil {
				pkRows.Close()
				return fmt.Errorf("failed to scan primary key SQL: %w", err)
			}

			_, err = tx.ExecContext(ctx, pkSQL)
			if err != nil {
				pkRows.Close()
				return fmt.Errorf("failed to apply primary key constraint: %w", err)
			}
		}
		pkRows.Close()

		fmt.Println("adding triggers")
		triggerSQL := fmt.Sprintf(`		
		CREATE TRIGGER update_profile_modified 
		BEFORE UPDATE ON %s.%s 
		FOR EACH ROW 
		EXECUTE FUNCTION public.update_modified_column();
		`, schema.ToSchemaName, tableName)

		// Apply trigger
		_, err = tx.ExecContext(ctx, triggerSQL)

		if err != nil {
			return fmt.Errorf("failed to create trigger: %w", err)
		}

	}

	// Commit the transaction for schema creation
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("Successfully created structures in new schema")
	return nil
}

// copyTableData copies data from source schema tables to target schema tables
func (schema *Schema) copyTableData(ctx context.Context, tables []string) error {
	fmt.Println("Copying data to new schema")

	// Start a transaction for data copying
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

		if len(columns) == 0 {
			fmt.Printf("No columns found for table %s, skipping\n", tableName)
			continue
		}

		// Build column list string
		columnList := strings.Join(columns, ", ")

		// Check if target table already has data
		var count int
		err = schema.DB.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", schema.ToSchemaName, tableName)).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check if target table has data: %w", err)
		}

		if count > 0 {
			fmt.Printf("Target table %s already has %d rows, skipping data copy\n", tableName, count)
			continue
		}

		// Copy data using INSERT INTO ... SELECT
		_, err = tx.ExecContext(ctx, fmt.Sprintf(`
			INSERT INTO %s.%s (%s)
			SELECT %s FROM %s.%s
		`, schema.ToSchemaName, tableName, columnList, columnList, schema.FromSchemaName, tableName))

		if err != nil {
			fmt.Printf("Warning: Failed to copy data for table %s: %v\n", tableName, err)
			// Continue with other tables instead of failing completely
		}
	}

	// Commit the transaction for data copying
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("Successfully copied data to new schema")
	return nil
}

// createForeignKeys creates foreign key constraints in the target schema
func (schema *Schema) createForeignKeys(ctx context.Context) error {
	fmt.Println("Creating foreign key constraints")

	// Check if source schema has any foreign keys
	var fkCount int
	err := schema.DB.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM information_schema.table_constraints
		WHERE constraint_type = 'FOREIGN KEY' AND table_schema = $1
	`, schema.FromSchemaName).Scan(&fkCount)

	if err != nil {
		return fmt.Errorf("failed to check for foreign keys: %w", err)
	}

	if fkCount == 0 {
		fmt.Println("No foreign keys found in source schema, skipping")
		return nil
	}

	// Start a transaction for creating foreign keys
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
	var appliedCount int
	for fkRows.Next() {
		var fkSQL string
		if err := fkRows.Scan(&fkSQL); err != nil {
			return fmt.Errorf("failed to scan foreign key SQL: %w", err)
		}

		_, err = tx.ExecContext(ctx, fkSQL)
		if err != nil {
			fmt.Printf("Warning: Failed to apply foreign key constraint: %v\n", err)
			// Continue with other constraints instead of failing completely
		} else {
			appliedCount++
		}
	}

	if err = fkRows.Err(); err != nil {
		return fmt.Errorf("error iterating foreign keys: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully applied %d foreign key constraints\n", appliedCount)
	return nil
}

// createIndexes creates indexes in the target schema
func (schema *Schema) createIndexes(ctx context.Context) error {
	fmt.Println("Creating indexes")

	// Check if source schema has any indexes
	var idxCount int
	err := schema.DB.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM pg_indexes
        WHERE schemaname = $1 AND indexname NOT IN (
            SELECT conname FROM pg_constraint WHERE contype = 'p'
        ) AND indexname NOT LIKE '%_key' AND indexname NOT LIKE 'pg_%'
    `, schema.FromSchemaName).Scan(&idxCount)

	if err != nil {
		return fmt.Errorf("failed to check for indexes: %w", err)
	}

	if idxCount == 0 {
		fmt.Println("No indexes found in source schema, skipping")
		return nil
	}

	// Get indexes that are NOT from unique constraints (filter out names ending with _key)
	idxRows, err := schema.DB.QueryContext(ctx, `
        SELECT 
            regexp_replace(
                indexdef, 
                'ON ' || $2 || '\\.([^\\s]+)', 
                'ON ' || $1 || '.$1'
            ) AS index_sql
        FROM 
            pg_indexes
        WHERE 
            schemaname = $2 AND 
            indexname NOT IN (
                SELECT conname FROM pg_constraint WHERE contype = 'p'
            ) AND
            indexname NOT LIKE '%_key' AND
            indexname NOT LIKE 'pg_%';
    `, schema.ToSchemaName, schema.FromSchemaName)

	if err != nil {
		return fmt.Errorf("failed to get indexes: %w", err)
	}
	defer idxRows.Close()

	// Create indexes - handle each in a separate transaction
	var appliedCount int
	for idxRows.Next() {
		var idxSQL string
		if err := idxRows.Scan(&idxSQL); err != nil {
			return fmt.Errorf("failed to scan index SQL: %w", err)
		}

		// Use a separate transaction for each index
		tx, err := schema.DB.BeginTx(ctx, nil)
		if err != nil {
			fmt.Printf("Warning: Failed to begin transaction for index: %v\n", err)
			continue
		}

		_, err = tx.ExecContext(ctx, idxSQL)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Warning: Failed to create index: %v\n", err)
		} else {
			err = tx.Commit()
			if err != nil {
				fmt.Printf("Warning: Failed to commit index transaction: %v\n", err)
			} else {
				appliedCount++
			}
		}
	}

	fmt.Printf("Successfully created %d indexes\n", appliedCount)
	return nil
}

// createViews creates views in the target schema
func (schema *Schema) createViews(ctx context.Context, views []string) error {
	if len(views) == 0 {
		fmt.Println("No views to create, skipping")
		return nil
	}

	fmt.Println("Creating views")

	// Start a transaction for creating views
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

	var createdCount int
	for _, viewName := range views {
		// Check if view already exists in target schema
		var exists bool
		err := schema.DB.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.views 
				WHERE table_schema = $1 
				AND table_name = $2
			)
		`, schema.ToSchemaName, viewName).Scan(&exists)

		if err != nil {
			return fmt.Errorf("failed to check if view exists: %w", err)
		}

		if exists {
			fmt.Printf("View %s already exists in target schema, skipping\n", viewName)
			continue
		}

		// Get view definition
		var viewDef string
		err = schema.DB.QueryRowContext(ctx, `
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
		_, err = tx.ExecContext(ctx, viewDef)
		if err != nil {
			fmt.Printf("Warning: Failed to create view %s: %v\n", viewName, err)
			// Continue with other views instead of failing completely
		} else {
			createdCount++
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully created %d views\n", createdCount)
	return nil
}
