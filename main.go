package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"

	"github.com/chespinoza/goliscan/commands"
	"github.com/chespinoza/goliscan/config"
)

func main() {
	var err error

	defer func() {
		r := recover()
		if r != nil {
			err, ok := r.(*commands.AppError)
			if ok {
				fmt.Println("Exiting with error:\n\t", err.Error.Error())
				os.Exit(1)
			}
			panic(r)
		}
	}()

	version := config.GetVersion()
	cli.VersionPrinter = version.Printer

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "License scanner"
	app.Version = version.ShortInfo()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Tomasz Maczukin",
			Email: "tomasz@maczukin.pl",
		},
	}

	app.Commands = commands.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		fmt.Printf("Command '%s' not found\n", command)
	}

	err = app.Run(os.Args)
	if err != nil {
		commands.ThrowError(err)
	}
}
