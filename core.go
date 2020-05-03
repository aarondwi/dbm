package dbm

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type conf struct {
	Dialect  string `yaml:"dialect"`
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// Generate directory, and all its necessary files
// for dbm to use later
func Init(dirname string) {
	err := os.MkdirAll(dirname+"/src", 'd')
	if err != nil {
		log.Fatalf("Failed generating src directory: %v", err)
	}

	c := conf{
		Dialect:  "postgres/mysql/mariadb/mongo",
		Url:      "connection url to your database",
		Database: "Database to be written to",
		Username: "username to use",
		Password: "Password of the username",
	}
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("Failed generating conf.yaml content: %v", err)
	}

	err = ioutil.WriteFile(dirname+"/conf.yaml", []byte(string(d)), 0700)
	if err != nil {
		log.Fatalf("Failed creating conf.yaml file: %v", err)
	}

	log.Println("Successfully generate dbm directory")
}
