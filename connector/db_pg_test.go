package connector

import (
	"dbm/schema"
	"os"
	"testing"
)

var db DbPostgres

func testSetupAndTeardown(m *testing.M) int {
	conf := schema.Conf{
		Dialect:  "postgres",
		Host:     "127.0.0.1",
		Port:     5432,
		Username: "dbm",
		Password: "dbm",
		Database: "dbm",
		Sslmode:  "disable",
	}
	db.Init(conf)
	defer db.Close()
	return m.Run()
}
func TestMain(m *testing.M) {
	os.Exit(testSetupAndTeardown(m))
}

func TestFailedConnection(t *testing.T) {
	var dbTest DbPostgres
	confTest := schema.Conf{
		Dialect:  "postgres",
		Host:     "127.0.10.10",
		Port:     5432,
		Username: "notxist",
		Password: "notxist",
		Database: "notxist",
		Sslmode:  "notxist",
	}
	err := dbTest.Init(confTest)
	if err == nil {
		t.Fatalf("Should fail creating connection, but it is not")
	}
}

func TestCreateAndDropLogTable(t *testing.T) {
	err := db.CreateLogTable()
	if err != nil {
		t.Fatalf("Failed Creating dbm_logs: %v", err)
	}
	err = db.DropLogTable()
	if err != nil {
		t.Fatalf("Failed Dropping dbm_logs: %v", err)
	}
}

func TestCreateAndDropLogTableFailed(t *testing.T) {
	err := db.DropLogTable()
	if err == nil {
		t.Fatal("Should fail dropping dbm_logs, but it is not")
	}
	db.CreateLogTable()
	err = db.CreateLogTable()
	if err == nil {
		t.Fatalf("Should fail creating dbm_logs, but it is not")
	}
	db.DropLogTable()
}

func TestInsertAndListAndDeleteLogs(t *testing.T) {
	db.CreateLogTable()
	defer db.DropLogTable()

	filenames := []string{"file1", "file2", "file3", "file4", "file5"}
	err := db.InsertLogs(filenames)
	if err != nil {
		t.Fatalf("Failed inserting into dbm_logs: %v", err)
	}

	result, error := db.ListAlreadyUp()
	if error != nil {
		t.Fatalf("Failed retrieving logs from dbm_logs: %v", err)
	}
	if len(result) != len(filenames) {
		t.Fatalf("Inserted and Retrieved Data Mismatch: %v", err)
	}

	var res string
	res, error = db.GetLastLog()
	if error != nil {
		t.Fatalf("Failed retrieving last log from dbm_logs: %v", err)
	}
	if res != filenames[len(filenames)-1] {
		t.Fatalf("Inserted and Retrieved Data Mismatch: %v", err)
	}

	err = db.DeleteLog(filenames[len(filenames)-1])
	if err != nil {
		t.Fatalf("Failed deleting from dbm_logs: %v", err)
	}

	result, error = db.ListAlreadyUp()
	if error != nil {
		t.Fatalf("Failed retrieving logs from dbm_logs: %v", error)
	}
	if len(result) != len(filenames)-1 {
		t.Fatalf("Inserted and Retrieved Data Mismatch: %v", error)
	}
}

func TestInsertAndListAndDeleteLogsFailed(t *testing.T) {
	filenames := []string{"file1", "file2", "file3", "file4", "file5"}

	err := db.InsertLogs(filenames)
	if err == nil {
		t.Fatalf("It should fail inserting but it is not")
	}

	_, error := db.ListAlreadyUp()
	if error == nil {
		t.Fatalf("It should fail listing but it is not")
	}

	_, error = db.GetLastLog()
	if error == nil {
		t.Fatalf("It should fail retrieving last log but it is not")
	}

	err = db.DeleteLog("random string")
	if err == nil {
		t.Fatalf("It should fail deleting but it is not")
	}
}

func TestBlindExec(t *testing.T) {
	err := db.BlindExec("SELECT 1+1")
	if err != nil {
		t.Fatalf("Failed Executing Blindly: %v", err)
	}
}

func TestBlindExecFailed(t *testing.T) {
	err := db.BlindExec("SELECT 1+'halo'")
	if err == nil {
		t.Fatalf("BlindExec should fail but it is not")
	}
}
