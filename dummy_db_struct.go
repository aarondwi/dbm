package main

import (
	"dbm/schema"
	"errors"
)

// dummyDB is used for mocking DbAccess interface only
type dummyDB struct{}

func (d *dummyDB) Init(conf schema.Conf)       {}
func (d *dummyDB) BlindExec(stmt string) error { return nil }
func (d *dummyDB) CreateLogTable() error       { return nil }
func (d *dummyDB) DropLogTable() error         { return nil }
func (d *dummyDB) InsertLogs(filenames []string) error {
	return nil
}
func (d *dummyDB) DeleteLog(filename string) error { return nil }
func (d *dummyDB) GetLastLog() (string, error)     { return "somefile", nil }
func (d *dummyDB) ListAlreadyUp() ([]string, error) {
	return []string{"somefile"}, nil
}
func (d *dummyDB) Close() {}

type dummyDBFail struct{}

func (d *dummyDBFail) Init(conf schema.Conf) {}
func (d *dummyDBFail) BlindExec(stmt string) error {
	return errors.New("some errors for test")
}
func (d *dummyDBFail) CreateLogTable() error {
	return errors.New("some errors for test")
}
func (d *dummyDBFail) DropLogTable() error {
	return errors.New("some errors for test")
}
func (d *dummyDBFail) InsertLogs(filenames []string) error {
	return errors.New("some errors for test")
}
func (d *dummyDBFail) DeleteLog(filename string) error {
	return errors.New("some errors for test")
}
func (d *dummyDBFail) GetLastLog() (string, error) {
	return "", errors.New("some errors for test")
}
func (d *dummyDBFail) ListAlreadyUp() ([]string, error) {
	return nil, errors.New("some errors for test")
}
func (d *dummyDBFail) Close() {}
