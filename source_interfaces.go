package dbm

// SourceFormat is the main interface
// for OS R/W operation
// split into interface mainly for testing purpose
type SourceFormat interface {
	GenerateDirectory(dirname string) error
	GenerateSrcfile(filename string) error
	ReadSrcfileContent(filename string) (*Srcfile, error)
	ReadFromSrcDir() ([]string, error)
}
