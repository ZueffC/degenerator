package commands

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type FileStruct struct {
	Meta map[string]interface{}
	Body string
}

func getFileMetadata(filepath string) (*FileStruct, error) {
	var buf bytes.Buffer
	data := FileStruct{}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	file, err := os.ReadFile(filepath)
	if err != nil {
		println("Can't open file " + filepath + " \n Please, check files permission!")
		return nil, errors.New("Can't read file by provided path")
	}

	context := parser.NewContext()
	content := string(file)
	if err := markdown.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		panic(err)
		return nil, errors.New("Can't get data from markdown, please check file")
	}

	data.Meta = meta.Get(context)
	data.Body = buf.String()

	return &data, nil
}

func GenerateCommand(path_t string) error {
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
		temp, error := getFileMetadata(absPath + "/" + string(e.Name()))
		if error != nil {
			panic(error.Error())
		}

		fmt.Println(temp.Meta["Title"])
	}

	return nil
}
