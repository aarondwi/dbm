package main

import (
	"dbm/connector"
	"dbm/filehandler"
	"dbm/schema"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Init generates directory, and all its necessary files
// for dbm to use later
func Init(sf filehandler.SourceFormat, dirname string) error {
	err := sf.GenerateDirectory(dirname)
	if err != nil {
		log.Printf("Failed generating dbm directory: %v", err)
		return err
	}

	log.Println("Successfully generate dbm directory")
	return nil
}

// GenerateSrcfile creates a yaml file, that includes generated Up and Down attributes
// Expected to be called on the same level as "src/" folder
func GenerateSrcfile(sf filehandler.SourceFormat, filename string) error {
	err := sf.GenerateSrcfile(filename)
	if err != nil {
		log.Printf("Failed Generating src file: %v", err)
		return err
	}

	return nil
}

// Setup initiate logs table/store
// for dbm to track what has been up-ed and what has not
func Setup(db connector.DbAccess) error {
	err := db.CreateLogTable()
	if err != nil {
		log.Printf("Failed generating log table: %v", err)
		return err
	}
	return nil
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
func Status(sf filehandler.SourceFormat, db connector.DbAccess) error {
	files, err := sf.ReadFromSrcDir()
	if err != nil {
		log.Printf("Failed reading src directory: %v", err)
		return err
	}

	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		log.Printf("Failed retrieving logs from db: %v", err)
		return err
	}

	for _, f := range files {
		status := "up"
		if !stringInSlice(f, alreadyUp) {
			status = "down"
		}
		log.Printf("%s : %s", f, status)
	}
	return nil
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
		log.Printf("Failed Unmarshalling %s: %v", filename, err)
		return nil, err
	}

	return s, nil
}

// Up applies one or more additions to database
// if want to apply all notYetUp, pass empty string to filename
func Up(db connector.DbAccess, filename string) {
	files, err := ioutil.ReadDir("src")
	if err != nil {
		log.Printf("Failed reading src directory: %v", err)
		return
	}

	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		log.Printf("Failed retrieving logs from db: %v", err)
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
		log.Printf("File not found: %v", err)
	}
}

// Down remove one addition to database
// different from Up, Down without empty filename only deletes one latest from logs
func Down(db connector.DbAccess, filename string) {
	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		log.Printf("Failed retrieving logs from db: %v", err)
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
		log.Printf("File not found in logs: %v", err)
	}
}
