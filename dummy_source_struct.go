package main

import (
	"dbm/schema"
	"errors"
)

type dummySource struct{}

func (d *dummySource) GenerateDirectory(dirname string) error { return nil }
func (d *dummySource) GenerateSrcfile(filename string) error  { return nil }
func (d *dummySource) ReadSrcfileContent(filename string) (*schema.Srcfile, error) {
	return &schema.Srcfile{
		Up:   "hello",
		Down: "World",
	}, nil
}
func (d *dummySource) ReadFromSrcDir() ([]string, error) {
	return []string{"suatufile"}, nil
}

type dummySourceFail struct{}

func (d *dummySourceFail) GenerateDirectory(dirname string) error {
	return errors.New("some error for test")
}
func (d *dummySourceFail) GenerateSrcfile(filename string) error {
	return errors.New("some error for test")
}
func (d *dummySourceFail) ReadSrcfileContent(filename string) (*schema.Srcfile, error) {
	return nil, errors.New("some error for test")
}
func (d *dummySourceFail) ReadFromSrcDir() ([]string, error) {
	return nil, errors.New("some error for test")
}
