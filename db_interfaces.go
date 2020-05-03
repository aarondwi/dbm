package dbm

type DbAccess interface {
	Init()
	CreateLogTable() error
	DropLogTable() error // mostly for testing, or if used as a library
	InsertLogs(filenames []string) error
	DeleteLogs(filenames []string) error
	ListAlreadyUp() (filename []string, err error)
	Close()
}
