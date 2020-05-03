package dbm

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type DbPostgresConfig struct {
	Host     string
	Port     int16
	Username string
	Password string
	Database string
	Sslmode  string
}

type DbPostgres struct {
	db *sql.DB
}

func (d *DbPostgres) Init(conf DbPostgresConfig) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port,
		conf.Username, conf.Password,
		conf.Database, conf.Sslmode,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed creating connection: %v", err)
	}

	d.db = db
}
func (d *DbPostgres) CreateLogTable() error {
	_, err := d.db.Exec(`
		CREATE TABLE dbm_logs (
			filename text primary key,
			created_at timestamptz not null default now()
		)
	`)
	if err != nil {
		log.Fatalf("Failed creating table dbm_logs: %v", err)
		return err
	}
	return nil
}

func (d *DbPostgres) DropLogTable() error {
	_, err := d.db.Exec(`DROP TABLE dbm_logs`)
	if err != nil {
		log.Fatalf("Failed dropping table dbm_logs: %v", err)
		return err
	}
	return nil
}

func (d *DbPostgres) InsertLogs(filenames []string) error {
	params := make([]string, len(filenames))
	args := make([]interface{}, len(filenames))
	c := 1
	for i, s := range filenames {
		params[i] = fmt.Sprintf("($%d)", c)
		args[i] = s
		c += 1
	}
	
	stmt := "INSERT INTO dbm_logs(filename) VALUES " + 
		strings.Join(params, ",")
	_, err := d.db.Exec(stmt, args...)
	if err != nil {
		log.Fatalf("Failed inserting into dbm_logs: %v", err)
		return err
	}
	return nil
}

func (d *DbPostgres) DeleteLogs(filenames []string) error {
	params := make([]string, len(filenames))
	args := make([]interface{}, len(filenames))
	c := 1
	for i, s := range filenames {
		params[i] = fmt.Sprintf("$%d", c)
		args[i] = s
		c += 1
	}

	stmt := "DELETE FROM dbm_logs WHERE filename IN ( " + 
		strings.Join(params, ",") + ")"
	_, err := d.db.Exec(stmt, args...)
	if err != nil {
		log.Fatalf("Failed deleting dbm_logs: %v", err)
		return err
	}
	return nil
}

func (d *DbPostgres) ListAlreadyUp() ([]string, error) {
	rows, err := d.db.Query(`SELECT filename FROM dbm_logs ORDER BY filename`)
	if err != nil {
		log.Fatalf("Failed retrieving dbm_logs: %v", err)
	}
	var result []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		result = append(result, filename)
	}
	return result, nil
}

func (d *DbPostgres) Close() {
	d.db.Close()
}
