package main

import (
	"html/template"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	val := Config{
		Package: "pkg",
		Services: []ServiceInfo{
			{
				ServiceName:      "Abc",
				ServiceNameLower: "abc",
			},
		},
	}
	tp := template.Must(template.New("").Parse(outTemplate))
	err := tp.Execute(os.Stdout, val)
	checkErr(err)
}
