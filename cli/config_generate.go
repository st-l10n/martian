// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/st-10n/martian/assets"
)

func main() {
	err := vfsgen.Generate(assets.Config, vfsgen.Options{
		PackageName:  "cli",
		BuildTags:    "!dev",
		VariableName: "Config",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
