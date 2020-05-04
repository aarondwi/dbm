package dbm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// Srcfile is the format of the source file
// in which the user will fill with `up` and `down` quereis
type Srcfile struct {
	Up   string `yaml:"up"`
	Down string `yaml:"down"`
}

// Source is the one and only implementation
// for SourceFormat interface
// It fully handles the handling to OS R/W operation
type Source struct{}

// GenerateDirectory creates a directory with given `dirname`
// includes the src subfolder in it
// and generates dummy conf.yaml file
func (source *Source) GenerateDirectory(dirname string) error {
	err := os.MkdirAll(filepath.Join(dirname, "src"), 'd')
	if err != nil {
		log.Fatalf("Failed generating src directory: %v", err)
		return err
	}

	c := &Conf{
		Dialect:  "postgresql/mysql/mariadb",
		Host:     "Host of your database",
		Port:     5432,
		Database: "Database to be written to",
		Username: "username to use",
		Password: "Password of the username",
		Sslmode:  "Whether to use ssl",
	}
	d, _ := yaml.Marshal(&c)

	err = ioutil.WriteFile(
		filepath.Join(dirname, "conf.yaml"),
		[]byte(string(d)), 0700)
	if err != nil {
		log.Fatalf("Failed creating conf.yaml file: %v", err)
		return err
	}

	return nil
}

// ReadFromSrcDir returns all sourcefile names
// in the src directory
func (source *Source) ReadFromSrcDir() ([]string, error) {
	files, err := ioutil.ReadDir("src")
	if err != nil {
		log.Fatalf("Failed reading from src directory: %v", err)
		return nil, err
	}

	var result []string
	for _, f := range files {
		result = append(result, f.Name())
	}

	return result, nil
}

// ReadSrcfileContent reads the yaml file specified by filename
// umarshal the yaml and returns the Srcfile structs
// taken from https://stackoverflow.com/questions/30947534/how-to-read-a-yaml-file
func (source *Source) ReadSrcfileContent(filename string) (*Srcfile, error) {
	yamlFile, err := ioutil.ReadFile(filepath.Join("src", filename))
	if err != nil {
		log.Printf("Failed read %s: %v", filename, err)
		return nil, err
	}

	s := &Srcfile{}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		log.Fatalf("Failed Unmarshalling %s: %v", filename, err)
		return nil, err
	}

	return s, nil
}

// GenerateSrcfile creates skeleton file
// that user can fill later with `up` and `down` queries
func (source *Source) GenerateSrcfile(filename string) error {
	s := &Srcfile{
		Up:   "Add feature, such as table, index, etc",
		Down: `To retract the result of "Up"`,
	}
	d, _ := yaml.Marshal(&s)

	filename = fmt.Sprintf("%d-%s.yaml",
		int32(time.Now().Unix()), filename)

	err := ioutil.WriteFile(filepath.Join("src", filename),
		[]byte(string(d)), 0700)
	if err != nil {
		log.Fatalf("Failed creating conf.yaml file: %v", err)
		return err
	}
	return nil
}
