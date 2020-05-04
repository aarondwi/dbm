package main

import "dbm/schema"

// DummyDB is used for mocking DbAccess interface only
type DummyDB struct{}

func (d *DummyDB) Init(conf schema.Conf) {}
func (d *DummyDB) BlindExec(stmt string) error {
	return nil
}
func (d *DummyDB) CreateLogTable() error {
	return nil
}
func (d *DummyDB) DropLogTable() error {
	return nil
}
func (d *DummyDB) InsertLogs(filenames []string) error {
	return nil
}
func (d *DummyDB) DeleteLog(filename string) error {
	return nil
}
func (d *DummyDB) GetLastLog() (string, error) {
	return "LastLog", nil
}
func (d *DummyDB) ListAlreadyUp() ([]string, error) {
	return []string{"FirstLog", "MiddleLog", "LastLog"}, nil
}
func (d *DummyDB) Close() {}
