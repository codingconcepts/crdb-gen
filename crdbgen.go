package main

import (
	"context"

	"flag"
	"fmt"
	"log"
	"os"

	"github.com/codingconcepts/crdb-gen/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	url := flag.String("url", "postgres://root@localhost:26257/defaultdb", "Connection URL, of the form: postgresql://[user[:passwd]@]host[:port]/[db][?parameters...]")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		os.Exit(2)
	}

	db, err := pgxpool.New(context.Background(), *url)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	tables, err := fetchTables(db)
	if err != nil {
		log.Fatalf("error fetching tables: %v", err)
	}

	for _, table := range tables {
		fmt.Printf("%s\n", table.Name)
		for _, column := range table.Columns {
			fmt.Printf("\t%s (%s)\n", column.Name, column.DataType)
		}
	}
}

func fetchTables(pool *pgxpool.Pool) ([]model.Table, error) {
	// Tables.
	const tableStmt = `SELECT table_name FROM information_schema.tables
					   WHERE is_insertable_into = 'YES'`

	rows, err := pool.Query(context.Background(), tableStmt)
	if err != nil {
		return nil, fmt.Errorf("fetching tables: %w", err)
	}

	var tables []model.Table
	for rows.Next() {
		var t model.Table
		if err = rows.Scan(&t.Name); err != nil {
			return nil, fmt.Errorf("scanning table: %w", err)
		}
		tables = append(tables, t)
	}

	// Columns.
	for i, table := range tables {
		if tables[i].Columns, err = fetchColumns(pool, &table); err != nil {
			return nil, fmt.Errorf("adding columns: %w", err)
		}
	}

	return tables, nil
}

func fetchColumns(pool *pgxpool.Pool, table *model.Table) ([]model.Column, error) {
	const tableStmt = `SELECT column_name, crdb_sql_type FROM information_schema.columns
					   WHERE table_name = $1`

	rows, err := pool.Query(context.Background(), tableStmt, table.Name)
	if err != nil {
		return nil, fmt.Errorf("fetching columns: %w", err)
	}

	var columns []model.Column
	for rows.Next() {
		var c model.Column
		if err = rows.Scan(&c.Name, &c.DataType); err != nil {
			return nil, fmt.Errorf("scanning column: %w", err)
		}
		columns = append(columns, c)
	}

	return columns, nil
}
