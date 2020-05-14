package main

import (
	"os"

	"github.com/aarondwi/dbm/connector"
	"github.com/aarondwi/dbm/filehandler"
	cli "github.com/jawher/mow.cli"
)

var sf = &filehandler.Source{}
var pgdb = &connector.DbPostgres{}

func main() {
	app := cli.App("dbm", "ORM-style database-schema migration tools, minus the ORM bloats")
	app.Version("version", "0.1.0")

	app.Command("init", "generate dbm directory", func(cmd *cli.Cmd) {
		dirname := cmd.StringArg("DIRNAME", "", "directory name to be generated")
		cmd.Action = func() {
			Init(sf, *dirname)
		}
	})

	app.Command("generate", "generate yaml file inside src directory", func(cmd *cli.Cmd) {
		filename := cmd.StringArg("FILENAME", "", "directory name to be generated")
		cmd.Action = func() {
			GenerateSrcfile(sf, *filename)
		}
	})

	app.Command("setup", "create log table in the chosen dbms", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conf, _ := ReadConfigFile(sf)
			err := pgdb.Init(*conf)
			if err == nil {
				Setup(pgdb)
			}
		}
	})

	app.Command("status", "list status for all src file", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conf, _ := ReadConfigFile(sf)
			err := pgdb.Init(*conf)
			if err == nil {
				Status(sf, pgdb)
			}
		}
	})

	app.Command("up", "apply src file to database", func(cmd *cli.Cmd) {
		filename := cmd.StringOpt("FILENAME", "", "file to be applied")
		cmd.Action = func() {
			conf, _ := ReadConfigFile(sf)
			err := pgdb.Init(*conf)
			if err == nil {
				Up(sf, pgdb, *filename)
			}
		}
	})

	app.Command("down", "un-apply src file from database", func(cmd *cli.Cmd) {
		filename := cmd.StringOpt("FILENAME", "", "file to be un-applied")
		cmd.Action = func() {
			conf, _ := ReadConfigFile(sf)
			err := pgdb.Init(*conf)
			if err == nil {
				Down(sf, pgdb, *filename)
			}
		}
	})

	app.Run(os.Args)
}
