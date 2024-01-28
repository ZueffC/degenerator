package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ListCommand(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Panicln("Can't find specified path: " + err.Error())
	}

	entries, err := os.ReadDir(absPath)

	for _, e := range entries {
		fmt.Println(e.Name())
	}
}
