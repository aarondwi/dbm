package main

import (
	"dbm/connector"
	"dbm/schema"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// Init generates directory, and all its necessary files
// for dbm to use later
func Init(dirname string) {
	err := os.MkdirAll(filepath.Join(dirname, "src"), 'd')
	if err != nil {
		log.Fatalf("Failed generating src directory: %v", err)
		return
	}

	c := &schema.Conf{
		Dialect:  "postgresql/mysql/mariadb",
		Host:     "Host of your database",
		Port:     5432,
		Database: "Database to be written to",
		Username: "username to use",
		Password: "Password of the username",
		Sslmode:  "Whether to use ssl",
	}
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("Failed generating conf.yaml content: %v", err)
		return
	}

	err = ioutil.WriteFile(
		filepath.Join(dirname, "conf.yaml"),
		[]byte(string(d)), 0700)
	if err != nil {
		log.Fatalf("Failed creating conf.yaml file: %v", err)
		return
	}

	log.Println("Successfully generate dbm directory")
}

// Generate creates a yaml file, that includes generated Up and Down attributes
// Expected to be called on the same level as "src/" folder
func Generate(filename string) {
	s := schema.Srcfile{
		Up:   "Add feature, such as table, index, etc",
		Down: `To retract the result of "Up"`,
	}
	d, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatalf("Failed generating srcfile %s: %v", filename, err)
		return
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

// Setup initiate logs table/store
// for dbm to track what has been up-ed and what has not
func Setup(db connector.DbAccess) {
	db.CreateLogTable()
}

// taken from
// https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python#15323988
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Status check the statuses of all files
// in src directory
func Status(db connector.DbAccess) {
	files, err := ioutil.ReadDir("src")
	if err != nil {
		log.Fatalf("Failed reading src directory: %v", err)
		return
	}

	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		log.Fatalf("Failed retrieving logs from db: %v", err)
		return
	}

	for _, f := range files {
		var status string
		if stringInSlice(f.Name(), alreadyUp) {
			status = "up"
		} else {
			status = "down"
		}
		log.Printf("%s : %s", f.Name(), status)
	}
}

// taken from
// https://stackoverflow.com/questions/30947534/how-to-read-a-yaml-file
func readYamlFileFromSrcDir(filename string) (*schema.Srcfile, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Failed read %s: %v", filename, err)
		return nil, err
	}

	s := &schema.Srcfile{}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		log.Fatalf("Failed Unmarshalling %s: %v", filename, err)
		return nil, err
	}

	return s, nil
}

// Up applies one or more additions to database
// if want to apply all notYetUp, pass empty string to filename
func Up(db connector.DbAccess, filename string) {
	files, err := ioutil.ReadDir("src")
	if err != nil {
		log.Fatalf("Failed reading src directory: %v", err)
		return
	}

	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		log.Fatalf("Failed retrieving logs from db: %v", err)
		return
	}

	var notYetUp []string
	for _, f := range files {
		if !stringInSlice(f.Name(), alreadyUp) {
			notYetUp = append(notYetUp, f.Name())
		}
	}

	if filename == "" {
		for _, nyu := range notYetUp {
			s, err := readYamlFileFromSrcDir(nyu)
			if err != nil {
				return
			}
			db.BlindExec(s.Up)
		}
		db.InsertLogs(notYetUp)
	} else if stringInSlice(filename, notYetUp) {
		s, err := readYamlFileFromSrcDir(filename)
		if err != nil {
			return
		}
		db.BlindExec(s.Up)
		db.InsertLogs([]string{filename})
	} else {
		log.Fatalf("File not found: %v", err)
	}
}

// Down remove one addition to database
// different from Up, Down without empty filename only deletes one latest from logs
func Down(db connector.DbAccess, filename string) {
	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		log.Fatalf("Failed retrieving logs from db: %v", err)
		return
	}

	if filename == "" {
		targetFilename, err := db.GetLastLog()
		s, err := readYamlFileFromSrcDir(targetFilename)
		if err != nil {
			return
		}
		db.BlindExec(s.Down)
		db.DeleteLog(filename)
	} else if stringInSlice(filename, alreadyUp) {
		s, err := readYamlFileFromSrcDir(filename)
		if err != nil {
			return
		}
		db.BlindExec(s.Down)
		db.DeleteLog(filename)
	} else {
		log.Fatalf("File not found in logs: %v", err)
	}
}
