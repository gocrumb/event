package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)

	fmt.Println(os.Args)

	types := flag.String("type", "", "")
	output := flag.String("output", "event.go", "")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, ".")
	}

	g := Generator{}
	g.types = strings.Split(*types, ",")

	var (
		pkg *Package
		err error
	)
	if len(args) == 1 && isDirectory(args[0]) {
		pkg, err = parsePackageDir(args[0])
	} else {
		pkg, err = parsePackageFiles(args)
	}
	if err != nil {
		log.Fatalf("parsing package: %s", err)
	}
	g.pkg = pkg

	src, err := g.generate()
	if err != nil {
		log.Fatalf("generating code: %s", err)
	}

	err = ioutil.WriteFile(filepath.Join(g.pkg.dir, *output), src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
