package embed

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const DefaultInput = "/web/templates"
const DefaultOutput = "/internal/template/templates.go"

func Generate(root, input, output string, skip []string) {
	check := func(err error) {
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}

	abs, err := filepath.Abs(root)
	check(err)

	input = path.Join(abs, input)
	output = path.Join(abs, output)

	if _, err := os.Stat(input); os.IsNotExist(err) {
		check(err)
	}
	fmt.Printf("Input\t\t`%s`\n", input)

	if _, err := os.Stat(output); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(output), os.ModePerm)
		check(err)
	}
	fmt.Printf("Output\t\t`%s`\n", output)

	paths, err := filePaths(input)
	check(err)

	templates, err := templates(input, paths, skip)
	check(err)
	for _, t := range templates {
		fmt.Printf("Template\tGenerate=%t\t%s\n", t.Generate, t.Title)
	}

	outputFile, err := os.Create(output)
	check(err)

	defer func() {
		err = outputFile.Close()
		check(err)
	}()

	builder := &bytes.Buffer{}

	err = tmpl.Execute(builder, templates)
	check(err)

	data, err := format.Source(builder.Bytes())
	check(err)

	err = ioutil.WriteFile(output, data, os.ModePerm)
	check(err)

	fmt.Println("Done")
}
