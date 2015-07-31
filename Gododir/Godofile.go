package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	. "gopkg.in/godo.v1"
)

var mime = map[string]string{
	".css": "text/css",
	".js":  "application/javascript",
}

func main() {
	Godo(tasks)
}

func tasks(p *Project) {
	Env = `GOPATH=.vendor::$GOPATH`

	p.Task("default", D{"build"})

	p.Task("jsmin", func(c *Context) error {

		b, err := ioutil.ReadFile("templates/scripts.html")
		if err != nil {
			return err
		}

		reBlock := regexp.MustCompile(`<!-- build:scripts-->(?s)(.*)<!-- endbuild-->`)
		block := reBlock.FindStringSubmatch(string(b))

		if len(block) == 0 {
			return errors.New("Build block was not found")
		}

		reSources := regexp.MustCompile(`src="(.*)"`)
		files := reSources.FindAllStringSubmatch(block[1], -1)

		var filenames []string
		for _, f := range files {
			filenames = append(filenames, f[1])
		}

		err = os.MkdirAll("static/dist/js", os.ModePerm)
		if err != nil {
			return err
		}

		output := "static/dist/js/app.min.js"
		return doMinify(filenames, output)
	})

	p.Task("cssmin", func(c *Context) error {

		b, err := ioutil.ReadFile("templates/styles.html")
		if err != nil {
			return err
		}

		reBlock := regexp.MustCompile(`<!-- build:styles-->(?s)(.*)<!-- endbuild-->`)
		block := reBlock.FindStringSubmatch(string(b))

		if len(block) == 0 {
			return errors.New("Build block was not found")
		}

		reSources := regexp.MustCompile(`href="(.*)"`)
		files := reSources.FindAllStringSubmatch(block[1], -1)

		var filenames []string
		for _, f := range files {
			filenames = append(filenames, f[1])
		}

		err = os.MkdirAll("static/dist/css", os.ModePerm)
		if err != nil {
			return err
		}
		output := "static/dist/css/app.min.css"
		return doMinify(filenames, output)
	})

	p.Task("build", D{"jsmin", "cssmin"}, func() error {
		return Run("go build")
	}).Watch("**/*.go")

	p.Task("server", func() {
		Start("main.go")
	}).Watch("**/*.{go,html}")

	p.Task("views", func() error {
		return Run("razor templates")
	}).Watch("templates/**/*.go.html")

	p.Task("lint", func() {
		Run("golint .")
		Run("gofmt -w -s .")
		Run("go vet .")
	})
}

func doMinify(filenames []string, output string) error {
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)
	m.AddFunc("text/css", css.Minify)

	err := os.MkdirAll("static/dist/js", os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, input := range filenames {

		in, err := os.Open(input)
		if err != nil {
			return err
		}
		defer in.Close()
		ext := filepath.Ext(input)

		mediatype := mime[ext]

		if err := m.Minify(mediatype, out, in); err != nil {
			if err == minify.ErrNotExist {
				io.Copy(out, in)
			} else {
				return err
			}
		}
	}
	return nil
}
