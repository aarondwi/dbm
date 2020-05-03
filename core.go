package dbm

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type conf struct {
	Dialect  string `yaml:"dialect"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Sslmode string `yaml:"sslmode"`
}

// Generate directory, and all its necessary files
// for dbm to use later
func Init(dirname string) {
	err := os.MkdirAll(filepath.Join(dirname, "src"), 'd')
	if err != nil {
		log.Fatalf("Failed generating src directory: %v", err)
	}

	c := conf{
		Dialect:  "postgresql/mysql/mariadb",
		Host:      "Host of your database",
		Port:      "Port to use",
		Database: "Database to be written to",
		Username: "username to use",
		Password: "Password of the username",
		Sslmode: "Whether to use ssl",
	}
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("Failed generating conf.yaml content: %v", err)
	}

	err = ioutil.WriteFile(
		filepath.Join(dirname, "conf.yaml"),
		[]byte(string(d)), 0700)
	if err != nil {
		log.Fatalf("Failed creating conf.yaml file: %v", err)
	}

	log.Println("Successfully generate dbm directory")
}

type srcfile struct {
	Up   string `yaml:"up"`
	Down string `yaml:"down"`
}

// Creates a yaml file, that includes generated Up and Down attributes
// Expected to be called on the same level as "src/" folder
func Generate(filename string) {
	s := srcfile{
		Up:   "Add feature, such as table, index, etc",
		Down: `To retract the result of "Up"`,
	}
	d, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatalf("Failed generating srcfile %s: %v", filename, err)
	}

	filename = fmt.Sprintf(
		"%d-%s.yaml",
		int32(time.Now().Unix()),
		filename)

	err = ioutil.WriteFile(
		filepath.Join("src", filename),
		[]byte(string(d)), 0700)
	if err != nil {
		log.Fatalf("Failed creating conf.yaml file: %v", err)
	}
}
