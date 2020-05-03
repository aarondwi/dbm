package dbm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	dirname := "test"
	os.Mkdir(dirname, 'd')
	defer os.RemoveAll(dirname)
	Init(dirname)

	_, err := os.Stat(filepath.Join(dirname, "/src"))
	if err != nil {
		t.Fatalf(`Failed Creating "src" directory`)
	}

	_, err = os.Stat(filepath.Join(dirname, "/conf.yaml"))
	if err != nil {
		t.Fatalf(`Failed Creating "conf.yaml"`)
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
		t.Fatalf("Failed generating srcfile")
	}
}
