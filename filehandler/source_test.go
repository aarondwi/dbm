package filehandler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aarondwi/dbm/schema"

	"gopkg.in/yaml.v2"
)

var source Source

func TestGenerateDirectory(t *testing.T) {
	dirname := "test"
	source.GenerateDirectory(dirname)
	defer os.RemoveAll(dirname)

	_, err := os.Stat(filepath.Join(dirname, "/src"))
	if err != nil {
		t.Fatalf(`Failed Creating "src" directory: %v`, err)
	}

	_, err = os.Stat(filepath.Join(dirname, "/conf.yaml"))
	if err != nil {
		t.Fatalf(`Failed Creating "conf.yaml": %v`, err)
	}
}

func TestGenerateSrcfile(t *testing.T) {
	dirname := "src"
	os.Mkdir(dirname, 'd')
	defer os.RemoveAll(dirname)

	filename := "CreateTableDummy"
	source.GenerateSrcfile(filename)

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

func TestGenerateSrcfileFail(t *testing.T) {
	filename := "CreateTableDummy"
	err := source.GenerateSrcfile(filename)
	if err == nil {
		t.Fatalf("should fail generating src file but it is not")
	}
}

func TestReadSrcfileContent(t *testing.T) {
	dirname := "src"
	os.Mkdir(dirname, 'd')
	defer os.RemoveAll(dirname)

	s := &schema.Srcfile{
		Up:   "hello",
		Down: "World",
	}
	d, _ := yaml.Marshal(&s)
	filename := fmt.Sprintf("%d-%s.yaml",
		int32(time.Now().Unix()), "CreateTableDummy")

	ioutil.WriteFile(filepath.Join("src", filename),
		[]byte(string(d)), 0700)

	result := &schema.Srcfile{}
	result, err := source.ReadSrcfileContent(filename)
	if err != nil {
		t.Fatalf("Failed generating srcfile content: %v", err)
	}
	if result.Up != s.Up || result.Down != s.Down {
		t.Fatalf(
			"The Content of Up and Down are different: \n"+
				"Expected Up: %s"+
				"Received Up: %s"+
				"Expected Up: %s"+
				"Received Up: %s"+
				"Error: %v", s.Up, result.Up, s.Down, result.Down, err)
	}
}

func TestReadSrcfileContentFail(t *testing.T) {
	dirname := "src"
	os.Mkdir(dirname, 'd')
	defer os.RemoveAll(dirname)
	ioutil.WriteFile(filepath.Join("src", "dummyfile"),
		[]byte(string("not a yaml file")), 0700)

	_, err := source.ReadSrcfileContent("nonexistentfile.txt")
	if err == nil {
		t.Fatalf("should fail because file not exist, but it is not")
	}
	_, err = source.ReadSrcfileContent("dummyfile")
	if err == nil {
		t.Fatalf("should fail because not a yaml file")
	}
}

func TestReadConfigfileContent(t *testing.T) {
	s := &schema.Conf{
		Dialect:  "test",
		Host:     "test",
		Port:     1000,
		Username: "test",
		Password: "test",
		Database: "test",
		Sslmode:  "test"}
	d, _ := yaml.Marshal(&s)
	ioutil.WriteFile("conf.yaml", []byte(string(d)), 0700)
	defer os.Remove("conf.yaml")

	result := &schema.Conf{}
	result, err := source.ReadConfigFile()
	if err != nil {
		t.Fatalf("Failed reading config: %v", err)
	}
	if result.Dialect != "test" ||
		result.Host != "test" ||
		result.Port != 1000 ||
		result.Username != "test" ||
		result.Password != "test" ||
		result.Database != "test" ||
		result.Sslmode != "test" {
		t.Fatalf(
			"The Content is different: \n"+
				"Expected: \n%v\n"+
				"Received: \n%v\n", s, result)
	}
}

func TestReadConfigfileContentFail(t *testing.T) {
	_, err := source.ReadConfigFile()
	if err == nil {
		t.Fatalf("should fail because file not exist, but it is not")
	}

	ioutil.WriteFile("conf.yaml", []byte(string("<html/>")), 0700)
	defer os.Remove("conf.yaml")
	_, err = source.ReadConfigFile()
	if err == nil {
		t.Fatalf("should fail because not a yaml file")
	}
}

func TestReadFromSrcDir(t *testing.T) {
	dirname := "src"
	os.Mkdir(dirname, 'd')
	defer os.RemoveAll(dirname)

	s := &schema.Srcfile{
		Up:   "hello",
		Down: "World",
	}
	d, _ := yaml.Marshal(&s)
	filename := fmt.Sprintf("%d-%s.yaml",
		int32(time.Now().Unix()), "CreateTableDummy")

	err := ioutil.WriteFile(filepath.Join("src", filename),
		[]byte(string(d)), 0700)
	if err != nil {
		t.Fatalf("Failed generating mock src file: %v", err)
	}

	var result []string
	result, err = source.ReadFromSrcDir()
	if err != nil {
		t.Fatalf("Failed reading from src directory: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("Different number of files received")
	}
}

func TestReadFromSrcDirFail(t *testing.T) {
	_, err := source.ReadFromSrcDir()
	if err == nil {
		t.Fatalf("should fail reading from src dir but it is not")
	}
}
