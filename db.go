package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB should be constructed with NewDB.
type DB struct {
	client *sql.DB
}

// NewDB conditionally initializes our database for demo purposes, returning a DB with
// the methods intended for use by other layers.
func NewDB(filepath string, initData bool) DB {
	client := InitDB(filepath)
	if initData {
		InitTables(client)
		InitReferenceData(client)
	}

	db := DB{client}
	return db
}

func InitDB(filepath string) *sql.DB {
	client, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic("Failed to create or open db")
	}
	return client
}

// InitTables is provided for demo convenience. A real app shouldn't do this.
func InitTables(db *sql.DB) {
	tables := []string{`
	CREATE TABLE IF NOT EXISTS instrument(
		instrument_id INTEGER PRIMARY KEY,
		description TEXT NOT NULL
	);
	`, `
	CREATE TABLE IF NOT EXISTS run_instrument(
		run_id INTEGER PRIMARY KEY,
		instrument_id INTEGER NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(instrument_id) REFERENCES instrument(instrument_id)
	);
	`, `
	CREATE TABLE IF NOT EXISTS run_sample(
		sample_id INTEGER,
		run_id INTEGER,
		FOREIGN KEY(run_id) REFERENCES run_instrument(run_id)
		PRIMARY KEY(sample_id, run_id)
	);
	`}

	for _, t := range tables {
		fmt.Println(t)
		_, err := db.Exec(t)
		if err != nil {
			panic(err)
		}
	}
}

// InitReferenceData is provided for demo convenience. A real app shouldn't do this.
func InitReferenceData(db *sql.DB) {
	instruments := [3]string{"Instrument 1", "Instrument 2", "Instrument 3"}
	sql := `
	INSERT INTO instrument(
		description
	) VALUES (?);
	`

	statement, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for _, item := range instruments {
		_, err := statement.Exec(item)
		if err != nil {
			panic(err)
		}
	}
}

// SaveRun returns the newly generated run ID.
func (db DB) SaveRun(run Run) (int, error) {
	tx, err := db.client.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Println(err)
	}
	sql := `
	INSERT INTO run_instrument(
		instrument_id
	) VALUES (?);
	`

	runStatement, err := tx.Prepare(sql)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return 0, err
	}
	defer runStatement.Close()

	_, err = runStatement.Exec(run.InstrumentID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return 0, err
	}

	// Use sqlite's somewhat awkward last_insert_rowid() for retrieving our new run_id.
	sql = "SELECT last_insert_rowid() AS id"

	rows, err := tx.Query(sql)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var runID int
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&runID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	sql = `
	INSERT INTO run_sample(
		sample_id,
		run_id
	) VALUES (?, ?);
	`

	sampleStatement, err := tx.Prepare(sql)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return 0, err
	}
	defer sampleStatement.Close()

	for _, s := range run.Samples {
		_, err = sampleStatement.Exec(s.ID, runID)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	return runID, nil
}
