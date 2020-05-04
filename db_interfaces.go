package dbm

// DbAccess is the main interface for our database
// to be properly used by `core`, a database connector should implement all of these below
type DbAccess interface {
	Init(conf Conf)
	BlindExec(stmt string) error // meaning we trust our users to make it right
	CreateLogTable() error
	DropLogTable() error // mostly for testing, or if used as a library
	InsertLogs(filenames []string) error
	GetLastLog() (string, error)
	DeleteLog(filename string) error
	ListAlreadyUp() (filename []string, err error)
	Close()
}
