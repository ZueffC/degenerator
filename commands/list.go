package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ListCommand(path_t string) {
	path, _ := filepath.Abs(path_t)
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Panicln("Can't find specified path: " + err.Error())
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		panic(err.Error())
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}
}
