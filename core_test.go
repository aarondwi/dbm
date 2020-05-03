package dbm

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	dirname := "test"
	os.Mkdir(dirname, 'd')
	Init(dirname)

	_, err := os.Stat(dirname + "/src")
	if err != nil {
		t.Fatalf(`Failed Creating "src" directory`)
	}

	_, err = os.Stat(dirname + "/conf.yaml")
	if err != nil {
		t.Fatalf(`Failed Creating "conf.yaml"`)
	}

	os.RemoveAll(dirname)
}
