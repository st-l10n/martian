// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

// Config contains project configuration.
var cfg http.FileSystem = http.Dir("config")

func main() {
	err := vfsgen.Generate(cfg, vfsgen.Options{
		PackageName:  "cli",
		BuildTags:    "!dev",
		VariableName: "Config",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
