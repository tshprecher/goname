package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"golang.org/x/tools/go/loader"
	"log"
	"os/exec"
	"strings"
	"unicode"
)

var (
	targets []renameTarget = make([]renameTarget, 0, 16)

	// flags
	list    bool
	verbose bool
)

func init() {
	flag.BoolVar(&list, "l", false, "list gorename targets instead of renaming")
	flag.BoolVar(&verbose, "v", false, "verbose")
}

// renameTarget represents a rename operation that can be fed into gorename
type renameTarget struct {
	from string
	to   string
}

// rename takes a variable name and returns the name properly formatted. If no
// change is necessary, the return value is the same as the input value.
func rename(name string) string {
	if name == "" {
		return ""
	}

	terms := strings.Split(name, "_")
	if len(terms) == 1 {
		return name
	}

	for t := range terms {
		rt := []rune(terms[t])
		startCap := true
		if t == 0 {
			startCap = unicode.IsUpper(rt[0])
		}
		for r := range rt {
			if r == 0 {
				rt[r] = unicode.ToUpper(rt[r])
				continue
			}
			rt[r] = unicode.ToLower(rt[r])
		}
		if t == 0 && !startCap {
			rt[0] = unicode.ToLower(rt[0])
		}
		terms[t] = string(rt)
	}
	return strings.Join(terms, "")
}

// loadProgram loads a list of packages in the GOPATH
func loadProgram(pkgs []string) (*loader.Program, error) {
	conf := loader.Config{
		Build: &build.Default,
	}
	for _, pkg := range pkgs {
		conf.ImportWithTests(pkg)
	}
	return conf.Load()
}

// handleFile finds the improperly formatted variable names and
// adds them to the global slice of rename targets
func handleFile(file *ast.File, pkgPrefix string) {
	sp := file.Scope
	for _, o := range sp.Objects {
		if o.Kind != ast.Con && o.Kind != ast.Var {
			continue
		}

		formatted := rename(o.Name)
		if formatted != o.Name {
			targets = append(targets, renameTarget{fmt.Sprintf("\"%s\".%s", pkgPrefix, o.Name), formatted})
		}
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	log.SetPrefix("")
	if len(args) == 0 {
		log.Fatal("missing package argument")
	}
	if len(args) > 1 {
		log.Fatal("only one package argument allowed")
	}

	pkgPrefix := args[0]
	prgm, err := loadProgram([]string{pkgPrefix})
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, pkginfo := range prgm.Imported {
		for _, astfile := range pkginfo.Files {
			handleFile(astfile, pkgPrefix)
		}
	}

	if list {
		for _, t := range targets {
			fmt.Printf("%s\t%s\n", t.from, t.to)
		}
	} else {
		cmd := exec.Command("gorename")
		err := cmd.Run()
		if err == exec.ErrNotFound {
			log.Fatalf("gorename not found in $PATH. Run 'go install golang.org/x/tools/cmd/gorename' to install.")
		}

		for _, t := range targets {
			if verbose {
				log.Printf("renaming target %s to %s", t.from, t.to)
			}
			cmd = exec.Command("gorename", "-from", t.from, "-to", t.to)
			err = cmd.Run()
			if err != nil {
				log.Fatalf("error replacing '%s': %s", t.from, err.Error())
			}
		}
	}
}
