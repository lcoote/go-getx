package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/desal/dsutil"
	"github.com/desal/go-getx/getx"
	"github.com/desal/gocmd"
	"github.com/desal/richtext"
	"github.com/jawher/mow.cli"
)

//Example config:
//a/hats=http://server/special_repos/hats.git
//a/([^/]+)=http://server/repos/$1.git
//b/([^/]+)=http://other/repos/$1.git

func main() {
	app := cli.App("go-getx", "go get extended")
	app.Spec = "[-d] [-v] [-i] [-f | -u] [-t] [--goflags] [PKG...]"

	var (
		dependencies = app.BoolOpt("d deps-only", false, "Do not fetch named packages, only their dependencies")
		verbose      = app.BoolOpt("v verbose", false, "Verbose output")
		veryverbose  = app.BoolOpt("vv veryverbose", false, "Very Verbose output (outputs executed commansd)")
		install      = app.BoolOpt("i install", false, "Install all fetched packages (will continue if package fails to compile)")
		fetch        = app.BoolOpt("f fetch-missing", false, "Performs a deep search for any missing dependencies and fetches them")
		update       = app.BoolOpt("u update", false, "Updates package, and all transisitive depnediencs where possible")
		tests        = app.BoolOpt("t tests", false, "Fetches tests for the named packages")
		buildFlags   = app.StringOpt("goflags", "", "Additional flags to parse to go install (e.g. '-tags netgo')")

		pkgs = app.StringsArg("PKG", nil, "Packages")
	)

	app.Action = func() {
		if len(*pkgs) == 0 {
			app.PrintHelp()
			os.Exit(0)
		}

		ruleSet, err := getx.LoadRulesFromFile(filepath.Join(dsutil.UserHomeDir(), ".go-getx-map"))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		format := richtext.New()
		flags := []getx.Flag{}
		goFlags := []gocmd.Flag{}

		if *fetch {
			flags = append(flags, getx.DeepScan)
		} else if *update {
			flags = append(flags, getx.Update)
		}

		if *install {
			flags = append(flags, getx.Install)
		}

		if *verbose {
			flags = append(flags, getx.Verbose)
		} else if *veryverbose {
			flags = append(flags, getx.Verbose)
			flags = append(flags, getx.CmdVerbose)
			goFlags = append(goFlags, gocmd.Verbose)
		}

		goPath, err := gocmd.EnvGoPath()
		if err != nil {
			format.ErrorLine("%s", err)
			os.Exit(1)
		}

		ctx := getx.New(format, goPath, ruleSet, *buildFlags, flags...)
		for _, pkg := range *pkgs {
			ctx.Get(".", pkg, *dependencies, *tests)
		}

		if *install {
			goCtx := gocmd.New(format, goPath, *buildFlags, goFlags...)
			ok := true
			for _, pkg := range *pkgs {
				err := goCtx.Install(".", pkg)
				if err != nil {
					ok = false
					format.ErrorLine("Failed to install %s: %s", pkg, err.Error())
				}
			}

			if !ok {
				os.Exit(1)
			}
		}
	}

	app.Run(os.Args)
}
