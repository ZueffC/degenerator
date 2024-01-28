package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/cristalhq/acmd"
	cmd "github.com/zueffc/degenerator/commands"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func loggerInit() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Can't create config.ini file. Please check read/write permissions!")
	}

	Info = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	Warn = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	Error = log.New(file, "ERROR: ", log.Ldate|log.Ltime)

	log.SetOutput(file)
}

func main() {
	loggerInit()
	parseINI()

	commands := []acmd.Command{
		{
			Name:        "list",
			Description: "get names of all .md files from a specific directory",
			ExecFunc: func(ctx context.Context, args []string) error {
				path, _ := filepath.Abs(getPostsPath())
				cmd.ListCommand(path)
				return nil
			},
		},
	}

	cmdRunner := acmd.RunnerOf(commands, acmd.Config{
		AppName:        "Degenerator",
		AppDescription: "Degenerator is a tool that provides foolish ability to convert markdown to html and vice resa",
		Version:        "0.0.1",
	})

	if err := cmdRunner.Run(); err != nil {
		log.Fatalln("Application was crashed during the unresolvable error: " + err.Error())
		cmdRunner.Exit(err)
	} else {
		Info.Println("The application was successfully started")
	}
}
