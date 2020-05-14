package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/aarondwi/dbm/connector"
	"github.com/aarondwi/dbm/filehandler"
	"github.com/aarondwi/dbm/schema"
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

// ReadConfigFile read and parse the conf.yaml file
func ReadConfigFile(sf filehandler.SourceFormat) (*schema.Conf, error) {
	r, err := sf.ReadConfigFile()
	if err != nil {
		log.Printf("Failed reading conf.yaml: %v", err)
		return nil, err
	}
	return r, nil
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
		fmt.Printf("%s : %s\n", f, status)
	}
	return nil
}

// Up applies one or more additions to database
// if want to apply all notYetUp, pass empty string to filename
func Up(sf filehandler.SourceFormat, db connector.DbAccess, filename string) error {
	files, err := sf.ReadFromSrcDir()
	if err != nil {
		return err
	}

	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		return err
	}

	var notYetUp []string
	for _, f := range files {
		if !stringInSlice(f, alreadyUp) {
			notYetUp = append(notYetUp, f)
		}
	}

	if filename == "" {
		for _, nyu := range notYetUp {
			s, err := sf.ReadSrcfileContent(nyu)
			if err != nil {
				return err
			}
			db.BlindExec(s.Up)
		}
		db.InsertLogs(notYetUp)
	} else if stringInSlice(filename, alreadyUp) {
		return errors.New("file is already applied")
	} else if stringInSlice(filename, notYetUp) {
		s, err := sf.ReadSrcfileContent(filename)
		if err != nil {
			return err
		}
		db.BlindExec(s.Up)
		db.InsertLogs([]string{filename})
	} else {
		return errors.New("file not found")
	}
	return nil
}

// Down remove one addition to database
// different from Up, Down without empty filename only deletes one latest from logs
func Down(sf filehandler.SourceFormat, db connector.DbAccess, filename string) error {
	alreadyUp, err := db.ListAlreadyUp()
	if err != nil {
		return err
	}

	if filename == "" {
		targetFilename, err := db.GetLastLog()
		s, err := sf.ReadSrcfileContent(targetFilename)
		if err != nil {
			return err
		}
		db.BlindExec(s.Down)
		db.DeleteLog(targetFilename)
	} else if stringInSlice(filename, alreadyUp) {
		s, err := sf.ReadSrcfileContent(filename)
		if err != nil {
			return err
		}
		db.BlindExec(s.Down)
		db.DeleteLog(filename)
	} else {
		return errors.New("file not found in dbm_logs")
	}
	return nil
}
