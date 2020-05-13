package connector

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/aarondwi/dbm/schema"

	// this module directly use postgres
	_ "github.com/lib/pq"
)

// DbPostgres is the implementation
// of DbAccess interface for postgres
type DbPostgres struct {
	db *sql.DB
}

// Init creates the connection for target DB
func (d *DbPostgres) Init(conf schema.Conf) error {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port,
		conf.Username, conf.Password,
		conf.Database, conf.Sslmode,
	)
	// we silent the error
	// because either way it won't throw database error
	// see https://github.com/lib/pq/issues/581
	db, _ := sql.Open("postgres", psqlInfo)
	err := db.Ping()
	if err != nil {
		log.Printf("Failed creating connection: %v", err)
		return err
	}

	d.db = db
	return nil
}

// BlindExec blindly executes (shocking, right?)
// the query in the files
// this is the EXPECTED behavior, as we want to allow people to
// write it in their database dialect directly
func (d *DbPostgres) BlindExec(stmt string) error {
	_, err := d.db.Exec(stmt)
	if err != nil {
		log.Printf("Failed executing BlindExec: %v", err)
		return err
	}
	return nil
}

// CreateLogTable generates the table
// to store logs in the database
// it will be named dbm_logs
func (d *DbPostgres) CreateLogTable() error {
	_, err := d.db.Exec(`
		CREATE TABLE dbm_logs (
			filename text primary key,
			created_at timestamptz not null default now()
		)
	`)
	if err != nil {
		log.Printf("Failed creating table dbm_logs: %v", err)
		return err
	}
	return nil
}

// DropLogTable removes dbm_logs table
func (d *DbPostgres) DropLogTable() error {
	_, err := d.db.Exec(`DROP TABLE dbm_logs`)
	if err != nil {
		log.Printf("Failed dropping table dbm_logs: %v", err)
		return err
	}
	return nil
}

// InsertLogs adds the filenames to dbm_logs table
func (d *DbPostgres) InsertLogs(filenames []string) error {
	params := make([]string, len(filenames))
	args := make([]interface{}, len(filenames))
	c := 1
	for i, s := range filenames {
		params[i] = fmt.Sprintf("($%d)", c)
		args[i] = s
		c++
	}

	stmt := "INSERT INTO dbm_logs(filename) VALUES " +
		strings.Join(params, ",")
	_, err := d.db.Exec(stmt, args...)
	if err != nil {
		log.Printf("Failed inserting into dbm_logs: %v", err)
		return err
	}
	return nil
}

// DeleteLog removes the filename from the dbm_logs table
func (d *DbPostgres) DeleteLog(filename string) error {
	stmt := "DELETE FROM dbm_logs WHERE filename = $1"
	_, err := d.db.Exec(stmt, filename)
	if err != nil {
		log.Printf("Failed deleting %s from dbm_logs: %v", filename, err)
		return err
	}
	return nil
}

// GetLastLog retrieves the last log from dbm_logs table
func (d *DbPostgres) GetLastLog() (string, error) {
	var filename string
	err := d.db.QueryRow(`SELECT filename FROM dbm_logs ORDER BY filename DESC LIMIT 1`).Scan(&filename)
	if err != nil {
		log.Printf("Failed retrieving from dbm_logs: %v", err)
		return "", err
	}
	return filename, nil
}

// ListAlreadyUp retrieves all filenames in dbm_logs table
// filenames listed here mean those that has been applied to the database
func (d *DbPostgres) ListAlreadyUp() ([]string, error) {
	rows, err := d.db.Query(`SELECT filename FROM dbm_logs ORDER BY filename`)
	if err != nil {
		log.Printf("Failed retrieving dbm_logs: %v", err)
		return nil, err
	}
	var result []string
	for rows.Next() {
		var filename string
		rows.Scan(&filename)
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		result = append(result, filename)
	}
	return result, nil
}

// Close the connection to the database
func (d *DbPostgres) Close() {
	d.db.Close()
}
