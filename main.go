package main

import "github.com/mightymatth/arcli/cmd"

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
