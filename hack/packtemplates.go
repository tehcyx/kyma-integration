package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Reads all *.html files in the template folder
// and encodes them as strings literals in templates.go
func main() {
	tmplPath := flag.String("p", "./templates", "template path (Optional, default: './templates')")
	variablePrefix := flag.String("a", "TMPL", "prefix for variables in target file (Optional, default: 'TMPL')")

	flag.Parse()

	if *tmplPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *variablePrefix == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fs, errDir := ioutil.ReadDir(*tmplPath)
	if errDir != nil {
		fmt.Println(fmt.Errorf("error occured: %s", errDir))
	}
	out, errData := os.Create("./internal/tmpl/data.go")
	if errData != nil {
		fmt.Println(fmt.Errorf("error occured: %s", errData))
	}
	out.Write([]byte("// THIS FILE IS AUTO-GENERATED, DO NOT EDIT \npackage tmpl \n\nvar TMPLMap = map[string]string{\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".html") {
			out.Write([]byte(fmt.Sprintf("\"%s\": ", normalizeName(strings.TrimSuffix(f.Name(), ".html")))))
			f, errOpen := os.Open(fmt.Sprintf("%s/%s", *tmplPath, f.Name()))
			if errOpen != nil {
				fmt.Println(fmt.Errorf("error occured: %s", errOpen))
			}
			out.Write([]byte("`"))
			io.Copy(out, f)
			out.Write([]byte("`,\n"))
		}
	}
	out.Write([]byte("}\n"))
}

func normalizeName(name string) string {
	tmp := strings.ReplaceAll(name, "_", "")
	return strings.ReplaceAll(tmp, "-", "")
}