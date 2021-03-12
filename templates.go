package embed

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
)

type Template struct {
	Name     string
	Title    string
	Data     string
	Generate bool
}

type Templates map[string]Template

func templates(root string, paths map[string]string, skip []string) (Templates, error) {
	templates := make(Templates)

	for path, _ := range paths {
		name := strings.TrimPrefix(path, root)
		name = strings.TrimPrefix(name, "/")
		name = strings.TrimSuffix(name, ".gohtml")
		name = strings.ToLower(name)
		name = camelCase(name)

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		data, err := bytesToString(bytes)
		if err != nil {
			return nil, err
		}

		template := Template{
			Name:     name,
			Title:    strings.Title(name),
			Data:     data,
			Generate: true,
		}

		for _, s := range skip {
			if template.Title == s {
				template.Generate = false
			}
		}

		templates[name] = template
	}

	return templates, nil
}

func bytesToString(b []byte) (string, error) {
	builder := strings.Builder{}

	for _, v := range b {
		_, err := builder.WriteString(fmt.Sprintf("%d,", int(v)))
		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}

func camelCase(s string) string {
	concat := func(words []string) string {
		builder := strings.Builder{}

		for i, w := range words {
			if i > 0 {
				w = strings.Title(w)
			}
			builder.WriteString(w)
		}

		return builder.String()
	}

	specialChars := ""

	for _, r := range s {
		if unicode.IsLetter(r) == false && unicode.IsNumber(r) == false {
			specialChars += string(r)
		}
	}

	for _, c := range specialChars {
		if strings.Contains(s, string(c)) {
			s = concat(strings.Split(s, string(c)))
		}
	}

	return s
}
