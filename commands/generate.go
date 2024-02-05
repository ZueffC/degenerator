package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	cp "github.com/otiai10/copy"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type FileStruct struct {
	Meta     map[string]interface{}
	Body     string
	BlogName string
}

type IndexPage struct {
	BlogName string
	Posts    []string
}

var index IndexPage

// Data getters
func getFileData(filepath string) (*FileStruct, error) {
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

func getFileName(file *FileStruct) string {
	title := fmt.Sprint(file.Meta["Title"])
	title = strings.ReplaceAll(strings.ToLower(title), " ", "-") //

	return title
}

func getPostPartial(templateName string, file *FileStruct) (string, error) {
	var postBuff bytes.Buffer

	if len(templateName) == 0 || templateName == "default" {
		templateName = "default"
	}

	file.Meta["Link"] = getFileName(file) + ".html"

	path, _ := filepath.Abs("templates/" + templateName + "/partials/post-card.html")
	partialBytes, _ := os.ReadFile(path)
	partialTemplate, _ := template.New("post-card").Parse(string(partialBytes))
	partialTemplate.Execute(&postBuff, file)

	return postBuff.String(), nil
}

// Section where we make some HTML
func makeIndexHTML(templateName string, index IndexPage, blog_name string) error {
	var indexBuff bytes.Buffer
	index.BlogName = blog_name

	if len(templateName) == 0 || templateName == "default" {
		templateName = "default"
	}

	headerPath, _ := filepath.Abs("templates/" + templateName + "/partials/header.html")
	footerPath, _ := filepath.Abs("templates/" + templateName + "/partials/footer.html")
	indexPath, _ := filepath.Abs("templates/" + templateName + "/index.html")

	indexBytes, _ := os.ReadFile(indexPath)
	headerBytes, _ := os.ReadFile(headerPath)
	footerBytes, _ := os.ReadFile(footerPath)

	indexPage := string(indexBytes)
	indexPage = strings.Replace(indexPage, "{{.Header}}", string(headerBytes), 1)
	indexPage = strings.Replace(indexPage, "{{.Footer}}", string(footerBytes), 1)
	indexPage = strings.ReplaceAll(indexPage, "{{.BlogName}}", blog_name)

	partialTemplate, _ := template.New("index").Parse(indexPage)
	partialTemplate.Execute(&indexBuff, index)

	os.RemoveAll("release/index.html")
	os.Mkdir("release", fs.FileMode(os.O_CREATE))
	os.WriteFile("release/index.html", indexBuff.Bytes(), fs.FileMode(os.O_RDWR))

	return nil
}

func makePostHTML(templateName string, data *FileStruct, blog_name string) error {
	var postBuffer bytes.Buffer

	if len(templateName) == 0 || templateName == "default" {
		templateName = "default"
	}

	headerPath, _ := filepath.Abs("templates/" + templateName + "/partials/header.html")
	footerPath, _ := filepath.Abs("templates/" + templateName + "/partials/footer.html")
	postPath, _ := filepath.Abs("templates/" + templateName + "/post.html")

	postBytes, _ := os.ReadFile(postPath)
	headerBytes, _ := os.ReadFile(headerPath)
	footerBytes, _ := os.ReadFile(footerPath)

	postPage := string(postBytes)
	postPage = strings.Replace(postPage, "{{.Header}}", string(headerBytes), 1)
	postPage = strings.Replace(postPage, "{{.Footer}}", string(footerBytes), 1)
	postPage = strings.ReplaceAll(postPage, "{{.BlogName}}", blog_name)

	partialTemplate, _ := template.New("index").Parse(postPage)
	partialTemplate.Execute(&postBuffer, data)

	filename := getFileName(data) + ".html"
	os.MkdirAll("release/posts", os.ModePerm)
	_ = os.WriteFile("release/posts/"+filename, postBuffer.Bytes(), fs.FileMode(os.O_RDWR))

	cp.Copy("templates/"+templateName+"/static", "release/static")
	cp.Copy("templates/"+templateName+"/about.html", "release/about.html")
	return nil
}

// Global avaliable function can be called from other modules
func GenerateCommand(path_t string, blog_name string) error {
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
		data, error := getFileData(absPath + "/" + string(e.Name()))
		if error != nil {
			panic(error.Error())
		}

		postCardHTML, _ := getPostPartial("", data)
		index.Posts = append(index.Posts, postCardHTML)
		makePostHTML("", data, blog_name)
	}

	_ = makeIndexHTML("", index, blog_name)

	return nil
}
