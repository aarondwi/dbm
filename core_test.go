package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var dbMock DummyDB

func TestInit(t *testing.T) {
	dirname := "test"
	os.Mkdir(dirname, 'd')
	defer os.RemoveAll(dirname)
	Init(dirname)

	_, err := os.Stat(filepath.Join(dirname, "/src"))
	if err != nil {
		t.Fatalf(`Failed Creating "src" directory: %v`, err)
	}

	_, err = os.Stat(filepath.Join(dirname, "/conf.yaml"))
	if err != nil {
		t.Fatalf(`Failed Creating "conf.yaml": %v`, err)
	}
}

func TestGenerate(t *testing.T) {
	dirname := "src"
	os.Mkdir(dirname, 'd')
	defer os.RemoveAll(dirname)

	filename := "CreateTableDummy"
	Generate(filename)

	files, err := ioutil.ReadDir("src")
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, f := range files {
		if strings.Contains(f.Name(), filename) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Failed generating srcfile: %v", err)
	}
}

func TestUp(t *testing.T) {

}

func TestDown(t *testing.T) {

}
